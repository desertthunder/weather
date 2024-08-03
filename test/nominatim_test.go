package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	osm "github.com/desertthunder/weather/internal/nominatim" // osm is an alias for nominatim (openstreetmap)
)

func TestParams(t *testing.T) {
	t.Run("String", func(t *testing.T) {

		t.Run("empty", func(t *testing.T) {
			p := osm.Params{}

			s := p.String()

			if s != "" {
				t.Errorf("Expected empty string, got %s", s)
			}
		})

		t.Run("q", func(t *testing.T) {
			p := osm.Params{
				Q: "Austin",
			}

			got := p.String()
			want := "q=Austin&format=jsonv2&limit=25"

			if got != want {
				t.Errorf("Expected %s, got %s", want, got)
			}
		})
	})
}

func TestNominatimClient(t *testing.T) {
	t.Run("Setters", func(t *testing.T) {
		client := osm.Client()
		t.Run("SetURL", func(t *testing.T) {
			want := "https://nominatim.openstreetmap.org"
			client.SetURL(want)

			got := client.BaseURL()
			if got != want {
				t.Errorf("Expected base URL to be %s, got %s", want, got)
			}
		})

		t.Run("SetParams", func(t *testing.T) {
			want := "q=Austin"
			wantParams := osm.Params{
				Q: "Austin",
			}

			client.SetParams(wantParams)
			got := client.GetParams()

			if !strings.Contains(got.String(), want) {
				t.Errorf("Expected params to include %s, got %s", want, got.String())
			}
		})

		t.Run("SetUserAgent", func(t *testing.T) {
			want := "new_user_agent"

			client.SetUserAgent("new_user_agent")

			got := client.UserAgent()

			if got != want {
				t.Errorf("Expected user agent to be %s, got %s", want, got)
			}
		})
	})

	t.Run("Search", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`[{"place_id": 305086781, "licence": "Data © OpenStreetMap contributors, ODbL 1.0. http://osm.org/copyright", "osm_type": "relation", "osm_id": 113314, "lat": "30.2711286", "lon": "-97.7436995", "category": "boundary", "type": "administrative", "place_rank": 16, "importance": 0.6494265717866669, "addresstype": "city", "name": "Austin", "display_name": "Austin, Travis County, Texas, United States", "boundingbox": ["30.0985133", "-97.9367663", "30.5166255", "-97.5605288"]}]`))
		}))

		defer server.Close()

		client := osm.Client()

		client.SetURL(server.URL)
		p := osm.Params{
			Q: "Austin",
		}

		client.SetParams(p)

		results := client.Search()

		if len(results) < 1 {
			fmt.Println(results)
			t.Errorf("Expected 1 result, got %d", len(results))
		}
	})

	t.Run("Geocode", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`[{"place_id": 312908827, "licence": "Data © OpenStreetMap contributors, ODbL 1.0. http://osm.org/copyright", "osm_type": "relation", "osm_id": 237385, "lat": "47.6038321", "lon": "-122.330062", "category": "boundary", "type": "administrative", "place_rank": 16, "importance": 0.6729791735643788, "addresstype": "city", "name": "Seattle", "display_name": "Seattle, King County, Washington, United States", "boundingbox": ["47.4810022", "47.7341354", "-122.4596960", "-122.2244330"]}]`))
		}))

		defer server.Close()

		client := osm.Client()

		client.SetURL(server.URL)

		t.Run("ByPoint", func(t *testing.T) {
			lat := 47.6062
			lon := -122.3321

			city, err := client.GeocodeByPoint(lat, lon)

			if err != nil {
				t.Errorf("Expected no error, got %s", err.Error())
			}

			if !strings.Contains(city.Name, "Seattle") {
				t.Errorf("Expected city name to be Seattle, got %s", city.Name)
			}

		})

		t.Run("ByCity", func(t *testing.T) {
			city, err := client.GeocodeByCity("Seattle")

			if err != nil {
				t.Errorf("Expected no error, got %s", err.Error())
			}

			if !strings.Contains(city.Name, "Seattle") {
				t.Errorf("Expected city name to be Seattle, got %s", city.Name)
			}
		})
	})
}
