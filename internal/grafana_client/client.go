package grafana_client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/ahatojli4/grafana-helper/internal/cache"
	"github.com/ahatojli4/grafana-helper/internal/grafana_client/entities"
)

type Client struct {
	host, auth, sessionId string
	client                *http.Client
	cache                 *cache.Cache
}

func New(host, auth string, c *cache.Cache) *Client {
	return &Client{
		host:   host,
		auth:   auth,
		client: &http.Client{},
		cache:  c,
	}
}

// Search /api/search/
func (c *Client) Search() ([]entities.SearchDashboard, error) {
	u := (&url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   "/api/search/",
		RawQuery: (url.Values{
			"limit": []string{"5000"},
		}).Encode(),
	}).String()
	var raw []byte
	raw, ok := c.cache.Get(u)
	if !ok {
		req, err := http.NewRequest("GET", u, nil)
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
		raw, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		c.cache.Set(u, raw)
	}
	result := make([]entities.SearchDashboard, 1000)
	err := json.Unmarshal(raw, &result)

	return result, err
}

// DashboardByUid /api/dashboards/uid/:uid
func (c *Client) DashboardByUid(uid string) (*entities.Dashboard, error) {
	u := (&url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   "/api/dashboards/uid/" + uid,
	}).String()
	var raw []byte
	raw, ok := c.cache.Get(u)
	if !ok {
		req, err := http.NewRequest("GET", u, nil)
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
		raw, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		c.cache.Set(u, raw)
	}
	result := entities.Dashboard{}
	err := json.Unmarshal(raw, &result)

	return &result, err
}
