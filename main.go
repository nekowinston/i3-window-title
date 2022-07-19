package main

import (
	"fmt"
	"go.i3wm.org/i3/v4"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strings"
)

type Config struct {
	DefaultIcon string `yaml:"default_icon"`
	Padding     int    `yaml:"padding"`
	Capitalize  bool   `yaml:"capitalize"`

	Mappings []struct {
		Class string `yaml:"class"`
		Title string `yaml:"title"`
		Icon  string `yaml:"icon"`
	}
}

var config Config

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
	// one space is added here, to ensure that NF render correctly
	// and spaces after the icon are added, depending on the padding
	padding := " " + strings.Repeat(" ", config.Padding)
	// shorthands
	class := event.Container.WindowProperties.Class
	title := event.Container.WindowProperties.Title
	// default to the default icon
	icon := config.DefaultIcon + padding

	var name string

	for _, mapConf := range config.Mappings {
		if strings.ToLower(mapConf.Class) == strings.ToLower(class) {
			// check if the wm class matches a mapping
			// if there's an icon set, prepend it
			if len(mapConf.Icon) > 0 {
				icon = mapConf.Icon + padding
			}
			// use the mapped title
			if len(mapConf.Title) > 0 {
				name = mapConf.Title
			}
			break
		} else {
			// if there's no mapping for the wm class, use the native title
			name = title
			// and finally, capitalize the title if set in the config
			if config.Capitalize {
				name = strings.Title(name)
			}
		}
	}

	v := fmt.Sprintf("%s%s", icon, name)
	fmt.Println(v)
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
