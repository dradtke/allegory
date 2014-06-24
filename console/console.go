// Package console implements an in-game console for logging and running commands.
package console

import (
	"container/list"
	"fmt"
	al "github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/gopher/bus"
	"github.com/dradtke/gopher/config"
	"github.com/dradtke/gopher/graphics"
	"io"
	"os"
	"unicode"
)

const (
	SUMMARY_LENGTH = 30
	BLINK_SPEED    = 0.6
)

var (
	blinker   *al.Timer
	cmd       string
	color_map map[string]al.Color
	is_blunk  bool
	log       list.List
	visible   bool
)

type message struct {
	level string
	text  string
}

func (m message) String() string {
	return fmt.Sprintf("[%s] %s", m.level, m.text)
}

func (m message) Line() graphics.Line {
	return graphics.Line{Text: m.String(), Color: color_map[m.level]}
}

func getSummary() (sum []message) {
	sum = make([]message, 0, SUMMARY_LENGTH)
	for e, i := log.Back(), 0; e != nil && i < SUMMARY_LENGTH; e, i = e.Prev(), i+1 {
		sum = append(sum, e.Value.(message))
	}
	return
}

func Debug(msg string) {
	log.PushBack(message{level: "DEBUG", text: msg})
}

func Debugf(msg string, v ...interface{}) {
	log.PushBack(message{level: "DEBUG", text: fmt.Sprintf(msg, v...)})
}

func Info(msg string) {
	log.PushBack(message{level: "INFO", text: msg})
}

func Infof(msg string, v ...interface{}) {
	log.PushBack(message{level: "INFO", text: fmt.Sprintf(msg, v...)})
}

func Error(msg string) {
	log.PushBack(message{level: "ERROR", text: msg})
}

func Errorf(msg string, v ...interface{}) {
	log.PushBack(message{level: "ERROR", text: fmt.Sprintf(msg, v...)})
}

func Fatal(msg string) {
	log.PushBack(message{level: "FATAL", text: msg})
}

func Fatalf(msg string, v ...interface{}) {
	log.PushBack(message{level: "FATAL", text: fmt.Sprintf(msg, v...)})
}

func Init(eventQueue *al.EventQueue) {
	var err error
	if blinker, err = al.CreateTimer(BLINK_SPEED); err != nil {
		panic(err)
	}
	eventQueue.Register(blinker)
	blinker.Start()

	color_map = map[string]al.Color{
		"DEBUG": al.MapRGB(0, 0, 255),
		"INFO":  al.MapRGB(0, 255, 0),
		"ERROR": al.MapRGB(255, 0, 0),
		"FATAL": al.MapRGB(255, 0, 0),
	}
}

func Render() {
	if !visible {
		return
	}
	sum := getSummary()
	lines := make([]graphics.Line, len(sum))
	for i, msg := range sum {
		lines[i] = msg.Line()
	}
	graphics.RenderConsole(lines, cmd, is_blunk)
}

func HandleEvent(ev interface{}) bool {
	switch e := ev.(type) {

	case al.KeyDownEvent:
		switch e.KeyCode() {

		case al.KEY_F12:
			visible = !visible
			return true
		}

	case al.KeyCharEvent:
		if visible {
			switch e.KeyCode() {

			case al.KEY_BACKSPACE:
				backspaceCmd()
				return true

			case al.KEY_ENTER:
				submitCmd()
				return true

			default:
				unichar := rune(e.Unichar())
				if unicode.IsPrint(unichar) {
					cmd += string(unichar)
					return true
				}
			}
		}

	case al.TimerEvent:
		if e.Source() == blinker {
			is_blunk = !is_blunk
			return true
		}
	}

	return false
}

func Save() {
	if f, err := os.Create(config.CONSOLE_FILE); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	} else {
		for e := log.Front(); e != nil; e = e.Next() {
			io.WriteString(f, e.Value.(message).String()+"\n")
		}
		f.Close()
	}
}

func Toggle() {
	visible = !visible
}

func backspaceCmd() {
	if cmd == "" {
		return
	}
	cmd = cmd[:len(cmd)-1]
}

func submitCmd() {
	if cmd == "" {
		return
	}
	bus.Signal(bus.ConsoleCommandEvent, cmd)
	cmd = ""
}
