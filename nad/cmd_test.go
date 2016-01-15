package nad

import (
	"testing"
)

func TestCmdString(t *testing.T) {
	cmd := Cmd{Variable: "Power", Operator: "=", Value: "On"}
	expected := "Main.Power=On"
	actual := cmd.String()
	if actual != expected {
		t.Fatalf("Expected %q, got %q", expected, actual)
	}
}

func TestCmdDelimited(t *testing.T) {
	cmd := Cmd{Variable: "Power", Operator: "=", Value: "On"}
	expected := "\nMain.Power=On\n"
	actual := cmd.Delimited()
	if actual != expected {
		t.Fatalf("Expected %q, got %q", expected, actual)
	}
}

func TestParseCmd(t *testing.T) {
	var tests = []struct {
		in  string
		out Cmd
	}{
		{"Main.Power=On\n", Cmd{Variable: "Power", Operator: "=", Value: "On"}},
		{"Main.Model?", Cmd{Variable: "Model", Operator: "?", Value: ""}},
		{"Main.Model=C356", Cmd{Variable: "Model", Operator: "=", Value: "C356"}},
		{"main.volume+", Cmd{Variable: "volume", Operator: "+", Value: ""}},
		{"Main.Source=DISC/MDC", Cmd{Variable: "Source", Operator: "=", Value: "DISC/MDC"}},
		{"Main.Tape1?", Cmd{Variable: "Tape1", Operator: "?", Value: ""}},
	}
	for _, tt := range tests {
		cmd, err := ParseCmd(tt.in)
		if err != nil {
			t.Fatal(err)
		}
		if cmd != tt.out {
			t.Errorf("Expected %q, got %q", tt.out, cmd)
		}
	}
}

func TestCmdValid(t *testing.T) {
	var tests = []struct {
		in  Cmd
		out bool
	}{
		{Cmd{Variable: "Model", Operator: "?"}, true},
		{Cmd{Variable: "Mute", Operator: "?"}, true},
		{Cmd{Variable: "Power", Operator: "?"}, true},
		{Cmd{Variable: "power", Operator: "?"}, true},
		{Cmd{Variable: "Source", Operator: "?"}, true},
		{Cmd{Variable: "SpeakerA", Operator: "?"}, true},
		{Cmd{Variable: "SpeakerB", Operator: "?"}, true},
		{Cmd{Variable: "Tape1", Operator: "?"}, true},
		{Cmd{Variable: "Mute", Operator: "=", Value: "On"}, true},
		{Cmd{Variable: "Power", Operator: "=", Value: "On"}, true},
		{Cmd{Variable: "SpeakerA", Operator: "=", Value: "On"}, true},
		{Cmd{Variable: "SpeakerB", Operator: "=", Value: "On"}, true},
		{Cmd{Variable: "Tape1", Operator: "=", Value: "On"}, true},
		{Cmd{Variable: "Source", Operator: "=", Value: "CD"}, true},
		{Cmd{Variable: "Source", Operator: "=", Value: "TUNER"}, true},
		{Cmd{Variable: "Source", Operator: "=", Value: "VIDEO"}, true},
		{Cmd{Variable: "Source", Operator: "=", Value: "DISC/MDC"}, true},
		{Cmd{Variable: "Source", Operator: "=", Value: "TAPE2"}, true},
		{Cmd{Variable: "Source", Operator: "=", Value: "AUX"}, true},
		{Cmd{Variable: "Volume", Operator: "+"}, true},
		{Cmd{Variable: "Volume", Operator: "-"}, true},
		{Cmd{Variable: "Volume", Operator: "?"}, false},
		{Cmd{Variable: "foo", Operator: "?"}, false},
		{Cmd{Variable: "Model", Operator: "=", Value: "On"}, false},
		{Cmd{Variable: "Source", Operator: "=", Value: "On"}, false},
		{Cmd{Variable: "Volume", Operator: "=", Value: "On"}, false},
		{Cmd{Variable: "Source", Operator: "=", Value: "foo"}, false},
		{Cmd{Variable: "Power", Operator: "=", Value: "foo"}, false},
	}
	for _, tt := range tests {
		if valid := tt.in.Valid(); valid != tt.out {
			t.Errorf("Expected %t, got %t", tt.out, valid)
		}
	}
}

func TestCmds(t *testing.T) {
	cmds := Cmds()
	cmds[3] = "foo"
	if Cmds()[3] == "foo" {
		t.Errorf("Expected %q, got %q", "Main.Mute=Off", Cmds()[3])
	}
}
