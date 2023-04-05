package cloudfinder


import(
	"github.com/projectdiscovery/cdncheck"
	"github.com/projectdiscovery/dnsx/libs/dnsx"
	"net"
	"net/url"
	"github.com/bobesa/go-domain-util/domainutil"


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

func Resolver(string) []net.IP {
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


func cloudfinder(target string, verbose bool, wg *sync.WaitGroup) {
	var host string

	client, err := cdncheck.NewWithCache()
	if err != nil {
			log.Fatal(err)
	}
	if isURL(target) {
		htok = ParseUrlTokens(target)
		host = htok.Subdomain + "." + htok.Domain
	} else {
		host = target
	}

	ip := net.ParseIP(host)
	ips := []net.IP{}
	if ip != nil {
		ips = append(ips, ip)
	} else {
		//ips = append(ips, Resolver(host)...)
		ips = append(ips, Resolver(host))
	}

	for _, ip := range ips {
		found, result, err := client.Check(ip); 
		if found && err == nil {
			log.Printf("hostname: %s provider: %s\n", host, result)
			return true
		}
	}
	return false

}



