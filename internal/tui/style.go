package tui

import "github.com/charmbracelet/lipgloss"

var primaryColor = lipgloss.Color("69")   // blue
var subtleColor  = lipgloss.Color("245")  // gray
var accentColor  = lipgloss.Color("75")   // teal
var dangerColor  = lipgloss.Color("204")  // red-ish

var lipgloss_style = struct {
	HeaderBorder lipgloss.Style
	Footer       lipgloss.Style
	TableHeader  lipgloss.Style
	RowEven      lipgloss.Style
	RowOdd       lipgloss.Style
	RowSelected  lipgloss.Style
	Separator    lipgloss.Style
	Label        lipgloss.Style
	Value        lipgloss.Style
	Title        lipgloss.Style
	Error        lipgloss.Style
	Hint         lipgloss.Style
}{
	HeaderBorder: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1).
		Foreground(primaryColor).
		Bold(true),

	Footer: lipgloss.NewStyle().
		Foreground(subtleColor).
		Italic(true),

	TableHeader: lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Padding(0, 1),

	RowEven: lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Padding(0, 1),

	RowOdd: lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).
		Padding(0, 1),

	RowSelected: lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(primaryColor).
		Bold(false).
		Padding(0, 1),

	Separator: lipgloss.NewStyle().
		Foreground(subtleColor),

	Label: lipgloss.NewStyle().
		Foreground(subtleColor).
		Bold(false),

	Value: lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")),

	Title: lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Padding(0, 1),

	Error: lipgloss.NewStyle().
		Foreground(dangerColor).
		Bold(true),

	Hint: lipgloss.NewStyle().
		Foreground(subtleColor).
		Italic(true),
}
