package pkg

import (
	"github.com/esonhugh/k8spider/define"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
)

var _ StaticTaskRunner[*dns.Envelope] = (*AxfrTask)(nil)

type AxfrTask struct {
	Target    string
	DnsServer string
	Records   []define.Record
	lock      sync.Mutex
}

func (a *AxfrTask) Generator(tasks chan *dns.Envelope) {
	t := new(dns.Transfer)
	m := new(dns.Msg)
	m.SetAxfr(a.Target)
	ch, err := t.In(m, a.DnsServer)
	if err != nil {
		log.Fatalf("Transfer failed: %v", err)
	}
	for rr := range ch {
		tasks <- rr
	}
}

func (a *AxfrTask) Solver(task *dns.Envelope) {
	for _, r := range task.RR {
		a.AddRecord(define.Record{
			SvcDomain: r.Header().Name,
			Extra:     strings.Join(strings.Split(r.String(), "\t"), " "),
		})
	}
	log.Debugf("Record: %v", task.RR)
}

func (a *AxfrTask) AddRecord(r define.Record) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.Records = append(a.Records, r)
}

// default target should be zone
func DumpAXFR(target string, dnsServer string) []define.Record {
	tasker := &AxfrTask{Target: target, DnsServer: dnsServer}
	RunStatic[*dns.Envelope](tasker, Threads)
	return tasker.Records
}
