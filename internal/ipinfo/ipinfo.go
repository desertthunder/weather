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

type IPInfoClient struct {
	BaseURL string
	Token   string
}

func (c *IPInfoClient) SetToken(token string) {
	c.Token = token
}

func (i *IPInfoClient) SetURL(url string) {
	i.BaseURL = url
}

func (i *IPInfoResponse) Point() (float64, float64) {
	coords := strings.Split(i.Location, ",")

	lat, _ := strconv.ParseFloat(coords[0], 64)
	lon, _ := strconv.ParseFloat(coords[1], 64)

	return lat, lon
}

func NewIPInfoClient(token string) *IPInfoClient {
	return &IPInfoClient{Token: token, BaseURL: baseURL}
}

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
	if *ipaddr != "" {
		valid = utils.ValidateIPAddress(*ipaddr)

		uri.Path = fmt.Sprintf("/%s", *ipaddr)
	}

	if !valid {
		err := errors.New("invalid IP address")
		fmt.Printf("Invalid IP address: %s\n", *ipaddr)
		return ipinfo, err
	}

	uri.Query().Add("token", c.Token)

	rsp, err := http.Get(uri.String())

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

	err = json.Unmarshal(data, &ipinfo)

	if err != nil {
		fmt.Printf("Failed to unmarshal response body: %s\n", err.Error())
	}

	return ipinfo, err
}
