package main

import (
	"fmt"
	"os"

	"github.com/miekg/dns"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprint(os.Stderr, "Invalid number of parameter.\n Currect usage: dns <name>\n")
		os.Exit(1)
	}
	name := os.Args[1]
	answer, err := resolveDNS(name)
	if err == nil {
		for _, record := range answer {
			fmt.Println(record)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Failed to resolve %s\n", name)
		os.Exit(1)
	}
}
func resolveDNS(name string) ([]dns.RR, error) {
	nameserver := "192.168.1.1"

	c := new(dns.Client)

	for {

		m := new(dns.Msg)
		m.SetQuestion(dns.Fqdn(name), dns.TypeA)

		// DNS request to IP in nameserver ( dns resolver)
		fmt.Printf("Quering %s about %s\n", nameserver, name)
		resp, _, err := c.Exchange(m, fmt.Sprintf("%s:53", nameserver))
		if err != nil {
			return nil, err
		}

		// Answer Section
		if len(resp.Answer) > 0 {

			if cname, ok := resp.Answer[0].(*dns.CNAME); ok {
				return resolveDNS(cname.Target)
			}
			return resp.Answer, nil
		}

		// if additional sections exits, look for in it to find the next-level namesever, if it doesnot error out.
		found := false

		for _, rr := range resp.Extra {

			record, ok := rr.(*dns.A)
			if ok {
				nameserver = record.A.String()
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("break in resolution")
		}

		if len(resp.Extra) == 0 && len(resp.Ns) != 0 {
			ns := resp.Ns[0].(*dns.NS)
			nsIP, err := resolveDNS(ns.Ns)
			if err != nil {
				return nil, fmt.Errorf("break in resolution")
			}
			nameserver = nsIP[0].(*dns.A).A.String()
		}
	}

}
