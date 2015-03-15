package nad

import (
	"fmt"
	"regexp"
)

const prefix = "Main"

var cmdExp = regexp.MustCompile("^" + prefix + "\\." +
	"(Model|Mute|Power|Source|Speaker[A-B]|Tape1|Volume)" +
	"([=+-?])" +
	"([A-Za-z0-9]+|[+-]\\d+)?\\r$")

type Cmd struct {
	Variable string
	Operator string
	Value    string
}

func (c *Cmd) String() string {
	return fmt.Sprintf("%s.%s%s%s", prefix, c.Variable, c.Operator, c.Value)
}

func (c *Cmd) Delimited() string {
	return fmt.Sprintf("\r%s\r", c.String())
}

func (c *Cmd) Valid() bool {
	// All variables support querying using ?
	switch c.Variable {
	case "Model", "Mute", "Power", "Source", "SpeakerA", "SpeakerB",
		"Tape1", "Volume":
		if c.Operator == "?" && c.Value == "" {
			return true
		}
	}

	// Variables which support On/Off toggling
	switch c.Variable {
	case "Mute", "Power", "SpeakerA", "SpeakerB", "Tape1":
		if c.Operator == "=" && (c.Value == "On" || c.Value == "Off") {
			return true
		}
	}

	// Valid sources
	if c.Variable == "Source" && c.Operator == "=" {
		switch c.Value {
		case "CD", "Tuner", "Video", "Disc", "Ipod", "Tape2", "Aux":
			return true
		default:
			return false
		}
	}

	// Volume adjustment
	return c.Variable == "Volume" &&
		(c.Operator == "+" || c.Operator == "-") &&
		c.Value == ""
}

func ParseCmd(s string) (Cmd, error) {
	m := cmdExp.FindAllStringSubmatch(s, -1)
	if len(m) == 0 || len(m[0]) < 4 {
		return Cmd{}, fmt.Errorf("failed to parse command")
	}
	return Cmd{
		Variable: m[0][1],
		Operator: m[0][2],
		Value:    m[0][3],
	}, nil
}
