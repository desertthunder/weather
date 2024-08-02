package ipinfo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/desertthunder/weather/internal/nws"
	"github.com/desertthunder/weather/internal/utils"
)

const (
	baseURL string = "https://ipinfo.io"
)

type IPInfoResponse struct {
	IP           string `json:"ip"`
	Hostname     string `json:"hostname"`
	City         string `json:"city"`
	Region       string `json:"region"`
	Country      string `json:"country"`
	Location     string `json:"loc"`
	Organization string `json:"org"`
	Postal       string `json:"postal"`
	Timezone     string `json:"timezone"`
}

// Client for the https://ipinfo.io API.
type IPInfoClient struct {
	// Base URL for the API. Defaults to https://ipinfo.io but
	// can be overridden for tests.
	BaseURL string
	// Token for the API. Required for requests.
	Token string
	// Logger for the client.
	Log *log.Logger
}

// BuildCity converts an IPInfoResponse to a City object.
func (r IPInfoResponse) BuildCity() nws.City {
	lat, lon := r.Point()

	return nws.City{
		Name: r.City,
		Lat:  lat,
		Long: lon,
	}
}

func (c *IPInfoClient) SetToken(token string) {
	c.Token = token
}

func (i *IPInfoClient) SetURL(url string) {
	i.BaseURL = url
}

// Setter for the logger.
func (i *IPInfoClient) SetLogger(logger *log.Logger) {
	i.Log = logger
}

func (i *IPInfoResponse) Point() (float64, float64) {
	coords := strings.Split(i.Location, ",")

	lat, _ := strconv.ParseFloat(coords[0], 64)
	lon, _ := strconv.ParseFloat(coords[1], 64)

	return lat, lon
}

// Call the IPInfo API to geolocate a given IP address.
//
// If no IP address is provided, no param is passed to the API, which means the
// client's IP address is used.
func (c *IPInfoClient) Geolocate(ipaddr *string) (IPInfoResponse, error) {
	ipinfo := IPInfoResponse{}

	if c.Token == "" {
		err := errors.New("IPInfo token is required")

		return ipinfo, err
	}

	uri, err := url.ParseRequestURI(c.BaseURL)

	if err != nil {
		fmt.Printf("Failed to parse URL: %s\n", err.Error())

		return ipinfo, err
	}

	valid := true

	if ipaddr != nil {
		valid = utils.ValidateIPAddress(*ipaddr)

		uri.Path = fmt.Sprintf("/%s", *ipaddr)
	}

	if !valid {
		err := errors.New("invalid IP address")
		fmt.Printf("Invalid IP address: %s\n", *ipaddr)
		return ipinfo, err
	}

	query := uri.Query()
	query.Add("token", c.Token)

	withQuery := fmt.Sprintf("%s?%s", uri.String(), query.Encode())

	rsp, err := http.Get(withQuery)

	if err != nil {
		fmt.Printf("Request to %s failed with error: %s\n", uri, err.Error())

		return ipinfo, err
	}

	defer rsp.Body.Close()

	data, err := io.ReadAll(rsp.Body)

	if err != nil {
		fmt.Printf("Failed to parse response body: %s\n", err.Error())

		return ipinfo, err
	}

	err = ipinfo.Validate(data)

	return ipinfo, err
}

func (r *IPInfoResponse) Validate(data []byte) error {
	s := utils.GetRawJSON(data)

	if strings.Contains(s, "bogon") {
		return errors.New("IP address is private (may be local)")
	} else {
		err := json.Unmarshal(data, &r)

		return err
	}
}

func NewIPInfoClient(token string) *IPInfoClient {
	return &IPInfoClient{Token: token, BaseURL: baseURL}
}
