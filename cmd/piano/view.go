package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"monks.co/piano-alone/data"
	"monks.co/piano-alone/game"
)

var (
	yellow  = lipgloss.Color("#B58900")
	orange  = lipgloss.Color("#CB4B16")
	red     = lipgloss.Color("#DC322F")
	magenta = lipgloss.Color("#D33682")
	violet  = lipgloss.Color("#6C71C4")
	blue    = lipgloss.Color("#268BD2")
	cyan    = lipgloss.Color("#2AA198")
	green   = lipgloss.Color("#859900")

	xxxLight = lipgloss.Color("#FDF6E3") // base3
	xxLight  = lipgloss.Color("#EEE8D5") // base2
	xLight   = lipgloss.Color("#93A1A1") // base1
	light    = lipgloss.Color("#839496") // base0
	dark     = lipgloss.Color("#657B83") // base00
	xDark    = lipgloss.Color("#586E75") // base01
	xxDark   = lipgloss.Color("#073642") // base02
	xxxDark  = lipgloss.Color("#002B36") // base03

	modalStyle = lipgloss.NewStyle().
			Background(xxxLight).
			Align(lipgloss.Center, lipgloss.Center)
	modalHeaderStyle = lipgloss.NewStyle().
				Background(xxxLight).
				Foreground(xxxDark).
				Bold(true).
				Underline(true).
				Blink(true)
	modalDismisserStyle = lipgloss.NewStyle().
				Background(xxxLight).
				Foreground(xDark)

	pageStyle = lipgloss.NewStyle().
			Background(xxxDark).
			Foreground(xxxLight)

	contentStyle = lipgloss.NewStyle().
			Background(xxxDark).
			Foreground(xxxLight).
			Padding(1, 4)
	headerStyle = lipgloss.NewStyle().
			Background(xxDark).
			Foreground(orange).
			Align(lipgloss.Center).
			Bold(true)

	boxStyle = lipgloss.NewStyle().
			Background(xxxDark).
			Foreground(xxxLight).
			Border(lipgloss.NormalBorder()).
			BorderBackground(xxxDark).
			BorderForeground(dark).
			Padding(0, 1)
	focusedBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderBackground(xxxDark).
			BorderForeground(xxxLight).
			Padding(0, 1)
	boxHeaderStyle = lipgloss.NewStyle().
			Background(xxxDark).
			Foreground(orange).
			Bold(true)

	menuItemStyle = lipgloss.NewStyle().
			Background(xxxDark).
			Foreground(xxxLight)
	focusedMenuItemStyle = lipgloss.NewStyle().
				Background(xxxDark).
				Foreground(orange)

	statusbarStyle = lipgloss.NewStyle().
			Background(xxDark).
			Foreground(xxxLight)
	versionStyle = lipgloss.NewStyle().
			Background(xxDark).
			Foreground(xxxLight).
			Padding(0, 1).
			Width(lipgloss.Width(data.CurrentVersion) + 2)
	msgStyle = lipgloss.NewStyle().
			Background(xxDark).
			Foreground(xxxLight).
			Align(lipgloss.Right).
			Padding(0, 1)

	buttonStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#93A1A1")).
			Foreground(xxxDark).
			Padding(0, 3).
			MarginTop(1)
	activeButtonStyle = buttonStyle.Copy().
				Background(lipgloss.Color("#FDF6E3")).
				Foreground(xxDark).
				MarginRight(2).
				Underline(true)
)

func (m model) View() string {
	return zone.Scan(m.view())
}

func (m model) view() string {
	if m.modal != "" {
		return m.viewModal()
	}

	return joinVertical(
		m.viewHeader(),
		lipgloss.JoinHorizontal(lipgloss.Center,
			m.viewMenu(),
			m.viewContent(),
		),
		m.viewStatusbar(),
	)
}

func (m model) viewHeader() string {
	return headerStyle.Copy().Width(m.width).Render("LIFE ONLINE: Piano Telephone")
}

