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
	actual, err := ParseCmd("Main.Power=On\r")
	if err != nil {
		t.Fatal(err)
	}
	expected := Cmd{Variable: "Power", Operator: "=", Value: "On"}
	if expected != actual {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
	actual, err = ParseCmd("Main.Model?\r")
	if err != nil {
		t.Fatal(err)
	}
	expected = Cmd{Variable: "Model", Operator: "?", Value: ""}
	if expected != actual {
		t.Errorf("Expected %q, got %q", expected, actual)
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
	assertTrue(Cmd{Variable: "Source", Operator: "?"})
	assertTrue(Cmd{Variable: "SpeakerA", Operator: "?"})
	assertTrue(Cmd{Variable: "SpeakerB", Operator: "?"})
	assertTrue(Cmd{Variable: "Tape1", Operator: "?"})
	assertTrue(Cmd{Variable: "Volume", Operator: "?"})
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
	assertTrue(Cmd{Variable: "Source", Operator: "=", Value: "Tuner"})
	assertTrue(Cmd{Variable: "Source", Operator: "=", Value: "Video"})
	assertTrue(Cmd{Variable: "Source", Operator: "=", Value: "Disc"})
	assertTrue(Cmd{Variable: "Source", Operator: "=", Value: "Ipod"})
	assertTrue(Cmd{Variable: "Source", Operator: "=", Value: "Tape2"})
	assertTrue(Cmd{Variable: "Source", Operator: "=", Value: "Aux"})
	assertFalse(Cmd{Variable: "Source", Operator: "=", Value: "foo"})
	assertFalse(Cmd{Variable: "Power", Operator: "=", Value: "foo"})

	assertTrue(Cmd{Variable: "Volume", Operator: "+"})
	assertTrue(Cmd{Variable: "Volume", Operator: "-"})
}
