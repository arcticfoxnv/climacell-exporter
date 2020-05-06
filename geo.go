package main

import (
	"github.com/codingsince1985/geo-golang/openstreetmap"
)

func LookupCityCoords(city string) (float64, float64, error) {
	geocoder := openstreetmap.Geocoder()
	location, err := geocoder.Geocode(city)
	if err != nil {
		return 0, 0, err
	}

	return location.Lat, location.Lng, nil
}
