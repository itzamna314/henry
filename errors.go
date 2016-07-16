package henry

import (
	"fmt"
)

type UsageError struct {
	UsageMessage string
}

func (u UsageError) Error() string {
	return fmt.Sprintf("Handled illegal user input.  Expected: '%s'", u.UsageMessage)
}
