package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/empijei/cli/lg"
	"github.com/gopherjs/websocket"
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 2 {
		lg.Failure("Please specify a host to connect to")
	}

	host := flag.Arg(0)
	port, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		lg.Failure(err)
	}
	if port == 0 {
		port = 80
	}
	ws, err := websocket.Dial(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		lg.Failure(err.Error())
	}
	defer func() {
		_ = ws.Close()
	}()
	go func() {
		_, err := io.Copy(ws, os.Stdin)
		if err != nil {
			lg.Failure(err)
		}
	}()
	_, err = io.Copy(os.Stdout, ws)
	if err != nil {
		lg.Failure(err)
	}
}
