package main

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Action struct {
	Command string
	Mode    string
}

// TODO: user should be able to define if we do additional filtering
// and if cmd should be one shot
type Completer struct {
	Name    string
	Command string
	Mode    string
	Help    string
	Action  Action
}

func getRegisteredCompleters() (map[string]Completer, error) {
	t := make(map[string][]Completer)

	data, err := os.ReadFile("pkr.yaml")
	if err != nil {
		return nil, errors.Wrap(err, "reading file")
	}

	err = yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling yaml")
	}

	cptrs := make(map[string]Completer)
	cp, ok := t["completers"]
	if ok {
		for _, c := range cp {
			cptrs[c.Name] = c
		}
	}

	return cptrs, nil
}
