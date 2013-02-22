package lotto

import (
	"appengine"
	"appengine/datastore"
	"regexp"
	"time"
)

const (
	NUMBERS_RE = `(\d{0,2})\s*`
)

var (
	numbersRe *regexp.Regexp = regexp.MustCompile(NUMBERS_RE)
)

// Nubers to check 
type LuckyNumbers struct {
	Results   string // space separated integers
	ExpDate   string
	StartDate string
	Plus      bool
}

// Return entity record with specify key. Should be only one record in db.
func GetLuckyNumbers(c appengine.Context) (nbrs *LuckyNumbers) {
	nbrs = &LuckyNumbers{}
	datastore.Get(c, datastore.NewKey(c, "Numbers", "", 1, nil), nbrs)
	return
}

// Return numbers as slice of ints
func (l LuckyNumbers) GetNumbers() (nbrs []int) {
	nbrs, _ = GetNumbers(l.Results)
	return
}

// parse and return number valid end date
func (l LuckyNumbers) GetExpDate() (date time.Time, err error) {
	date, err = parseDate(l.ExpDate)
	return
}

// parse and return number valid start date
func (l LuckyNumbers) GetStartDate() (date time.Time, err error) {
	date, err = parseDate(l.StartDate)
	return
}

// Results
type LottoResults struct {
	Plus     bool
	LottoNo  string
	LottoWin string
	PlusNo   string
	PlusWin  string
	LuckyNo  string
	Date     string
}

// If exists return Key of entity record from collection Results with specify date.
func GetLottoResult(c appengine.Context, date string) (recordKey *datastore.Key, err error) {
	var key *datastore.Key
	q := datastore.NewQuery("Results").
		KeysOnly().Limit(1).
		Filter("Date =", date)
	for t := q.Run(c); ; {
		key, err = t.Next(nil)
		if err == datastore.Done {
			err = nil
			break
		}
		recordKey = key
		if err != nil {
			return
		}
	}
	return
}

// Get entite record with last Date from Result collection.
func GetLastLottoResutls(c appengine.Context) (result *LottoResults, err error) {
	result = &LottoResults{}
	q := datastore.NewQuery("Results").
		Limit(1).
		Order("-Date")
	for t := q.Run(c); ; {
		_, err = t.Next(result)
		if err == datastore.Done {
			err = nil
			break
		}
		if err != nil {
			return
		}
	}
	return
}

// Return slice of all entities from Results collection.
// Results are sorted by Date
func GetLottoResults(c appengine.Context) (results []*LottoResults, err error) {
	q := datastore.NewQuery("Results").Order("Date")
	for t := q.Run(c); ; {
		result := &LottoResults{}
		_, err = t.Next(result)
		if err == datastore.Done {
			err = nil
			break
		}
		if err != nil {
			return
		}
		results = append(results, result)
	}
	return
}
