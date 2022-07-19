package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.i3wm.org/i3/v4"
	"strings"
)

type Config struct {
	DefaultIcon string `mapstructure:"default_icon"`
	Padding     int
	Capitalize  bool

	Workspace struct {
		Enabled bool
		Icon    string
		Title   string
	}
	Mappings []struct {
		Class string
		Title string
		Icon  string
	}
}

var config Config

func init() {
	viper.SetConfigName("window_titles")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("$XDG_CONFIG_HOME")
	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath(".")

	viper.SetDefault("default_icon", "ﬓ")
	viper.SetDefault("padding", 1)
	viper.SetDefault("capitalize", false)
	viper.SetDefault("workspace.enabled", true)
	viper.SetDefault("workspace.icon", "ﬓ")
	viper.SetDefault("workspace.title", "Desktop")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("Could not read config file:", err)
	}

	// update on config file change
	viper.OnConfigChange(func(e fsnotify.Event) {
		viper.Unmarshal(&config)
	})
	viper.WatchConfig()
	fmt.Println(config)
}

func main() {
	subscription := i3.Subscribe(i3.WindowEventType, i3.WorkspaceEventType)

	for subscription.Next() {
		event := subscription.Event()
		switch event.(type) {
		case *i3.WindowEvent:
			getActiveWindowTitle(event.(*i3.WindowEvent))

			// NB: unsure if this is the best way to handle this:
			// this skips to the next event, which is the WorkspaceEvent
			// otherwise the default mapping would show up
			subscription.Next()
		case *i3.WorkspaceEvent:
			// only handle this event if the default title is set
			if config.Workspace.Enabled {
				showWorkspaceMapping(event.(*i3.WorkspaceEvent))
			}
		}
	}
}

func getActiveWindowTitle(event *i3.WindowEvent) {
	// shorthands
	class := event.Container.WindowProperties.Class
	title := event.Container.WindowProperties.Title
	icon := config.DefaultIcon

	var name string

	for _, mapConf := range config.Mappings {
		if strings.ToLower(mapConf.Class) == strings.ToLower(class) {
			// check if the wm class matches a mapping
			// if there's an icon set, prepend it
			if len(mapConf.Icon) > 0 {
				icon = mapConf.Icon
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
	printOutput(icon, name)
}

func showWorkspaceMapping(event *i3.WorkspaceEvent) {
	if len(event.Current.Nodes) == 0 {
		printOutput(config.Workspace.Icon, config.Workspace.Title)
	}
}

func printOutput(icon string, title string) {
	// one space is added here, to ensure that NF render correctly
	// and spaces after the icon are added, depending on the padding
	padding := " " + strings.Repeat(" ", config.Padding)
	icon = icon + padding
	v := fmt.Sprintf("%s%s", icon, title)
	fmt.Println(v)
}
