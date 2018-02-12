package db

import (
	"encoding/csv"
	"fmt"
	"github.com/iizotop/baseweb/utils"
	"gopkg.in/mgo.v2/bson"
	"io"
	"mongo_kml/model"
	"os"
)

func ReadCSV(fileName string) (meteoUnitList []model.MeteoUnit) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0644)

	if err != nil {
		fmt.Printf("Can`t open file %s, error: %s\n", fileName, err.Error())
		return
	}
	defer file.Close()

	r := csv.NewReader(file)

	meteoUnitList = make([]model.MeteoUnit, 0, 124)

	first := true

	for {
		record, err := r.Read()

		if first {
			first = false
			continue
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if len(record) < 14 {
			continue
		}

		meteoUnit := model.MeteoUnit{}

		meteoUnit.Location = record[0]
		meteoUnit.Station = record[1]
		meteoUnit.CodeID = record[2]
		meteoUnit.Site = record[3]
		meteoUnit.Dist = record[4]
		meteoUnit.Start = utils.ToInt(record[5])
		meteoUnit.End = utils.ToInt(record[6])

		var lon, lat float64
		lat = utils.ToFloat64(record[7])
		lon = utils.ToFloat64(record[8])

		meteoUnit.Point = model.Point{
			Type:        "Point",
			Coordinates: [2]float64{lon, lat},
		}
		meteoUnit.Source = record[9]
		meteoUnit.STA = record[10]
		meteoUnit.Height = utils.ToFloat64(record[11])
		meteoUnit.Bar = utils.ToFloat64(record[12])
		meteoUnit.WMO = record[13]

		meteoUnitList = append(meteoUnitList, meteoUnit)
	}
	return
}

func InsertMeteo(meteoList []model.MeteoUnit) (err error) {

	db, def := getDatabase()
	defer def()

	for _, meteo := range meteoList {

		err = db.C("meteoStations").Insert(meteo)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	return
}

func FindNearestStation(point model.Point) (meteo *model.MeteoUnit) {

	fmt.Println("get here")

	db, def := getDatabase()
	defer def()

	query := bson.M{
		"point": bson.M{
			"$near": bson.M{
				"$geometry": point,
			},
		},
	}

	err := db.C("meteoStations").Find(query).One(&meteo)

	if err != nil {
		fmt.Println(err)
		fmt.Println("error")
	}

	return
}

func SetMeteoCode() {

	db, def := getDatabase()
	defer def()

	results := make([]model.GeoKml, 1024)

	query := bson.M{
		"$or": []bson.M{
			{"meteoID": ""},
			{"meteoID": bson.M{"$exists": false}},
		},
	}

	err := db.C("geoTile").Find(query).All(&results)

	if err != nil {
		fmt.Printf("Can`t get all records from database, error: %s\n", err.Error())
	}

	for _, res := range results {

		if centerPoint, err := res.Geometry.Center(); err == nil {

			fmt.Println(centerPoint)

			if nearestMeteo := FindNearestStation(centerPoint); nearestMeteo != nil {

				fmt.Println(nearestMeteo)

				meteoCodeID := nearestMeteo.CodeID

				if meteoCodeID != "" {

					err := db.C("geoTile").Update(bson.M{
						"_id": res.ID,
					},
						bson.M{
							"$set": bson.M{
								"meteoID": meteoCodeID,
							},
						})

					if err != nil {

						fmt.Printf("Can`t insert codeID, error: %s\n", err.Error())
					}
				}
			}
		}

	}
}
