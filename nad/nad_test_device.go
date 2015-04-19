package nad

import "fmt"

type testPort struct {
	reply chan string
}

func (p *testPort) Close() (err error) {
	return
}

func (p *testPort) Read(b []byte) (n int, err error) {
	r := []byte(<-p.reply)
	copy(b, r)
	return len(r), nil
}

func (p *testPort) Write(b []byte) (n int, err error) {
	cmd := string(b)
	switch cmd {
	case "\rMain.Model?\r":
		p.reply <- "Main.Model=C356\n"
	case "\rMain.Mute=On\r":
		p.reply <- "\nMain.Mute=On\n"
	case "\rMain.Mute=Off\r":
		p.reply <- "\n\nMain.Mute=Off\n"
	case "\rMain.Power=On\r":
		p.reply <- "\n\n\nMain.Power=On\n"
	case "\rMain.Power=Off\r":
		p.reply <- "\nMain.Power=Off\n"
	case "\rMain.Source=CD\r":
		p.reply <- "\nMain.Source=CD\n"
	case "\rMain.Source=TUNER\r":
		p.reply <- "\nMain.Source=TUNER\n"
	case "\rMain.Source=VIDEO\r":
		p.reply <- "\nMain.Source=VIDEO\n"
	case "\rMain.Source=DISC/MDC\r":
		p.reply <- "\nMain.Source=DISC/MDC\n"
	case "\rMain.Source=TAPE2\r":
		p.reply <- "\nMain.Source=TAPE2\n"
	case "\rMain.Source=AUX\r":
		p.reply <- "\nMain.Source=AUX\n"
	case "\rMain.Source?\r":
		p.reply <- "\nMain.Source=CD\n"
	case "\rMain.SpeakerA=On\r":
		p.reply <- "\nMain.SpeakerA=On\n"
	case "\rMain.SpeakerA=Off\r":
		p.reply <- "\nMain.SpeakerA=Off\n"
	case "\rMain.SpeakerB=On\r":
		p.reply <- "\nMain.SpeakerB=On\n"
	case "\rMain.SpeakerB=Off\r":
		p.reply <- "\nMain.SpeakerB=Off\n"
	case "\rMain.Tape1=On\r":
		p.reply <- "\nMain.Tape1=On\n"
	case "\rMain.Tape1=Off\r":
		p.reply <- "\nMain.Tape1=Off\n"
	case "\rMain.Volume+\r":
		p.reply <- "\nMain.Volume+\n"
	case "\rMain.Volume-\r":
		p.reply <- "\nMain.Volume-\n"
	default:
		return 0, fmt.Errorf("unknown command: %q", cmd)
	}
	return len(b), nil
}

// NewTestClient creates a client that communicates with a simulated amp.
func NewTestClient() *Client {
	reply := make(chan string, 1)
	port := &testPort{reply: reply}
	return &Client{port: port, EnableVolume: true}
}
