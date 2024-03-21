package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
	"github.com/inconshreveable/go-update"
	zone "github.com/lrstanley/bubblezone"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/baseurl"
	"monks.co/piano-alone/data"
	"monks.co/piano-alone/game"
	"monks.co/piano-alone/songs"
)

type model struct {
	role    string
	baseURL baseurl.BaseURL
	inbox   <-chan *game.Message

	ws    *websocket.Conn
	state *game.State
	midi  []byte
	log   []*game.Message

	width  int
	height int

	modal    string
	quitting string

	menuSelectionIndex int
	contentInFocus     bool

	conductorButtonIndex int

	midiOutPorts     midi.OutPorts
	midiOutPortIndex int

	connection    string
	latestVersion string

	output string
}

type (
	msgTick               time.Time
	msgQuit               string
	msgDismissModal       string
	msgGotMIDIOutputPorts midi.OutPorts
	msgVersion            string
	msgGotWSClient        *websocket.Conn
	msgGotWSMessage       *game.Message
)

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.checkVersion,
		m.checkMIDIOutPorts,
		m.connect,
		m.tick,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case msgTick:
		return m, m.tick

	case msgGotWSClient:
		m.ws = msg
		return m, m.acceptMessage

	case msgGotWSMessage:
		m.log = append(m.log, msg)
		switch msg.Type {
		case game.MessageTypeInitialState:
			m.state = game.StateFromBytes(msg.Data)

		case game.MessageTypeBroadcastPhase:
			phase := game.PhaseFromBytes(msg.Data)
			m.state.Phase = phase

		case game.MessageTypeConductorConnected:
			m.state.ConductorIsConnected = true
		case game.MessageTypeConductorDisconnected:
			m.state.ConductorIsConnected = false

		case game.MessageTypeDisklavierConnected:
			m.state.DisklavierIsConnected = true
		case game.MessageTypeDisklavierDisconnected:
			m.state.DisklavierIsConnected = false

		case game.MessageTypeBroadcastConnectedPlayer:
			player := game.PlayerFromBytes(msg.Data)
			m.state.Players[player.Fingerprint] = player

		case game.MessageTypeBroadcastDisconnectedPlayer:
			fingerprint := string(msg.Data)
			m.state.Players[fingerprint].ConnectionState = game.ConnectionStateDisconnected

		case game.MessageTypeBroadcastControllerModal:
			m.modal = string(msg.Data)

		case game.MessageTypeBroadcastSubmittedTrack:
			fingerprint := string(msg.Data)
			if _, ok := m.state.Players[fingerprint]; !ok {
				m.state.Players[fingerprint] = &game.Player{}
			}
			m.state.Players[fingerprint].HasSubmitted = true

		case game.MessageTypeBroadcastCombinedTrack:
			m.midi = msg.Data
			return m, tea.Batch(
				m.playSong,
				m.acceptMessage,
			)
		}

		return m, m.acceptMessage

	case msgQuit:
		m.output = string(msg)
		return m, tea.Quit

	case msgGotMIDIOutputPorts:
		m.midiOutPorts = midi.OutPorts(msg)
		for i, p := range m.midiOutPorts {
			if p.String() == "Komplete Audio 6" {
				m.midiOutPortIndex = i
				break
			}
		}
		return m, nil

	case msgVersion:
		m.latestVersion = string(msg)
		return m, nil

	case tea.MouseMsg:
		if msg.Action != tea.MouseActionRelease {
			return m, nil
		}
		for i, label := range m.menu() {
			if zone.Get(label).InBounds(msg) {
				m.contentInFocus = false
				m.menuSelectionIndex = i
				return m, nil
			}
		}
		for i, port := range m.midiOutPorts {
			if zone.Get(port.String()).InBounds(msg) {
				m.contentInFocus = true
				m.midiOutPortIndex = i
				return m, nil
			}
		}

		if zone.Get("Test MIDI Output").InBounds(msg) {
			m.contentInFocus = true
			return m, m.testMIDI
		}

		if zone.Get("Advance").InBounds(msg) {
			return m, m.advance
		} else if zone.Get("Restart").InBounds(msg) {
			return m, m.restart
		}

		if zone.Get("Modal").InBounds(msg) {
			m.modal = ""
			return m, nil
		}

		if zone.Get("Menu").InBounds(msg) {
			m.contentInFocus = false
			return m, nil
		} else if zone.Get("Content").InBounds(msg) {
			m.contentInFocus = true
			return m, nil
		}

	case tea.KeyMsg:
		if m.quitting != "" {
			if msg.String() == m.quitting {
				return m, tea.Quit
			} else {
				m.quitting = ""
				return m, nil
			}
		}

		if m.modal != "" {
			m.modal = ""
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c":
			m.quitting = msg.String()
			return m, nil

		case "esc", "q":
			if m.contentInFocus {
				m.contentInFocus = false
				return m, nil
			} else {
				m.quitting = msg.String()
			}

		case "u", "U":
			if m.latestVersion != "" && data.CurrentVersion != m.latestVersion {
				return m, m.updateBinary
			}

		case "tab":
			m.contentInFocus = !m.contentInFocus
			return m, nil

		case "enter":
			if !m.contentInFocus {
				m.contentInFocus = true
				return m, nil
			}
			switch m.selectedMenuItem() {
			case "MIDI Output Test":
				return m, m.testMIDI
			case "Performance Status":
				switch m.selectedConductorButton() {
				case "Advance":
					return m, m.advance
				case "Restart":
					return m, m.restart
				}
			}

		case "k", "up":
			if !m.contentInFocus {
				m.menuSelectionIndex -= 1
				if m.menuSelectionIndex < 0 {
					m.menuSelectionIndex = len(m.menu()) - 1
				}
				return m, nil
			}
			switch m.selectedMenuItem() {
			case "MIDI Configuration":
				m.midiOutPortIndex -= 1
				if m.midiOutPortIndex < 0 {
					m.midiOutPortIndex = len(m.midiOutPorts) - 1
				}
				return m, nil
			case "Performance Status":
				buttonCount := len(m.conductorButtons())
				if buttonCount > 1 {
					m.conductorButtonIndex -= 1
					if m.conductorButtonIndex < 0 {
						m.conductorButtonIndex = buttonCount - 1
					}
				}
			}
		case "j", "down":
			if !m.contentInFocus {
				m.menuSelectionIndex += 1
				if m.menuSelectionIndex >= len(m.menu()) {
					m.menuSelectionIndex = 0
				}
				return m, nil
			}
			switch m.selectedMenuItem() {
			case "MIDI Configuration":
				m.midiOutPortIndex += 1
				if m.midiOutPortIndex >= len(m.midiOutPorts) {
					m.midiOutPortIndex = 0
				}
				return m, nil
			case "Performance Status":
				buttonCount := len(m.conductorButtons())
				if buttonCount > 1 {
					m.conductorButtonIndex += 1
					if m.conductorButtonIndex > buttonCount+1 {
						m.conductorButtonIndex = 0
					}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m model) menu() []string {
	switch m.role {
	case "conductor":
		return []string{
			"Performance Status",
			"Message Log",
		}
	default:
		return []string{
			"Performance Status",
			"MIDI Configuration",
			"MIDI Output Test",
			"Message Log",
		}
	}
}

func (m model) conductorButtons() []string {
	if m.role != "conductor" {
		return []string{}
	}
	switch m.state.Phase.Type {
	case game.GamePhaseUninitialized:
		return []string{"Advance"}
	case game.GamePhaseLobby:
		return []string{"Advance"}
	case game.GamePhaseHero:
		return []string{"Advance", "Restart"}
	case game.GamePhaseProcessing:
		return []string{}
	case game.GamePhasePlayback:
		return []string{}
	case game.GamePhaseDone:
		return []string{"Restart"}
	default:
		return []string{}
	}
}

func (m model) selectedConductorButton() string {
	buttons := m.conductorButtons()
	if len(buttons) == 0 {
		return ""
	}
	return buttons[m.conductorButtonIndex%len(buttons)]
}

func (m model) selectedMenuItem() string {
	return m.menu()[m.menuSelectionIndex]
}

func (m model) checkVersion() tea.Msg {
	resp, err := http.Get(m.baseURL.Rest(data.PathLatestClientVersion))
	if err != nil {
		return msgQuit(fmt.Errorf("error finding latest version: %s", err).Error())
	}
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		return msgQuit(fmt.Errorf("error finding latest version: %s", resp.Status).Error())
	}
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return msgQuit(fmt.Errorf("error finding latest version: %s", err).Error())
	}
	return msgVersion(bs)
}

