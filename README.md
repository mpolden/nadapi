# nadapi

[![Build Status](https://travis-ci.org/martinp/nadapi.png)](https://travis-ci.org/martinp/nadapi)

A tool that implements the RS-232 protocol for NAD amplifiers. The protocol
documentation can be found [here](http://nadelectronics.com/software).

The tool has two modes:

* server - Serves a REST API where commands can be sent to the amplifier over HTTP.
* send - Sends a command directly to the amplifier and prints the response.

It's developed primarily for the NAD C356BEE, but should work for other NAD
amplifiers that have a RS-232 port.

## Usage

```
$ nadapi -h
Usage:
  nadapi [OPTIONS] <list | send | server>

Help Options:
  -h, --help  Show this help message

Available commands:
  list    List commands
  send    Send command
  server  Start API server
```
