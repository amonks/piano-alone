package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"monks.co/piano-alone/data"
	"monks.co/piano-alone/game"
)

var (
	highlight  = lipgloss.Color("#CB4B16")
	highlight2 = lipgloss.Color("#859900")

	faded       = lipgloss.Color("#657b83")
	background  = lipgloss.Color("#002B36")
	background2 = lipgloss.Color("#073642")
	foreground  = lipgloss.Color("#FDF6E3")

	backgroundStyle = lipgloss.NewStyle().
			Background(background)

	pageStyle = lipgloss.NewStyle().
			Background(background).
			Foreground(foreground)

	contentStyle = lipgloss.NewStyle().
			Background(background).
			Foreground(foreground).
			Padding(1, 4)
	headerStyle = lipgloss.NewStyle().
			Background(background2).
			Foreground(highlight).
			Align(lipgloss.Center).
			Bold(true)

	boxStyle = lipgloss.NewStyle().
			Background(background).
			Foreground(foreground).
			Border(lipgloss.NormalBorder()).
			BorderBackground(background).
			BorderForeground(faded).
			Padding(0, 1)
	focusedBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderBackground(background).
			BorderForeground(foreground).
			Padding(0, 1)
	boxHeaderStyle = lipgloss.NewStyle().
			Background(background).
			Foreground(highlight).
			Bold(true)

	menuItemStyle = lipgloss.NewStyle().
			Background(background).
			Foreground(foreground)
	focusedMenuItemStyle = lipgloss.NewStyle().
				Background(background).
				Foreground(highlight)

	statusbarStyle = lipgloss.NewStyle().
			Background(background2).
			Foreground(foreground)
	versionStyle = lipgloss.NewStyle().
			Background(background2).
			Foreground(foreground).
			Padding(0, 1).
			Width(lipgloss.Width(data.CurrentVersion) + 2)
	msgStyle = lipgloss.NewStyle().
			Background(background2).
			Foreground(foreground).
			Align(lipgloss.Right).
			Padding(0, 1)

	buttonStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#93A1A1")).
			Foreground(background).
			Padding(0, 3).
			MarginTop(1)
	activeButtonStyle = buttonStyle.Copy().
				Background(lipgloss.Color("#FDF6E3")).
				Foreground(background2).
				MarginRight(2).
				Underline(true)
)

func (m model) View() string {
	menuStyle := boxStyle.Copy().Height(m.height - 4)
	if !m.contentInFocus {
		menuStyle = menuStyle.BorderForeground(lipgloss.Color("#FFFFFF"))
	}
	menu := menuStyle.Render(m.viewMenu())
	menuWidth := lipgloss.Width(menu)

	contentStyle := boxStyle.Copy().Height(m.height - 4).Width(m.width - menuWidth - 4)
	if m.contentInFocus {
		contentStyle = contentStyle.BorderForeground(lipgloss.Color("#FFFFFF"))
	}
	content := contentStyle.Render(m.viewContent())

	return joinVertical(
		headerStyle.Copy().Width(m.width).Render("LIFE ONLINE: Piano Telephone"),
		lipgloss.JoinHorizontal(lipgloss.Center,
			menu,
			content,
		),
		statusbarStyle.Copy().Width(m.width).Render(m.viewStatusbar()),
	)
}

func (m model) viewMenu() string {
	content := strings.Builder{}
	for _, item := range menu {
		if item == menu[m.menuSelectionIndex] {
			content.WriteString(focusedMenuItemStyle.Render(item) + "\n")
		} else {
			content.WriteString(menuItemStyle.Render(item) + "\n")
		}
	}
	return content.String()
}

func (m model) viewMidiOutPorts() string {
	var midiOutPorts strings.Builder
	if len(m.midiOutPorts) == 0 {
		midiOutPorts.WriteString("none found")
	} else {
		for _, p := range m.midiOutPorts {
			if m.midiOutPorts[m.midiOutPortIndex].String() == p.String() {
				midiOutPorts.WriteString("* " + p.String() + "\n")
			} else {
				midiOutPorts.WriteString("  " + p.String() + "\n")
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
		selectedPortMessage = "selected port: " + m.midiOutPorts[m.midiOutPortIndex].String()
	}
	return joinVertical(
		boxHeaderStyle.Render("MIDI Output Test"),
		selectedPortMessage,
		buttonStyle.Render("Test Midi Output"),
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
			"Waiting for server operator to start the performance.",
			fmt.Sprintf("Please direct attendees to %s", m.baseURL.Rest("/")),
		)
	case game.GamePhaseLobby:
		return joinVertical(
			boxHeaderStyle.Render("Performance Status"),
			"Waiting for players to join.",
			fmt.Sprintf("Please direct attendees to %s", m.baseURL.Rest("/")),
			fmt.Sprintf("Connected players: %d", len(m.state.Players)),
			fmt.Sprintf("Next phase begins in: %s", time.Until(m.state.Phase.Exp)),
		)

	case game.GamePhaseHero:
		return joinVertical(
			boxHeaderStyle.Render("Performance Status"),
			"Players are playing.",
			"",
			fmt.Sprintf("Connected players: %d", len(m.state.Players)),
			fmt.Sprintf("Next phase begins in: %s", time.Until(m.state.Phase.Exp)),
		)

	case game.GamePhaseProcessing:
		return joinVertical(
			boxHeaderStyle.Render("Performance Status"),
			"Processing MIDI from players.",
			"",
			fmt.Sprintf("Connected players: %d", len(m.state.Players)),
			"",
		)

	case game.GamePhasePlayback:
		return joinVertical(
			boxHeaderStyle.Render("Performance Status"),
			"Playing combined MIDI on disklavier.",
			"",
			fmt.Sprintf("Connected players: %d", len(m.state.Players)),
			"",
		)

	case game.GamePhaseDone:
		return joinVertical(
			boxHeaderStyle.Render("Performance Status"),
			"The performance is over.",
			"",
			fmt.Sprintf("Connected players: %d", len(m.state.Players)),
			"",
		)

	default:
		return joinVertical(
			boxHeaderStyle.Render("Performance Status"),
			"Unknown state. Something is wrong.",
		)
	}
}

func (m model) viewContent() string {
	switch menu[m.menuSelectionIndex] {
	case "Performance Status":
		return m.viewPerformanceStatus()
	case "MIDI Configuration":
		return m.viewMidiOutPorts()
	case "MIDI Output Test":
		return m.viewMidiOutputTest()
	default:
		return "?"
	}
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
	return lipgloss.JoinHorizontal(lipgloss.Top,
		versionStyle.Render(data.CurrentVersion),
		msgStyle.Copy().Width(m.width-lipgloss.Width(data.CurrentVersion)-2).Render(msg),
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
	style := lipgloss.NewStyle().Width(maxWidth).Background(background)
	for i, it := range items {
		builtItems[i] = style.Render(it)
	}
	return strings.Join(builtItems, "\n")
}
