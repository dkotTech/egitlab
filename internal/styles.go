package internal

import (
	"github.com/charmbracelet/lipgloss"
)

type (
	// StatusProps a structure for status properties
	StatusProps struct {
		Emoji string
		Style lipgloss.Style
	}
)

var (
	// StatusStyles a map that combines emojis and styles
	StatusStyles = map[string]StatusProps{
		"created": {
			Emoji: "🆕",
			Style: lipgloss.NewStyle().Foreground(lipgloss.Color("#bf87e6")), // Purple
		},
		"pending": {
			Emoji: "⌛",
			Style: lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")), // Orange
		},
		"preparing": {
			Emoji: "⚙",
			Style: lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")), // Bright Yellow
		},
		"waiting_for_resource": {
			Emoji: "⏳",
			Style: lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")), // Cyan
		},
		"running": {
			Emoji: "🏃",
			Style: lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")), // Bright Lime Green
		},
		"running_with_artifacts": {
			Emoji: "📦",
			Style: lipgloss.NewStyle().Foreground(lipgloss.Color("#4682B4")), // Steel Blue
		},
		"success": {
			Emoji: "✅",
			Style: lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")), // Bright Green
		},
		"failed": {
			Emoji: "⛔",
			Style: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")), // Red
		},
		"failed_with_warnings": {
			Emoji: "⚠",
			Style: lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA07A")), // Light Salmon
		},
		"canceled": {
			Emoji: "🚫",
			Style: lipgloss.NewStyle().Foreground(lipgloss.Color("#DC143C")), // Crimson
		},
		"skipped": {
			Emoji: "⏭",
			Style: lipgloss.NewStyle().Foreground(lipgloss.Color("#D3D3D3")), // Light Gray
		},
		"manual": {
			Emoji: "✋",
			Style: lipgloss.NewStyle().Foreground(lipgloss.Color("#00CED1")), // Dark Turquoise
		},
		"scheduled": {
			Emoji: "⏰",
			Style: lipgloss.NewStyle().Foreground(lipgloss.Color("#1E90FF")), // Dodger Blue
		},
		"deployed": {
			Emoji: "🚀",
			Style: lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")), // Bright Lime Green
		},
	}

	TableBorderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#660066"))

	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	CursorStyle  = FocusedStyle
	NoStyle      = lipgloss.NewStyle()
	HelpStyle    = BlurredStyle
)
