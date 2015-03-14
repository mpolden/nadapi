package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/martinp/nadapi/api"
	//"github.com/martinp/nadapi/nad"
	"log"
	"os"
)

func main() {
	var opts struct {
		Listen string `short:"l" long:"listen" description:"Listen address" value-name:"ADDR" default:":8080"`
		Device string `short:"d" long:"device" description:"Path to serial device" value-name:"FILE" required:"true"`
	}
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		os.Exit(1)
	}

	// nad, err := nad.New(opts.Device)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	api := api.API{}

	log.Printf("Listening on %s", opts.Listen)
	if err := api.ListenAndServe(opts.Listen); err != nil {
		log.Fatal(err)
	}
}
