package dns

import (
	"fmt"
        "errors"
	"github.com/miekg/dns"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/heartbeat/monitors"
	"github.com/elastic/beats/heartbeat/reason"
	"net"
)

func newDNSMonitorHostJob(
        nameserver string,
	question string,
	config *Config,
) (monitors.Job, error) {
	typ := config.Name
	jobName := fmt.Sprintf("%v@%v@%v", typ, question, nameserver)

	fields := common.MapStr{
		"nameserver":  nameserver,
		"question": question,
	}

	return monitors.MakeSimpleJob(jobName, typ, func() (common.MapStr, error) {
		event, err := execQuery(nameserver, question)
		if event == nil {
			event = common.MapStr{}
		}
		event.Update(fields)
		return event, err
	}), nil
}

func execQuery(nameserver string, question string)(common.MapStr, reason.Reason){

     host, port, port_err := net.SplitHostPort(nameserver)

     if port_err != nil {
     	nameserver += ":53"
	port = "53" 
     }

     m1 := new(dns.Msg)
     m1.Id = dns.Id()
     m1.RecursionDesired = true
     m1.Question = make([]dns.Question, 1)
     m1.Question[0] = dns.Question{question, dns.TypeANY, dns.ClassINET}
     c := new(dns.Client)
     in, rtt, err := c.Exchange(m1, nameserver)

     event := common.MapStr{
     	      "response": common.MapStr{
              		"in":   in,
	      },
	      "nameserver": nameserver,
	      "question": question,
	      "rtt":  rtt,
	      "dst_host": host,
	      "port": port,
     }

     if len(in.Answer) == 0 {
     	resp_err := errors.New("Zero Answers")
        return event, reason.IOFailed(resp_err)
     }

     if err != nil {
        return event, reason.IOFailed(err)
     }

     return event, nil
}
