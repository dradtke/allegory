package allegory

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func readStdin() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		_stdin <- scanner.Text()
	}
}

func ParseAssignment(str string) (name, value string, err error) {
	eq := strings.Index(str, "=")
	if eq == -1 {
		return "", "", errors.New("not an assignment")
	}
	return strings.TrimSpace(str[:eq]), strings.TrimSpace(str[eq+1:]), nil
}
