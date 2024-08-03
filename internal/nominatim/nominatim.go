// Nominatim API client for geocoding and reverse geocoding.
//
// Endpoints:
//
// Base URL: https://nominatim.openstreetmap.org
//
// /search - search OSM objects by name or type
// /reverse - search OSM object by their location
// /lookup - look up address details for OSM objects by their ID
// /status - query the status of the server
package nominatim

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/desertthunder/weather/internal/nws"
)

type Formats = string
type Endpoints = string

const (
	Json    Formats = "json"
	JsonV2  Formats = "jsonv2"
	GeoJson Formats = "geojson"
	GeoCode Formats = "geocodejson"
)

const (
	Search  Endpoints = "search"
	Reverse Endpoints = "reverse"
	Lookup  Endpoints = "lookup"
	Status  Endpoints = "status"
)

// Note: We need this to be overridden in tests and for future
// work in which a local instance of Nominatim is used.
const BaseURL string = "https://nominatim.openstreetmap.org"

// User-Agent for testing purposes.
const UserAgent string = "geocast-desertthunder@github.com"

// struct Nominatim represents the Nominatim API client.
type Nominatim struct {
	baseURL   string
	params    Params
	userAgent string
}

type nominatimSearchResult struct {
	PlaceID     int      `json:"place_id"`
	Licence     string   `json:"licence"`
	OSM_Type    string   `json:"osm_type"`
	OSM_ID      int      `json:"osm_id"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	Category    string   `json:"category"`
	Type        string   `json:"type"`
	PlaceRank   int      `json:"place_rank"`
	Importance  float64  `json:"importance"`
	AddressType string   `json:"addresstype"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	BoundingBox []string `json:"boundingbox"`
}

type NominatimSearchResponse = []nominatimSearchResult

type Params struct {
	// free form search query string
	Q           string
	Format      Formats
	Limit       int
	NameDetails bool
}

func (n *Nominatim) SetURL(url string) {
	n.baseURL = url
}

func (n Nominatim) BaseURL() string {
	return n.baseURL
}

func (n *Nominatim) SetParams(params Params) {
	n.params = params
}

func (n *Nominatim) SetUserAgent(ua string) {
	n.userAgent = ua
}

func (n *Nominatim) GetParams() Params {
	return n.params
}

func (n *Nominatim) UserAgent() string {
	return n.userAgent
}

func (p Params) String() string {
	qs := ""

	if p.Q == "" {
		return qs
	}

	qs = fmt.Sprintf("q=%s", p.Q)

	if p.Format == "" {
		p.Format = JsonV2
	}

	qs = fmt.Sprintf("%s&format=%s", qs, p.Format)

	if p.Format == JsonV2 || p.Format == Json {
		p.Limit = 25
	}

	qs = fmt.Sprintf("%s&limit=%d", qs, p.Limit)

	if p.NameDetails {
		qs = fmt.Sprintf("%s&namedetails=1", qs)
	}

	return qs
}

func (n *Nominatim) getRequest(endpoint Endpoints) ([]byte, error) {
	uri := n.baseURL

	if endpoint == Search {
		uri = fmt.Sprintf("%s/%s?%s", uri, endpoint, n.params.String())
	}

	client := &http.Client{}

	req, _ := http.NewRequest("GET", uri, nil)

	req.Header.Set("User-Agent", n.userAgent)

	rsp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Request to %s failed with error: %s\n", uri, err.Error())

		return nil, err
	}

	data, err := io.ReadAll(rsp.Body)

	if err != nil {
		fmt.Printf("Failed to read response body: %s\n", err.Error())

		return nil, err
	}

	return data, nil
}

func handleSimpleError(err error) {
	fmt.Printf("Error: %s\n", err)
}

func (n *Nominatim) Search() NominatimSearchResponse {
	d, err := n.getRequest(Search)

	if err != nil {
		handleSimpleError(err)

		return NominatimSearchResponse{}
	}

	rsp := NominatimSearchResponse{}
	json.Unmarshal(d, &rsp)

	return rsp
}

func (n *Nominatim) GeocodeByPoint(lat, lon float64) (*nws.City, error) {
	n.SetParams(Params{
		Q: fmt.Sprintf("%f,%f", lat, lon),
	})

	results := n.Search()

	if len(results) == 0 {
		return nil, errors.New("no results found for the provided point")
	}

	result := results[0]

	city := nws.BuildCity(result.DisplayName, result.Lat, result.Lon)

	return &city, nil
}

func (n *Nominatim) GeocodeByCity(c string) (*nws.City, error) {
	n.SetParams(Params{
		Q: c,
	})

	results := n.Search()

	if len(results) == 0 {
		return nil, errors.New("no results found for the provided city name")
	}

	result := results[0]

	city := nws.BuildCity(result.DisplayName, result.Lat, result.Lon)

	return &city, nil
}

func Init() *Nominatim {
	return &Nominatim{
		baseURL:   BaseURL,
		params:    Params{},
		userAgent: UserAgent,
	}
}

func Client() *Nominatim {
	return Init()
}
