package main

import (
	"bytes"
	"io"
	"log"
	"os/exec"
	"strings"
)

// help completions should not be selectable
// TODO: Implement non selectable completions
func helpCompletions(rc map[string]Completer) []string {
	outs := []string{}

	for k, v := range rc {
		outs = append(outs, k+" - "+v.Help)
	}

	return outs
}

// TODO: should add ability to cache results
func getCompletions(input string, rc map[string]Completer) []string {
	if len(input) == 0 {
		return helpCompletions(rc)
	}

	splits := strings.Split(input, " ")
	completions := []string{}

	if completer, ok := rc[splits[0]]; ok {
		switch completer.Mode {
		case "stdin":
			completions = commandOutputWithStdin(completer.Command, strings.Join(splits[1:], " "))
		case "args":
			completions = commandOutputWithArgs(completer.Command, strings.Join(splits[1:], " "))
		}
	}

	return completions
}

func commandOutputWithStdin(command, input string) []string {
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
		return []string{"[ERROR] " + command + ": " + err.Error()} // ability to color it (maybe red)
	}

	if len(stdout.String()) > 0 {
		return strings.Split(strings.Trim(stdout.String(), "\n"), "\n")
	}

	return []string{}
}

func commandOutputWithArgs(command, input string) []string {
	args := append(strings.Split(command, " ")[1:], strings.Split(strings.TrimSpace(input), " ")...)
	cmd := exec.Command(strings.Split(command, " ")[0], args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return []string{"[ERROR] " + command + ": " + err.Error()} // ability to color it (maybe red)
	}

	if len(stdout.String()) > 0 {
		return strings.Split(strings.Trim(stdout.String(), "\n"), "\n")
	}

	return []string{}
}
