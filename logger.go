package allegory

import (
	"fmt"
	"github.com/synful/term"
	"os"
)

func Debug(format string, v ...interface{}) {
	term.White(os.Stdout, "[DEBUG] "+fmt.Sprintf(format, v...)+"\n")
}

func Info(format string, v ...interface{}) {
	term.LightGreen(os.Stdout, " [INFO] "+fmt.Sprintf(format, v...)+"\n")
}

func Error(format string, v ...interface{}) {
	term.Red(os.Stderr, "[ERROR] "+fmt.Sprintf(format, v...)+"\n")
}
