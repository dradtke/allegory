// Package console implements an in-game console for logging and running commands.
package console

import (
	"container/list"
	"fmt"
	"github.com/dradtke/allegory/bus"
	"github.com/dradtke/allegory/config"
	"github.com/dradtke/allegory/graphics"
	"github.com/dradtke/go-allegro/allegro"
	"io"
	"os"
	"unicode"
)

const (
	SUMMARY_LENGTH = 30
	BLINK_SPEED    = 0.6
)

var (
	blinker   *allegro.Timer
	cmd       string
	color_map map[string]allegro.Color
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

func Init(eventQueue *allegro.EventQueue) {
	var err error
	if blinker, err = allegro.CreateTimer(BLINK_SPEED); err != nil {
		panic(err)
	}
	eventQueue.Register(blinker)
	blinker.Start()

	color_map = map[string]allegro.Color{
		"DEBUG": allegro.MapRGB(0, 0, 255),
		"INFO":  allegro.MapRGB(0, 255, 0),
		"ERROR": allegro.MapRGB(255, 0, 0),
		"FATAL": allegro.MapRGB(255, 0, 0),
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

	case allegro.KeyDownEvent:
		switch e.KeyCode() {

		case allegro.KEY_F12:
			visible = !visible
			return true
		}

	case allegro.KeyCharEvent:
		if visible {
			switch e.KeyCode() {

			case allegro.KEY_BACKSPACE:
				backspaceCmd()
				return true

			case allegro.KEY_ENTER:
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

	case allegro.TimerEvent:
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
