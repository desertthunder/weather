package nws

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/log"
)

const baseURL string = "https://api.weather.gov"

type WeatherClient struct {
	baseURL string
	Log     *log.Logger
	logger  *log.Logger
}

func (c *WeatherClient) SetURL(url string) {
	c.baseURL = url
}

func (c *WeatherClient) BaseURL() string {
	return c.baseURL
}

func (c *WeatherClient) SetLogger(logger *log.Logger) {
	c.logger = logger
	c.Log = logger
}

func (c *WeatherClient) GetWeather(city City) (*ForecastAPIResponse, error) {
	uri := city.OfficeURL()

	rsp, err := http.Get(uri)

	if err != nil {
		fmt.Printf("Request to %s failed with error: %s\n", uri, err.Error())

		c.logger.Error(fmt.Sprintf("Request to %s failed with error: %s", uri, err.Error()))

		return nil, err
	}

	defer rsp.Body.Close()

	data, err := io.ReadAll(rsp.Body)

	if err != nil {
		c.logger.Error(fmt.Sprintf("Failed to read response body: %s", err.Error()))
	}

	office := ForecastOfficeAPIResponse{}

	err = json.Unmarshal(data, &office)

	forecastURL := office.ForecastURL()

	c.logger.Debug(fmt.Sprintf("Found: %s", forecastURL))

	if err != nil {
		c.logger.Error(fmt.Sprintf("Failed to unmarshal response body: %s", err.Error()))

		return nil, err
	}

	rsp, err = http.Get(forecastURL)

	if err != nil {
		c.logger.Error(fmt.Sprintf("Request to %s failed with error: %s", forecastURL, err.Error()))

		return nil, err
	}

	defer rsp.Body.Close()

	data, err = io.ReadAll(rsp.Body)

	if err != nil {
		c.logger.Error(fmt.Sprintf("Failed to read response body: %s", err.Error()))

		return nil, err
	}

	fc := ForecastAPIResponse{}

	err = json.Unmarshal(data, &fc)

	if err != nil {
		c.logger.Error(fmt.Sprintf("Failed to unmarshal response body: %s", err.Error()))

		return nil, err
	}

	return &fc, nil
}

func NewWeatherClient() *WeatherClient {
	return &WeatherClient{baseURL: baseURL}
}
