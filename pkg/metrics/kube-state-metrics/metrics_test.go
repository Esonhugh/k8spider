package kube_state_metrics

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg/metrics"
)

func TestMetrics(t *testing.T) {
	f, err := os.Open("./metrics")
	if err != nil {
		t.Fatalf("open file failed: %v", err)
		t.Fail()
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	rule := DefaultMatchRules()
	if err := rule.Compile(); err != nil {
		t.Fatalf("compile rule failed: %v", err)
		t.Fail()
	}

	output, err := os.OpenFile("./output.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("open output file failed: %v", err)
	}
	defer output.Close()
	var rx []*metrics.MetricMatcher
	for scanner.Scan() {
		line := scanner.Text()
		res, err := rule.Match(line)
		if err != nil {
			continue
		} else {
			t.Logf("matched: %s", res.DumpString())
			_, _ = output.WriteString(res.DumpString() + "\n")
			rx = append(rx, res.CopyData())
		}
	}
}

func TestConvertToResource(t *testing.T) {
	output, err := os.OpenFile("./output.txt", os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		t.Fatalf("open output file failed: %v", err)
	}
	defer output.Close()
	var rules []*metrics.MetricMatcher
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		line := scanner.Text()
		var r *metrics.MetricMatcher
		e := json.Unmarshal([]byte(line), &r)
		if e != nil {
			t.Logf("unmarshal failed: %v", e)
			continue
		}
		rules = append(rules, r)
	}
	var res define.ResourceList = metrics.ConvertToResource(rules)
	res.Print(os.Stderr)
}
