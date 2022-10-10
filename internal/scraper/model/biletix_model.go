package model

import "time"

// BiletixResponse is a representation of a Biletix Events Response
type BiletixResponse struct {
	Response struct {
		NumFound int `json:"numFound"`
		Start    int `json:"start"`
		Docs     []struct {
			Region          string    `json:"region"`
			Vote            int       `json:"vote"`
			Citycount       string    `json:"citycount"`
			Status          string    `json:"status"`
			LinkURL         string    `json:"link_url"`
			ImageURL        string    `json:"image_url"`
			Liveevent       bool      `json:"liveevent"`
			Venue           string    `json:"venue"`
			Svenue          string    `json:"svenue"`
			Type            string    `json:"type"`
			City            string    `json:"city"`
			ID              string    `json:"id"`
			Venuecount      string    `json:"venuecount"`
			Category        string    `json:"category"`
			Start           time.Time `json:"start"`
			Description     string    `json:"description"`
			Subcategory     string    `json:"subcategory"`
			Detail          []string  `json:"detail"`
			Event           []string  `json:"event"`
			EventSuggestion string    `json:"EventSuggestion,omitempty"`
			Name            string    `json:"name"`
			Sname           string    `json:"sname"`
			NameWs          string    `json:"name_ws"`
			Group           string    `json:"group"`
			Parent          string    `json:"parent"`
			Venuecode       string    `json:"venuecode"`
			National        bool      `json:"national"`
			End             time.Time `json:"end"`
			Version         int64     `json:"_version_"`
			Eventcount      int       `json:"eventcount"`
			CityID          string    `json:"city_id"`
			Timestamp       time.Time `json:"timestamp"`
			Artist          []string  `json:"artist"`
			GroupSuggestion string    `json:"GroupSuggestion,omitempty"`
		} `json:"docs"`
	} `json:"response"`
}
