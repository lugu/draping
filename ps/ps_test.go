package ps

import (
	"fmt"
	"io"
	"bytes"
	"testing"
)

func TestScanln(t *testing.T) {
	var in = ` 0.0  0.0 56919 tmux
0.2  0.0 38558 zsh
0.0  0.0 75110 ps
0.0  0.0 89110 tmux
`
	reader := bytes.NewBuffer([]byte(in))
 	var cpu, mem float32
	var pid, comm string
	n, err := fmt.Fscanln(reader, &cpu, &mem, &pid, &comm)
	if n != 4 {
		t.Errorf("1: failed to parse 4 value: %d (%s)", n, err)
	} else if err != nil {
		t.Errorf("1: failed to parse: %s", err)
	}
	n, err = fmt.Fscanln(reader, &cpu, &mem, &pid, &comm)
	if n != 4 {
		t.Errorf("2: failed to parse 4 value: %d (%s)", n, err)
	} else if err != nil {
		t.Errorf("2: failed to parse: %s", err)
	}
	n, err = fmt.Fscanln(reader, &cpu, &mem, &pid, &comm)
	if n != 4 {
		t.Errorf("3: failed to parse 4 value: %d (%s)", n, err)
	} else if err != nil {
		t.Errorf("3: failed to parse: %s", err)
	}
	n, err = fmt.Fscanln(reader, &cpu, &mem, &pid, &comm)
	if n != 4 {
		t.Errorf("4: failed to parse 4 value: %d (%s)", n, err)
	} else if err != nil {
		t.Errorf("4: failed to parse: %s", err)
	}
	n, err = fmt.Fscanln(reader, &cpu, &mem, &pid, &comm)
	if n != 0 {
		t.Errorf("5: failed to parse 0 value: %d", n)
	} else if err != io.EOF {
		t.Errorf("5: failed to parse: %s", err)
	}
}

func TestPS(t *testing.T) {
	stats, err := PS()
	if err != nil {
		t.Fatal(err)
	}
	if len(stats) == 0 {
		t.Fatal("no stats")
	}
}
