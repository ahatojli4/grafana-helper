package grafana_client

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/ahatojli4/grafana-helper/internal/grafana_client/entities"
)

type Client struct {
	host, auth string
	client     *http.Client
}

func New(host, auth string) *Client {
	if !(strings.HasPrefix(host, "https://") || strings.HasPrefix(host, "http://")) {
		host = "https://" + host
	}

	return &Client{
		host:   host,
		auth:   auth,
		client: &http.Client{},
	}
}

// Search /api/search/
func (c *Client) Search() ([]entities.SearchDashboard, error) {
	req, err := http.NewRequest("GET", c.host+"/api/search/?limit=5000", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.auth)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := make([]entities.SearchDashboard, 1000)
	err = json.Unmarshal(raw, &result)

	return result, err
}

// DashboardByUid /api/dashboards/uid/:uid
func (c *Client) DashboardByUid(uid string) (*entities.Dashboard, error) {
	req, err := http.NewRequest("GET", c.host+"/api/dashboards/uid/"+uid, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.auth)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := entities.Dashboard{}
	err = json.Unmarshal(raw, &result)

	return &result, err
}
