package entities

type Panel struct {
	Title      string `json:"title"`
	Datasource string `json:"datasource,omitempty"`
	Targets    []struct {
		Expr           string `json:"expr"`
		Format         string `json:"format"`
		Hide           bool   `json:"hide"`
		IntervalFactor int    `json:"intervalFactor"`
		LegendFormat   string `json:"legendFormat"`
		RefId          string `json:"refId"`
	} `json:"targets,omitempty"`
}

type list string

func (l *list) UnmarshalJSON(bytes []byte) error {
	*l = list(bytes)

	return nil
}

type Dashboard struct {
	Meta struct {
		Url string `json:"url"`
	} `json:"meta"`
	Dashboard struct {
		Panels     []Panel `json:"panels"`
		Templating struct {
			List list `json:"list"` // just because awesome grafana sometimes has object in this place instead of array
		} `json:"templating"`
		Title string `json:"title"`
		Uid   string `json:"uid"`
	} `json:"dashboard"`
}
