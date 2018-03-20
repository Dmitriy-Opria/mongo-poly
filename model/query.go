package model

import (
	"github.com/iizotop/flurosat/utils"
	"net/http"
	"strings"
	"time"
)

type KmlQuery struct {
	Point   [2]float64   `json:"point"`
	Polygon [][2]float64 `json:"polygon"`
}

type KmlAnswer struct {
	MD5   []string `json:"md5"`
	Empty bool     `json:"-"`
}

type WeatherQuery struct {
	MD5       string  `json:"md5"`
	CodeID    string  `json:"code_id"`
	MonthList []Month `json:"month_list"`
	DocType   string  `json:"doc_type"`
	QueryType int     `json:"query_type"`
}

const (
	DAYDEGREE = 1
	WEATHER   = 2
)

func GetRequestParams(r *http.Request, queryType int) (query *WeatherQuery, statusCode int) {

	query = new(WeatherQuery)
	query.QueryType = queryType

	query.MD5 = r.URL.Query().Get("hash")
	query.CodeID = r.URL.Query().Get("code")

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	query.DocType = r.URL.Query().Get("type")

	if query.MD5 == "" && query.CodeID == "" {
		statusCode = http.StatusBadRequest
		return
	}

	if from == "" || to == "" {
		statusCode = http.StatusBadRequest
		return
	}

	monthList := GetMonthList(from, to)

	if len(monthList) < 1 {
		statusCode = http.StatusBadRequest
		return
	}

	if queryType == WEATHER {
		if query.DocType != "csv" && query.DocType != "xlsx" && query.DocType != "json" {
			statusCode = http.StatusBadRequest
			return
		}
	}

	return
}

func GetMonthList(from, to string) (monthList []Month) {

	fr := strings.Split(from, "-")
	t := strings.Split(to, "-")

	if len(fr) < 2 || len(t) < 2 {
		return
	}

	frYear := utils.ToInt(fr[0])
	frMonth := utils.ToInt(fr[1])

	toYear := utils.ToInt(t[0])
	toMonth := utils.ToInt(t[1])

	if frYear == 0 || frMonth == 0 || toYear == 0 || toMonth == 0 {
		return
	}

	if toYear-frYear < 0 {
		return
	}

	if toYear-frYear == 0 {
		if toMonth-frMonth < 0 {
			return
		}
	}

	yearMin := time.Now().Year() - 2
	yearMax := time.Now().Year()

	if frYear < yearMin {
		frYear = yearMin
	}

	if toYear > yearMax {
		toYear = yearMax
		toMonth = int(time.Now().Month())
	}

	switch toYear - frYear {
	case 0:
		for monthIndex := frMonth; monthIndex <= toMonth; monthIndex++ {

			month := Month{
				Year:  toYear,
				Month: monthIndex,
			}
			monthList = append(monthList, month)
		}
		return
	case 1:
		for monthIndex := frMonth; monthIndex <= 12; monthIndex++ {

			month := Month{
				Year:  frYear,
				Month: monthIndex,
			}
			monthList = append(monthList, month)
		}
		for monthIndex := 1; monthIndex <= toMonth; monthIndex++ {

			month := Month{
				Year:  toYear,
				Month: monthIndex,
			}
			monthList = append(monthList, month)
		}
		return
	case 2:
		for monthIndex := frMonth; monthIndex <= 12; monthIndex++ {

			month := Month{
				Year:  frYear,
				Month: monthIndex,
			}
			monthList = append(monthList, month)
		}
		for monthIndex := 1; monthIndex <= 12; monthIndex++ {

			month := Month{
				Year:  frYear + 1,
				Month: monthIndex,
			}
			monthList = append(monthList, month)
		}
		for monthIndex := 1; monthIndex <= toMonth; monthIndex++ {

			month := Month{
				Year:  toYear,
				Month: monthIndex,
			}
			monthList = append(monthList, month)
		}
		return
	}
	return
}
