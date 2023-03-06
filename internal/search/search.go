package search

import (
	"fmt"
	"sync"

	entities2 "github.com/ahatojli4/ghelper/internal/entities"
	"github.com/ahatojli4/ghelper/internal/grafana_client"
	"github.com/ahatojli4/ghelper/internal/grafana_client/entities"
)

type search struct {
	client           *grafana_client.Client
	searchDashboards []entities.SearchDashboard

	inflightRequests chan struct{}
	findChannel      chan *combinedResponse
	errChannel       chan error
}

type combinedResponse struct {
	sd *entities.SearchDashboard
	d  *entities.Dashboard
}

func New(client *grafana_client.Client) *search {
	dd, err := client.Search()
	if err != nil {
		panic(err)
	}

	return &search{
		client:           client,
		searchDashboards: dd,
		inflightRequests: make(chan struct{}, 10),
		findChannel:      make(chan *combinedResponse),
		errChannel:       make(chan error),
	}
}

func (s *search) Find(metric string, result chan<- *entities2.ResultDashboard) {
	go s.Check(metric, result)
	go func() {
		for err := range s.errChannel {
			fmt.Println("Error: " + err.Error())
		}
	}()
	wg := sync.WaitGroup{}
	for _, sd := range s.searchDashboards {
		s.inflightRequests <- struct{}{}
		wg.Add(1)
		go func(sd *entities.SearchDashboard) {
			defer wg.Done()
			defer func() {
				<-s.inflightRequests
			}()
			d, err := s.client.DashboardByUid(sd.Uid)
			if err != nil {
				s.errChannel <- err
				return
			}
			s.findChannel <- &combinedResponse{
				sd: sd,
				d:  d,
			}
		}(&sd)
	}
	wg.Wait()
	close(s.findChannel)

	return
}

func (s *search) Check(metric string, resultChannel chan<- *entities2.ResultDashboard) {
	for r := range s.findChannel {
		res := &entities2.ResultDashboard{
			Title:     r.sd.Title,
			Url:       r.d.Meta.Url,
			Panels:    make([]entities2.ResultPanel, 0),
			Variables: make([]entities2.ResultVariable, 0),
		}
		found := false
		for _, panel := range r.d.Dashboard.Panels {
			if checkPanelDatasource(&panel) && checkPanelMetric(&panel, metric) {
				found = true
				res.Panels = append(res.Panels, entities2.ResultPanel{Title: panel.Title})
			}
		}

		if checkVarDatasource(string(r.d.Dashboard.Templating.List)) && checkVarMetric(string(r.d.Dashboard.Templating.List), metric) {
			found = true
			res.Variables = append(res.Variables, entities2.ResultVariable{Title: "True"})
		}
		if found {
			resultChannel <- res
		}
	}
	close(resultChannel)
}
