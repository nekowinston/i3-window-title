package main

import (
	"fmt"
	"go.i3wm.org/i3/v4"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

type MapConf []struct {
	Class string `yaml:"class"`
	Title string `yaml:"title"`
	Icon  string `yaml:"icon"`
}

var config MapConf

func main() {
	loadConfig()

	subscription := i3.Subscribe(i3.WindowEventType)

	for subscription.Next() {
		event := subscription.Event()
		switch event.(type) {
		case *i3.WindowEvent:
			getActiveWindowTitle(event.(*i3.WindowEvent))
		}
	}
}

func getActiveWindowTitle(event *i3.WindowEvent) {
	title := event.Container.WindowProperties.Class
	var icon string

	for _, mapConf := range config {
		if mapConf.Class == title {
			icon = mapConf.Icon
			title = mapConf.Title
		}
	}

	fmt.Printf(`%s  %s`, icon, title)
	fmt.Println()
}

func loadConfig() {
	// get XDG_CONFIG_HOME, fallback to $HOME/.config
	home := os.Getenv("XDG_CONFIG_HOME")
	if len(home) == 0 {
		home = os.Getenv("HOME") + "/.config"
	}

	r, err := os.Open(path.Join(home, "window_titles.yml"))
	if err != nil {
		fmt.Println("Could not open config file:", err)
	}
	defer r.Close()

	decoder := yaml.NewDecoder(r)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Could not decode config file:", err)
	}
}
