package exit

import (
	"fmt"
	"os"
)

func WithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
