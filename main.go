package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"mongo-poly/config"
	"mongo-poly/db"
	"mongo-poly/model"
	"mongo-poly/x_token"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	conf := config.Get()
	db.MongoInit(conf.Mongodb)

	go db.UpdateWeatherData()

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

	r.Post("/", meteoPointFinder)
	//r.Get("/meteosave", meteoSaver)
	r.Get("/period", weather)
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

	db.SetMeteoCode()

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

func meteoPointFinder(w http.ResponseWriter, r *http.Request) {

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
/*
func meteoSaver(w http.ResponseWriter, r *http.Request) {

*/
/*
	token := r.Header.Get("X-Token")

	if !x_token.CheckValid(token) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}*/ /*


	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	monthList := getMonthList(from, to)

	if len(monthList) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.SaveRangeWeather(monthList)
}
*/

//http://localhost:3000/weather?hash=testhash&code=testcode&from=2017-01&to=2018-01&type=testdoctype
func weather(w http.ResponseWriter, r *http.Request) {

	query, statusCode := model.GetRequestParams(r, model.WEATHER)

	if statusCode == 0 {
		if query.CodeID == "" {

			var err error
			query.CodeID, err = db.GetCodeIDByMD5(query.MD5)
			if err != nil {
				w.Write([]byte("ERR"))
				w.WriteHeader(statusCode)
				return
			}
		}
		db.GetWeatherResponse(w, *query)
	} else {
		w.Write([]byte("ERR"))
		w.WriteHeader(statusCode)
	}
}

//http://localhost:300/daydeg?hash=testhash&code=testcode&from=2017-01&to=2018-01
func dayDegree(w http.ResponseWriter, r *http.Request) {

	query, statusCode := model.GetRequestParams(r, model.DAYDEGREE)

	if statusCode == 0 {
		if query.CodeID == "" {

			var err error
			query.CodeID, err = db.GetCodeIDByMD5(query.MD5)
			if err != nil {
				w.Write([]byte("ERR"))
				w.WriteHeader(statusCode)
				return
			}
		}
		if dayDegreeList := db.GetMonthDayDegree(*query); len(dayDegreeList) > 0 {
			enc := json.NewEncoder(w)
			enc.Encode(dayDegreeList)
			w.Header().Set("Content-Type", "application/json")
		}
	} else {
		w.Write([]byte("ERR"))
		w.WriteHeader(statusCode)
	}
}
