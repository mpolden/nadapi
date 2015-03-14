package nad

import (
	"testing"
)

type Port struct {
	cmds chan []byte
}

func (p *Port) Close() (err error) {
	return
}

func (p *Port) Read(b []byte) (n int, err error) {
	cmd := <-p.cmds
	copy(b, cmd)
	return len(cmd), nil
}

func (p *Port) Write(b []byte) (n int, err error) {
	cmd := string(b)
	switch cmd {
	case "\rfoo?\r":
		p.cmds <- []byte("bar\r")
	case "\rMain.Model?\r":
		p.cmds <- []byte("C356\r")
	case "\rMain.Mute=On\r":
		p.cmds <- []byte("Main.Mute=On\r")
	case "\rMain.Mute=Off\r":
		p.cmds <- []byte("Main.Mute=Off\r")
	default:
		panic("unknown command: %s")
	}
	return len(b), nil
}

func newNAD() NAD {
	cmds := make(chan []byte, 1)
	port := &Port{cmds: cmds}
	return NAD{port: port}
}

func TestSend(t *testing.T) {
	nad := newNAD()
	b, err := nad.Send("foo?")
	if err != nil {
		t.Fatal(err)
	}
	actual := string(b)
	if expected := "bar"; actual != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, actual)
	}
}

func TestModel(t *testing.T) {
	nad := newNAD()
	actual, err := nad.Model()
	if err != nil {
		t.Fatal(err)
	}
	if expected := "C356"; actual != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, actual)
	}
}

func TestMuteEnable(t *testing.T) {
	nad := newNAD()
	actual, err := nad.Mute(true)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Mute=On"; actual != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, actual)
	}
}

func TestMuteDisable(t *testing.T) {
	nad := newNAD()
	actual, err := nad.Mute(false)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "Main.Mute=Off"; actual != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, actual)
	}
}
