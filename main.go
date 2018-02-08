package main

import (
	"mongo_kml/config"
	"mongo_kml/db"
	"mongo_kml/model"
	//"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
)

func main() {

	r := chi.NewRouter()

	r.Post("/", kmlFinder)
	http.ListenAndServe(":3000", r)

	westLon, eastLon, northLat, southLat := 10.0, 20.0, 10.0, 20.0

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
		ID:  bson.ObjectId("Just_test_01"),
		MD5:      "test",
		Size:     0,
		Geometry: poly,
	}*/

	conf := config.Get()
	db.MongoInit(conf.Mongodb)

	//db.InsertPolygon(geoKml)

	/*point := model.Point{
		Type:        "Point",
		Coordinates: [2]float64{middleLon, middleLat},
	}*/
	//db.FindPointInPolygon(point)

	westLon, eastLon, northLat, southLat = 12.0, 18.0, 12.0, 18.0

	poly := model.Polygon{
		Type: "Polygon",
		Coordinates: [][][2]float64{
			{
				{westLon, northLat},
				{eastLon, northLat},
				{eastLon, southLat},
				{westLon, southLat},
				{westLon, northLat},
			},
		},
	}
	db.FindPolygonInPolygon(poly)
}

func kmlFinder(w http.ResponseWriter, r *http.Request) {

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
