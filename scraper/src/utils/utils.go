package utils

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"
)

type Category int
type DatePeriod int
type City int
type Provider int

func (c Category) String() string {
	return [...]string{"MUSIC", "ART", "SPORT", "FAMILIY"}[c]
}

const (
	MUSIC Category = iota
	ART
)

func (dp DatePeriod) String() string {
	return [...]string{"thisweek", "next14days"}[dp]
}

const (
	THIS_WEEK DatePeriod = iota
	NEXT14_DAYS
)

func (c City) String() string {
	return [...]string{"İstanbul"}[c]
}

const (
	ISTANBUL City = iota
)

func (p Provider) String() string {
	return [...]string{"BILETIX", "PASSO", "KULTURIST"}[p]
}

const (
	BILETIX Provider = iota
	PASSO
	KULTURIST
)

// FailOnError returns a log based on given error and message
func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// GetQueryParameters returns a queryString based on given category and city
// In order to send a request to Biletix, it should be added to query params to url
// Example this uri returns the August events in Istanbul.
// https://www.biletix.com/solr/tr/select/?start=0&rows=1300&q=*:*&fq=start%3A%5B2022-08-01T00%3A00%3A00Z+TO+2022-08-31T00%3A00%3A00Z%2B1DAY%5D&wt=json&indent=true&facet=true&facet.field=category&facet.field=venuecode&facet.field=region&facet.field=subcategory&facet.mincount=1
func GetQueryParameters(category Category, city City) string {

	date1, date2 := GetDatesForBiletix()

	sb := strings.Builder{}
	sb.WriteString("?start=0&rows=1300&q=*:*&fq=")

	startDateString := fmt.Sprintf("start:[%s TO %s+1DAY]", date1, date2)
	sb.WriteString(url.QueryEscape(startDateString))
	sb.WriteString("&fq=city:")
	cityString := fmt.Sprintf(`"%s"`, city)
	sb.WriteString(url.QueryEscape(cityString))
	sb.WriteString("&wt=json&indent=true&facet=true&facet.field=category&facet.field=venuecode&facet.field=region&facet.field=subcategory&facet.mincount=1")
	return sb.String()
}

// GetDates returns  start date and end date which on current month
func GetDates() (string, string) {
	currentTimestamp := time.Now().UTC()
	currentYear, currentMonth, _ := currentTimestamp.Date()
	currentLocation, _ := time.LoadLocation("Europe/Istanbul")

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := time.Date(currentYear, currentMonth+1, 0, 0, 0, 0, 0, currentLocation)

	return firstOfMonth.Format("2006-01-02T00:00:00.000Z"), lastOfMonth.Format("2006-01-02T00:00:00.000Z")
}

// GetDatesForBiletix returns start date and end date which on current month for Biletix
func GetDatesForBiletix() (string, string) {
	currentTimestamp := time.Now().UTC()
	currentYear, currentMonth, _ := currentTimestamp.Date()
	currentLocation, err := time.LoadLocation("Europe/Istanbul")
	if err != nil {
		log.Errorf("location err: %v", err)
	}

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := time.Date(currentYear, currentMonth+1, 0, 0, 0, 0, 0, currentLocation)

	return firstOfMonth.Format("2006-01-02T00:00:00Z"), lastOfMonth.Format("2006-01-02T00:00:00Z")
}

func GetExJSON[T any](path string, t *T) T {
	jsonFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, t)

	defer jsonFile.Close()
	return *t
}
