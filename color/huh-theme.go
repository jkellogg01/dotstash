package color

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const (
	buttonPaddingHorizontal = 2
	buttonPaddingVertical   = 0
)

func ThemeBase() *huh.Theme {
	var t huh.Theme

	t.FieldSeparator = lipgloss.NewStyle().SetString("\n\n")

	button := lipgloss.NewStyle().
		Padding(buttonPaddingVertical, buttonPaddingHorizontal).
		MarginRight(1)

	// Focused styles.
	t.Focused.Base = lipgloss.NewStyle().PaddingLeft(1).BorderStyle(lipgloss.ThickBorder()).BorderLeft(true)
	t.Focused.Card = lipgloss.NewStyle().PaddingLeft(1)
	t.Focused.ErrorIndicator = lipgloss.NewStyle().SetString(" *")
	t.Focused.ErrorMessage = lipgloss.NewStyle().SetString(" *")
	t.Focused.SelectSelector = lipgloss.NewStyle().SetString("> ")
	t.Focused.NextIndicator = lipgloss.NewStyle().MarginLeft(1).SetString("→")
	t.Focused.PrevIndicator = lipgloss.NewStyle().MarginRight(1).SetString("←")
	t.Focused.MultiSelectSelector = lipgloss.NewStyle().SetString("> ")
	t.Focused.SelectedPrefix = lipgloss.NewStyle().SetString("[•] ")
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().SetString("[ ] ")
	t.Focused.FocusedButton = button.Foreground(Black).Background(Purple)
	t.Focused.BlurredButton = button.Foreground(White).Background(Black)
	t.Focused.TextInput.Placeholder = lipgloss.NewStyle().Foreground(BrightBlack)
	t.Focused.Title = lipgloss.NewStyle().Foreground(Purple).Bold(true)
	t.Focused.Description = lipgloss.NewStyle().Faint(true)

	t.Help = help.New().Styles

	// Blurred styles.
	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.MultiSelectSelector = lipgloss.NewStyle().SetString("  ")
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	return &t
}
