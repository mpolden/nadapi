package nad

import "fmt"

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
	if cmd.Operator == "?" {
		r, ok := p.state[cmd.Variable]
		if !ok {
			return 0, fmt.Errorf("missing initial value for: %s", cmd.String())
		}
		p.reply <- r.Delimited()
	} else {
		p.state[cmd.Variable] = cmd
		p.reply <- cmd.Delimited()
	}
	return len(b), nil
}

func newTestPort() *testPort {
	reply := make(chan string, 1)
	state := make(map[string]Cmd)
	// Initial state
	state["Model"] = Cmd{Variable: "Model", Operator: "=", Value: "C356"}
	state["Mute"] = Cmd{Variable: "Mute", Operator: "=", Value: "Off"}
	state["Power"] = Cmd{Variable: "Power", Operator: "=", Value: "Off"}
	state["Speakera"] = Cmd{Variable: "Speakera", Operator: "=", Value: "On"}
	state["Speakerb"] = Cmd{Variable: "Speakerb", Operator: "=", Value: "Off"}
	state["Tape1"] = Cmd{Variable: "Tape1", Operator: "=", Value: "Off"}
	state["Source"] = Cmd{Variable: "Source", Operator: "=", Value: "CD"}
	return &testPort{reply: reply, state: state}
}

// NewTestClient creates a client that communicates with a simulated amp.
func NewTestClient() *Client {
	port := newTestPort()
	return &Client{port: port, EnableVolume: true}
}
