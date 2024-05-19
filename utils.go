package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

type Completion struct {
	text       string
	action     Action
	failure    bool
	selectable bool
}

// help completions should not be selectable
// TODO: Implement non selectable completions
func helpCompletions(rc map[string]Completer) []Completion {
	outs := []Completion{}

	for k, v := range rc {
		outs = append(outs, Completion{text: k + " - " + v.Help, selectable: false})
	}

	return outs
}

// TODO: should add ability to cache results
func getCompletions(input string, rc map[string]Completer) []Completion {
	if len(input) == 0 {
		return helpCompletions(rc)
	}

	splits := strings.Split(input, " ")
	completions := []Completion{}

	if completer, ok := rc[splits[0]]; ok {
		var comps []string
		var err error

		switch completer.Mode {
		case "stdin":
			comps, err = commandOutputWithStdin(completer.Command, strings.Join(splits[1:], " "))
		case "args":
			comps, err = commandOutputWithArgs(completer.Command, strings.Join(splits[1:], " "))
		default:
			comps, err = commandOutputWithArgs(completer.Command, "")
		}

		if err != nil {
			completions = append(
				completions,
				Completion{
					text:       "[ERROR] " + err.Error(),
					action:     Action{Command: completer.Action.Command, Mode: completer.Action.Mode},
					selectable: false,
					failure:    true,
				})

			return completions
		}

		for _, comp := range comps {
			completions = append(
				completions,
				Completion{
					text:       comp,
					action:     Action{Command: completer.Action.Command, Mode: completer.Action.Mode},
					selectable: true,
				})
		}
	}

	return completions
}

func commandOutputWithStdin(command, input string) ([]string, error) {
	cmd := exec.Command(strings.Split(command, " ")[0], strings.Split(command, " ")[1:]...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal("creating stdin pipe:", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal("starting command:", err)
	}

	if _, err := io.WriteString(stdinPipe, input); err != nil {
		log.Fatal("writing to stdin:", err)
	}

	stdinPipe.Close()

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("%s: %s", command, err.Error())
	}

	if len(stdout.String()) > 0 {
		return strings.Split(strings.Trim(stdout.String(), "\n"), "\n"), nil
	}

	return []string{}, nil
}

func commandOutputWithArgs(command, input string) ([]string, error) {
	args := strings.Split(command, " ")[1:]
	if len(input) > 0 {
		args = append(args, strings.Split(strings.TrimSpace(input), " ")...)
	}
	cmd := exec.Command(strings.Split(command, " ")[0], args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s: %s", command, err.Error())
	}

	if len(stdout.String()) > 0 {
		return strings.Split(strings.Trim(stdout.String(), "\n"), "\n"), nil
	}

	return []string{}, nil
}
