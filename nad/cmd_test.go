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
