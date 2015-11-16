package system

import (
	"strconv"
	"strings"

	"github.com/aelsabbahy/GOnetstat"
	"github.com/codegangsta/cli"
)

type Port interface {
	Port() string
	Exists() (interface{}, error)
	Listening() (interface{}, error)
	IP() (interface{}, error)
}

type DefPort struct {
	port     string
	sysPorts map[string]GOnetstat.Process
}

func NewDefPort(port string, system *System, c *cli.Context) Port {
	p := normalizePort(port)
	return &DefPort{
		port:     p,
		sysPorts: system.Ports(c),
	}
}

func splitPort(fullport string) (network, port string) {
	split := strings.SplitN(fullport, ":", 2)
	if len(split) == 2 {
		return split[0], split[1]
	}
	return "tcp", fullport

}

func normalizePort(fullport string) string {
	net, addr := splitPort(fullport)
	return net + ":" + addr
}

func (p *DefPort) Port() string {
	return p.port
}

func (p *DefPort) Exists() (interface{}, error) { return p.Listening() }

func (p *DefPort) Listening() (interface{}, error) {
	if _, ok := p.sysPorts[p.port]; ok {
		return true, nil
	}
	return false, nil
}

func (p *DefPort) IP() (interface{}, error) {
	return p.sysPorts[p.port].Ip, nil
}

func GetPorts(lookupPids bool, c *cli.Context) map[string]GOnetstat.Process {
	ports := make(map[string]GOnetstat.Process)
	var net string
	var netstat []GOnetstat.Process
	if c.GlobalBool("ipv6") == false {
		netstat = GOnetstat.Tcp(lookupPids)
		net = "tcp"
		for _, entry := range netstat {
			if entry.State == "LISTEN" {
				port := strconv.FormatInt(entry.Port, 10)
				ports[net+":"+port] = entry
			}
		}
		netstat = GOnetstat.Udp(lookupPids)
		net = "udp"
		for _, entry := range netstat {
			port := strconv.FormatInt(entry.Port, 10)
			ports[net+":"+port] = entry
		}
	}
	if c.GlobalBool("ipv4") == false {
		netstat = GOnetstat.Tcp6(lookupPids)
		net = "tcp6"
		for _, entry := range netstat {
			if entry.State == "LISTEN" {
				port := strconv.FormatInt(entry.Port, 10)
				ports[net+":"+port] = entry
			}
		}
		netstat = GOnetstat.Udp6(lookupPids)
		net = "udp6"
		for _, entry := range netstat {
			port := strconv.FormatInt(entry.Port, 10)
			ports[net+":"+port] = entry
		}
	}
	return ports
}
