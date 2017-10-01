package miles

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

type errReader struct{}

var errRead = errors.New("err reader")

func (r errReader) Read(buf []byte) (int, error) {
	return 0, errRead
}

type errWriter struct{}

var errWrite = errors.New("err writer")

func (w errWriter) Write(buf []byte) (int, error) {
	return 0, errWrite
}

func TestOptionsChoose(t *testing.T) {

	tt := []struct {
		name string
		opt  Options
		ans  string
		err  error
	}{
		{
			name: "nil reader",
			opt:  Options{},
			err:  errNilFrom,
		},
		{
			name: "give up for wrong input even if default is provided and empty allowed",
			opt: Options{
				From:        strings.NewReader("x\nx\nx\na"),
				MaxAttempts: 2,
				Allowed:     []string{"a"},
				Default:     "a",
				AllowEmpty:  true,
			},
			err: errGiveUp,
		},
		{
			name: "custom max attempts",
			opt: Options{
				From:        strings.NewReader("x\nx\nx\nx\nx\nx\na"),
				MaxAttempts: 7,
				Allowed:     []string{"a"},
			},
			ans: "a",
		},
		{
			name: "valid option",
			opt: Options{
				From:    strings.NewReader("a"),
				Allowed: []string{"a"},
			},
			ans: "a",
		},
		{
			name: "valid option on second attempt",
			opt: Options{
				From:    strings.NewReader("b\na"),
				Allowed: []string{"a"},
			},
			ans: "a",
		},
		{
			name: "valid option with different case",
			opt: Options{
				From:    strings.NewReader("A"),
				Allowed: []string{"a"},
			},
			ans: "a",
		},
		{
			name: "allowed empty",
			opt: Options{
				From:       strings.NewReader(""),
				AllowEmpty: true,
			},
			ans: "",
		},
		{
			name: "empty not allowed",
			opt: Options{
				From:        strings.NewReader("\n\n\n"),
				MaxAttempts: 2,
			},
			err: errGiveUp,
		},

		{
			name: "default value",
			opt: Options{
				From:    strings.NewReader(""),
				Default: "a",
			},
			ans: "a",
		},
		{
			name: "reader error",
			opt: Options{
				From: errReader{},
			},
			err: errRead,
		},
		{
			name: "err writer",
			opt: Options{
				To:     errWriter{},
				Prompt: "PROMPT",
			},
			err: errWrite,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ans, err := tc.opt.Choose()
			if err != nil && tc.err == nil {
				t.Errorf("expected err to be nil; was %v", err)
			}
			if err == nil && tc.err != nil {
				t.Errorf("expected err to be %v; was nil", tc.err)
			}
			if ans != tc.ans {
				t.Errorf("expected ans to be %s; was %s", tc.ans, ans)
			}
		})
	}
}

func TestOptionsChooseWithPrompt(t *testing.T) {
	tt := []struct {
		name   string
		buf    *bytes.Buffer
		opt    Options
		err    error
		prompt string
	}{
		{
			name: "options appended to prompt",
			buf:  bytes.NewBufferString(""),
			opt: Options{
				From:    strings.NewReader("x\na"),
				Prompt:  "PROMPT",
				Allowed: []string{"a", "b"},
				Default: "a",
			},
			prompt: "PROMPT [A/b]: PROMPT [A/b]: ",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.opt.To = tc.buf
			_, err := tc.opt.Choose()
			if tc.err != nil && err != tc.err {
				t.Error(err)
			}
			prompt := tc.buf.String()
			if prompt != tc.prompt {
				t.Errorf("expected %s; was %s", tc.prompt, prompt)
			}
		})
	}
}

func ExampleOptions_Choose() {
	opt := Options{
		From:    strings.NewReader("Y"),
		Prompt:  "Is it cool?",
		Allowed: []string{"y", "n"},
		Default: "y",
	}
	_, err := opt.Choose()
	if err != nil {
		// Handle error
	}
	// The first value returned (here discarded) would contain the chosen value
	// Output: Is it cool? [Y/n]:
}
