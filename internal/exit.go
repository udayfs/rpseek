package internal

import (
	"fmt"
	"os"
	"strings"
)

func ExitOnError(message ...string) {
	fmt.Fprintln(os.Stderr, "Error:", strings.Join(message, " "))
	os.Exit(1)
}
