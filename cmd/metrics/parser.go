package metrics

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"strings"

	cmdx "github.com/esonhugh/k8spider/cmd"
	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg/metrics"
	"github.com/esonhugh/k8spider/pkg/metrics/coredns"
	kube_state_metrics "github.com/esonhugh/k8spider/pkg/metrics/kube-state-metrics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	cmdx.RootCmd.AddCommand(MetricsCmd)
}

var MetricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "parse kube stat metrics to readable resource, use: k8spider metrics <file/url> <file/url> ...",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			return
		}

		for _, From := range args {
			log.Debugf("parse metrics from %v", From)
			rule := append(kube_state_metrics.DefaultMatchRules(), coredns.CoreDNSMatchRules()...)
			if err := rule.Compile(); err != nil {
				log.Fatalf("compile rule failed: %v", err)
			}
			log.Debugf("compiled rules completed, start to get resource \n")

			ot := output()

			var r io.Reader
			if strings.HasPrefix("http://", From) || strings.HasPrefix("https://", From) {
				resp, err := http.Get(From)
				if err != nil {
					log.Fatalf("get metrics from %v failed: %v", From, err)
				}
				defer resp.Body.Close()
				r = resp.Body
			} else {
				f, err := os.OpenFile(From, os.O_RDONLY, 0666)
				if err != nil {
					log.Fatalf("open file %v failed: %v", From, err)
				}
				defer f.Close()
				r = f
			}
			log.Debugf("start to parse metrics line by line\n")

			var rx []*metrics.MetricMatcher

			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				line := scanner.Text()
				res, err := rule.Match(line)
				if err != nil {
					continue
				} else {
					log.Debugf("matched: %s", res.DumpString())
					rx = append(rx, res)
				}
			}
			if err := scanner.Err(); err != nil {
				log.Warnf("scan metrics failed and break out, reason: %v", err)
			}
			var res define.ResourceList = metrics.ConvertToResource(rx)
			log.Debugf("parse metrics completed, start to print result\n")

			res.Print(ot)
		}

	},
}

func output() io.WriteCloser {
	if cmdx.Opts.OutputFile != "" {
		f, err := os.OpenFile(cmdx.Opts.OutputFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Warnf("create output file failed: %v", err)
			return nil
		}
		return f
	} else {
		return os.Stdout
	}
}
