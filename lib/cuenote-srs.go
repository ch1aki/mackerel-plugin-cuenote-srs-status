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

	"github.com/jessevdk/go-flags"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

// CuenoteSrsStatPlugin mackerel plugin for CuenoteSrsStat
type CuenoteSrsStatPlugin struct {
	Prefix           string
	Tempfile         string
	Host             string
	User             string
	Password         string
	EnableGroupStats bool
}

// GraphDefinition interface for mackerelplugin
func (c CuenoteSrsStatPlugin) GraphDefinition() map[string]mp.Graphs {
	graphDef := map[string]mp.Graphs{
		"cuenote-srs.queue_total": {
			Label: "Cuenote SR-S Queue Total Status",
			Unit:  "float",
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
		"success",
		"failure",
		"dnsfail",
		"exclusion",
		"bounce_unique",
		"canceled",
		"expired",
		"deferral",
		"dnsdeferral",
		"connfail",
		"bounce",
		"exception",
	}

	for _, t := range types {
		graphdef["cuenote-srs.queue_group_"+t] = mp.Graphs{
			Label: "Cuenote SR-S Queue Group Status Delivering",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "*", Label: "%1", Diff: false},
			},
		}
	}

	return graphdef
}

func (c CuenoteSrsStatPlugin) fetchStat(type_ string) (io.Reader, error) {
	p := url.Values{}
	p.Add("cmd", "get_stat")
	p.Add("type", type_)
	u := url.URL{Scheme: "https", Host: c.Host, Path: "api", RawQuery: p.Encode()}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.User, c.Password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return nil, errors.New("Forbidden")
	}

	return resp.Body, nil
}

// FetchMetrics interface for mackerelplugin
func (c CuenoteSrsStatPlugin) FetchMetrics() (map[string]float64, error) {
	statRet := make(map[string]float64)
	body, err := c.fetchStat("now_total")
	if err != nil {
		return nil, err
	}
	statRet, err = c.parseNowTotal(body)
	if err != nil {
		return nil, err
	}

	if c.EnableGroupStats {
		body, err := c.fetchStat("now_total")
		if err != nil {
			return nil, err
		}
		groupStat, err := c.parseNowGroup(body)
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

	r := bufio.NewReader(body)
	for _, m := range [...]string{"delivering", "undelivered", "resend"} {
		line, _, err := r.ReadLine()
		if err != nil {
			return nil, errors.New("cannot get values")
		}
		re := regexp.MustCompile(m + "\t([0-9]+)")
		res := re.FindStringSubmatch(string(line))
		if res == nil || len(res) != 2 {
			return nil, errors.New("cannot get values")
		}
		stat[m], err = strconv.ParseFloat(res[1], 64)
		if err != nil {
			return nil, errors.New("cannot get values")
		}
	}

	return stat, nil
}

func (c CuenoteSrsStatPlugin) parseNowGroup(body io.Reader) (map[string]float64, error) {
	stat := make(map[string]float64)
	re := regexp.MustCompile(`#?(\S+)\t(\S+)\t([0-9]+)`)

	reader := bufio.NewReader(body)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		res := re.FindStringSubmatch(string(line))
		if res == nil || len(res) != 4 {
			return nil, errors.New("cannnot parse responce")
		}

		stat["queue_group_"+res[2]+"."+res[1]], err = strconv.ParseFloat(res[3], 64)
		if err != nil {
			return nil, errors.New("cannot get values")
		}
	}

	return stat, nil
}

type options struct {
	User             string `short:"u" long:"user" description:"Cuenote SR-S username"`
	Password         string `short:"p" long:"password" description:"Cuenote SR-S password"`
	Host             string `short:"H" long:"host" description:"Cuenote SR-S hostname (e.g. srsXXXX.cuenote.jp)"`
	Prefix           string `long:"prefix" default:"cuenote-srs-stat" description:"metric key prefix"`
	Tempfile         string `long:"template" description:"Tempfile name"`
	EnableGroupStats bool   `long:"group-stats" description:"Enable Grouped status"`
}

// Do the plugin
func Do() {
	var opts options
	_, err := flags.ParseArgs(&opts, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	c := CuenoteSrsStatPlugin{
		Prefix:           opts.Prefix,
		Host:             opts.Host,
		User:             opts.User,
		Password:         opts.Password,
		EnableGroupStats: opts.EnableGroupStats,
	}
	helper := mp.NewMackerelPlugin(c)
	helper.Tempfile = opts.Tempfile
	helper.Run()
}
