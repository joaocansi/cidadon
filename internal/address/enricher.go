package address

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Enricher interface {
	Enrich(ctx context.Context, address *Address) error
}

type NominatimEnricher struct {
	Client *http.Client
	URL    string
}

func NewNominatimAddressEnricher(client *http.Client, url string) *NominatimEnricher {
	return &NominatimEnricher{
		Client: client,
		URL:    url,
	}
}

const baseEnrichAddressUrl = "%s/reverse?lat=%f&lon=%f&format=json"

type Response struct {
	Address struct {
		Suburb      string `json:"suburb"`
		City        string `json:"city"`
		State       string `json:"state"`
		Region      string `json:"address"`
		Postcode    string `json:"postcode"`
		Country     string `json:"country"`
		CountryCode string `json:"country_code"`
	} `json:"address"`
}

func (ae *NominatimEnricher) Enrich(ctx context.Context, address *Address) error {
	u, err := url.Parse(ae.URL + "/reverse")
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set("lat", fmt.Sprintf("%f", address.CoordinateX))
	q.Set("lon", fmt.Sprintf("%f", address.CoordinateY))
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set(
		"User-Agent",
		"cidadon/1.0 (https://your-domain.com; contact@your-domain.com)",
	)

	resp, err := ae.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("address enrichment returned status %d", resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("decode address response: %w", err)
	}

	address.City = response.Address.City
	address.State = response.Address.State
	address.Neighborhood = response.Address.Suburb
	address.Postcode = response.Address.Postcode

	return nil
}
