# i3-window-title

A small helper application to display the focused window's title in polybar, or similar.

It uses i3wm's official Go event subscription, so it updates instantly without polling.

![Preview](assets/preview.gif)

## Configuration

The config file is YAML format, located in `~/.config/window_titles.yml`:

```yaml
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
