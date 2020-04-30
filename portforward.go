package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
)

var (
	forwards = make(map[int]fwdinfo)
)

type fwdinfo struct {
	name     string
	port     int
	listener net.Listener
	// timeout  time.Time
}

func newportfwd(lport int, name string, dport int) error {
	fi, found := forwards[lport]
	if found {
		println("Already forwarding from my own port", lport, "to", fi.name, "port", fi.port, " - closing listener")
		err := fi.listener.Close()
		if err != nil {
			println("Error terminating existing listener:", err)
			return err
		}
	}

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(lport))
	if err != nil {
		println("Error opening listener from my own port", lport, "to", name, "port", dport)
		return err
	}
	fi = fwdinfo{
		name,
		dport,
		ln,
	}
	println("Now forwarding from my own port", lport, "to", fi.name, "port", fi.port)
	forwards[lport] = fi
	go handlenewconnections(fi)
	return nil
}

func handlenewconnections(fi fwdinfo) error {
	for {
		conn, err := fi.listener.Accept()
		if err != nil {
			return err
		}
		go handleconnection(fi, conn)
	}
}

func handleconnection(fi fwdinfo, conn net.Conn) error {
	fmt.Println("Connecting to", fi.name, "port", fi.port)
	remote, err := net.Dial("tcp", fi.name+":"+strconv.Itoa(fi.port))
	if err != nil {
		fmt.Println("Error connecting to remote:", err)
		return err
	}
	go io.Copy(remote, conn)
	go io.Copy(conn, remote)
	return nil //  Best effort
}

func terminatelistener(fi fwdinfo) error {
	return fi.listener.Close()
}
