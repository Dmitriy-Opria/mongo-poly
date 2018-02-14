package main

import (
	"encoding/json"
	//"gopkg.in/mgo.v2/bson"
	"github.com/go-chi/chi"
	"github.com/iizotop/baseweb/utils"
	"mongo_kml/config"
	"mongo_kml/db"
	"mongo_kml/model"
	"mongo_kml/x_token"
	"net/http"
	"strings"
	"time"
	"fmt"
)

func main() {
	r := chi.NewRouter()


	conf := config.Get()
	db.MongoInit(conf.Mongodb)
	r.Post("/", kmlFinder)
	r.Get("/meteosave", meteoSaver)
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

	if from == "" || to == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fr := strings.Split(from, "-")
	t := strings.Split(to, "-")

	if len(fr) < 2 || len(t) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	frYear := utils.ToInt(fr[0])
	frMonth := utils.ToInt(fr[1])

	toYear := utils.ToInt(t[0])
	toMonth := utils.ToInt(t[1])

	if frYear == 0 || frMonth == 0 || toYear == 0 || toMonth == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if toYear-frYear < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if toYear-frYear == 0 {
		if toMonth-frMonth < 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	monthList := getMonthList(frYear, frMonth, toYear, toMonth)

	fmt.Println(monthList)

	db.SaveRangeWeather(monthList)
}

func getMonthList(fromYear, fromMonth, toYear, toMonth int) (monthList []model.Month) {

	yearMin := time.Now().Year() - 2
	yearMax := time.Now().Year()

	if fromYear < yearMin {
		fromYear = yearMin
	}

	if toYear > yearMax {
		toYear = yearMax
		toMonth = int(time.Now().Month())
	}

	switch toYear - fromYear {
	case 0:
		for monthIndex := fromMonth; monthIndex <= toMonth; monthIndex++ {

			month := model.Month{
				Year:  toYear,
				Month: monthIndex,
			}
			monthList = append(monthList, month)
		}
		return
	case 1:
		for monthIndex := fromMonth; monthIndex <= 12; monthIndex++ {

			month := model.Month{
				Year:  fromYear,
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
		for monthIndex := fromMonth; monthIndex <= 12; monthIndex++ {

			month := model.Month{
				Year:  fromYear,
				Month: monthIndex,
			}
			monthList = append(monthList, month)
		}
		for monthIndex := 1; monthIndex <= 12; monthIndex++ {

			month := model.Month{
				Year:  fromYear + 1,
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
