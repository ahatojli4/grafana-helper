package entities

type ResultPanel struct {
	Title string
}

type ResultVariable struct {
	Title string
}

type ResultDashboard struct {
	Title     string
	Url       string
	Panels    []ResultPanel
	Variables []ResultVariable
}
