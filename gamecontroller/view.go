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

	indicatorStyleOn = lipgloss.NewStyle().
				Background(xxxDark).
				Foreground(green).
				Bold(true)
	indicatorStyleOff = lipgloss.NewStyle().
				Background(xxxDark).
				Foreground(red).
				Bold(true)

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
			Padding(0, 3)
	activeButtonStyle = buttonStyle.Copy().
				Background(lipgloss.Color("#FDF6E3")).
				Foreground(xxDark).
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
	return zone.Mark("Modal", modalStyle.Copy().Width(m.width).Height(m.height).Render(
		modalHeaderStyle.Render(m.modal)+"\n"+modalDismisserStyle.Render("press any key to dismiss"),
	))
}

func (m model) viewMenu() string {
	content := strings.Builder{}
	for _, item := range m.menu() {
		if item == m.menu()[m.menuSelectionIndex] {
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
	switch m.menu()[m.menuSelectionIndex] {
	case MenuItemStartPerformance:
		content = m.viewStartPerformance()
	case MenuItemPerformanceStatus:
		content = m.viewPerformanceStatus()
	case MenuItemMIDIConfiguration:
		content = m.viewMidiOutPorts()
	case MenuItemMIDIOutputTest:
		content = m.viewMidiOutputTest()
	case MenuItemDebug:
		content = m.viewDebug()
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
	outPorts := make([]string, len(m.midiOutPorts))
	for i, p := range m.midiOutPorts {
		outPorts[i] = p.String()
	}

	return joinVertical(
		boxHeaderStyle.Render("MIDI Out Ports"),
		renderList(outPorts, m.midiOutPortIndex, "none found"),
	)
}

func renderList(items []string, selectedIndex int, alt string) string {
	if len(items) == 0 {
		return alt
	}
	out := make([]string, len(items))
	for i, label := range items {
		indicator := "  "
		if i == selectedIndex {
			indicator = "• "
		}
		out[i] = indicator + label
	}
	return joinVertical(out...)
}

func renderButton(label string, isActive bool) string {
	buttonStyle := buttonStyle
	if isActive {
		buttonStyle = activeButtonStyle
	}
	return zone.Mark(label, buttonStyle.Render(label))
}

func (m model) viewMidiOutputTest() string {
	hasMidiOutPort := m.midiOutPortIndex < len(m.midiOutPorts)
	selectedPortMessage := "no MIDI out ports found"
	if hasMidiOutPort {
		selectedPortMessage = "MIDI Port: " + m.midiOutPorts[m.midiOutPortIndex].String()
	}
	return joinVertical(
		boxHeaderStyle.Render(MenuItemMIDIOutputTest),
		selectedPortMessage,
		"",
		renderButton("Test MIDI Output", m.contentInFocus && hasMidiOutPort),
	)
}

func (m model) viewPerformanceStatus() string {
	if m.state == nil {
		return "nil state"
	}

	switch m.state.Phase {
	case game.GamePhaseUninitialized:
		return joinVertical(
			boxHeaderStyle.Render(MenuItemPerformanceStatus),
			m.viewSelectedMIDIOutPort(),
			"Waiting for server operator to start the performance.",
			fmt.Sprintf("Please direct attendees to %s", m.baseURL.Rest("/")),
			fmt.Sprintf("Connected players: %d", m.state.CountConnectedPlayers()),
			"",
			m.viewConnectionStatus(),
			m.viewConductorButtons(),
		)

	case game.GamePhaseLobby:
		return joinVertical(
			boxHeaderStyle.Render(MenuItemPerformanceStatus),
			m.viewSelectedMIDIOutPort(),
			"Waiting for players to join.",
			fmt.Sprintf("Please direct attendees to %s", m.baseURL.Rest("/")),
			fmt.Sprintf("Connected players: %d", m.state.CountConnectedPlayers()),
			"",
			m.viewConnectionStatus(),
			m.viewConductorButtons(),
		)

	case game.GamePhaseHero:
		return joinVertical(
			boxHeaderStyle.Render(MenuItemPerformanceStatus),
			m.viewSelectedMIDIOutPort(),
			"Players are playing.",
			"",
			fmt.Sprintf("Connected players: %d", m.state.CountConnectedPlayers()),
			fmt.Sprintf("Submitted tracks: %d", m.state.CountSubmittedTracks()),
			m.viewConnectionStatus(),
			m.viewConductorButtons(),
		)

	case game.GamePhaseProcessing:
		return joinVertical(
			boxHeaderStyle.Render(MenuItemPerformanceStatus),
			m.viewSelectedMIDIOutPort(),
			"Processing MIDI from players.",
			"",
			fmt.Sprintf("Connected players: %d", m.state.CountConnectedPlayers()),
			fmt.Sprintf("Submitted tracks: %d", m.state.CountSubmittedTracks()),
			m.viewConnectionStatus(),
			m.viewConductorButtons(),
		)

	case game.GamePhasePlayback:
		return joinVertical(
			boxHeaderStyle.Render(MenuItemPerformanceStatus),
			m.viewSelectedMIDIOutPort(),
			"Playing combined MIDI on disklavier.",
			"",
			fmt.Sprintf("Connected players: %d", m.state.CountConnectedPlayers()),
			fmt.Sprintf("Submitted tracks: %d", m.state.CountSubmittedTracks()),
			m.viewConnectionStatus(),
			m.viewConductorButtons(),
		)

	case game.GamePhaseDone:
		return joinVertical(
			boxHeaderStyle.Render(MenuItemPerformanceStatus),
			m.viewSelectedMIDIOutPort(),
			"The performance is over.",
			"",
			fmt.Sprintf("Connected players: %d", m.state.CountConnectedPlayers()),
			fmt.Sprintf("Submitted tracks: %d", m.state.CountSubmittedTracks()),
			"",
			m.viewConnectionStatus(),
			m.viewConductorButtons(),
		)

	default:
		return joinVertical(
			boxHeaderStyle.Render(MenuItemPerformanceStatus),
			m.viewSelectedMIDIOutPort(),
			"Unknown state. Something is wrong.",
			"",
			fmt.Sprintf("Connected players: %d", m.state.CountConnectedPlayers()),
			fmt.Sprintf("Submitted tracks: %d", m.state.CountSubmittedTracks()),
			m.viewConnectionStatus(),
			m.viewConductorButtons(),
		)
	}
}

func (m model) viewConnectionStatus() string {
	return joinVertical(
		renderStatus("disklavier", m.state.DisklavierIsConnected),
		renderStatus("conductor", m.state.ConductorIsConnected),
	)
}

func (m model) viewStartPerformance() string {
	if m.state == nil {
		return ""
	}
	switch m.state.Phase {
	case game.GamePhaseUninitialized, game.GamePhaseDone:
		return joinVertical(
			boxHeaderStyle.Render(MenuItemStartPerformance),
			m.viewScheduledPerformances(),
		)
	default:
		return joinVertical(
			boxHeaderStyle.Render(MenuItemStartPerformance),
			"Performance has started.",
		)
	}
}

func (m model) viewScheduledPerformances() string {
	var ps []string
	for _, p := range m.scheduledPerformances {
		ps = append(ps, p.Date.String())
	}
	return renderList(ps, m.scheduledPerformanceIndex, "none found")
}

func (m model) viewConductorButtons() string {
	buttons := m.conductorButtons()
	button := m.selectedConductorButton()
	if len(buttons) == 0 {
		return ""
	}
	output := []string{""}
	for _, b := range buttons {
		output = append(output, renderButton(b, b == button))
	}
	return joinVertical(output...)
}

func renderStatus(name string, isConnected bool) string {
	if isConnected {
		return indicatorStyleOn.Render("•") + pageStyle.Render(" "+name)
	} else {
		return indicatorStyleOff.Render("•") + pageStyle.Render(" "+name)
	}
}

func (m model) viewSelectedMIDIOutPort() string {
	if len(m.midiOutPorts) <= 0 {
		return ""
	}
	return fmt.Sprintf("MIDI Port: %s", m.midiOutPorts[m.midiOutPortIndex].String())
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

func (m model) viewDebug() string {
	lines := []string{
		boxHeaderStyle.Render("State"),
		fmt.Sprintf("Configuration: %t", m.state.Configuration != nil),
		fmt.Sprintf("Phase: %s", m.state.Phase.String()),
		"",
		boxHeaderStyle.Render("Log"),
	}
	lines = append(lines, m.debug...)
	return joinVertical(lines...)

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
