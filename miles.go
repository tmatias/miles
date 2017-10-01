// Package miles provides helpers for getting user input.
package miles

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var errNilFrom = errors.New("miles: nil reader")
var errGiveUp = errors.New("miles: gave up")

// Options defines a way to ask the user to choose from a set of possible values.
type Options struct {
	// From is the reader to read the option from.
	From io.Reader
	// To is the writer to write the prompt to.
	To io.Writer
	// Prompt is the text to show to the user before waiting for his input.
	Prompt string
	// Allowed options for the user to choose from (will be appended to Prompt).
	Allowed []string
	// Default value returned if AllowEmpty is true.
	Default string
	// AllowEmpty allows the user to proceed with an empty value.
	AllowEmpty bool
	// MaxAttempts before giving up.
	MaxAttempts int
}

// Choose asks the user to choose from a set of possible values.
//
// If the user enters an invalid option, he will be prompted again until
// MaxAttempts is reached (five times, unless specified otherwise).
//
// If we reach MaxAttempts and the user never enters a valid option, then an
// error is returned, even if a Default is provided and AllowEmpty is true.
//
// When AllowEmpty is set to true and the user enters a blank value, a blank
// value is returned (empty spaces befored and after the entered value are
// always discarded).
//
// If the user enters a blank value and a Default is provided, then the Default
// value will be returned (even if AllowEmpty is not set to true).
//
// The prompt presented to the user can be customized by setting the Prompt
// property. A string describing the valid options (with the default option
// highlighted) is appended to the Prompt set before printing. For example:
//
//     opt := miles.Options{
//         Prompt: "Choose",
//         Allowed: []string{"a", "b", "c"},
//         Default: "a",
//     }
//     chosen, err := opt.Choose()
//
// would generate the prompt:
//
//     Choose [A/b/c]:
//
// If no Prompt is specified, nothing is printed.
//
// If From or To are not specified, os.Stdin and os.Stdout are used.
//
// Options are always case-insensitive. When building the prompt, the default
// value (if provided) will be uppercased and all others lowercased. The
// returned value will always be lowercased.
func (o Options) Choose() (string, error) {
	if o.From == nil {
		return "", errNilFrom
	}
	var ans string
	attempts := 0
	maxAttempts := o.MaxAttempts
	if maxAttempts < 1 {
		maxAttempts = 5
	}
	from := o.From
	if from == nil {
		from = os.Stdin
	}
	to := o.To
	if to == nil {
		to = os.Stdout
	}
	s := bufio.NewScanner(from)
	for {
		if o.Prompt != "" {
			fmt.Fprint(to, buildPrompt(o.Prompt, o.Default, o.Allowed...))
		}
		s.Scan()
		err := s.Err()
		if err != nil {
			return "", err
		}
		ans = strings.ToLower(strings.Trim(s.Text(), " \n"))
		if ans == "" && (o.AllowEmpty || o.Default != "") {
			return o.Default, nil
		}
		for _, v := range o.Allowed {
			if ans == v {
				return ans, nil
			}
		}
		attempts++
		if attempts >= maxAttempts {
			return "", errGiveUp
		}
	}
}

func buildPrompt(prompt string, def string, allowed ...string) string {
	options := ""
	sep := ""
	for _, v := range allowed {
		if v == def {
			v = strings.ToUpper(v)
		} else {
			v = strings.ToLower(v)
		}
		options = fmt.Sprintf("%s%s%s", options, sep, v)
		sep = "/"
	}
	return fmt.Sprintf("%s [%s]: ", prompt, options)
}
