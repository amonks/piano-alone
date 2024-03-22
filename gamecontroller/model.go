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

	ws    *websocket.Conn
	state *game.State
	midi  []byte
	log   []*game.Message
	debug []string

	scheduledPerformances     []*game.Performance
	scheduledPerformanceIndex int

	width  int
	height int

	modal    string
	quitting string

	menuSelectionIndex int
	contentInFocus     bool

	conductorButtonIndex int

	midiOutPorts     midi.OutPorts
	midiOutPortIndex int

	latestVersion string

	output string
}

type (
	msgGotScheduledPerformances []*game.Performance
	msgTick                     time.Time
	msgQuit                     string
	msgGotMIDIOutputPorts       midi.OutPorts
	msgVersion                  string
	msgGotWSClient              *websocket.Conn
	msgGotWSMessage             *game.Message
	msgStartedPerformance       string
)

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.checkVersion,
		m.checkMIDIOutPorts,
		m.connect,
		m.tick,
		m.checkScheduledPerformances,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case msgTick, tea.KeyMsg, tea.MouseMsg, tea.WindowSizeMsg:
	case msgGotWSClient:
		m.debug = append(m.debug, "msgGotWSClient")
	case msgGotWSMessage:
		m.debug = append(m.debug, msg.Type.String())
	case msgVersion:
		m.debug = append(m.debug, "msgVersion: "+string(msg))
	case msgGotMIDIOutputPorts:
		m.debug = append(m.debug, fmt.Sprintf("msgGotMIDIOutputPorts: %d ports", len(msg)))
	case msgGotScheduledPerformances:
		m.debug = append(m.debug, fmt.Sprintf("msgGotScheduledPerformances: %d performances", len(msg)))
	case msgStartedPerformance:
		m.debug = append(m.debug, fmt.Sprintf("msgStartedPerformance: %s", string(msg)))
	default:
		m.debug = append(m.debug, fmt.Sprintf("%+v", msg))
	}

	switch msg := msg.(type) {

	case msgTick:
		return m, m.tick

	case msgGotWSClient:
		m.ws = msg
		return m, m.acceptMessage

	case msgGotWSMessage:
		m.log = append(m.log, msg)
		switch msg.Type {
		case game.MessageTypeState:
			m.state = game.StateFromBytes(msg.Data)
			return m, m.acceptMessage

		case game.MessageTypeBroadcastPhase:
			phase := game.PhaseFromBytes(msg.Data)
			m.state.Phase = phase
			return m, m.acceptMessage

		case game.MessageTypeConductorConnected:
			m.state.ConductorIsConnected = true
			return m, m.acceptMessage
		case game.MessageTypeConductorDisconnected:
			m.state.ConductorIsConnected = false
			return m, m.acceptMessage

		case game.MessageTypeDisklavierConnected:
			m.state.DisklavierIsConnected = true
			return m, m.acceptMessage
		case game.MessageTypeDisklavierDisconnected:
			m.state.DisklavierIsConnected = false
			return m, m.acceptMessage

		case game.MessageTypeBroadcastConnectedPlayer:
			player := game.PlayerFromBytes(msg.Data)
			m.state.Players[player.Fingerprint] = player
			return m, m.acceptMessage

		case game.MessageTypeBroadcastDisconnectedPlayer:
			m.state.Players[string(msg.Data)].ConnectionState = game.ConnectionStateDisconnected
			return m, m.acceptMessage

		case game.MessageTypeBroadcastControllerModal:
			m.modal = string(msg.Data)
			return m, m.acceptMessage

		case game.MessageTypeBroadcastSubmittedTrack:
			fingerprint := string(msg.Data)
			if _, ok := m.state.Players[fingerprint]; !ok {
				m.state.Players[fingerprint] = &game.Player{}
			}
			m.state.Players[fingerprint].HasSubmitted = true
			return m, m.acceptMessage

		case game.MessageTypeSendRenditionToDisklavier:
			m.midi = msg.Data
			return m, tea.Batch(
				m.playSong,
				m.acceptMessage,
			)

		default:
			return m, m.acceptMessage
		}

	case msgGotScheduledPerformances:
		m.scheduledPerformances = []*game.Performance(msg)
		return m, nil

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

		return m, nil

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
				return m, nil
			}

		case "u", "U":
			if m.latestVersion != "" && data.CurrentVersion != m.latestVersion {
				return m, m.updateBinary
			}
			return m, nil

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

			case "Start Performance":
				m.menuSelectionIndex = 1
				m.contentInFocus = false
				return m, m.startPerformance

			case "Performance Status":
				switch m.selectedConductorButton() {

				case "Advance":
					return m, m.advance

				case "Restart":
					return m, m.restart

				default:
					return m, nil
				}
			}
		case "k", "up":
			if !m.contentInFocus {
				m.menuSelectionIndex = decrement(m.menuSelectionIndex, len(m.menu()))
				return m, nil
			}
			switch m.selectedMenuItem() {

			case "Start Performance":
				m.scheduledPerformanceIndex = decrement(m.scheduledPerformanceIndex, len(m.scheduledPerformances))
				return m, nil

			case "MIDI Configuration":
				m.midiOutPortIndex = decrement(m.midiOutPortIndex, len(m.midiOutPorts))
				return m, nil

			case "Performance Status":
				m.conductorButtonIndex = decrement(m.conductorButtonIndex, len(m.conductorButtons()))
				return m, nil
			}
		case "j", "down":
			if !m.contentInFocus {
				m.menuSelectionIndex = increment(m.menuSelectionIndex, len(m.menu()))
				return m, nil
			}
			switch m.selectedMenuItem() {

			case "Start Performance":
				m.scheduledPerformanceIndex = increment(m.scheduledPerformanceIndex, len(m.scheduledPerformances))
				return m, nil

			case "MIDI Configuration":
				m.midiOutPortIndex = increment(m.midiOutPortIndex, len(m.midiOutPorts))
				return m, nil

			case "Performance Status":
				m.conductorButtonIndex = increment(m.conductorButtonIndex, len(m.conductorButtons()))
				return m, nil

			default:
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	default:
		return m, nil
	}

	return m, nil
}

