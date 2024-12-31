package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v3"

	"github.com/neilsf/sitemongoose/internal/pkg/monitor"
)

const APP_VERSION = "0.1.0"

var app struct {
	Start struct {
		ConfigFilePath string `type:"existingfile" name:"config-file" help:"The path to the configuration file, including the filename." short:"c" default:"config.yaml"`
	} `cmd:"" help:"Starts monitoring." default:"1"`
	Version struct {
	} `cmd:"" help:"Prints the program version."`
}

var appConfig struct {
	Monitors []monitor.Monitor
}

func getExePath() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err.Error())
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		log.Fatal(err.Error())
	}
	return exePath
}

func readConfig(path string) {
	if path == "" {
		exePath := getExePath()
		path = filepath.Join(filepath.Dir(exePath), "config.yaml")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	if yaml.Unmarshal([]byte(data), &appConfig) != nil {
		log.Fatalf("Error parsing config file")
	}
}

func validateConfig() {
	var monitorNames = make(map[string]bool)
	for _, m := range appConfig.Monitors {
		if _, ok := monitorNames[m.Name]; ok {
			log.Fatalf("Duplicate monitor name: %s", m.Name)
		}
		monitorNames[m.Name] = true
		if ok, err := m.Validate(); !ok {
			log.Fatalf("Invalid monitor configuration: %s", err.Error())
		}
		for _, event := range m.Events {
			if ok, err := event.Validate(); !ok {
				log.Fatalf("Invalid event configuration: %s", err.Error())
			}
			for _, alert := range event.Alerts {
				if ok, err := alert.Validate(); !ok {
					log.Fatalf("Invalid alert configuration: %s", err.Error())
				}
			}
		}
	}
}

func main() {
	ctx := kong.Parse(&app, kong.Name("sitemongoose"), kong.Description("A simple website monitoring tool (version: ${version})"), kong.Vars{"version": APP_VERSION})
	switch ctx.Command() {
	case "start":
		readConfig(app.Start.ConfigFilePath)
		validateConfig()
		moni := monitor.GetService()
		for _, m := range appConfig.Monitors {
			moni.AddMonitor(m)
		}
		moni.Start()
	case "version":
		fmt.Println(APP_VERSION)
	default:
		panic(ctx.Command())
	}
}
