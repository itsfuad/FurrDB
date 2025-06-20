package repl

import (
	"bufio"
	"bytes"
	"furr/internal/db"
	"strings"
	"testing"
)

func TestReplCoreLogic(t *testing.T) {
	input := "SET foo bar\nGET foo\nEXISTS foo\nDEL foo\nEXISTS foo\nEXIT\n"
	r := strings.NewReader(input)
	w := &bytes.Buffer{}

	// Patch the REPL to use our reader/writer
	repl := NewTestRepl(r, w)
	repl.Run()

	out := w.String()
	if !strings.Contains(out, "bar") {
		t.Errorf("expected output to contain 'bar', got: %s", out)
	}
	if !strings.Contains(out, "1") || !strings.Contains(out, "0") {
		t.Errorf("expected output to contain '1' and '0', got: %s", out)
	}
}

// NewTestRepl and Run for testable REPL

type TestRepl struct {
	in  *strings.Reader
	out *bytes.Buffer
}

func NewTestRepl(in *strings.Reader, out *bytes.Buffer) *TestRepl {
	return &TestRepl{in, out}
}

func (r *TestRepl) Run() {
	scanner := bufio.NewScanner(r.in)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		tokens := strings.Fields(line)
		if len(tokens) == 0 {
			continue
		}
		cmd := strings.ToUpper(tokens[0])
		args := tokens[1:]
		if cmd == "EXIT" {
			r.out.WriteString("bye!\n")
			return
		}
		handler, ok := db.Commands[cmd]
		if !ok {
			r.out.WriteString("ERR unknown command\n")
			continue
		}
		result, err := handler(args)
		if err != nil {
			r.out.WriteString("ERR " + err.Error() + "\n")
			continue
		}
		r.out.WriteString(result + "\n")
	}
}
