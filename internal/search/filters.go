package search

import (
	"grafana-helper/internal/grafana_client/entities"
	"strings"
)

func checkPanelDatasource(p *entities.Panel) bool {
	return p.Datasource == "prometheus"
}

func checkPanelMetric(p *entities.Panel, metric string) bool {
	for _, t := range p.Targets {
		if strings.Contains(t.Expr, metric) {
			return true
		}
	}

	return false
}

func checkVarDatasource(v string) bool {
	return strings.Contains(v, "prometheus")
}

func checkVarMetric(v string, metric string) bool {
	return strings.Contains(v, metric)
}
