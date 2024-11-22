package comands

import (
	"fmt"
	"strings"

	"github.com/dkotTech/egitlab/internal"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
	"github.com/zalando/go-keyring"
)

func NewSetCredentialsCommand() *cli.Command {
	var credentials internal.Credentials

	return &cli.Command{
		Name:    "set-credentials",
		Aliases: []string{"sc"},
		Usage:   "set a credential",
		Flags:   credentials.CliFlags(),
		Action: func(cCtx *cli.Context) error {
			err := credentials.Parse(cCtx)
			if err != nil {
				return err
			}

			if _, err := tea.NewProgram(initialModel(credentials)).Run(); err != nil {
				return err
			}

			return nil
		},
	}
}

type model struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode

	credentials internal.Credentials
}

func initialModel(credentials internal.Credentials) model {
	m := model{
		inputs: []textinput.Model{
			func() textinput.Model {
				t := textinput.New()
				t.Cursor.Style = internal.CursorStyle

				t.Cursor.SetMode(cursor.CursorBlink)
				t.PromptStyle = internal.FocusedStyle
				t.TextStyle = internal.FocusedStyle
				t.Focus()
				t.Placeholder = "Gitlab token"
				t.EchoMode = textinput.EchoNormal
				t.EchoCharacter = '•'

				return t
			}(),
		},
		credentials: credentials,
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				err := m.saveCredentials()
				if err != nil {
					panic(err)
				}
				return m, tea.Quit
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = internal.FocusedStyle
					m.inputs[i].TextStyle = internal.FocusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = internal.NoStyle
				m.inputs[i].TextStyle = internal.NoStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := fmt.Sprintf("[ %s ]", internal.BlurredStyle.Render("Submit"))
	if m.focusIndex == len(m.inputs) {
		button = internal.FocusedStyle.Render("[ Submit ]")
	}

	b.WriteString(fmt.Sprintf("\n\n%s\n\n", button))
	b.WriteString(internal.HelpStyle.Render("q: exit\n"))

	return b.String()
}

func (m *model) saveCredentials() error {
	return keyring.Set(m.credentials.Service, m.credentials.User, m.inputs[0].Value())
}
