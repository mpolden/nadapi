package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/martinp/nadapi/api"
	"github.com/martinp/nadapi/nad"
	"log"
	"os"
)

type opts struct {
	Device       string `short:"d" long:"device" description:"Path to serial device" value-name:"FILE" required:"true"`
	EnableVolume bool   `short:"x" long:"volume" description:"Allow volume adjustment. Use with caution!"`
}

type server struct {
	opts
	Listen string `short:"l" long:"listen" description:"Listen address" value-name:"ADDR" default:":8080"`
}

func (s *server) Execute(args []string) error {
	client, err := nad.New(s.Device)
	if err != nil {
		return err
	}
	client.EnableVolume = s.EnableVolume
	api := api.New(client)
	log.Printf("Listening on %s", s.Listen)
	if err := api.ListenAndServe(s.Listen); err != nil {
		return err
	}
	return nil
}

type send struct {
	opts
	Args struct {
		Command string `positional-arg-name:"<command>" description:"Command to send"`
	} `positional-args:"yes" required:"yes"`
}

func (s *send) Execute(args []string) error {
	client, err := nad.New(s.Device)
	if err != nil {
		return err
	}
	client.EnableVolume = s.EnableVolume
	reply, err := client.SendString(s.Args.Command)
	if err != nil {
		return err
	}
	fmt.Println(reply)
	return nil
}

type list struct{}

func (l *list) Execute(args []string) error {
	for _, cmd := range nad.Commands() {
		fmt.Println(cmd)
	}
	return nil
}

func main() {
	p := flags.NewParser(nil, flags.Default)
	var server server
	if _, err := p.AddCommand("server", "Start API server",
		"REST API for NAD amplifier.", &server); err != nil {
		log.Fatal(err)
	}
	var send send
	if _, err := p.AddCommand("send", "Send command",
		"Send command to NAD amplifier.", &send); err != nil {
		log.Fatal(err)
	}

	var list list
	if _, err := p.AddCommand("list", "List commands",
		"List commands accepted by NAD amplifier.", &list); err != nil {
		log.Fatal(err)
	}

	if _, err := p.Parse(); err != nil {
		os.Exit(1)
	}
}
