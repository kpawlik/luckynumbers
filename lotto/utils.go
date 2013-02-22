package lotto

import (
	"appengine"
	"appengine/urlfetch"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	LOTTO_SIZE  = 6
	MIN_WIN_NO  = 2
	DATE_FORMAT = `02-01-2006`
)

type intSet struct {
	set map[int]bool
}

func NewIntSet() *intSet {
	return &intSet{make(map[int]bool)}
}

func NewIntSetFrom(ints []int) *intSet {
	s := NewIntSet()
	for i := 0; i < len(ints); i++ {
		s.set[ints[i]] = true
	}
	return s
}

func (s *intSet) Add(i int) {
	s.set[i] = true
}
func (s *intSet) Remove(i int) {
	if s.set[i] {
		delete(s.set, i)
	}
}

func (s intSet) Has(i int) bool {
	return s.set[i]
}

func (s intSet) Items() []int {
	result := make([]int, 0, len(s.set))
	for k, v := range s.set {
		if v {
			result = append(result, k)
		}
	}
	return result

}
func (s *intSet) Diff(other *intSet) *intSet {
	result := NewIntSet()
	for k, _ := range s.set {
		if other.Has(k) {
			result.Add(k)
		}
	}
	return result
}

func GetNumbers(str string) (nbrs []int, err error) {
	var (
		no int
	)
	nbrs = make([]int, LOTTO_SIZE)
	found := numbersRe.FindAllString(str, -1)
	if found == nil {
		err = errors.New(fmt.Sprintf("Error parsing expresion '%s' t numbers ", str))
		return
	}
	if foundLen := len(found); foundLen != LOTTO_SIZE {
		err = errors.New(fmt.Sprintf("Bad slice length after parsing %s != %s", foundLen, LOTTO_SIZE))
		return
	}
	for i := 0; i < len(found); i++ {
		if no, err = strconv.Atoi(strings.TrimSpace(found[i])); err != nil {
			return
		}
		nbrs[i] = no
	}
	return
}

func numbersToString(ints []int) string {
	strs := make([]string, LOTTO_SIZE)
	for i := 0; i < len(ints); i++ {
		strs = append(strs, strconv.FormatInt(int64(ints[i]), 10))
	}
	return strings.Join(strs, " ")
}

func parseDate(strDate string) (date time.Time, err error) {
	dd := strings.Split(strDate, "-")
	if ddLen := len(dd); ddLen != 3 || ddLen < 3 {
		err = errors.New(fmt.Sprintf("Bad data string '%s'", strDate))
		return
	}
	y := dd[2]
	if len(y) == 2 {
		y = fmt.Sprintf("20%s", y)
	}
	m, d := dd[1], dd[0]
	dateS := fmt.Sprintf("%s-%s-%s", d, m, y)
	date, err = time.Parse(DATE_FORMAT, dateS)
	return
}

func getPageBody(c appengine.Context, url string) (body []byte, err error) {
	var (
		resp *http.Response
	)
	client := urlfetch.Client(c)
	if resp, err = client.Get(url); err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func sortedDiff(one []int, other []int) []int {
	oneSet := NewIntSetFrom(one)
	otherSet := NewIntSetFrom(other)
	res := oneSet.Diff(otherSet).Items()
	sort.Ints(res)
	return res
}

func PrinNe(args ...interface{}) string {
	str, _ := args[0].(string)
	if strings.TrimSpace(str) == "" {
		return "-"
	}
	return str
}
func PrintWin(args ...interface{}) string {
	str, _ := args[0].(string)
	fmt.Println(str, len(strings.Split(strings.TrimSpace(str), " ")))
	if len(strings.Split(strings.TrimSpace(str), " ")) > MIN_WIN_NO {
		return "win"
	}
	return ""
}
