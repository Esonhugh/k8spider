package all

import (
	"net"
	"os"
	"strings"
	"sync"

	command "github.com/esonhugh/k8spider/cmd"
	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg"
	"github.com/esonhugh/k8spider/pkg/mutli"
	"github.com/esonhugh/k8spider/pkg/post"
	"github.com/esonhugh/k8spider/pkg/printer"
	"github.com/esonhugh/k8spider/pkg/scanner"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Opts struct {
	OnlyService bool
}

func init() {
	command.RootCmd.AddCommand(AllCmd)
	AllCmd.PersistentFlags().BoolVarP(&Opts.OnlyService, "only-service", "O", false, "only dump service cidr")
}

var AllCmd = &cobra.Command{
	Use: "all",
	Aliases: []string{
		"a",
	},
	Short: "all is a tool to discover k8s services and available ip in subnet",
	Run: func(cmd *cobra.Command, args []string) {
		if command.Opts.Cidr == "" {
			log.Warn("cidr is required")
			return
		}
		// Wildcard
		records := scanner.DumpWildCard(command.Opts.Zone)
		if records != nil {
			printer.PrintResult(records, command.Opts.OutputFile)
		}
		// AXFR Dumping
		records, err := scanner.DumpAXFR(dns.Fqdn(command.Opts.Zone), "ns.dns."+command.Opts.Zone+":53")
		if err == nil {
			printer.PrintResult(records, command.Opts.OutputFile)
		} else {
			log.Errorf("Transfer failed: %v", err)
		}

		// Service Discovery
		ipNets, err := pkg.ParseStringToIPNet(command.Opts.Cidr)
		if err != nil {
			log.Warnf("ParseStringToIPNet failed: %v", err)
			return
		}
		podNets, err := pkg.ParseStringToIPNet(command.Opts.PodCidr)
		if err != nil {
			log.Warnf("ParseStringToIPNet failed: %v", err)
			return
		}

		var finalRecord define.Records
		if Opts.OnlyService {
			finalRecord = OnlyService(ipNets, command.Opts.ThreadingNum)
		} else {
			finalRecord = RunMultiThread(ipNets, podNets, command.Opts.ThreadingNum)
		}
		printer.PrintResult(finalRecord, command.Opts.OutputFile)

		PostRun(finalRecord, command.Opts.OutputFile)
	},
}

func mergeRecords(cs ...<-chan define.Record) chan define.Record {
	out := make(chan define.Record)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan define.Record) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func RunMultiThread(net, pod *net.IPNet, count int) (finalRecord define.Records) {
	scan := mutli.ScanAll(net, count)
	scan2 := mutli.ScanNeighborSvc(pod, count)
	for r := range mergeRecords(scan, scan2) {
		finalRecord = append(finalRecord, r)
	}
	return
}

func OnlyService(net *net.IPNet, count int) (finalRecord define.Records) {
	scan := mutli.ScanAll(net, count)
	for r := range scan {
		finalRecord = append(finalRecord, r)
	}
	return
}

func PostRun(finalRecord define.Records, file string) {
	if finalRecord == nil || len(finalRecord) == 0 {
		return
	}
	var f *os.File = nil
	if file != "" {
		var err error
		f, err = os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Warnf("OpenFile failed: %v", err)
		}
		defer f.Close()
	}
	writeString := func(s string) {
		if f != nil {
			_, _ = f.WriteString(s)
		}
	}
	writeString("// Extracted information under \n")

	{ // Namespaces
		log.Info("Extract Namespaces: ")
		writeString("// Extract Namespaces:\n")
		list := post.RecordsDumpNameSpace(finalRecord, command.Opts.Zone)
		for _, ns := range list {
			log.Infof("Namespace: %s", ns)
		}
		writeString("// Namespace: [" + strings.Join(list, ",") + "]\n")
	}

	{ // service
		log.Info("Extract Service: ")
		writeString("// Extract Service: \n")
		list := post.RecordsDumpFullService(finalRecord, command.Opts.Zone)
		for _, svc := range list {
			if strings.Contains(svc, "metrics") {
				if strings.Contains(svc, "kube-state-metrics") {
					log.Warnf("Checkout service %v, which maybe contains cluster metrics information", svc)
				} else {
					log.Warnf("Checkout service %v, which maybe contains apps metrics information", svc)
				}
			}
			if strings.Contains(svc, "ingress-nginx-controller-admission") ||
				(strings.Contains(svc, "admission") &&
					(strings.Contains(svc, "nginx-ingress") || strings.Contains(svc, "ingress-nginx"))) {
				log.Warnf("Checkout service %v, which maybe effects CVE-2025-1974 nginx ingress controllor ", svc)
			}
			log.Infof("Service: %s", svc)
			writeString("// \t\t" + svc + "\n")
		}
	}

	{ // service maps with pods and service ips
		log.Info("Possible Pod and service ip maps")
		writeString("// Service maps: \n")
		maps := post.PodServiceMap(finalRecord, command.Opts.Zone)
		for svc, ips := range maps {
			log.Infof("Service: %s\n\tips: [%s]", svc, strings.Join(ips, ","))
			writeString("// Service: " + svc + "\n//\tips:\t" + strings.Join(ips, "\n//\t\t") + "\n")
		}
	}
}
