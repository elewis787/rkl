package note

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func NoteCmd() *cobra.Command {
	noteCmd := &cobra.Command{
		Use: "note",
		RunE: func(cmd *cobra.Command, args []string) error {
			return chat()
		},
	}
	return noteCmd

}

func chat() error {
	if !term.IsTerminal(0) || !term.IsTerminal(1) {
		return fmt.Errorf("stdin/stdout should be terminal")
	}
	oldState, err := term.MakeRaw(0)
	if err != nil {
		return err
	}
	defer term.Restore(0, oldState)
	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	term := term.NewTerminal(screen, "")
	term.SetPrompt(string(term.Escape.Red) + "> " + string(term.Escape.Reset))

	rePrefix := string(term.Escape.Cyan) + "Human says:" + string(term.Escape.Reset)

	for {
		line, err := term.ReadLine()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if line == "" {
			continue
		}
		fmt.Fprintln(term, rePrefix, line)
	}
}
