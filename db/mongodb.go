package db

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"mongo_kml/config"
	"mongo_kml/model"
	"strings"
)

var (
	mgoSession *mgo.Session
	mgoDBName  string = "test"
)

func MongoInit(mgoConfig config.MongodbConfig) {

	userPassword := ""

	if mgoConfig.User != "" {
		userPassword = fmt.Sprintf("%s:%s@", mgoConfig.User, mgoConfig.Password)
	}

	url := fmt.Sprintf("mongodb://%s%s/%s", userPassword, mgoConfig.Host, mgoConfig.DBName)

	fmt.Println("MongoInit:", url)

	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Println("mgo.Dial:", err)
		return
	}

	mgoSession = session
	mgoDBName = mgoConfig.DBName

	initMgoIndexes()
}

func getSession() (*mgo.Session, func()) {
	s := mgoSession.Clone()
	return s, s.Close
}

func getDatabase() (*mgo.Database, func()) {
	s := mgoSession.Clone()
	return s.DB(mgoDBName), s.Close
}

func initMgoIndexes() {

	processedIndexes()
	geoTileIndexes()
	geoKmlIndexes()
	tasksIndexes()
	meteoIndex()
}

func tasksIndexes() {

	db, def := getDatabase()
	defer def()

	var err error
	var key mgo.Index

	tasks := db.C("tasks")

	key = mgo.Index{
		Key:      []string{"hash"},
		Unique:   true,
		DropDups: true,
	}
	if err = tasks.EnsureIndex(key); err != nil {
		fmt.Printf("tasks(%q): %#v\n", strings.Join(key.Key, "_"), err)
	}

	key = mgo.Index{
		Key: []string{"granule"},
	}
	if err = tasks.EnsureIndex(key); err != nil {
		fmt.Printf("tasks(%q): %#v\n", strings.Join(key.Key, "_"), err)
	}

	key = mgo.Index{
		Key: []string{"processed"},
	}
	if err = tasks.EnsureIndex(key); err != nil {
		fmt.Printf("tasks(%q): %#v\n", strings.Join(key.Key, "_"), err)
	}

	key = mgo.Index{
		Key: []string{"sensing_time"},
	}
	if err = tasks.EnsureIndex(key); err != nil {
		fmt.Printf("tasks(%q): %#v\n", strings.Join(key.Key, "_"), err)
	}

	key = mgo.Index{
		Key: []string{"update_time"},
	}
	if err = tasks.EnsureIndex(key); err != nil {
		fmt.Printf("tasks(%q): %#v\n", strings.Join(key.Key, "_"), err)
	}
}

func processedIndexes() {

	db, def := getDatabase()
	defer def()

	var err error
	var key mgo.Index

	processed := db.C("processed")

	key = mgo.Index{
		Key:      []string{"md5", "type", "day"},
		Unique:   true,
		DropDups: true,
	}
	if err = processed.EnsureIndex(key); err != nil {
		fmt.Printf("processed(%q): %#v\n", strings.Join(key.Key, "_"), err)
	}
}

func geoTileIndexes() {

	db, def := getDatabase()
	defer def()

	var err error
	var key mgo.Index

	geoTile := db.C("geoTile")

	key = mgo.Index{
		Key:      []string{"granuleID"},
		Unique:   true,
		DropDups: true,
	}
	if err = geoTile.EnsureIndex(key); err != nil {
		fmt.Printf("geoTile(%q): %#v\n", strings.Join(key.Key, "_"), err)
	}

	key = mgo.Index{
		Key:  []string{"$2dsphere:geometry"},
		Bits: 26,
	}
	if err = geoTile.EnsureIndex(key); err != nil {
		fmt.Printf("geoTile(%q): %#v\n", strings.Join(key.Key, "_"), err)
	}

	key = mgo.Index{
		Key: []string{"mgrsTile"},
	}
	if err = geoTile.EnsureIndex(key); err != nil {
		fmt.Printf("geoTile(%q): %#v\n", strings.Join(key.Key, "_"), err)
	}
}

func geoKmlIndexes() {

	db, def := getDatabase()
	defer def()

	var err error
	var key mgo.Index

	geoKml := db.C("geoKml")

	key = mgo.Index{
		Key:      []string{"md5", "size"},
		Unique:   true,
		DropDups: true,
	}
	if err = geoKml.EnsureIndex(key); err != nil {
		fmt.Printf("geoKml(%q): %#v\n", strings.Join(key.Key, "_"), err)
	}

	key = mgo.Index{
		Key:  []string{"$2dsphere:geometry"},
		Bits: 26,
	}
	if err = geoKml.EnsureIndex(key); err != nil {
		fmt.Printf("geoKml(%q): %#v\n", strings.Join(key.Key, "_"), err)
	}
}

func meteoIndex() {
	db, def := getDatabase()
	defer def()

	var err error
	var key mgo.Index

	meteoStations := db.C("meteoStations")

	key = mgo.Index{
		Key:  []string{"$2dsphere:point"},
		Bits: 26,
	}
	if err = meteoStations.EnsureIndex(key); err != nil {
		fmt.Printf("meteoStations(%q): %#v\n", strings.Join(key.Key, "_"), err)
	}
}

/*func FindTileByGeometry(geometry model.Polygon) (tiles []model.GeoTile) {
	db, def := getDatabase()
	defer def()
	query := bson.M{
		"geometry": bson.M{
			"$geoIntersects": bson.M{
				"$geometry": geometry,
			},
		},
	}
	err := db.C("geoTile").Find(query).All(&tiles) // .Select(geoTileFields)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}*/

func InsertPolygon(poly model.GeoKml) {

	db, def := getDatabase()
	defer def()

	if centerPoint, err := poly.Geometry.Center(); err != nil {

		meteoStation := FindNearestStation(centerPoint)
		poly.MeteoCodeID = meteoStation.CodeID
	}

	err := db.C("geoTile").Insert(poly)
	if err != nil {
		fmt.Println(err)
		return
	}
	return

}

func FindPolygonInPolygon(poly model.Polygon) (answer model.KmlAnswer) {

	db, def := getDatabase()
	defer def()

	query := bson.M{
		"geometry": bson.M{
			"$geoIntersects": bson.M{
				"$geometry": poly,
			},
		},
	}
	intersectedFields := make([]model.GeoKml, 0, 16)

	err := db.C("geoTile").Find(query).All(&intersectedFields)

	if err != nil {
		answer.Empty = true
		fmt.Println(err)
		return
	} else {
		answer.MD5 = make([]string, 0, len(intersectedFields))
		if len(intersectedFields) == 0 {
			answer.Empty = true
			return
		}
		for _, field := range intersectedFields {
			answer.MD5 = append(answer.MD5, field.MD5)
		}
	}
	return

}
func FindPointInPolygon(point model.Point) (answer model.KmlAnswer) {

	db, def := getDatabase()
	defer def()

	query := bson.M{
		"geometry": bson.M{
			"$geoIntersects": bson.M{
				"$geometry": point,
			},
		},
	}
	intersectedFields := make([]model.GeoKml, 124)

	err := db.C("geoTile").Find(query).All(&intersectedFields)
	if err != nil {
		answer.Empty = true
		fmt.Println(err)
		return
	} else {
		answer.MD5 = make([]string, 0, len(intersectedFields))
		if len(intersectedFields) == 0 {
			answer.Empty = true
			return
		}
		for _, field := range intersectedFields {
			answer.MD5 = append(answer.MD5, field.MD5)
		}
	}
	return

}
