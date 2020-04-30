package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/miekg/dns"
)

var (
	forwarder string
	localip   string
	bindip    string
	bindport  int
)

func main() {
	flag.StringVar(&basedomain, "basedomain", "portfwd", "DNS tld for asking for a dynamic port forward")
	flag.StringVar(&localip, "localip", "", "Local IP to return to DNS responses (local port forwarding address)")
	flag.StringVar(&bindip, "bindip", "0.0.0.0", "Bind IP to listen for DNS requests")
	flag.IntVar(&bindport, "bindport", 53, "Port to listen for DNS requests")
	flag.StringVar(&forwarder, "forwarder", "", "Upstream server and port (8.8.8.8:53) to forward all non-answerable requests to")
	flag.Parse()

	for i := 0; i < len(flag.Args()); i++ {
		switch flag.Arg(i) {
		case "help":
			fallthrough
		case "-h":
			flag.Usage()
		}
	}

	if localip == "" {
		fmt.Println("You need to at least assign a local IP address to respond with")
		fmt.Println()
		flag.Usage()
		os.Exit(1)
	}

	var dnshandlerinstance dnshandler

	fmt.Println("DNS Port Forwarder 0.1 launched")

	bindaddr := bindip + ":" + strconv.Itoa(bindport)
	err := dns.ListenAndServe(bindaddr, "udp", dnshandlerinstance)
	if err != nil {
		fmt.Println("Can't start DNS listener on", bindaddr, ":", err)
	}
}