func (m model) viewModal() string {
	return modalStyle.Copy().Width(m.width).Height(m.height).Render(
		modalHeaderStyle.Render(m.modal) + "\n" + modalDismisserStyle.Render("press any key to dismiss"),
	)
}

func (m model) viewMenu() string {
	content := strings.Builder{}
	for _, item := range menu {
		if item == menu[m.menuSelectionIndex] {
			content.WriteString(zone.Mark(item, focusedMenuItemStyle.Render(item)) + "\n")
		} else {
			content.WriteString(zone.Mark(item, menuItemStyle.Render(item)) + "\n")
		}
	}
	menuStyle := boxStyle.Copy().Height(m.height - 4)
	if !m.contentInFocus {
		menuStyle = menuStyle.BorderForeground(lipgloss.Color("#FFFFFF"))
	}
	return zone.Mark("Menu", menuStyle.Render(content.String()))
}

func (m model) viewContent() string {
	var content string
	switch menu[m.menuSelectionIndex] {
	case "Performance Status":
		content = m.viewPerformanceStatus()
	case "MIDI Configuration":
		content = m.viewMidiOutPorts()
	case "MIDI Output Test":
		content = m.viewMidiOutputTest()
	case "Message Log":
		content = m.viewMessageLog()
	default:
		content = "?"
	}
	contentStyle := boxStyle.Copy().Height(m.height - 4).Width(m.width - menuWidth - 2)
	if m.contentInFocus {
		contentStyle = contentStyle.BorderForeground(lipgloss.Color("#FFFFFF"))
	}
	return zone.Mark("Content", contentStyle.Render(content))
}

func (m model) viewMidiOutPorts() string {
	var midiOutPorts strings.Builder
	if len(m.midiOutPorts) == 0 {
		midiOutPorts.WriteString("none found")
	} else {
		for _, p := range m.midiOutPorts {
			if m.midiOutPorts[m.midiOutPortIndex].String() == p.String() {
				midiOutPorts.WriteString(zone.Mark(p.String(), "* "+p.String()) + "\n")
			} else {
				midiOutPorts.WriteString(zone.Mark(p.String(), "  "+p.String()) + "\n")
			}
		}
	}

	return joinVertical(
		boxHeaderStyle.Render("MIDI Out Ports"),
		midiOutPorts.String(),
	)
}

func (m model) viewMidiOutputTest() string {
	hasMidiOutPort := m.midiOutPortIndex < len(m.midiOutPorts)
	buttonStyle := buttonStyle
	if m.contentInFocus && hasMidiOutPort {
		buttonStyle = activeButtonStyle
	}
	selectedPortMessage := "no MIDI out ports found"
	if hasMidiOutPort {
		selectedPortMessage = "MIDI Port: " + m.midiOutPorts[m.midiOutPortIndex].String()
	}
	return joinVertical(
		boxHeaderStyle.Render("MIDI Output Test"),
		selectedPortMessage,
		zone.Mark("Test MIDI Output", buttonStyle.Render("Test MIDI Output")),
	)
}

