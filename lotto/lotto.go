package lotto

import (
	"regexp"
	"strconv"
	"time"
)

const (
	LOTTO      = `<div class=\"glowna_wyniki_lotto\">[r\n\s.\w\<\>\"\=\/\-\+\!]*<div class=\"start\-wyniki_lotto\">`
	LOTTO_PLUS = `<div class=\"glowna_wyniki_lottoplus\">[r\n\s.\w\<\>\"\=\/\-\+\!]*<div class=\"start\-wyniki_mini\-lotto\">`
	LOTTO_DATE = `<div class=\"start-wyniki_lotto\">[r\n\s\.\w\<\>\"\=\/\-\+\!\:\;?\,]*</div>`
	URL        = "http://www.lotto.pl/"
)

var (
	reExp = map[string]*regexp.Regexp{}
)

func init() {
	strs := map[string]string{"lotto": LOTTO, "plus": LOTTO_PLUS, "date": LOTTO_DATE}
	for k, v := range strs {
		reExp[k] = regexp.MustCompile(v)
	}
}

type FetchResult struct {
	body  []byte
	Lotto []int
	Plus  []int
	Date  string
}

func GetResults(body []byte) (result *FetchResult) {
	result = &FetchResult{body: body}
	result.readLotto()
	result.readPlus()
	result.readDate()
	return result
}

func (r *FetchResult) readLotto() {
	re, _ := reExp["lotto"]
	r.Lotto = r.find(re)
}

func (r *FetchResult) readPlus() {
	re, _ := reExp["plus"]
	r.Plus = r.find(re)
}

func (r *FetchResult) readDate() {
	re, _ := reExp["date"]
	dText := re.Find(r.body)
	dateExp := regexp.MustCompile(`\d{2}-\d{2}-\d{2}`)
	r.Date = string(dateExp.Find(dText))
}

func (r *FetchResult) find(re *regexp.Regexp) (numbers []int) {
	numbers = make([]int, LOTTO_SIZE)
	lotto := re.Find(r.body)
	noExp := regexp.MustCompile(`\d+`)
	noBytes := noExp.FindAll(lotto, -1)
	for i := 0; i < len(noBytes); i++ {
		if no, err := strconv.Atoi(string(noBytes[i])); err == nil {
			numbers[i] = no
		}
	}
	return
}

func (r FetchResult) LottoStr() string {
	return numbersToString(r.Lotto[:])
}

func (r FetchResult) PlusStr() string {
	return numbersToString(r.Plus[:])
}

func (f FetchResult) GetDate() (date time.Time, err error) {
	date, err = parseDate(f.Date)
	return
}
