# Nominatim API Documentation

Nominatim is an open-source search engine used for geocoding and reverse geocoding.
It provides access to OpenStreetMap data and can be used to find location information based on an address or coordinates.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Endpoint Overview](#endpoint-overview)
3. [Usage](#usage)
   - [Forward Geocoding](#forward-geocoding)
   - [Reverse Geocoding](#reverse-geocoding)
   - [Additional Parameters](#additional-parameters)
4. [Rate Limits and Usage Policy](#rate-limits-and-usage-policy)
5. [Examples](#examples)
6. [Common Errors](#common-errors)
7. [Further Resources](#further-resources)

## Getting Started

To use the Nominatim API, you need to make HTTP GET requests to the appropriate endpoint.
No API key is required for basic usage, but be sure to check the usage policy for rate limits and usage restrictions.

## Endpoint Overview

- **Forward Geocoding**: Converts an address into geographical coordinates.
    - Endpoint: `https://nominatim.openstreetmap.org/search`
- **Reverse Geocoding**: Converts geographical coordinates into an address.
    - Endpoint: `https://nominatim.openstreetmap.org/reverse`

## Usage

### Forward Geocoding

To perform forward geocoding, use the `/search` endpoint with the following parameters:

- `q`: The query string (e.g., address, place name)
- `format`: The format of the output (`json`, `xml`, `html`, etc.)
- `addressdetails`: (Optional) Include a breakdown of the address details (0 or 1)
- `limit`: (Optional) Limit the number of returned results
- `countrycodes`: (Optional) Restrict search results to specific countries (comma-separated list of ISO 3166-1 alpha2 codes)

**Example Request:**

```http
GET https://nominatim.openstreetmap.org/search?q=1600+Amphitheatre+Parkway,+Mountain+View&format=json&limit=1
```

### Reverse Geocoding

To perform reverse geocoding, use the `/reverse` endpoint with the following parameters:

- `lat`: Latitude of the location
- `lon`: Longitude of the location
- `format`: The format of the output (`json`, `xml`, `html`, etc.)
- `zoom`: (Optional) Level of detail in the response
- `addressdetails`: (Optional) Include a breakdown of the address details (0 or 1)

**Example Request:**

```http
GET https://nominatim.openstreetmap.org/reverse?lat=37.423021&lon=-122.083739&format=json
```

### Additional Parameters

- `accept-language`: (Optional) Preferred language for the response
- `namedetails`: (Optional) Include additional details about the place (0 or 1)
- `extratags`: (Optional) Include additional tags about the place (0 or 1)

## Rate Limits and Usage Policy

Nominatim has rate limits to prevent abuse and ensure fair usage. The exact limits may vary, so refer to the [official documentation](https://nominatim.org/release-docs/latest/api/Overview/) for up-to-date information.

## Examples

### Forward Geocoding Example

```http
GET https://nominatim.openstreetmap.org/search?q=Empire+State+Building&format=json&addressdetails=1
```

### Reverse Geocoding Example

```http
GET https://nominatim.openstreetmap.org/reverse?lat=40.748817&lon=-73.985428&format=json&zoom=18&addressdetails=1
```

## Common Errors

- **403 Forbidden**: The server rejected the request. This may occur if the request rate limit is exceeded.
- **400 Bad Request**: The request was malformed. Check for missing or incorrect parameters.
- **500 Internal Server Error**: A problem occurred on the server side. This is usually temporary.

## Further Resources

- [Nominatim API Documentation](https://nominatim.org/release-docs/latest/api/Overview/)
- [OpenStreetMap Wiki](https://wiki.openstreetmap.org/wiki/Nominatim)

This basic documentation provides an overview of the Nominatim API, its endpoints, usage examples, and common errors. Be sure to refer to the [official Nominatim documentation](https://nominatim.org/release-docs/latest/api/Overview/) for more detailed information and updates.
