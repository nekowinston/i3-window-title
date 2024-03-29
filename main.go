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
		Class           string
		Title           string
		Icon            string
		ShowNativeTitle bool `mapstructure:"show_native_title"`
	}

	NativeTitleSeparators struct {
		Start string
		End   string
	} `mapstructure:"native_title_separators"`
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
	viper.SetDefault("native_title_separators.start", "[")
	viper.SetDefault("native_title_separators.end", "]")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
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
}

func main() {
	subscription := i3.Subscribe(i3.WindowEventType, i3.WorkspaceEventType)

	for subscription.Next() {
		event := subscription.Event()
		switch event.(type) {
		case *i3.WindowEvent:
			getActiveWindowTitle(event.(*i3.WindowEvent))
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

	// if the window was closed, check if it should display the workspace, if enabled
	if config.Workspace.Enabled && event.Change == "close" {
		totalNodes := len(event.Container.Nodes) + len(event.Container.FloatingNodes)

		// only do so if it was the last window in the workspace
		if totalNodes == 0 {
			printOutput(config.Workspace.Icon, config.Workspace.Title)
			return
		}
	}

	// check if the wm class matches a mapping
	for _, mapConf := range config.Mappings {
		if strings.ToLower(mapConf.Class) == strings.ToLower(class) {
			// if there's an icon set, prepend it
			if len(mapConf.Icon) > 0 {
				icon = mapConf.Icon
			}
			// use the mapped title
			if len(mapConf.Title) > 0 {
				name = mapConf.Title
				if mapConf.ShowNativeTitle {
					s := config.NativeTitleSeparators
					name = fmt.Sprintf("%s %s%s%s", name, s.Start, title, s.End)
				}
			} else {
				name = title
			}
			printOutput(icon, name)
			return
		}
	}

	// use the native WM title if no mapping is found
	name = title
	// capitalize the title if set in the config
	if config.Capitalize {
		name = strings.Title(name)
	}
	printOutput(icon, name)
}

func showWorkspaceMapping(event *i3.WorkspaceEvent) {
	totalNodes := len(event.Current.Nodes) + len(event.Current.FloatingNodes)
	focused := event.Current.Focused

	if totalNodes == 0 && focused {
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
