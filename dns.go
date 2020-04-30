package main

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"

	"github.com/miekg/dns"
)

var (
	basedomain = "portfwd"
)

type dnshandler struct{}

func (d dnshandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	var err error
	var remotehost string
	var remoteport, localport int

	err = errors.New("Unsupported query type")

	reg := regexp.MustCompile(`^(.+)\.([\d]+)\.([\d]+)\.` + basedomain + `\.$`)
	matches := reg.FindStringSubmatch(r.Question[0].Name)
	if len(matches) == 4 {
		// Our TLD
		if r.Question[0].Qtype == dns.TypeA && r.Question[0].Qclass == dns.ClassINET {

			err = nil // Now we assume there are no problems

			remotehost = matches[1]
			remoteport, err = strconv.Atoi(matches[2])
			if err == nil {
				localport, err = strconv.Atoi(matches[3])
			}
			if err == nil {
				newportfwd(localport, remotehost, remoteport)

				aRec := &dns.A{
					Hdr: dns.RR_Header{
						Name:   r.Question[0].Name,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    0,
					},
					A: net.ParseIP(localip).To4(), // Return our local address
				}
				m.Answer = append(m.Answer, aRec)
				fmt.Println("Responding with", localip, "for mapping", remotehost, "port", remoteport, "from local port", localport)
			}
		} else {
			// Just ignore everything other than IPv4 type A lookups
			m.Opcode = dns.RcodeNotImplemented
		}
	} else if forwarder != "" {
		// Forward the request upstream and use that as answer
		var um *dns.Msg
		ur := r.Copy()                        // Copy to new request
		um, err = dns.Exchange(ur, forwarder) // As forwarder
		if err != nil {
			fmt.Println("Problem quering upstream:", err)
			m.Opcode = dns.RcodeServerFailure
		} else {
			um.CopyTo(m) // Copy the answer to our own response
		}
	}

	w.WriteMsg(m)
}
