package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

// runCmd runs a command and returns its stdout as a string.
func runCmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s failed: %w", name, err)
	}
	return out.String(), nil
}

func main() {
	app := &cli.App{
		Name:  "gocat",
		Usage: "Concatenate files in a directory using ripgrep + optional fzf",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "fzf", Usage: "Use fzf to interactively pick files (press Tab to multi-select)"},
			&cli.BoolFlag{Name: "hidden", Usage: "Include hidden files (dotfiles)"},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return fmt.Errorf("please provide a folder path")
			}
			folder := c.Args().First()

			// ripgrep collects files, respecting .gitignore and ignore rules
			args := []string{"--files"}
			if c.Bool("hidden") {
				args = append(args, "--hidden")
			}
			args = append(args, folder)

			out, err := runCmd("rg", args...)
			if err != nil {
				return err
			}
			files := strings.Split(strings.TrimSpace(out), "\n")

			// optionally filter through fzf
			if c.Bool("fzf") {
				cmd := exec.Command("fzf", "--multi")
				cmd.Stdin = strings.NewReader(strings.Join(files, "\n"))

				var fzfOut bytes.Buffer
				cmd.Stdout = &fzfOut
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					return err
				}
				files = strings.Split(strings.TrimSpace(fzfOut.String()), "\n")
			}

			// concatenate files
			var builder strings.Builder
			for _, f := range files {
				if f == "" {
					continue
				}
				data, err := os.ReadFile(f)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Skip %s: %v\n", f, err)
					continue
				}
				builder.WriteString(fmt.Sprintf("------ %s ------\n", f))
				builder.Write(data)
				builder.WriteString("\n\n")
			}

			_, err = fmt.Print(builder.String())
			return err
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
