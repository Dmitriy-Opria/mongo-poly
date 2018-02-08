package main

import (
	"mongo_kml/model"
	"mongo_kml/db"
	"mongo_kml/config"
)

func main(){
	westLon, eastLon, northLat, southLat := 10.0, 20.0, 10.0, 20.0

	//middleLon := westLon + (eastLon - westLon)/2
	//middleLat := southLat - (southLat - northLat)/2

	/*poly := model.MultiPoint{
		Type: "MultiPolygon",
		Coordinates: [][2]float64{
			{westLon, northLat},
			{eastLon, northLat},
			{eastLon, southLat},
			{westLon, southLat},
			{westLon, northLat},
		},
	}*/

	conf := config.Get()
	db.MongoInit(conf.Mongodb)

	//db.InsertPolygon(poly)
	/*point := model.Point{
		Type: "Point",
		Coordinates: [2]float64{middleLon, middleLat},
	}
*/
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
	db.FindPointInPolygon(poly)

}

