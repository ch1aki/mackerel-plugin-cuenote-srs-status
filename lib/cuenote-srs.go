package mpcuenotesrs

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

// CuenoteSrsStatPlugin mackerel plugin for CuenoteSrsStat
type CuenoteSrsStatPlugin struct {
	Prefix              string
	Tempfile            string
	Host                string
	User                string
	Password            string
	EnableGroupStats    bool
	EnableDeliveryStats bool
}

// MetricKeyPrefix interface for PluginWithPrefix
func (c CuenoteSrsStatPlugin) MetricKeyPrefix() string {
	if c.Prefix == "" {
		c.Prefix = "cuenote-srs-stat"
	}
	return c.Prefix
}

// GraphDefinition interface for mackerelplugin
func (c CuenoteSrsStatPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(c.MetricKeyPrefix())

	graphDef := map[string]mp.Graphs{
		"queue_total": {
			Label: labelPrefix + " Queue Total Status",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "delivering", Label: "delivering", Diff: false, Stacked: false},
				{Name: "undelivered", Label: "undelivering", Diff: false, Stacked: false},
				{Name: "resend", Label: "resend", Diff: false, Stacked: false},
			},
		},
	}

	if c.EnableGroupStats {
		graphDef = c.addGraphDefGroup(graphDef)
	}

	return graphDef
}

func (c CuenoteSrsStatPlugin) addGraphDefGroup(graphdef map[string]mp.Graphs) map[string]mp.Graphs {
	types := [...]string{
		"delivering",
		"undelivered",
		"resend",
	}
	labelPrefix := strings.Title(c.MetricKeyPrefix())

	for _, t := range types {
		graphdef["queue_group."+t] = mp.Graphs{
			Label: labelPrefix + " Queue Group Status " + strings.Title(t),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "*", Label: "%1", Diff: false},
			},
		}
	}

	return graphdef
}

func (c CuenoteSrsStatPlugin) newRequest(reqType string) (*http.Request, error) {
	p := url.Values{}
	p.Add("cmd", "get_stat")
	p.Add("type", reqType)
	u := url.URL{Scheme: "https", Host: c.Host, Path: "api", RawQuery: p.Encode()}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.User, c.Password)
	req.Header.Set("User-Agent", "mackerel-plugin-cuenote-srs-status")

	return req, nil
}

// FetchMetrics interface for mackerelplugin
func (c CuenoteSrsStatPlugin) FetchMetrics() (map[string]float64, error) {
	statRet := make(map[string]float64)

	req, err := c.newRequest("now_total")
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return nil, errors.New("Forbidden")
	}

	statRet, err = c.parseNowTotal(resp.Body)
	if err != nil {
		return nil, err
	}

	if c.EnableGroupStats {
		reqGroup, err := c.newRequest("now_group")
		if err != nil {
			return nil, err
		}

		respGroup, err := http.DefaultClient.Do(reqGroup)
		if err != nil {
			return nil, err
		}

		groupStat, err := c.parseNowGroup(respGroup.Body)
		if err != nil {
			return nil, err
		}

		for k, v := range groupStat {
			statRet[k] = v
		}
	}

	return statRet, nil
}

func (c CuenoteSrsStatPlugin) parseNowTotal(body io.Reader) (map[string]float64, error) {
	stat := make(map[string]float64)
	re := regexp.MustCompile(`(delivering|undelivered|resend)\t([0-9]+)`)

	r := bufio.NewReader(body)
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		res := re.FindStringSubmatch(string(line))
		if res != nil && len(res) == 3 {
			stat[res[1]], err = strconv.ParseFloat(res[2], 64)
			if err != nil {
				return nil, errors.New("cannot get values")
			}
		}
	}

	return stat, nil
}

func (c CuenoteSrsStatPlugin) parseNowGroup(body io.Reader) (map[string]float64, error) {
	stat := make(map[string]float64)
	re := regexp.MustCompile(`#?(\S+)\t(delivering|undelivered|resend)\t([0-9]+)`)

	reader := bufio.NewReader(body)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		res := re.FindStringSubmatch(string(line))
		if res != nil && len(res) == 4 {
			stat["queue_group."+res[2]+"."+res[1]], err = strconv.ParseFloat(res[3], 64)
			if err != nil {
				return nil, errors.New("cannot get values")
			}
		}
	}

	return stat, nil
}

func (c CuenoteSrsStatPlugin) parseDeliveryGroup(body io.Reader) (map[string]float64, error) {
	stat := make(map[string]float64)
	re := regexp.MustCompile(`#?(\S+)\t(success|failure|deferral|dnsdeferral|connfail|exception|dnsfail|expired|canceled|bounce|exclusion)\t([0-9]+)`)

	reader := bufio.NewReader(body)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		res := re.FindStringSubmatch(string(line))
		if res != nil && len(res) == 4 {
			stat["delivery_group."+res[2]+"."+res[1]], err = strconv.ParseFloat(res[3], 64)
			if err != nil {
				return nil, errors.New("cannot get values")
			}
		}
	}

	return stat, nil
}

type options struct {
	User                string `short:"u" long:"user" required:"true" description:"Cuenote SR-S username"`
	Password            string `short:"p" long:"password" required:"true" description:"Cuenote SR-S password"`
	Host                string `short:"H" long:"host" required:"true" description:"Cuenote SR-S hostname (e.g. srsXXXX.cuenote.jp)"`
	Prefix              string `long:"prefix" description:"metric key prefix (default: cuenote-srs-stat)"`
	Tempfile            string `long:"tempfile" description:"Tempfile name"`
	EnableGroupStats    bool   `long:"group-stats" description:"Enable Grouped status (default: false)"`
	EnableDeliveryStats bool   `long:"delivery-stats" description:"Enable Delivery status (default: false)"`
}

// Do the plugin
func Do() {
	var opts options
	_, err := flags.ParseArgs(&opts, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	c := CuenoteSrsStatPlugin{
		Prefix:              opts.Prefix,
		Host:                opts.Host,
		User:                opts.User,
		Password:            opts.Password,
		EnableGroupStats:    opts.EnableGroupStats,
		EnableDeliveryStats: opts.EnableDeliveryStats,
	}
	helper := mp.NewMackerelPlugin(c)
	helper.Tempfile = opts.Tempfile
	helper.Run()
}
