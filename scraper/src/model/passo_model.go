package model

// PassoResponse is a representation of a Passo Events Response
type PassoResponse struct {
	TotalItemCount int `json:"totalItemCount"`
	ValueList      []struct {
		IsFollowing       bool   `json:"isFollowing"`
		ID                int    `json:"id"`
		Date              string `json:"date"`
		Name              string `json:"name"`
		HomePageImagePath string `json:"homePageImagePath"`
		HomeTeamName      string `json:"homeTeamName,omitempty"`
		EndDate           string `json:"endDate"`
		IsEntertainment   bool   `json:"isEntertainment"`
		EventType         int    `json:"eventType"`
		HashTagList       []struct {
			HashTagName string `json:"hashTagName"`
			RefEventID  int    `json:"refEventID"`
			HashTagID   int    `json:"hashTagId"`
		} `json:"hashTagList"`
		Priorty               int    `json:"priorty,omitempty"`
		VenueID               int    `json:"venueID"`
		VenueName             string `json:"venueName"`
		VenueSeoURL           string `json:"venueSeoUrl"`
		VenueSeoTitle         string `json:"venueSeoTitle"`
		VenueSeoDescription   string `json:"venueSeoDescription"`
		FirstRGB              string `json:"firstRGB"`
		LastRGB               string `json:"lastRGB"`
		SeoTitle              string `json:"seoTitle"`
		SeoDescription        string `json:"seoDescription"`
		SeoURL                string `json:"seoUrl"`
		HideDate              bool   `json:"hideDate"`
		ShowSecondSlider      bool   `json:"showSecondSlider"`
		SecondSliderSortOrder int    `json:"secondSliderSortOrder"`
		DontShowHomepage      bool   `json:"dontShowHomepage,omitempty"`
		ButtonText            string `json:"buttonText,omitempty"`
		IsShowWeb             bool   `json:"isShowWeb,omitempty"`
		IsShowMobile          bool   `json:"isShowMobile,omitempty"`
	} `json:"valueList"`
	IsError    bool `json:"isError"`
	ResultCode int  `json:"resultCode"`
}

type PassoRequestBody struct {
	CountRequired bool        `json:"CountRequired"`
	HastagID      interface{} `json:"HastagId"`
	CityID        interface{} `json:"CityId"`
	Date          interface{} `json:"date"`
	VenueID       interface{} `json:"VenueId"`
	StartDate     string      `json:"StartDate"`
	EndDate       string      `json:"EndDate"`
	LanguageID    int         `json:"LanguageId"`
	From          int         `json:"from"`
	Size          int         `json:"size"`
}
