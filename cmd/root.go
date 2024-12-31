package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/esonhugh/k8spider/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Opts = struct {
	Cidr    string
	PodCidr string

	DnsServer  string
	Zone       string
	OutputFile string
	Verbose    int

	MultiThreadingMode bool
	ThreadingNum       int

	SkipKubeDNSCheck bool

	FilterRules   []string
	FilterStrings []string
}{}

func defaultPodCidr() string {
	interfaces, _ := net.Interfaces()
	for _, i := range interfaces {
		if i.Name == "eth0" {
			addrs, _ := i.Addrs()
			if addrs != nil || len(addrs) > 0 {
				ip := strings.Split(addrs[0].String(), "/")[0]
				return fmt.Sprintf("%v/16", ip)
			}
		}
	}
	return "10.0.0.1/16"
}

func defaultCidr() string {
	if host := os.Getenv("KUBERNETES_SERVICE_HOST"); host != "" {
		return host + "/16"
	}
	return "10.96.0.1/16"
}

func init() {

	RootCmd.PersistentFlags().StringVarP(&Opts.Cidr, "cidr", "c", defaultCidr(), "cidr like: 192.168.0.1/16")
	RootCmd.PersistentFlags().StringVarP(&Opts.PodCidr, "pod-cidr", "p", defaultPodCidr(), "pod cidr list, watch out for the network interface name, default is eth0")

	RootCmd.PersistentFlags().StringVarP(&Opts.DnsServer, "dns-server", "d", "", "dns server")
	RootCmd.PersistentFlags().IntVarP(&pkg.DnsTimeout, "dns-timeout", "i", 2, "dns timeout")

	RootCmd.PersistentFlags().StringVarP(&Opts.Zone, "zone", "z", "cluster.local", "zone")

	RootCmd.PersistentFlags().StringVarP(&Opts.OutputFile, "output-file", "o", "", "output file")

	RootCmd.PersistentFlags().CountVarP(&Opts.Verbose, "verbose", "v", "log level (-v debug,-vv trace, info")

	RootCmd.PersistentFlags().BoolVarP(&Opts.MultiThreadingMode, "thread", "t", false, "multi threading mode, work pair with -n")
	RootCmd.PersistentFlags().IntVarP(&Opts.ThreadingNum, "thread-num", "n", 16, "threading num, default 16")

	RootCmd.PersistentFlags().BoolVarP(&Opts.SkipKubeDNSCheck, "skip-kube-dns-check", "k", false, "skip kube-dns check, force check if current environment is matched kube-dns schema")

	RootCmd.PersistentFlags().StringSliceVarP(&Opts.FilterRules, "filter-rules", "F", []string{}, "filter regexp rules")
	RootCmd.PersistentFlags().StringSliceVarP(&Opts.FilterStrings, "filter-strings", "f", []string{}, "filter contained strings")
}

var RootCmd = &cobra.Command{
	Use:   "k8spider",
	Short: "k8spider is a tool to discover k8s services",
	Long:  "k8spider is Powerful+Fast+Low Privilege Kubernetes service discovery tools via kubernetes DNS service. Currently supported service ip-port BruteForcing / AXFR Domain Transfer Dump / Coredns WildCard Dump / Pod Verified IP discovery\n\nTopics\n",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set Log Levels
		SetLogLevel(Opts.Verbose)
		// Set pkg global config
		pkg.Zone = Opts.Zone

		// debug print option data
		opt, _ := json.MarshalIndent(Opts, "", "  ")
		log.Tracef("Opts: %v", string(opt))

		if Opts.DnsServer != "" {
			pkg.NetResolver = pkg.WarpDnsServer(Opts.DnsServer)
		}
		for _, rules := range Opts.FilterRules {
			pkg.NetResolver.SetFilter(rules)
		}
		for _, rules := range Opts.FilterStrings {
			pkg.NetResolver.SetContainsFilter(rules)
		}
		// Check if current environment is a kubernetes cluster
		// If the command is whereisdns, which means DNS is not sure , so skip this check!
		// If SkipKubeDNSCheck is true, skip this check!
		if Opts.SkipKubeDNSCheck == false {
			if cmd.Use != "whereisdns" {
				if !pkg.CheckKubeDNS() {
					log.Warn("current environment is not a kubernetes cluster")
					os.Exit(1)
				}
			}
		} else {
			log.Tracef("kubernetes environment checking bypassed")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var LevelMap = []log.Level{
	log.InfoLevel,  // 0
	log.DebugLevel, // 1
	log.TraceLevel, // 2
}

func SetLogLevel(level int) {
	if level > 2 {
		level = 2
	}
	log.SetLevel(LevelMap[level])
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
