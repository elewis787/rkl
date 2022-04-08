package tui

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InitPromptModel struct {
	index  int
	inputs map[string]textinput.Model
	done   bool
	cmd    tea.Cmd
}

func NewInitPrompt() *InitPromptModel {
	homeDir, _ := os.UserHomeDir()
	historyFilePrompt := textinput.New()
	historyFilePrompt.Placeholder = homeDir + "/.history"
	historyFilePrompt.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: `#353C3B`, Dark: `#e5e5e5`})
	historyFilePrompt.Focus()
	return &InitPromptModel{
		inputs: map[string]textinput.Model{
			"history file path": historyFilePrompt,
		},
	}
}

func (i InitPromptModel) Init() tea.Cmd {
	return textinput.Blink
}

func (i InitPromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return i, tea.Quit
		case "enter":
			i.done = true
			return i, tea.Quit
		}
	}
	cmd := i.updateInputs(msg)
	return i, cmd
}

func (i InitPromptModel) View() string {
	output := strings.Builder{}
	if i.done {
		for _, v := range i.inputs {
			if v.Value() == "" {
				v.SetValue(v.Placeholder)
			}
			err := os.WriteFile("/tmp/dat1", []byte(v.Value()), 0644)
			if err != nil {
				output.WriteString(err.Error())
			}
		}
	} else {
		for k, v := range i.inputs {
			output.WriteString(k + "\n")
			output.WriteString(v.View())
		}
	}
	return output.String()
}

func (i *InitPromptModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	for k := range i.inputs {
		if i.inputs[k].Focused() {
			m, cmd := i.inputs[k].Update(msg)
			i.inputs[k] = m
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}
