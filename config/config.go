// Package config provides support for getting and setting configuration values.
package config

import (
	al "github.com/dradtke/go-allegro/allegro"
)

var (
	blank_color    al.Color
	fps            = 60
	display_width  = 640
	display_height = 480
    display_flags = al.WINDOWED
)

const CONSOLE_FILE = "build/console.txt"

func init() {
	blank_color = al.MapRGB(0, 0, 0)
}

func Fps() int {
	return fps
}

func SetFps(value int) {
	fps = value
}

func BlankColor() al.Color {
	return blank_color
}

func SetBlankColor(value al.Color) {
    blank_color = value
}

func DisplaySize() (w, h int) {
	return display_width, display_height
}

func SetDisplaySize(w, h int) {
	display_width = w
	display_height = h
}

func DisplayFlags() al.DisplayFlags {
    return display_flags
}

func SetDisplayFlags(value al.DisplayFlags) {
    display_flags = value
}

func GameName() string {
	return "Hello World"
}
