dnsportfwd

So you want to proxy stuff through another computer, but need an elegant way of doing this dynamically? dnsportfwd does the heavy lifing!

The idea is that you can open port forwards based on DNS queries, which makes it easy for you to use standard applications, but still get your packets routed through another device.

The scenario is:
- computer A running VPN to another site or similar, granting this computer acesss to services there
- another computer B, that needs access through the VPN

Computer A runs dnsportfwd, which listens for DNS requests on the portfwd TLD, and dynamically proxies TCP connections thorugh itself based on these requests
Computer B uses Computer A for DNS lookups

Normal internet access from computer B works as normal, but DNS requests are routed through computer A (split DNS also possible, if your router supports it)
When computer A is asked for something on the magic portfwd TLD, is maps this dynamically and responds with its own IP address
Computer B then thinks computer A is the host, and connects to this IP on the given port
Computer A then proxies this connection to the remote system

Accessing things from Computer B via Computer A is done by using a DNS query of this format:

hostname.tld.rport.lport.portfwd

Usage:

dnsportfwd -localip 1.2.3.4 -forwarder 8.8.8.8:53

So to get Google's unencrypted http frontpage (which only redirects to HTTPS) you could do:

$ curl http://www.google.com.80.8080.portfwd:8080

I've probably reinvented the wheel here, but at least it was fun to do :)
