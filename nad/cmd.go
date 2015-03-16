package nad

import (
	"fmt"
	"regexp"
)

const prefix = "Main"

var cmdExp = regexp.MustCompile("^" + prefix + "\\." +
	"(Model|Mute|Power|Source|Speaker[A-B]|Tape1|Volume)" +
	"([=+-?])" +
	"([A-Za-z0-9]+|[+-]\\d+)?\\r?$")

var commands = [...]string{
	"Main.Model?",
	"Main.Mute?",
	"Main.Mute=On",
	"Main.Mute=Off",
	"Main.Power?",
	"Main.Power=On",
	"Main.Power=Off",
	"Main.Source?",
	"Main.Source=CD",
	"Main.Source=Tuner",
	"Main.Source=Video",
	"Main.Source=Disc",
	"Main.Source=Ipod",
	"Main.Source=Tape2",
	"Main.Source=Aux",
	"Main.SpeakerA?",
	"Main.SpeakerA=On",
	"Main.SpeakerA=Off",
	"Main.SpeakerB?",
	"Main.SpeakerB=On",
	"Main.SpeakerB=Off",
	"Main.Tape1?",
	"Main.Tape1=On",
	"Main.Tape1=Off",
	"Main.Volume+",
	"Main.Volume-",
}

type Cmd struct {
	Variable string
	Operator string
	Value    string
}

func (c *Cmd) String() string {
	return fmt.Sprint(prefix, ".", c.Variable, c.Operator, c.Value)
}

func (c *Cmd) Delimited() string {
	return fmt.Sprint("\r", c.String(), "\r")
}

func (c *Cmd) Valid() bool {
	cmd := c.String()
	for _, c := range commands {
		if c == cmd {
			return true
		}
	}
	return false
}

func ParseCmd(s string) (Cmd, error) {
	m := cmdExp.FindAllStringSubmatch(s, -1)
	if len(m) == 0 || len(m[0]) < 4 {
		return Cmd{}, fmt.Errorf("could not parse command: %s", s)
	}
	return Cmd{
		Variable: m[0][1],
		Operator: m[0][2],
		Value:    m[0][3],
	}, nil
}

func Cmds() [26]string {
	return commands
}
