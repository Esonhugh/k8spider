package coredns

import (
	"bufio"
	"os"
	"testing"

	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg/metrics"
)

func TestCoreDNSMatchRules(t *testing.T) {
	t.Logf("start TestCoreDNSMatchRules")
	rule := CoreDNSMatchRules()

	f, err := os.Open("./test_coredns")
	if err != nil {
		t.Fatalf("open file failed: %v", err)
		t.Fail()
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	if err := rule.Compile(); err != nil {
		t.Fatalf("compile rule failed: %v", err)
		t.Fail()
	}

	var resList []*metrics.MetricMatcher
	for scanner.Scan() {
		line := scanner.Text()
		res, err := rule.Match(line)
		if err != nil {
			continue
		} else {
			resList = append(resList, res)
			t.Logf("matched: %s", res.DumpString())
		}
	}
	e := define.ConvertToResource(resList, MergeCoreDnsPlugin)
	for _, r := range e {
		t.Logf("resource: %s", r.JSON())
	}
}
