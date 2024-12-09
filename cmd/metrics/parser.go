package metrics

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"strings"

	cmdx "github.com/esonhugh/k8spider/cmd"
	"github.com/esonhugh/k8spider/pkg/metrics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var MetricOpt struct {
	From string
}

func init() {
	cmdx.RootCmd.AddCommand(MetricCmd)
	MetricCmd.PersistentFlags().StringVarP(&MetricOpt.From, "metric", "m", "", "metrics from (file / remote url)")

}

var MetricCmd = &cobra.Command{
	Use:   "metric",
	Short: "parse kube stat metrics to readable resource",
	Run: func(cmd *cobra.Command, args []string) {
		if MetricOpt.From == "" {
			return
		}
		log.Debugf("parse metrics from %v", MetricOpt.From)
		rule := metrics.DefaultMatchRules()
		if err := rule.Compile(); err != nil {
			log.Fatalf("compile rule failed: %v", err)
		}
		log.Debugf("compiled rules completed, start to get resource \n")

		ot := output()

		var r io.Reader
		if strings.HasPrefix("http://", MetricOpt.From) || strings.HasPrefix("https://", MetricOpt.From) {
			resp, err := http.Get(MetricOpt.From)
			if err != nil {
				log.Fatalf("get metrics from %v failed: %v", MetricOpt.From, err)
			}
			defer resp.Body.Close()
			r = resp.Body
		} else {
			f, err := os.OpenFile(MetricOpt.From, os.O_RDONLY, 0666)
			if err != nil {
				log.Fatalf("open file %v failed: %v", MetricOpt.From, err)
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
		var res metrics.ResourceList = metrics.ConvertToResource(rx)
		log.Debugf("parse metrics completed, start to print result\n")

		res.Print(ot)
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
