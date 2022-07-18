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

Move it to one of your `bin` directories, which is in your `$PATH`.

## Configuration

The config file is YAML format, located in `~/.config/window_titles.yml`:

```yaml
# the icon to fall back to, if there is no mapping
default_icon: ﬓ
# spaces to insert between the icon & window title
padding: 1

mappings:
# class is the WM_CLASS, which you can get via xprop.
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
