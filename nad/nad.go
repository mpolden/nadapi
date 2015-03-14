package nad

import (
	"bufio"
	"bytes"
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

func (n *Client) Send(cmd Cmd) ([]byte, error) {
	return n.SendString(cmd.Delimited())
}

func (n *Client) SendString(cmd string) ([]byte, error) {
	_, err := n.port.Write([]byte(cmd))
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(n.port)
	reply, err := reader.ReadBytes('\r')
	if err != nil {
		return nil, err
	}
	return bytes.TrimRight(reply, "\r"), nil
}

func (n *Client) Model() (string, error) {
	cmd := Cmd{Variable: "Model", Operator: "?"}
	b, err := n.Send(cmd)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (n *Client) enable(variable string, enable bool) (string, error) {
	cmd := Cmd{Variable: variable, Operator: "="}
	if enable {
		cmd.Value = "On"
	} else {
		cmd.Value = "Off"
	}
	b, err := n.Send(cmd)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (n *Client) Mute(enable bool) (string, error) {
	return n.enable("Mute", enable)
}

func (n *Client) Power(enable bool) (string, error) {
	return n.enable("Power", enable)
}

func (n *Client) SpeakerA(enable bool) (string, error) {
	return n.enable("SpeakerA", enable)
}

func (n *Client) SpeakerB(enable bool) (string, error) {
	return n.enable("SpeakerB", enable)
}

func (n *Client) Tape1(enable bool) (string, error) {
	return n.enable("Tape1", enable)
}

func (n *Client) Source(source Source) (string, error) {
	cmd := Cmd{Variable: "Source", Operator: "=", Value: string(source)}
	b, err := n.Send(cmd)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (n *Client) VolumeUp() (string, error) {
	cmd := Cmd{Variable: "Volume", Operator: "+"}
	b, err := n.Send(cmd)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (n *Client) VolumeDown() (string, error) {
	cmd := Cmd{Variable: "Volume", Operator: "-"}
	b, err := n.Send(cmd)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
