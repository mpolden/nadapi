package nad

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"io"
)

type NAD struct {
	port io.ReadWriteCloser
}

func New(device string) (NAD, error) {
	options := serial.OpenOptions{
		PortName:        device,
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
		ParityMode:      serial.PARITY_NONE,
	}
	port, err := serial.Open(options)
	if err != nil {
		return NAD{}, err
	}
	return NAD{port: port}, nil
}

func (n *NAD) Send(cmd string) ([]byte, error) {
	cmd = fmt.Sprintf("\r%s\r", cmd)
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

func (n *NAD) Model() (string, error) {
	b, err := n.Send("Main.Model?")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (n *NAD) Mute(enable bool) (string, error) {
	var cmd string
	if enable {
		cmd = "Main.Mute=On"
	} else {
		cmd = "Main.Mute=Off"
	}
	b, err := n.Send(cmd)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