func (m model) checkMIDIOutPorts() tea.Msg {
	ports := midi.GetOutPorts()
	return msgGotMIDIOutputPorts(ports)
}

func (m model) connect() tea.Msg {
	conn, _, err := websocket.DefaultDialer.Dial(m.baseURL.WS(data.PathControllerWS)+fmt.Sprintf("?role=%s", m.role), nil)
	if err != nil {
		return msgQuit(err.Error())
	}
	var msgType game.MessageType
	switch m.role {
	case "conductor":
		msgType = game.MessageTypeConductorConnected
	case "disklavier":
		msgType = game.MessageTypeDisklavierConnected
	default:
		panic("unknown role")
	}
	if err := conn.WriteMessage(
		websocket.BinaryMessage,
		game.NewMessage(msgType, m.role, nil).Bytes(),
	); err != nil {
		return msgQuit(err.Error())
	}
	return msgGotWSClient(conn)
}

func (m model) tick() tea.Msg {
	time.Sleep(time.Second / 10)
	return msgTick(time.Now())
}

func (m model) acceptMessage() tea.Msg {
	mt, bs, err := m.ws.ReadMessage()
	if err != nil {
		return msgQuit(fmt.Sprintf("msg read error: %s", err.Error()))
	}
	switch mt {
	case websocket.BinaryMessage:
		msg := game.MessageFromBytes(bs)
		return msgGotWSMessage(msg)
	default:
		return msgQuit("unexpected message type")
	}
}

func (m model) advance() tea.Msg {
	http.Post(m.baseURL.Rest(data.PathAdvance), "", nil)
	return nil
}

func (m model) restart() tea.Msg {
	http.Post(m.baseURL.Rest(data.PathRestart), "", nil)
	return nil
}

func (m model) testMIDI() tea.Msg {
	return m.playMIDI(songs.ExcerptBytes)
}

func (m model) playSong() tea.Msg {
	return m.playMIDI(m.midi)
}

func (m model) playMIDI(bs []byte) tea.Msg {
	outPort := m.midiOutPorts[m.midiOutPortIndex].String()
	out, err := midi.FindOutPort(outPort)
	if err != nil {
		return msgQuit("midi out port not found")
	}
	rd := bytes.NewReader(bs)
	smf.ReadTracksFrom(rd).Play(out)
	return nil
}

func (m model) updateBinary() tea.Msg {
	resp, err := http.Get(m.baseURL.Rest(data.PathLatestClientDownload))
	if err != nil {
		return msgQuit(err.Error())
	}
	defer resp.Body.Close()
	if err := update.Apply(resp.Body, update.Options{}); err != nil {
		if err := update.RollbackError(err); err != nil {
			return msgQuit(fmt.Sprintf("Failed to rollback from bad update: %s", err))
		}
	}
	return msgQuit("Update complete. Please run again.")
}
