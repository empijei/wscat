package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/empijei/cli/lg"
	"golang.org/x/net/websocket"
)

var lst = flag.Bool("l", false, "start in listen mode, requires a port")
var source = flag.String("s", "", "local source address (ip or hostname)")
var path = flag.String("p", "/", "the http path to listen on")
var domain = flag.String("d", "", "optional, the domain to filter for. (will check Host header)")
var hexdump = flag.Bool("x", false, "hexdump incoming and outgoing traffic")

//TODO specify proxy

func main() {
	flag.Parse()

	// Listen mode
	if *lst {
		if len(flag.Args()) < 1 {
			lg.Failure("No port specified")
		}
		port, err := strconv.Atoi(flag.Arg(0))
		if err != nil {
			lg.Failure(err)
		}
		err = listen(*source, *domain, port, *path)
		if err != nil {
			lg.Failure(err)
		}
		return
	}

	// Connect mode
	if len(flag.Args()) < 2 {
		lg.Failure("Please specify a host and a port to connect to")
	}
	host := flag.Arg(0)
	port, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		lg.Failure(err)
	}
	err = connect(host, port)
	if err != nil {
		lg.Failure(err)
	}
}

func listen(source string, domain string, port int, path string) error {
	http.Handle(domain+"/"+path, websocket.Handler(func(ws *websocket.Conn) {
		err := connectIO(ws, os.Stdout, os.Stdin)
		if err != nil {
			lg.Error(err)
		}
	}))
	err := http.ListenAndServe(source+":"+strconv.Itoa(port), nil)
	return err
}

func connect(host string, port int) error {
	if port == 0 {
		port = 80
	}
	ws, err := websocket.Dial(fmt.Sprintf("%s:%d", host, port), "", "http://localhost/")
	if err != nil {
		return err
	}
	defer func() {
		_ = ws.Close()
	}()
	err = connectIO(ws, os.Stdout, os.Stdin)
	return err
}

func connectIO(ws net.Conn, out io.Writer, in io.Reader) (err error) {
	var errin, errout error
	defer func() {
		if errin != nil {
			err = errin
		}
		if errout != nil {
			err = errout
		}
		return
	}()
	go func() {
		_, errin = io.Copy(ws, in)
	}()
	_, errout = io.Copy(out, ws)
	return err
}
