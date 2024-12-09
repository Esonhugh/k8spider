package metrics

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/elastic/go-grok"
	log "github.com/sirupsen/logrus"
)

var (
	MetricsPatterns = map[string]string{
		"METRIC_HEAD": "^%{WORD:metric_name}{",
		"LAST_LABEL":  `(%{DATA:last_label_name}="%{DATA:last_label_value}",?)`,
		"NUMBER_SCI":  "%{NUMBER}(e%{NUMBER})?",
		"METRIC_TAIL": `%{SPACE}%{NUMBER_SCI:metric_value}$`,
		"NORMAL_EXPR": `^%{WORD:metric_name}{((namespace="%{DATA:namespace}")|(%{DATA:last_label_name}=\"%{DATA:last_label_value}\")(,|))*}%{SPACE}%{NUMBER}(e%{NUMBER})?$`,
	}
	COMMON_MATCH_GROK = grok.New()
)

const COMMON_MATCH = `^%{WORD:name}{((namespace="%{DATA:namespace}")|(%{DATA:last_label_name}="%{DATA:last_label_value}")(,|))*}%{SPACE}%{NUMBER}(e%{NUMBER})?$`

type Label struct {
	Key   string `json:"key"`
	Value string `json:"value""`
}

type MetricMatcher struct {
	Type         string     `json:"type"`
	Header       string     `json:"-"`
	Labels       []Label    `json:"labels"`
	nameLabel    string     `json:"-"`
	grok         *grok.Grok `json:"-"`
	finalPattern string
	ptr          any `json:"-"`
}

func NewMetricMatcher(t string) *MetricMatcher {
	return &MetricMatcher{
		Type:      t,
		nameLabel: t,
		grok:      grok.New(),
		Labels:    make([]Label, 0),
	}
}

// example: ^%{WORD:name}{((namespace="%{DATA:ns}")|(%{DATA:last_label_name}="%{DATA:last_label_value}")(,|))*}%{SPACE}%{NUMBER}(e%{NUMBER})?$
//          ^%{WORD:name}{
//                        (
//                         (namespace="%{DATA:ns}")
//                                                 |
//                                                  (%{DATA:last_label_name}="%{DATA:last_label_value}")
//                                                                                                      (,|)
//                                                                                                          )
//                                                                                                           *
//                                                                                                            }%{SPACE}%{NUMBER}(e%{NUMBER})?$

func (mt *MetricMatcher) Compile() error {
	if err := mt.grok.AddPatterns(MetricsPatterns); err != nil {
		return err
	}
	if mt.Header == "" {
		mt.Header = `kube_` + mt.Type + `_info` // custom head
	}
	header := mt.Header
	var body []string
	for _, label := range mt.Labels {
		body = append(body, "("+label.Key+`=`+`"%{DATA:`+label.Key+`}",?)`)
	}
	tail := `%{METRIC_TAIL}`
	pattern := header + "{" + "(" + strings.Join(body, "|") + "|%{LAST_LABEL})*" + "}" + tail
	mt.finalPattern = pattern
	log.Debugln("generated pattern: ", pattern)
	return mt.grok.Compile(pattern, true)
}

func (mt *MetricMatcher) Match(target string) (res map[string]string, err error) {
	if !strings.HasPrefix(target, mt.Header) {
		log.Debugf("not match: %s", target)
		if COMMON_MATCH_GROK.MatchString(target) {
			res, err = COMMON_MATCH_GROK.ParseString(target)
			if err != nil {
				return nil, errors.Join(
					errors.New("match failed, in common expr"),
					err,
				)
			}
			return res, errors.Join(
				errors.New("match failed, in custom expr, but success in common expr, check ret result to get more detail"),
				err,
			)
		}
		return nil, errors.New("can't match")
	} else {
		res, err = mt.grok.ParseString(target)
		if err != nil || res == nil || len(res) == 0 {
			return nil, errors.New("match failed, no result found")
		}
		mt.setResult(res)
		return res, err
	}
}

func (mt *MetricMatcher) SetHeader(header string) *MetricMatcher {
	mt.Header = header
	return mt
}

func (mt *MetricMatcher) AddLabel(label string) *MetricMatcher {
	mt.Labels = append(mt.Labels, Label{
		Key:   label,
		Value: "",
	})
	return mt
}

func (mt *MetricMatcher) setResult(res map[string]string) {
	for i, label := range mt.Labels {
		if v, ok := res[label.Key]; ok {
			mt.Labels[i].Value = v
		} else {
			log.Debugf("label %s not found in result", label.Key) // keep it empty
		}
	}
}

func (mt *MetricMatcher) FindLabel(Key string) string {
	for _, label := range mt.Labels {
		if label.Key == Key {
			return label.Value
		}
	}
	return ""
}

func (mt *MetricMatcher) DumpString() string {
	b, _ := json.Marshal(mt)
	return string(b)
}

func (mt *MetricMatcher) SetNameLabel(n string) *MetricMatcher {
	mt.nameLabel = n
	return mt
}

func (mt *MetricMatcher) LabelNameOfName() string {
	return mt.nameLabel
}

func init() {
	err := COMMON_MATCH_GROK.AddPatterns(MetricsPatterns)
	if err != nil {
		log.Fatalf("add patterns failed: %v", err)
	}
	err = COMMON_MATCH_GROK.Compile(COMMON_MATCH, true)
	if err != nil {
		log.Fatalf("compile pattern failed: %v", err)
	}
}

func (mt *MetricMatcher) CopyData() *MetricMatcher {
	var newLabel []Label
	for _, label := range mt.Labels {
		newLabel = append(newLabel, label)
	}
	return &MetricMatcher{
		Type:      strings.Clone(mt.Type),
		Header:    strings.Clone(mt.Header),
		nameLabel: strings.Clone(mt.nameLabel),
		Labels:    newLabel,
	}
}
