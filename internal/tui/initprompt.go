package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/elewis787/rkl/internal/cfg"
)

const (
	histKey = `History File Path`
)

type InitPromptModel struct {
	inputs  map[string]textinput.Model
	done    bool
	cfgPath string
}

func NewInitPrompt(cfgPath string) *InitPromptModel {
	homeDir, _ := os.UserHomeDir()
	historyFilePrompt := textinput.New()
	historyFilePrompt.Placeholder = homeDir + "/.history"
	historyFilePrompt.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: `#353C3B`, Dark: `#e5e5e5`})
	historyFilePrompt.Focus()
	fmt.Println(cfgPath)
	return &InitPromptModel{
		cfgPath: cfgPath,
		inputs: map[string]textinput.Model{
			histKey: historyFilePrompt,
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
	// Write output file
	if i.done {
		v := i.inputs[histKey]
		if v.Value() == "" {
			v.SetValue(v.Placeholder)
		}
		config := &cfg.Config{
			HistoryFile: v.Value(),
		}
		err := cfg.ToFile(i.cfgPath, config)
		if err != nil {
			return err.Error()
		}
		return "Initalization complete! \n"
	}
	output := strings.Builder{}
	// Write input to screen
	for k, v := range i.inputs {
		output.WriteString(k + "\n")
		output.WriteString(v.View())
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
