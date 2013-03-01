package lotto

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"html/template"
	"net/http"
)

const (
	PIN       = "6666"
	LOTTO_URL = "http://www.lotto.pl"
)

var (
	tmpls *template.Template
)

func init() {
	var (
		err error
	)
	tmpls = template.New("tmpls").Funcs(template.FuncMap{"printNE": PrinNe, "printWin": PrintWin})
	if tmpls, err = tmpls.ParseFiles(
		"templates/header.html",
		"templates/style.css",
		"templates/index.html",
		"templates/change.html",
		"templates/error.html",
		"templates/check.html",
		"templates/stats.html",
	); err != nil {
		panic(err)
	}
	http.HandleFunc("/", check)
	http.HandleFunc("/change", change)
	http.HandleFunc("/store", store)
	http.HandleFunc("/stats", stats)
	http.HandleFunc("/check", check)
	http.HandleFunc("/favicon.ico", dump)
}

// Do nothing
func dump(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "")
}

func change(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	nbrs := GetLuckyNumbers(c)
	tmpls.ExecuteTemplate(w, "change", nbrs)
}

func store(w http.ResponseWriter, r *http.Request) {
	//check pseudo auth
	if r.FormValue("pin") != PIN {
		tmpls.ExecuteTemplate(w, "error", "Wrong PIN number")
		return
	}
	// validate numbers
	numbers := r.FormValue("numbers")
	if _, err := GetNumbers(numbers); err != nil {
		tmpls.ExecuteTemplate(w, "error", err)
		return
	}
	// get rest of fields
	expDate := r.FormValue("expiration")
	if _, err := parseDate(expDate); err != nil {
		tmpls.ExecuteTemplate(w, "error", fmt.Sprintf("Parsing date error: %v \n", err))
		return
	}
	startDate := r.FormValue("start")
	if _, err := parseDate(startDate); err != nil {
		tmpls.ExecuteTemplate(w, "error", fmt.Sprintf("Parsing date error: %v \n", err))
		return
	}
	plus := r.FormValue("plus") == "on"
	ln := &LuckyNumbers{Results: numbers, ExpDate: expDate, StartDate: startDate, Plus: plus}
	c := appengine.NewContext(r)
	_, err := datastore.Put(c, datastore.NewKey(c, "Numbers", "", 1, nil), ln)
	if err != nil {
		tmpls.ExecuteTemplate(w, "error", fmt.Sprintf("Storing error: %v \n", err))
		return
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func check(w http.ResponseWriter, r *http.Request) {

	var (
		key     *datastore.Key
		plusNo  string
		lres    *LottoResults
		body    []byte
		err     error
		plusWin []int
	)
	c := appengine.NewContext(r)
	// fetch remote site data
	if body, err = getPageBody(c, LOTTO_URL); err != nil {
		tmpls.ExecuteTemplate(w, "error", fmt.Sprintf("Error fetching page %s, %v", LOTTO_URL, err))
		return
	}
	// parse site and get data
	results := GetResults(body)
	// get lucky no from db
	lucky := GetLuckyNumbers(c)
	// check dates 
	expDate, _ := lucky.GetExpDate()
	startDate, _ := lucky.GetStartDate()
	lottoDate, _ := results.GetDate()
	strLottoDate := lottoDate.Format("2006-01-02")
	if lottoDate.After(expDate) || lottoDate.Before(startDate) {
		if lres, err = GetLastLottoResutls(c); err != nil {
			tmpls.ExecuteTemplate(w, "error", fmt.Sprintf("Get last result error: %v \n", err))
			return
		}
	} else {
		// check win
		luckyNbrs := lucky.GetNumbers()
		lottWin := sortedDiff(luckyNbrs, results.Lotto)
		// create new rec
		if lucky.Plus {
			plusNo = results.PlusStr()
			plusWin = sortedDiff(luckyNbrs, results.Plus)
		}
		lres = &LottoResults{Date: strLottoDate,
			LottoNo:  results.LottoStr(),
			PlusNo:   plusNo,
			LuckyNo:  lucky.Results,
			Plus:     lucky.Plus,
			LottoWin: numbersToString(lottWin),
			PlusWin:  numbersToString(plusWin)}
		// check if record already in db, 
		if key, err = GetLottoResult(c, strLottoDate); err != nil {
			tmpls.ExecuteTemplate(w, "error", fmt.Sprintf("Get results error: %v \n", err))
			return
		}
		// generate new key for insert
		if key == nil {
			key = datastore.NewIncompleteKey(c, "Results", nil)
		}
		// insert new rec or update existing
		if _, err = datastore.Put(c, key, lres); err != nil {
			tmpls.ExecuteTemplate(w, "error", fmt.Sprintf("Storing error: %v \n", err))
			return
		}
	}
	// display results
	tmpls.ExecuteTemplate(w, "check", lres)
}

func stats(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if results, err := GetLottoResults(c); err != nil {
		tmpls.ExecuteTemplate(w, "error", err)
		return
	} else {
		tmpls.ExecuteTemplate(w, "stats", results)
	}
}
