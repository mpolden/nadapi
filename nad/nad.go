package nad

import (
	"bufio"
	"fmt"
	"io"
	"sync"

	"github.com/pkg/term"
)

// Client reprensents a client to the amplifier.
type Client struct {
	port         io.ReadWriteCloser
	mu           sync.Mutex
	EnableVolume bool
}

// New creates a new client to the amplifier, using device for communication.
func New(device string) (*Client, error) {
	// From RS-232 Protocol for NAD Products v2.02:
	//
	// All communication should be done at a rate of 115200 bps with 8 data
	// bits, 1 stop bit and no parity bits. No flow control should be
	// performed.
	port, err := term.Open(device, term.Speed(115200), term.RawMode)
	if err != nil {
		return nil, err
	}
	return &Client{port: port}, nil
}

// Close closes the underlying device
func (n *Client) Close() error {
	return n.port.Close()
}

// SendCmd validates and sends the command cmd to the amplifier.
func (n *Client) SendCmd(cmd Cmd) (Reply, error) {
	// Check if volume adjustment is explicitly enabled. This check is done
	// because incorrect volume adjust might damage your amp, speakers
	// and/or cat.
	if cmd.Variable == "Volume" && !n.EnableVolume {
		return Reply{}, fmt.Errorf("volume adjustment is not enabled")
	}
	if !cmd.Valid() {
		return Reply{}, fmt.Errorf("invalid command: %s", cmd.String())
	}
	b, err := n.Send(cmd.Bytes())
	if err != nil {
		return Reply{}, err
	}
	return ParseReply(b)
}

// SendString parses, validates and sends the command s.
func (n *Client) SendString(s string) (string, error) {
	cmd, err := ParseCmd(s)
	if err != nil {
		return "", err
	}
	reply, err := n.SendCmd(cmd)
	if err != nil {
		return "", err
	}
	return reply.String(), nil
}

// Send sends cmd to the amplifier without any preprocessing or validation.
func (n *Client) Send(cmd []byte) ([]byte, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	reader := bufio.NewReader(n.port)
	// Discard any unread data before sending command
	if t, ok := n.port.(*term.Term); ok {
		if err := t.Flush(); err != nil {
			return nil, err
		}
	}
	if _, err := n.port.Write(cmd); err != nil {
		return nil, err
	}
	// Discard newlines
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		// Rewind when we hit non-newline
		if b != '\n' {
			if err := reader.UnreadByte(); err != nil {
				return nil, err
			}
			break
		}
	}
	b, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Model retrieves the amplifier model.
func (n *Client) Model() (Reply, error) {
	cmd := Cmd{Variable: "Model", Operator: "?"}
	return n.SendCmd(cmd)
}

func (n *Client) enable(variable string, enable bool) (Reply, error) {
	cmd := Cmd{Variable: variable, Operator: "="}
	if enable {
		cmd.Value = "On"
	} else {
		cmd.Value = "Off"
	}
	return n.SendCmd(cmd)
}

// Mute mutes the amplifier.
func (n *Client) Mute(enable bool) (Reply, error) {
	return n.enable("Mute", enable)
}

// Power turns the amplifier on/off.
func (n *Client) Power(enable bool) (Reply, error) {
	return n.enable("Power", enable)
}

// SpeakerA enables/disables output to speaker A.
func (n *Client) SpeakerA(enable bool) (Reply, error) {
	return n.enable("SpeakerA", enable)
}

// SpeakerB enables/disables output to speaker B.
func (n *Client) SpeakerB(enable bool) (Reply, error) {
	return n.enable("SpeakerB", enable)
}

// Tape1 enables/disables output to tape 1.
func (n *Client) Tape1(enable bool) (Reply, error) {
	return n.enable("Tape1", enable)
}

// Source sets the current audio source, specified by name
func (n *Client) Source(name string) (Reply, error) {
	cmd := Cmd{Variable: "Source", Operator: "=", Value: name}
	return n.SendCmd(cmd)
}

// VolumeUp increases volume.
func (n *Client) VolumeUp() (Reply, error) {
	cmd := Cmd{Variable: "Volume", Operator: "+"}
	return n.SendCmd(cmd)
}

// VolumeDown decreases volume.
func (n *Client) VolumeDown() (Reply, error) {
	cmd := Cmd{Variable: "Volume", Operator: "-"}
	return n.SendCmd(cmd)
}
