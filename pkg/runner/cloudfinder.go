package cloudfinder


import(
	"github.com/projectdiscovery/cdncheck"
	"github.com/projectdiscovery/dnsx/libs/dnsx"
	"github.com/bobesa/go-domain-util/domainutil"
	"log"
	"net"
	"sync"
	"strings"
)

type DomainTokens struct {
    Protocol string
	Subdomain string
	Domain string
	Tld string
}

func ParseUrlTokens(value string) (*DomainTokens){
	var d DomainTokens
    d.Protocol = domainutil.Protocol(value)
	d.Subdomain = domainutil.Subdomain(value)
	d.Domain = domainutil.Domain(value)
	d.Tld = domainutil.DomainSuffix(value)
	return &d
}

func isURL(candidate string) bool {
	return strings.Contains(candidate, "://")
}

func Resolver(name string) []net.IP {
	resolver, err := dnsx.New(dnsx.DefaultOptions)
	if err != nil {
			log.Fatal(err)
	}

	validIPs := []net.IP{}
	ips, err := resolver.Lookup(name)
	if err != nil {
		return validIPs
	}
	
	for _, ip := range ips {
		parsedIP := net.ParseIP(ip)
		if parsedIP.To4() == nil {
			continue
		}
		validIPs = append(validIPs, parsedIP)
	}
	return validIPs
}


func start(target string, verbose bool, wg *sync.WaitGroup) {
	var host string

	client, err := cdncheck.NewWithCache()
	if err != nil {
			log.Fatal(err)
	}
	if isURL(target) {
		htok := ParseUrlTokens(target)
		host = htok.Subdomain + "." + htok.Domain
	} else {
		host = target
	}

	ip := net.ParseIP(host)
	ips := []net.IP{}
	if ip != nil {
		ips = append(ips, ip)
	} else {
		//... fixes "Resolver(host) (value of type []"net".IP) as "net".IP value in argument to append"
		ips = append(ips, Resolver(host)...)
	}

	for _, ip := range ips {
		found, result, err := client.Check(ip); 
		if found && err == nil {
			log.Printf("hostname: %s provider: %s\n", host, result)
		}
	}
}



