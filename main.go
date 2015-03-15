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

	var cli struct {
		Command string `short:"c" long:"command" description:"Command to send" required:"true"`
	}

	p := flags.NewParser(&opts, flags.Default)
	serverCmd, err := p.AddCommand("server", "Serve API",
		"REST API for NAD amplifier.", &server)
	if err != nil {
		log.Fatal(err)
	}
	cliCmd, err := p.AddCommand("cli", "Send command",
		"Send command to a NAD amplifier.", &cli)
	if err != nil {
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

	api := api.New(nad)
	if p.Active == serverCmd {
		log.Printf("Listening on %s", server.Listen)
		if err := api.ListenAndServe(server.Listen); err != nil {
			log.Fatal(err)
		}
	} else if p.Active == cliCmd {
		reply, err := nad.SendString(cli.Command)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(reply)
	}
}
