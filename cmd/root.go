package cmd

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/esonhugh/k8spider/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Opts = struct {
	Cidr       string
	DnsServer  string
	SvcDomains []string
	Zone       string
	OutputFile string
	LogLevel   string
	Threads    int
}{}

func init() {
	RootCmd.PersistentFlags().StringVarP(&Opts.Cidr, "cidr", "c", os.Getenv("KUBERNETES_SERVICE_HOST")+"/16", "cidr like: 192.168.0.1/16")
	RootCmd.PersistentFlags().StringVarP(&Opts.DnsServer, "dns-server", "d", "", "dns server")
	RootCmd.PersistentFlags().StringSliceVarP(&Opts.SvcDomains, "svc-domains", "s", []string{}, "service domains, like: kubernetes.default,etcd.default don't add zone like svc.cluster.local")
	RootCmd.PersistentFlags().StringVarP(&Opts.Zone, "zone", "z", "cluster.local", "zone")
	RootCmd.PersistentFlags().StringVarP(&Opts.OutputFile, "output-file", "o", "", "output file")
	RootCmd.PersistentFlags().StringVarP(&Opts.LogLevel, "log-level", "l", "info", "log level")
	RootCmd.PersistentFlags().IntVarP(&Opts.Threads, "threads", "t", 4, "threads")
}

var RootCmd = &cobra.Command{
	Use:   "k8spider",
	Short: "k8spider is a tool to discover k8s services",
	Long:  "k8spider is a tool to discover k8s services",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		SetLogLevel(Opts.LogLevel)
		if Opts.DnsServer != "" {
			pkg.NetResolver = &net.Resolver{
				PreferGo: true,
				Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
					d := net.Dialer{}
					return d.DialContext(ctx, network, Opts.DnsServer)
				},
			}
		}
		if Opts.Threads < 1 {
			log.Fatalln("threads must be greater than 0")
		}
		pkg.Threads = Opts.Threads
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func SetLogLevel(level string) {
	switch level {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
