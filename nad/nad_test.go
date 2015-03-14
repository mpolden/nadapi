package nad

import (
	"fmt"
	"testing"
)

type Port struct {
	reply chan string
}

func (p *Port) Close() (err error) {
	return
}

func (p *Port) Read(b []byte) (n int, err error) {
	r := []byte(<-p.reply)
	copy(b, r)
	return len(r), nil
}

func (p *Port) Write(b []byte) (n int, err error) {
	cmd := string(b)
	switch cmd {
	case "\rfoo?\r":
		p.reply <- "bar\r"
	case "\rMain.Model?\r":
		p.reply <- "C356\r"
	case "\rMain.Mute=On\r":
		p.reply <- "Main.Mute=On\r"
	case "\rMain.Mute=Off\r":
		p.reply <- "Main.Mute=Off\r"
	case "\rMain.Power=On\r":
		p.reply <- "Main.Power=On\r"
	case "\rMain.Power=Off\r":
		p.reply <- "Main.Power=Off\r"
	case "\rMain.Source=CD\r":
		p.reply <- "Main.Source=CD\r"
	case "\rMain.SpeakerA=On\r":
		p.reply <- "Main.SpeakerA=On\r"
	case "\rMain.SpeakerA=Off\r":
		p.reply <- "Main.SpeakerA=Off\r"
	case "\rMain.SpeakerB=On\r":
		p.reply <- "Main.SpeakerB=On\r"
	case "\rMain.SpeakerB=Off\r":
		p.reply <- "Main.SpeakerB=Off\r"
	case "\rMain.Tape1=On\r":
		p.reply <- "Main.Tape1=On\r"
	case "\rMain.Tape1=Off\r":
		p.reply <- "Main.Tape1=Off\r"
	case "\rMain.Volume+\r":
		p.reply <- "Main.Volume=+1\r"
	case "\rMain.Volume-\r":
		p.reply <- "Main.Volume=-1\r"
	default:
		return 0, fmt.Errorf("unknown command: %q", cmd)
	}
	return len(b), nil
}

func newNAD() NAD {
	reply := make(chan string, 1)
	port := &Port{reply: reply}
	return NAD{port: port}
}

func TestSendString(t *testing.T) {
	nad := newNAD()
	b, err := nad.SendString("\rfoo?\r")
	if err != nil {
		t.Fatal(err)
	}
	actual := string(b)
	if expected := "bar"; actual != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestModel(t *testing.T) {
	nad := newNAD()
	actual, err := nad.Model()
	if err != nil {
		t.Fatal(err)
	}
	if expected := "C356"; actual != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestMuteEnable(t *testing.T) {
	nad := newNAD()
	actual, err := nad.Mute(true)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Mute=On"; actual != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestMuteDisable(t *testing.T) {
	nad := newNAD()
	actual, err := nad.Mute(false)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Mute=Off"; actual != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestPowerEnable(t *testing.T) {
	nad := newNAD()
	actual, err := nad.Power(true)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Power=On"; actual != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestPowerDisable(t *testing.T) {
	nad := newNAD()
	actual, err := nad.Power(false)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Power=Off"; actual != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestSource(t *testing.T) {
	nad := newNAD()
	actual, err := nad.Source(CD)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Source=CD"; actual != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestSpeakerAEnable(t *testing.T) {
	nad := newNAD()
	actual, err := nad.SpeakerA(true)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.SpeakerA=On"; actual != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestSpeakerADisable(t *testing.T) {
	nad := newNAD()
	actual, err := nad.SpeakerA(false)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.SpeakerA=Off"; actual != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestSpeakerBEnable(t *testing.T) {
	nad := newNAD()
	actual, err := nad.SpeakerB(true)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.SpeakerB=On"; actual != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestSpeakerBDisable(t *testing.T) {
	nad := newNAD()
	actual, err := nad.SpeakerB(false)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.SpeakerB=Off"; actual != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestTape1Enable(t *testing.T) {
	nad := newNAD()
	actual, err := nad.Tape1(true)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Tape1=On"; actual != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestTape1Disable(t *testing.T) {
	nad := newNAD()
	actual, err := nad.Tape1(false)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Tape1=Off"; actual != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestVolumeUp(t *testing.T) {
	nad := newNAD()
	actual, err := nad.VolumeUp()
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Volume=+1"; actual != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestVolumeDown(t *testing.T) {
	nad := newNAD()
	actual, err := nad.VolumeDown()
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Volume=-1"; actual != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}
