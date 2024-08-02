# Weather

A weather application made on top of the weather.gov and Nominatim/OpenStreetMap
APIs, written in Go (v1.22.4)

---

![Demo Webm](assets/demo.gif)

---

The demo video is generated using `vhs`. See [demo.tape](assets/demo.tape) for
the exact commands seen in the video.

## Usage

| Command   | Arguments | Flags    | Description                                            |
| --------- | --------- | -------- | ------------------------------------------------------ |
| `geocast` | none      | `-ip`    | geocode the IP address to get the city and state. If no IP is provided, the current IP is used. |
| `geocast` | `city`    | `-state` | geocode a city to get its latitude and longitude       |
| `geocast` | `code`    | none     | reverse geocode a latitude and longitude to get a city |

### Weather Commands

| Command   | Arguments | Flags    | Description                                            |
| --------- | --------- | -------- | ------------------------------------------------------ |
| `geocast` | `forecast`| none     | fetch the forecast                                     |
| `geocast` | `forecast`| `-city`  | fetch the forecast for a city                          |

## Data Sources

1. Geocoding

   - ipinfo
   - Nominatim/OpenStreetMap (osm)

2. Weather
   - weather.gov (US)

### Sample US Data

| City        | State | Latitude | Longitude |
| ----------- | ----- | -------- | --------- |
| Seattle     | WA    | 47.6062  | -122.3321 |
| Austin      | TX    | 30.2672  | -97.7431  |
| Cleveland   | OH    | 41.4993  | -81.6944  |
| Hartford    | CT    | 41.7658  | -72.6734  |
| Boston      | MA    | 42.3601  | -71.0589  |
| Los Angeles | CA    | 34.0522  | -118.2437 |
| Pittsburgh  | PA    | 40.4406  | -79.9959  |

## Updating Test Coverage

This project leverages a few tools executed by `coverage.py` to generate the
below coverage image, namely playwright and the built-in `go tool cover` command.

To ensure that you're able to generate the coverage image, you'll need to install
the `playwright` package and install the `chromium` browser.

Setup the virtual environment and install the dependencies:

```bash
virtualenv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

Then, install `playwright` with firefox and chromium:

```bash
playwright install
```
