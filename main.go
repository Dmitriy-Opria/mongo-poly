package main

import (
	"encoding/json"
	//"github.com/go-chi/chi"
	//"gopkg.in/mgo.v2/bson"
	"fmt"
	"mongo_kml/config"
	"mongo_kml/db"
	"mongo_kml/model"
	"mongo_kml/x_token"
	"net/http"
)

/*
func main() {

	point := model.Point{
		Type:        "Point",
		Coordinates: [2]float64{147.1398, -38.1016},
	}

	conf := config.Get()
	db.MongoInit(conf.Mongodb)

	db.FindNearestStation(point)
}
*/

/*
func main() {
	fileName := "./stations_codes.csv"
	meteolist := db.ReadCSV(fileName)


	conf := config.Get()
	db.MongoInit(conf.Mongodb)


	if err := db.InsertMeteo(meteolist); err != nil {
		fmt.Printf("Can`t insert meteo error: %s", err.Error())
	}
}*/

func main() {
	/*
		r := chi.NewRouter()

		r.Post("/", kmlFinder)
		http.ListenAndServe(":3000", r)*/

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

	conf := config.Get()
	db.MongoInit(conf.Mongodb)

	//db.InsertPolygon(geoKml)

	//db.SetMeteoCode()

	fileInfo := db.ReadWeaterFile("./IDCJDW3001.201712.csv")

	for _, v := range fileInfo.Days {

		fmt.Printf("%#v\n", v)
	}

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
