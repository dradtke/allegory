// Package config provides support for getting and setting configuration values.
package config

import (
	"github.com/dradtke/go-allegro/allegro"
)

// Default values.
var (
	title          string
	pkg_root       string
	blank_color    allegro.Color
	fps            = 60
	display_width  = 640
	display_height = 480
	display_flags  = allegro.WINDOWED
)

const CONSOLE_FILE = "build/console.txt"

func init() {
	blank_color = allegro.MapRGB(0, 0, 0)
}

func Fps() int {
	return fps
}

func SetFps(value int) {
	fps = value
}

func BlankColor() allegro.Color {
	return blank_color
}

func SetBlankColor(value allegro.Color) {
	blank_color = value
}

func DisplaySize() (w, h int) {
	return display_width, display_height
}

func SetDisplaySize(w, h int) {
	display_width = w
	display_height = h
}

func DisplayFlags() allegro.DisplayFlags {
	return display_flags
}

func SetDisplayFlags(value allegro.DisplayFlags) {
	display_flags = value
}

func SetWindowTitle(value string) {
	title = value
}

func WindowTitle() string {
	if title == "" {
		return "Untitled"
	} else {
		return title
	}
}

func SetPackageRoot(value string) {
	pkg_root = value
}

func PackageRoot() string {
	return pkg_root
}
