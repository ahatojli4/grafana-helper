package main

import (
	"embed"
	"encoding/json"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/ahatojli4/grafana-helper/internal/cache"
	"github.com/ahatojli4/grafana-helper/internal/entities"
	"github.com/ahatojli4/grafana-helper/internal/grafana_client"
	"github.com/ahatojli4/grafana-helper/internal/helper"
	"github.com/ahatojli4/grafana-helper/internal/search"
)

var (
	apiKey      = flag.String("api-key", "", "grafana api key")
	grafanaHost = flag.String("grafana-host", "grafana.rtty.in", "grafana host")
	port        = flag.String("port", "4142", "port")
)

//go:embed frontend/*
var content embed.FS

func main() {
	flag.Parse()

	c := cache.New(30 * time.Minute)
	srv := &Server{
		grafanaClient: grafana_client.New(*grafanaHost, entities.ApiKey(*apiKey).BearerHeader(), c),
		srv:           &http.Server{},
	}

	http.Handle("/search", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := r.URL.Query().Get("metric")
		w.Header().Set("Content-Type", "application/json")
		resCh := make(chan *entities.ResultDashboard, 10)
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			search.New(srv.grafanaClient).Find(m, resCh)
			if c.IsExpired() {
				c.Reset()
			}
		}()
		type result struct {
			Title  string `json:"title"`
			Url    string `json:"url"`
			Vars   string `json:"vars"`
			Panels string `json:"panels"`
		}
		resp := []result{}
		for res := range resCh {
			uniquePanelNames := helper.Uniq(helper.ConvertPanelsTo(res.Panels))
			uniqueVariables := helper.Uniq(helper.ConvertVariablesTo(res.Variables))
			resp = append(resp, result{
				Title: res.Title,
				Url: (&url.URL{
					Scheme: "https",
					Host:   *grafanaHost,
					Path:   res.Url,
				}).String(),
				Vars:   strings.Join(uniqueVariables, ", "),
				Panels: strings.Join(uniquePanelNames, ", "),
			})
		}
		b, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		wg.Wait()
		_, _ = w.Write(b)
	}))

	frontend, err := fs.Sub(content, "frontend")
	if err != nil {
		return
	}
	http.Handle("/", http.FileServer(http.FS(frontend)))
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

type Server struct {
	grafanaClient *grafana_client.Client
	srv           *http.Server
}
