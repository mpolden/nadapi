package nad

import (
	"fmt"
	"io"
	"strings"
)

type testPort struct {
	reply chan string
	state map[string]Cmd
}

func (p *testPort) Close() (err error) { return }

func (p *testPort) Read(b []byte) (n int, err error) {
	r := []byte(<-p.reply)
	copy(b, r)
	return len(r), nil
}

func (p *testPort) Write(b []byte) (n int, err error) {
	cmd, err := ParseCmd(string(b))
	if err != nil {
		return 0, err
	}
	if !cmd.Valid() {
		return 0, fmt.Errorf("invalid command: %s", cmd.String())
	}
	v := strings.ToLower(cmd.Variable)
	if cmd.Operator == "?" {
		r, ok := p.state[v]
		if !ok {
			return 0, fmt.Errorf("missing initial value for: %s", cmd.String())
		}
		p.reply <- r.Delimited()
	} else {
		p.state[v] = cmd
		p.reply <- cmd.Delimited()
	}
	return len(b), nil
}

func newTestPort(string) (io.ReadWriteCloser, error) {
	reply := make(chan string, 1)
	state := make(map[string]Cmd)
	// Initial state
	state["model"] = Cmd{Variable: "Model", Operator: "=", Value: "C356BEE"}
	state["mute"] = Cmd{Variable: "Mute", Operator: "=", Value: "Off"}
	state["power"] = Cmd{Variable: "Power", Operator: "=", Value: "Off"}
	state["speakera"] = Cmd{Variable: "SpeakerA", Operator: "=", Value: "On"}
	state["speakerb"] = Cmd{Variable: "SpeakerB", Operator: "=", Value: "Off"}
	state["tape1"] = Cmd{Variable: "Tape1", Operator: "=", Value: "Off"}
	state["source"] = Cmd{Variable: "Source", Operator: "=", Value: "CD"}
	return &testPort{reply: reply, state: state}, nil
}

// NewTestClient creates a client that communicates with a simulated amp.
func NewTestClient() *Client {
	device := &device{
		name:         "/dev/foo",
		openPort:     newTestPort,
		evalSymlinks: func(string) (string, error) { return "/dev/realfoo", nil },
	}
	return &Client{device: device, EnableVolume: true}
}
