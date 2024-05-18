package main

import (
	"bytes"
	"io"
	"log"
	"os/exec"
	"strings"
)

type Completer struct {
	Command string
	Mode    string
	Help    string
}

// TODO: should add ability to cache results
// TODO: this should be read from a config file
var registeredCompletors = map[string]Completer{
	"calc": {Command: "bc", Mode: "stdin", Help: "Calculates an expression using bc"},
}

// help completions should not be selectable
// TODO: Implement non selectable completions
func helpCompletions() []string {
	outs := []string{}

	for k, v := range registeredCompletors {
		outs = append(outs, k+" - "+v.Help)
	}

	return outs
}

func getCompletions(input string) []string {
	if len(input) == 0 {
		return helpCompletions()
	}

	splits := strings.Split(input, " ")
	completions := []string{}

	if completer, ok := registeredCompletors[splits[0]]; ok {
		switch completer.Mode {
		case "stdin":
			completions = commandOutputWithStdin(completer.Command, strings.Join(splits[1:], " "))
		}
	}

	return completions
}

func commandOutputWithStdin(command, input string) []string {
	cmd := exec.Command(command)
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