type MenuItem string

const (
	MenuItemStartPerformance  = "Start Performance"
	MenuItemPerformanceStatus = "Performance Status"
	MenuItemDebug             = "Debug"
	MenuItemMIDIConfiguration = "MIDI Configuration"
	MenuItemMIDIOutputTest    = "MIDI Output Test"
)

func (m model) menu() []string {
	switch m.role {
	case "conductor":
		return []string{
			MenuItemStartPerformance,
			MenuItemPerformanceStatus,
			MenuItemDebug,
		}
	default:
		return []string{
			MenuItemPerformanceStatus,
			MenuItemMIDIConfiguration,
			MenuItemMIDIOutputTest,
			MenuItemDebug,
		}
	}
}

func (m model) conductorButtons() []string {
	if m.role != "conductor" {
		return []string{}
	}
	switch m.state.Phase {
	case game.GamePhaseUninitialized:
		return []string{}
	case game.GamePhaseLobby:
		return []string{"Advance"}
	case game.GamePhaseHero:
		return []string{"Advance", "Restart"}
	case game.GamePhaseProcessing:
		return []string{}
	case game.GamePhasePlayback:
		return []string{"Advance"}
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
	defer resp.Body.Close()
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

func (m model) checkScheduledPerformances() tea.Msg {
	if m.role != "conductor" {
		return nil
	}
	resp, err := http.Get(m.baseURL.Rest(data.PathScheduledPerformances))
	if err != nil {
		return msgQuit(err.Error())
	}
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return msgQuit(err.Error())
	}
	ps := game.PerformancesFromBytes(bs)
	return msgGotScheduledPerformances(ps)
}

func (m model) startPerformance() tea.Msg {
	p := m.scheduledPerformances[m.scheduledPerformanceIndex]
	url := m.baseURL.Rest(data.PathBeginPerformance, "id", p.Configuration.PerformanceID)
	resp, err := http.Post(url, "", nil)
	if err != nil {
		return msgQuit(err.Error())
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return msgQuit(fmt.Sprintf("status code %d", resp.StatusCode))
	}
	return msgStartedPerformance(url)
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

func increment(i, n int) int {
	if i+1 >= n {
		return 0
	}
	return i + 1
}

func decrement(i, n int) int {
	if i-1 < 0 {
		return n - 1
	}
	return i - 1
}