func (m model) viewPerformanceStatus() string {
	if m.state == nil {
		return "nil state"
	}
	switch m.state.Phase.Type {
	case game.GamePhaseUninitialized:
		return joinVertical(
			boxHeaderStyle.Render("Performance Status"),
			m.viewSelectedMIDIOutPort(),
			"Waiting for server operator to start the performance.",
			fmt.Sprintf("Please direct attendees to %s", m.baseURL.Rest("/")),
			"",
			"",
		)
	case game.GamePhaseLobby:
		return joinVertical(
			boxHeaderStyle.Render("Performance Status"),
			m.viewSelectedMIDIOutPort(),
			"Waiting for players to join.",
			fmt.Sprintf("Please direct attendees to %s", m.baseURL.Rest("/")),
			fmt.Sprintf("Connected players: %d", m.state.CountConnectedPlayers()),
			"",
		)

	case game.GamePhaseHero:
		return joinVertical(
			boxHeaderStyle.Render("Performance Status"),
			m.viewSelectedMIDIOutPort(),
			"Players are playing.",
			"",
			fmt.Sprintf("Connected players: %d", m.state.CountConnectedPlayers()),
			fmt.Sprintf("Submitted tracks: %d", m.state.CountSubmittedTracks()),
		)

	case game.GamePhaseProcessing:
		return joinVertical(
			boxHeaderStyle.Render("Performance Status"),
			m.viewSelectedMIDIOutPort(),
			"Processing MIDI from players.",
			"",
			fmt.Sprintf("Connected players: %d", m.state.CountConnectedPlayers()),
			fmt.Sprintf("Submitted tracks: %d", m.state.CountSubmittedTracks()),
		)

	case game.GamePhasePlayback:
		return joinVertical(
			boxHeaderStyle.Render("Performance Status"),
			m.viewSelectedMIDIOutPort(),
			"Playing combined MIDI on disklavier.",
			"",
			fmt.Sprintf("Connected players: %d", m.state.CountConnectedPlayers()),
			fmt.Sprintf("Submitted tracks: %d", m.state.CountSubmittedTracks()),
		)

	case game.GamePhaseDone:
		return joinVertical(
			boxHeaderStyle.Render("Performance Status"),
			m.viewSelectedMIDIOutPort(),
			"The performance is over.",
			"",
			fmt.Sprintf("Connected players: %d", m.state.CountConnectedPlayers()),
			fmt.Sprintf("Submitted tracks: %d", m.state.CountSubmittedTracks()),
			"",
		)

	default:
		return joinVertical(
			boxHeaderStyle.Render("Performance Status"),
			m.viewSelectedMIDIOutPort(),
			"Unknown state. Something is wrong.",
			"",
			fmt.Sprintf("Connected players: %d", m.state.CountConnectedPlayers()),
			fmt.Sprintf("Submitted tracks: %d", m.state.CountSubmittedTracks()),
		)
	}
}

func (m model) viewSelectedMIDIOutPort() string {
	if len(m.midiOutPorts) <= 0 {
		return ""
	}
	return fmt.Sprintf("MIDI Port: %s", m.midiOutPorts[m.midiOutPortIndex].String())
}

func (m model) viewMessageLog() string {
	msgs := make([]string, len(m.log)+1)
	msgs[0] = boxHeaderStyle.Render("Message Log")
	for i, m := range m.log {
		msgs[i+1] = m.String()
	}
	return joinVertical(
		msgs...,
	)
}

func (m model) viewStatusbar() string {
	if m.quitting != "" {
		return msgStyle.Copy().Width(m.width).Render(
			fmt.Sprintf("press %s again to quit", m.quitting),
		)
	}

	var msg string
	switch m.latestVersion {
	case "":
		msg = "checking for updates"
	case data.CurrentVersion:
		msg = "latest version"
	default:
		msg = fmt.Sprintf("%s is newer, press u to update", m.latestVersion)
	}
	return statusbarStyle.Copy().Width(m.width).Render(
		lipgloss.JoinHorizontal(lipgloss.Top,
			versionStyle.Render(data.CurrentVersion),
			msgStyle.Copy().Width(m.width-lipgloss.Width(data.CurrentVersion)-2).Render(msg),
		),
	)
}

func joinVertical(items ...string) string {
	var maxWidth int
	for _, it := range items {
		if width := lipgloss.Width(it); width > maxWidth {
			maxWidth = width
		}
	}
	builtItems := make([]string, len(items))
	style := lipgloss.NewStyle().Width(maxWidth).Background(xxxDark)
	for i, it := range items {
		builtItems[i] = style.Render(it)
	}
	return strings.Join(builtItems, "\n")
}