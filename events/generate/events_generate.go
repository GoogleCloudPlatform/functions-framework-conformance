// This binary generates usable Go code from the data files.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	tmpl = template.Must(template.New("events").Parse(eventDataTemplate))
)

const (
	input             = "input"
	output            = "output"
	legacyType        = "legacy"
	cloudeventType    = "cloudevent"
	dataDir           = "generate/data"
	outputFile        = "events_data.go"
	eventDataTemplate = `// Code generated by events_generate.go. DO NOT EDIT.

package events

type EventData struct {
	LegacyEvent []byte
	CloudEvent []byte
}

type Event struct {
	Input EventData
	Output EventData
}

var Events = map[string]Event{
		{{ range $k, $v := . }}"{{ $k }}": Event{
			Input: EventData{
				LegacyEvent: {{ if $v.LegacyInput }}{{ $v.LegacyInput }}{{ else }} nil {{ end }},
				CloudEvent: {{ if $v.CloudEventInput }}{{ $v.CloudEventInput }}{{ else }} nil {{ end }},
			},
			Output: EventData{
				LegacyEvent: {{ if $v.LegacyOutput }}{{ $v.LegacyOutput }}{{ else }} nil {{ end }},
				CloudEvent: {{ if $v.CloudEventOutput }}{{ $v.CloudEventOutput }}{{ else }} nil {{ end }},
			},
		},
		{{ end }}
}
`
)

type eventData struct {
	LegacyInput      string
	LegacyOutput     string
	CloudEventInput  string
	CloudEventOutput string
}

func breakdownFileName(path string) (string, string) {
	// Must be a JSON file.
	if !strings.HasSuffix(path, ".json") {
		return "", ""
	}
	fileName := strings.TrimSuffix(path, ".json")

	var et, ft string
	if strings.HasSuffix(fileName, input) {
		ft = input
		fileName = strings.TrimSuffix(fileName, "-"+input)
	} else if strings.HasSuffix(fileName, output) {
		ft = output
		fileName = strings.TrimSuffix(fileName, "-"+output)
	}

	if strings.HasSuffix(fileName, legacyType) {
		et = legacyType
		fileName = strings.TrimSuffix(fileName, "-"+legacyType)
	} else if strings.HasSuffix(fileName, cloudeventType) {
		et = cloudeventType
		fileName = strings.TrimSuffix(fileName, "-"+cloudeventType)
	}

	return fileName, et + ft
}

func main() {
	events := make(map[string]*eventData)
	err := filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("accessing path %q: %v", path, err)
		}

		if info.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			// A source file not existing is fine.
			if os.IsNotExist(err) {
				return nil
			}
			return fmt.Errorf("reading file %q: %v", path, err)
		}

		if data == nil {
			return nil
		}

		name, t := breakdownFileName(info.Name())
		ed, ok := events[name]
		if !ok {
			ed = &eventData{}
			events[name] = ed
		}

		d := "[]byte(`" + string(data) + "`)"
		switch t {
		case legacyType + input:
			ed.LegacyInput = d
		case legacyType + output:
			ed.LegacyOutput = d
		case cloudeventType + input:
			ed.CloudEventInput = d
		case cloudeventType + output:
			ed.CloudEventOutput = d
		}

		return nil
	})
	if err != nil {
		log.Fatalf("walking %q: %v", dataDir, err)
	}

	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("creating event_data.go: %v", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, events); err != nil {
		log.Fatalf("executing template: %v", err)
	}

	return

}
