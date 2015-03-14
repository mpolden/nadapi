package nad

import (
	"fmt"
	"regexp"
)

const prefix = "Main"

var cmdExp = regexp.MustCompile("^" + prefix + "\\." +
	"(Model|Mute|Power|Source|Speaker[A-B]|Tape1|Volume)" +
	"([=+-?])" +
	"(On|Off|CD|Tuner|Video|Disc|Ipod|Tape2|Aux)?\\r$")

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

func ParseCmd(s string) (Cmd, error) {
	m := cmdExp.FindAllStringSubmatch(s, -1)
	if len(m) == 0 || len(m[0]) < 4 {
		return Cmd{}, fmt.Errorf("expected 4 submatches")
	}
	return Cmd{
		Variable: m[0][1],
		Operator: m[0][2],
		Value:    m[0][3],
	}, nil
}
