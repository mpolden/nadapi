package nad

import (
	"fmt"
	"regexp"
	"strings"
)

const prefix = "Main"

var (
	commands = [...]string{
		"Main.Model?",
		"Main.Mute?",
		"Main.Mute=On",
		"Main.Mute=Off",
		"Main.Power?",
		"Main.Power=On",
		"Main.Power=Off",
		"Main.Source?",
		"Main.Source=CD",
		"Main.Source=TUNER",
		"Main.Source=VIDEO",
		"Main.Source=DISC/MDC",
		"Main.Source=TAPE2",
		"Main.Source=AUX",
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
	cmdPattern = regexp.MustCompile("(?iU)^" + prefix + "\\." +
		"(.+)" +
		"([=+?-])" +
		"(.*)$")
)

// Cmd represents a command sent to the amplifier.
type Cmd struct {
	Variable string
	Operator string
	Value    string
}

// Reply represents an reply received from the amplifier. A reply has the same fields as a command.
type Reply struct{ Cmd }

// String formats command as a string.
func (c *Cmd) String() string {
	return fmt.Sprint(prefix, ".", c.Variable, c.Operator, c.Value)
}

// Bytes returns a command as bytes.
func (c *Cmd) Bytes() []byte {
	return []byte(c.Delimited())
}

// Delimited formats command before sending it to amplifier.
func (c *Cmd) Delimited() string {
	return fmt.Sprint("\n", c.String(), "\n")
}

// Valid returns true if command is a command accepted by amplifier.
func (c *Cmd) Valid() bool {
	cmd := strings.ToLower(c.String())
	for _, c := range commands {
		if strings.ToLower(c) == cmd {
			return true
		}
	}
	return false
}

// ParseCmd parses s into a command.
func ParseCmd(s string) (Cmd, error) {
	s = strings.Trim(s, "\r\n")
	m := cmdPattern.FindAllStringSubmatch(s, -1)
	if len(m) == 0 || len(m[0]) < 4 {
		return Cmd{}, fmt.Errorf("could not parse command: %q", s)
	}
	return Cmd{
		Variable: m[0][1],
		Operator: m[0][2],
		Value:    m[0][3],
	}, nil
}

// ParseReply parses b into a reply.
func ParseReply(b []byte) (Reply, error) {
	cmd, err := ParseCmd(string(b))
	if err != nil {
		return Reply{}, err
	}
	return Reply{cmd}, nil
}

// Cmds returns all valid commands.
func Cmds() [25]string {
	return commands
}
