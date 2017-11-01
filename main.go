package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/mpolden/nadapi/api"
	"github.com/mpolden/nadapi/nad"
)

type opts struct {
	Device       string `short:"d" long:"device" description:"Path to serial device" value-name:"FILE" default:"/dev/ttyUSB0"`
	EnableVolume bool   `short:"x" long:"volume" description:"Allow volume adjustment. Use with caution!"`
	Test         bool   `short:"t" long:"test" description:"Test mode. Sends commands to a simulated device."`
}

type serverCmd struct {
	opts
	Listen    string `short:"l" long:"listen" description:"Listen address" value-name:"ADDR" default:":8080"`
	StaticDir string `short:"s" long:"static" description:"Path to directory containing static assets" value-name:"DIR"`
}

func newClient(device string, test bool) (*nad.Client, error) {
	if test {
		return nad.NewTestClient(), nil
	}
	return nad.New(device)
}

func (s *serverCmd) Execute(args []string) error {
	client, err := newClient(s.Device, s.Test)
	if err != nil {
		return err
	}
	client.EnableVolume = s.EnableVolume
	api := api.New(client)
	api.StaticDir = s.StaticDir
	if strings.HasPrefix(s.Listen, ":") {
		log.Printf("Serving at http://0.0.0.0%s", s.Listen)
	} else {
		log.Printf("Serving at http://%s", s.Listen)
	}
	if err := http.ListenAndServe(s.Listen, api.Handler()); err != nil {
		return err
	}
	return nil
}

type sendCmd struct {
	opts
	Args struct {
		Command string `positional-arg-name:"<command>" description:"Command to send"`
	} `positional-args:"yes" required:"yes"`
}

func (s *sendCmd) Execute(args []string) error {
	client, err := newClient(s.Device, s.Test)
	if err != nil {
		return err
	}
	defer client.Close()
	client.EnableVolume = s.EnableVolume
	cmd := s.Args.Command
	if !strings.HasPrefix(strings.ToLower(cmd), "main.") {
		cmd = "Main." + cmd
	}
	reply, err := client.SendString(cmd)
	if err != nil {
		return err
	}
	fmt.Println(reply)
	return nil
}

type listCmd struct{}

func (l *listCmd) Execute(args []string) error {
	for _, c := range nad.Cmds() {
		fmt.Println(c)
	}
	return nil
}

func main() {
	p := flags.NewParser(nil, flags.Default)
	var server serverCmd
	if _, err := p.AddCommand("server", "Start API server",
		"REST API for NAD amplifier.", &server); err != nil {
		log.Fatal(err)
	}
	var send sendCmd
	if _, err := p.AddCommand("send", "Send command",
		"Send command to NAD amplifier.", &send); err != nil {
		log.Fatal(err)
	}

	var list listCmd
	if _, err := p.AddCommand("list", "List commands",
		"List commands accepted by NAD amplifier.", &list); err != nil {
		log.Fatal(err)
	}

	if _, err := p.Parse(); err != nil {
		os.Exit(1)
	}
}
