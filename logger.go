package allegory

import (
	"fmt"
	"github.com/synful/term"
	"os"
)

func Debug(value interface{}) {
	term.White(os.Stdout, "[DEBUG] "+toString(value)+"\n")
}

func Debugf(format string, v ...interface{}) {
	Debug(fmt.Sprintf(format, v...))
}

func Info(value interface{}) {
	term.LightGreen(os.Stdout, " [INFO] "+toString(value)+"\n")
}

func Infof(format string, v ...interface{}) {
	Info(fmt.Sprintf(format, v...))
}

func Error(value interface{}) {
	term.Red(os.Stderr, "[ERROR] "+toString(value)+"\n")
}

func Errorf(format string, v ...interface{}) {
	Error(fmt.Sprintf(format, v...))
}

func toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case error:
		return v.Error()
	default:
		return fmt.Sprintf("%v", v)
	}
}
