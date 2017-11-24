package dns

import (
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/heartbeat/monitors"
        "github.com/miekg/dns"
	"net"
	"strings"
	"fmt"
)

func init() {
	monitors.RegisterActive("dns", create)
}

var debugf = logp.MakeDebug("dns")

func create(
	info monitors.Info,
	cfg *common.Config,
) ([]monitors.Job, error) {
	config := defaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, err
	}

	var err error

	jobs := make([]monitors.Job, len(config.Questions) * len(config.NameServers))

	var (
	    index int
	    qtype uint16
	    isv6  bool
        )

	for _, nameserver := range config.NameServers {

	    host, port, port_err := net.SplitHostPort(nameserver)

            if port_err != nil {
               host = nameserver
               nameserver += ":53"
               port = "53"
            }
	    
	    if strings.Contains(host, ":") {
	       isv6 = true
	       fmt.Printf("host is ipv6[%v] nameserver[%v] isv6[%v]\n", host, nameserver, isv6)
	    }else{
	       isv6 = false
	       fmt.Printf("host is NOT ipv6[%v] nameserver[%v] isv6[%v]\n", host, nameserver, isv6)
	    }

	    if isv6 {
 	       fmt.Printf("isv6\n")
	    }

	    fmt.Printf("nameserver[%v] host[%v] port[%v]\n", nameserver, host, port)
	    for _, question := range config.Questions {

	    	query, qtypestr, qtype_err := net.SplitHostPort(question)

            	if qtype_err != nil {
               	    query = question
               	    qtype = dns.TypeA

            	}else{
		    if k, ok := dns.StringToType[strings.ToUpper(qtypestr)]; ok {
		       qtype = k
		    }else{
		       qtype = dns.TypeA
		    }
		}
	        fmt.Printf("    query[%v] qtype[%v]\n", query, qtype)
	    	jobs[index], err = newDNSMonitorHostJob(nameserver, host, port, isv6, query, qtype, &config)

	    	if err != nil {
	               return nil, err
	    	}
		index++
            }
        }

	return jobs, nil
}

