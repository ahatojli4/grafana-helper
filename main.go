package main

import (
	"flag"
	"fmt"
	"grafana-helper/internal/cmd"
	entities2 "grafana-helper/internal/entities"
	"grafana-helper/internal/grafana_client"
	"grafana-helper/internal/search"
	"net/url"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"
)

func main() {
	if runtime.GOOS != "darwin" {
		fmt.Println("Error: this util is only for macos")
		os.Exit(1)
	}

	if len(os.Args[1:]) == 0 {
		fmt.Println("usage: ghelper [-set-config] <metric_name>")
		return
	}

	var flagSetConfig = flag.Bool("set-config", false, "set config")
	flag.Parse()
	if *flagSetConfig {
		_, err := cmd.SetCfg()
		if err != nil {
			panic(err)
		}

		return
	}
	cfg := cmd.LoadConfig()
	client := grafana_client.New(cfg.Grafana.Host, cfg.Grafana.Auth.BasicHeader())
	s := search.New(client)
	resCh := make(chan *entities2.ResultDashboard, 10)
	go s.Find(os.Args[1:][0], resCh)
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	_, _ = fmt.Fprintf(w, "Title\tUrl\tIn Dashboard Vars\tPanels\n")
	for r := range resCh {
		uniquePanelNames := uniq(convertPanelsTo(r.Panels))
		uniqueVariables := uniq(convertVariablesTo(r.Variables))
		u := url.URL{
			Scheme: "https",
			Host:   cfg.Grafana.Host,
			Path:   r.Url,
		}

		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.Title, u.String(), strings.Join(uniqueVariables, ", "), strings.Join(uniquePanelNames, ", "))
	}
	_, _ = fmt.Fprintf(w, "Done")
	_ = w.Flush()
}

func uniq(s []string) []string {
	m := make(map[string]struct{}, len(s))
	for _, v := range s {
		m[v] = struct{}{}
	}
	res := make([]string, 0, len(m))
	for k := range m {
		res = append(res, k)
	}

	return res
}

func convertPanelsTo(tt []entities2.ResultPanel) []string {
	rr := make([]string, 0, len(tt))
	for _, t := range tt {
		rr = append(rr, t.Title)
	}

	return rr
}

func convertVariablesTo(tt []entities2.ResultVariable) []string {
	rr := make([]string, 0, len(tt))
	for _, t := range tt {
		rr = append(rr, t.Title)
	}

	return rr
}
