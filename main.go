package main

import (
	"encoding/json"
	//"gopkg.in/mgo.v2/bson"
	"github.com/go-chi/chi"
	"github.com/iizotop/baseweb/utils"
	"mongo-poly/config"
	"mongo-poly/db"
	"mongo-poly/model"
	"mongo-poly/x_token"
	"net/http"
	"strings"
	"time"
)

func main() {
	r := chi.NewRouter()

	conf := config.Get()
	db.MongoInit(conf.Mongodb)

	/*monthWeather := db.FindFieldWeather("test", 2017, 1)

	fmt.Println(monthWeather)

	for _, day := range monthWeather.Days {

		fmt.Printf("%#v\n", day)
		fmt.Println(day.Date)
	}*/

	/*from := time.Date(2016, 1,1,0,0,0,0, time.UTC)
	to := time.Date(2018, 5,1,0,0,0, 0,time.UTC)

	weatherList, err := db.GetPeriodWeather("test", from, to)

	db.WriteCSVWeather(weatherList)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else {
		fmt.Println(weatherList)
	}*/

	r.Post("/", kmlFinder)
	r.Get("/meteosave", meteoSaver)
	r.Get("/period", periodWeather)
	r.Get("/daydeg", dayDegree)
	http.ListenAndServe(":3000", r)

	//westLon, eastLon, northLat, southLat := 147.1397, 147.14, -38.1013, -38.1019

	//middleLon := westLon + (eastLon-westLon)/2
	//middleLat := southLat - (southLat-northLat)/2

	/*poly := model.MultiPolygon{
		Type: "MultiPolygon",
		Coordinates: [][][][2]float64{
			{
				{
					{westLon, northLat},
					{eastLon, northLat},
					{eastLon, southLat},
					{westLon, southLat},
					{westLon, northLat},
				},
			},
		},
	}
	geoKml := model.GeoKml{
		ID:       bson.ObjectId("Just_test_01"),
		MD5:      "test",
		Size:     0,
		Geometry: poly,
	}*/

	//db.InsertPolygon(geoKml)

	//db.SetMeteoCode()

	/*if err := db.GetWeather(2017, 12); err != nil {

		fmt.Println(err.Error())
	}*/

	/*point := model.Point{
		Type:        "Point",
		Coordinates: [2]float64{middleLon, middleLat},
	}*/
	//db.FindNearestStation(point)

	//db.FindPointInPolygon(point)

	/*	westLon, eastLon, northLat, southLat = 12.0, 18.0, 12.0, 18.0


		db.FindPolygonInPolygon(poly)*/
}

func kmlFinder(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("X-Token")

	if !x_token.CheckValid(token) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	kmlField := new(model.KmlQuery)

	json.NewDecoder(r.Body).Decode(kmlField)

	if len(kmlField.Point) > 0 {
		point := model.Point{
			Type:        "Point",
			Coordinates: kmlField.Point,
		}
		result := db.FindPointInPolygon(point)

		if result.Empty {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if res, err := json.Marshal(result); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.Write(res)
		}

	} else if len(kmlField.Polygon) > 0 {
		polygon := model.Polygon{
			Type: "Polygon",
			Coordinates: [][][2]float64{
				kmlField.Polygon,
			},
		}
		result := db.FindPolygonInPolygon(polygon)

		if result.Empty {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if res, err := json.Marshal(result); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.Write(res)
		}
	}

}

//	http://localhost:3000/meteosave?from=2016-01&to=2019-01
func meteoSaver(w http.ResponseWriter, r *http.Request) {

	/*
		token := r.Header.Get("X-Token")

		if !x_token.CheckValid(token) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}*/

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	monthList := getMonthList(from, to)

	if len(monthList) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.SaveRangeWeather(monthList)
}

func getMonthList(from, to string) (monthList []model.Month) {


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

			month := model.Month{
				Year:  toYear,
				Month: monthIndex,
			}
			monthList = append(monthList, month)
		}
		return
	case 1:
		for monthIndex := frMonth; monthIndex <= 12; monthIndex++ {

			month := model.Month{
				Year:  frYear,
				Month: monthIndex,
			}
			monthList = append(monthList, month)
		}
		for monthIndex := 1; monthIndex <= toMonth; monthIndex++ {

			month := model.Month{
				Year:  toYear,
				Month: monthIndex,
			}
			monthList = append(monthList, month)
		}
		return
	case 2:
		for monthIndex := frMonth; monthIndex <= 12; monthIndex++ {

			month := model.Month{
				Year:  frYear,
				Month: monthIndex,
			}
			monthList = append(monthList, month)
		}
		for monthIndex := 1; monthIndex <= 12; monthIndex++ {

			month := model.Month{
				Year:  frYear + 1,
				Month: monthIndex,
			}
			monthList = append(monthList, month)
		}
		for monthIndex := 1; monthIndex <= toMonth; monthIndex++ {

			month := model.Month{
				Year:  toYear,
				Month: monthIndex,
			}
			monthList = append(monthList, month)
		}
		return
	}
	return
}

//http://localhost:3000/period?hash=testhash&code=testcode&from=2017-01&to=2018-01&type=testdoctype
func periodWeather(w http.ResponseWriter, r *http.Request) {

	md5 := r.URL.Query().Get("hash")
	codeID := r.URL.Query().Get("code")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	docType := r.URL.Query().Get("type")

	if md5 == "" && codeID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if from == "" || to == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if docType == "csv" || docType == "xlsx" {
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	monthList := getMonthList(from, to)

	if len(monthList) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if codeID == "" {

		var err error

		codeID, err = db.GetCodeIDByMD5(md5)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	db.GetWeatherResponse(w, r,codeID, docType, monthList)
}

//http://localhost:300/daydeg?hash=testhash&code=testcode&from=2017-01&to=2018-01
func dayDegree(w http.ResponseWriter, r *http.Request) {

	md5 := r.URL.Query().Get("hash")
	codeID := r.URL.Query().Get("code")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if md5 == "" && codeID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if from == "" || to == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	monthList := getMonthList(from, to)

	if len(monthList) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if codeID == "" {

		var err error

		codeID, err = db.GetCodeIDByMD5(md5)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if dayDegreeList := db.GetMonthDayDegree(codeID, monthList); len(dayDegreeList) > 0 {
		enc := json.NewEncoder(w)
		enc.Encode(dayDegreeList)
		w.Header().Set("Content-Type", "application/json")
	}
}