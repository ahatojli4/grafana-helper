package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ahatojli4/grafana-helper/internal/entities"
	"github.com/ahatojli4/grafana-helper/internal/grafana_client"
	"github.com/ahatojli4/grafana-helper/internal/helper"
	"github.com/ahatojli4/grafana-helper/internal/search"
)

var (
	apiKey      = flag.String("api-key", "eyJrIjoidW1MNVNtUXRyVVNRRUZSN0pyc1V6ZHEzMFF5clV0NGsiLCJuIjoiZC5ib3Jpc2V2aWNodXMiLCJpZCI6MX0=", "grafana api key")
	session     = flag.String("session", "7e6726ddb14a7abbc04e9e90372f2fee", "grafana session")
	grafanaHost = flag.String("grafana-host", "grafana.rtty.in", "grafana host")
	port        = flag.String("port", "4141", "port")
)

func main() {
	flag.Parse()

	srv := &Server{
		grafanaClient: grafana_client.New(
			*grafanaHost,
			entities.ApiKey(*apiKey).BearerHeader(),
			*session,
		),
		srv: &http.Server{},
	}

	dir, err := os.ReadDir("./frontend/build/")
	if err != nil {
		fmt.Println(err)
		return
	}
	for i, entry := range dir {
		fmt.Println(i, entry.Name())
	}
	http.Handle("/search", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := r.URL.Query().Get("metric")
		w.Header().Set("Content-Type", "application/json")
		resCh := make(chan *entities.ResultDashboard, 10)
		go search.New(srv.grafanaClient).Find(m, resCh)
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
		_, _ = w.Write(b)
		w.WriteHeader(http.StatusOK)
	}))
	http.Handle("/", http.FileServer(http.Dir("./frontend/build"))) // todo: embed frontend
	log.Fatal(http.ListenAndServe(":"+*port, nil))

	fmt.Println(srv)
}

type Server struct {
	grafanaClient *grafana_client.Client
	srv           *http.Server
}
