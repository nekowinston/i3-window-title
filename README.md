<h3 align="center">
  i3-window-title
</h3>

<p align="center">
  <img alt="License: MIT" src="https://img.shields.io/github/license/nekowinston/i3-window-title">
  <img alt="Go version" src="https://img.shields.io/github/go-mod/go-version/nekowinston/i3-window-title">
  <img alt="Downloads" src="https://img.shields.io/github/downloads/nekowinston/i3-window-title/total">
  <img alt="goreleaser" src="https://github.com/nekowinston/i3-window-title/actions/workflows/release.yml/badge.svg">
</p>

A small helper application to display the focused window's title in polybar, or similar.

It uses i3wm's official Go event subscription, so it updates instantly without polling.

![Preview](assets/preview.gif)

## Installation

Download from [GitHub releases](https://github.com/nekowinston/i3-window-title/releases),
and move it to one of your `bin` directories, which is in your `$PATH`.

## Configuration

The config file is YAML format, located in `~/.config/window_titles.yml`:

```yaml
# the icon to fall back to, if there is no mapping
default_icon: ﬓ
# spaces to insert between the icon & window title
padding: 1
# auto-capitalize window titles, when they are not manually set.
capitalize: true

# show workspace info, when no window is focused.
workspace:
  # turning this off will always show the last window title
  enabled: true 
  # icon and title mapping for the workspace placeholder
  icon: ﬓ
  title: Desktop

# mappings for each application
mappings:
# class is the WM_CLASS, which you can get via xprop.
# NB: you can ignore casing, it doesn't matter if lower/uppercase
- class: org.wezfurlong.wezterm
  title: WezTerm
  icon: 
- class: firefox
  title: Firefox
  icon: 
- class: Pcmanfm
  title: PCManFM
  icon: 
```

## Usage (polybar)

```ini
[module/window-title]
type = custom/script
exec = i3-window-title

; important! this is required to only show the last printed line in the bar
tail = true
```
