package nad

import (
	"bufio"
	"github.com/pkg/term"
	"io"
)

type Source string

const (
	CD    Source = "CD"
	Tuner Source = "Tuner"
	Video Source = "Video"
	Disc  Source = "Disc"
	Ipod  Source = "Ipod"
	Tape2 Source = "Tape2"
	Aux   Source = "Aux"
)

type Client struct {
	port io.ReadWriteCloser
}

func New(device string) (Client, error) {
	// From RS-232 Protocol for NAD Products v2.02:
	//
	// All communication should be done at a rate of 115200 bps with 8 data
	// bits, 1 stop bit and no parity bits. No flow control should be
	// performed.
	port, err := term.Open(device, term.Speed(115200))
	if err != nil {
		return Client{}, err
	}
	return Client{port: port}, nil
}

func (n *Client) SendCmd(cmd Cmd) (Cmd, error) {
	b, err := n.Send([]byte(cmd.Delimited()))
	if err != nil {
		return Cmd{}, err
	}
	return ParseCmd(string(b))
}

func (n *Client) Send(cmd []byte) ([]byte, error) {
	if _, err := n.port.Write(cmd); err != nil {
		return nil, err
	}
	reader := bufio.NewReader(n.port)
	b, err := reader.ReadBytes('\r')
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (n *Client) Model() (Cmd, error) {
	cmd := Cmd{Variable: "Model", Operator: "?"}
	return n.SendCmd(cmd)
}

func (n *Client) enable(variable string, enable bool) (Cmd, error) {
	cmd := Cmd{Variable: variable, Operator: "="}
	if enable {
		cmd.Value = "On"
	} else {
		cmd.Value = "Off"
	}
	return n.SendCmd(cmd)
}

func (n *Client) Mute(enable bool) (Cmd, error) {
	return n.enable("Mute", enable)
}

func (n *Client) Power(enable bool) (Cmd, error) {
	return n.enable("Power", enable)
}

func (n *Client) SpeakerA(enable bool) (Cmd, error) {
	return n.enable("SpeakerA", enable)
}

func (n *Client) SpeakerB(enable bool) (Cmd, error) {
	return n.enable("SpeakerB", enable)
}

func (n *Client) Tape1(enable bool) (Cmd, error) {
	return n.enable("Tape1", enable)
}

func (n *Client) Source(source Source) (Cmd, error) {
	cmd := Cmd{Variable: "Source", Operator: "=", Value: string(source)}
	return n.SendCmd(cmd)
}

func (n *Client) VolumeUp() (Cmd, error) {
	cmd := Cmd{Variable: "Volume", Operator: "+"}
	return n.SendCmd(cmd)
}

func (n *Client) VolumeDown() (Cmd, error) {
	cmd := Cmd{Variable: "Volume", Operator: "-"}
	return n.SendCmd(cmd)
}
