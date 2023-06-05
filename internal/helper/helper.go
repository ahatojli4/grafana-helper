package helper

import "github.com/ahatojli4/grafana-helper/internal/entities"

func Uniq(s []string) []string {
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

func ConvertPanelsTo(tt []entities.ResultPanel) []string {
	rr := make([]string, 0, len(tt))
	for _, t := range tt {
		rr = append(rr, t.Title)
	}

	return rr
}

func ConvertVariablesTo(tt []entities.ResultVariable) []string {
	rr := make([]string, 0, len(tt))
	for _, t := range tt {
		rr = append(rr, t.Title)
	}

	return rr
}
