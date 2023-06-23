package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/ahatojli4/grafana-helper/internal/cache"
	"github.com/ahatojli4/grafana-helper/internal/cmd"
	entities2 "github.com/ahatojli4/grafana-helper/internal/entities"
	"github.com/ahatojli4/grafana-helper/internal/grafana_client"
	"github.com/ahatojli4/grafana-helper/internal/helper"
	"github.com/ahatojli4/grafana-helper/internal/search"
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
	c := cache.New(30 * time.Minute)
	c.Load()
	client := grafana_client.New(cfg.Grafana.Host, cfg.Grafana.Auth.BasicHeader(), c)
	s := search.New(client)
	resCh := make(chan *entities2.ResultDashboard, 10)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		s.Find(os.Args[1:][0], resCh)
		c.Store()
	}()
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	_, _ = fmt.Fprintf(w, "Title\tUrl\tIn Dashboard Vars\tPanels\n")
	for r := range resCh {
		uniquePanelNames := helper.Uniq(helper.ConvertPanelsTo(r.Panels))
		uniqueVariables := helper.Uniq(helper.ConvertVariablesTo(r.Variables))
		u := url.URL{
			Scheme: "https",
			Host:   cfg.Grafana.Host,
			Path:   r.Url,
		}

		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.Title, u.String(), strings.Join(uniqueVariables, ", "), strings.Join(uniquePanelNames, ", "))
	}
	wg.Wait()

	_, _ = fmt.Fprintf(w, "Done")
	_ = w.Flush()
}
