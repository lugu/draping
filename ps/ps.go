package ps

import (
	"io"
	"bytes"
	"io/ioutil"
	"bufio"
	"fmt"
	"os/exec"
)

type PidStatus struct {
	PID int
	CPU float32
	Mem float32
	Command string
}

type Stats map[string][]PidStatus

func PS() ([]PidStatus, error) {

	cmd := exec.Command("ps", "-ewwo", "%cpu,%mem,pid,comm")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(stdout)

	// $ ps -ewwo %cpu,%mem,pid,comm
	// %CPU %MEM   PID COMMAND
	//  0.0  0.0 56919 tmux
	//  0.2  0.0 38558 zsh
	//  0.0  0.0 75110 ps
	//  0.0  0.0 89110 tmux
	reader := bufio.NewReader(bytes.NewBuffer(buf))
	reader.ReadLine()
	stats := make([]PidStatus, 0)
	for {
		var status PidStatus
		n, err := fmt.Fscanln(reader,
			&status.CPU, &status.Mem, &status.PID, &status.Command)
		if err == io.EOF {
			break
		} else if n != 4 || err != nil {
			return nil, fmt.Errorf("reading %d values %s: %#v",
				n, err, stats)
		}
		stats = append(stats, status)
	}
	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	return stats, nil
}
