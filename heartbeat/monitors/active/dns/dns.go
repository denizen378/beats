package dns

import (
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/heartbeat/monitors"
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
	var index int
	
	for _, nameserver := range config.NameServers {
	    for _, question := range config.Questions {
	    	jobs[index], err = newDNSMonitorHostJob(nameserver, question, &config)
	    	if err != nil {
	               return nil, err
	    	}
		index++
            }
        }

	return jobs, nil
}

