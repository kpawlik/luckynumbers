package lotto

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"net/http"
)

func init() {
	//http.HandleFunc("/fix", fix)
}

func fix(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		key *datastore.Key
	)
	c := appengine.NewContext(r)
	results := make(map[*datastore.Key]*LottoResults)
	q := datastore.NewQuery("Results").Order("Date")
	for t := q.Run(c); ; {
		result := &LottoResults{}
		key, err = t.Next(result)
		if err == datastore.Done {
			err = nil
			break
		}
		if err != nil {
			return
		}
		results[key] = result
	}
	for key, rec := range results {
		d := rec.Date
		if len(d) > 8 {
			continue
		}
		date, _ := parseDate(d)
		nd := date.Format(`2006-01-02`)
		rec.Date = nd
		if _, err = datastore.Put(c, key, rec); err != nil {
			tmpls.ExecuteTemplate(w, "error", fmt.Sprintf("Storing error: %v \n", err))
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
