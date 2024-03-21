package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type GooglePhotoJsonPayload struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	ImageViews   string `json:"imageViews"`
	CreationTime struct {
		Timestamp string `json:"timestamp"`
		Formatted string `json:"formatted"`
	} `json:"creationTime"`
	PhotoTakenTime struct {
		Timestamp string `json:"timestamp"`
		Formatted string `json:"formatted"`
	} `json:"photoTakenTime"`
	GeoData struct {
		Latitude      float64 `json:"latitude"`
		Longitude     float64 `json:"longitude"`
		Altitude      float64 `json:"altitude"`
		LatitudeSpan  float64 `json:"latitudeSpan"`
		LongitudeSpan float64 `json:"longitudeSpan"`
	} `json:"geoData"`
	GeoDataExif struct {
		Latitude      float64 `json:"latitude"`
		Longitude     float64 `json:"longitude"`
		Altitude      float64 `json:"altitude"`
		LatitudeSpan  float64 `json:"latitudeSpan"`
		LongitudeSpan float64 `json:"longitudeSpan"`
	} `json:"geoDataExif"`
	URL                string `json:"url"`
}

type GooglePhotoJsonFile struct {
	Path *string
	Data *GooglePhotoJsonPayload
}

func (g *GooglePhotoJsonFile) GetPayload() error {
	dat, err := os.ReadFile(*g.Path)
	if err != nil {
		return fmt.Errorf("unable to read file: %w", err)
	}

	g.Data = &GooglePhotoJsonPayload{}
	err = json.Unmarshal(dat, g.Data)
	if err != nil {
		return fmt.Errorf("unable to unmarshal json for file %s: %w", *g.Path, err)
	}

	return nil
}