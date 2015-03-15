package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/martinp/nadapi/api"
	"github.com/martinp/nadapi/nad"
	"log"
	"os"
)

func main() {
	var opts struct {
		Device       string `short:"d" long:"device" description:"Path to serial device" value-name:"FILE" required:"true"`
		EnableVolume bool   `short:"x" long:"volume" description:"Allow volume adjustment. Use with caution!"`
	}

	var server struct {
		Listen string `short:"l" long:"listen" description:"Listen address" value-name:"ADDR" default:":8080"`
	}

	var send = struct {
		Args struct {
			Command string `positional-arg-name:"<command>" description:"Command to send"`
		} `positional-args:"yes" required:"yes"`
	}{}

	p := flags.NewParser(&opts, flags.Default)
	if _, err := p.AddCommand("server", "Start API server",
		"REST API for NAD amplifier.", &server); err != nil {
		log.Fatal(err)
	}
	if _, err := p.AddCommand("send", "Send command",
		"Send command to NAD amplifier.", &send); err != nil {
		log.Fatal(err)
	}

	if _, err := p.Parse(); err != nil {
		os.Exit(1)
	}

	nad, err := nad.New(opts.Device)
	if err != nil {
		log.Fatal(err)
	}
	nad.EnableVolume = opts.EnableVolume

	switch p.Active.Name {
	case "server":
		api := api.New(nad)
		log.Printf("Listening on %s", server.Listen)
		if err := api.ListenAndServe(server.Listen); err != nil {
			log.Fatal(err)
		}
	case "send":
		reply, err := nad.SendString(send.Args.Command)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(reply)
	default:
		log.Fatal("unknown subcommand")
	}
}
