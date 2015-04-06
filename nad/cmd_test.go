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
	expected := "\rMain.Power=On\r"
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
	assertTrue := func(cmd Cmd) {
		if !cmd.Valid() {
			t.Errorf("Expected true for %q", cmd)
		}
	}
	assertFalse := func(cmd Cmd) {
		if cmd.Valid() {
			t.Errorf("Expected false for %q", cmd)
		}
	}
	assertTrue(Cmd{Variable: "Model", Operator: "?"})
	assertTrue(Cmd{Variable: "Mute", Operator: "?"})
	assertTrue(Cmd{Variable: "Power", Operator: "?"})
	assertTrue(Cmd{Variable: "power", Operator: "?"})
	assertTrue(Cmd{Variable: "Source", Operator: "?"})
	assertTrue(Cmd{Variable: "SpeakerA", Operator: "?"})
	assertTrue(Cmd{Variable: "SpeakerB", Operator: "?"})
	assertTrue(Cmd{Variable: "Tape1", Operator: "?"})
	assertFalse(Cmd{Variable: "Volume", Operator: "?"})
	assertFalse(Cmd{Variable: "foo", Operator: "?"})

	assertFalse(Cmd{Variable: "Model", Operator: "=", Value: "On"})
	assertTrue(Cmd{Variable: "Mute", Operator: "=", Value: "On"})
	assertTrue(Cmd{Variable: "Power", Operator: "=", Value: "On"})
	assertFalse(Cmd{Variable: "Source", Operator: "=", Value: "On"})
	assertTrue(Cmd{Variable: "SpeakerA", Operator: "=", Value: "On"})
	assertTrue(Cmd{Variable: "SpeakerB", Operator: "=", Value: "On"})
	assertTrue(Cmd{Variable: "Tape1", Operator: "=", Value: "On"})
	assertFalse(Cmd{Variable: "Volume", Operator: "=", Value: "On"})

	assertTrue(Cmd{Variable: "Source", Operator: "=", Value: "CD"})
	assertTrue(Cmd{Variable: "Source", Operator: "=", Value: "TUNER"})
	assertTrue(Cmd{Variable: "Source", Operator: "=", Value: "VIDEO"})
	assertTrue(Cmd{Variable: "Source", Operator: "=", Value: "DISC/MDC"})
	assertTrue(Cmd{Variable: "Source", Operator: "=", Value: "TAPE2"})
	assertTrue(Cmd{Variable: "Source", Operator: "=", Value: "AUX"})
	assertFalse(Cmd{Variable: "Source", Operator: "=", Value: "foo"})
	assertFalse(Cmd{Variable: "Power", Operator: "=", Value: "foo"})

	assertTrue(Cmd{Variable: "Volume", Operator: "+"})
	assertTrue(Cmd{Variable: "Volume", Operator: "-"})
}

func TestCmds(t *testing.T) {
	cmds := Cmds()
	cmds[3] = "foo"
	if Cmds()[3] == "foo" {
		t.Errorf("Expected %q, got %q", "Main.Mute=Off", Cmds()[3])
	}
}
