package model

import "gopkg.in/mgo.v2/bson"

type (
	GeoKmlIndex struct {
		MD5  string `bson:"md5"`
		Size int    `bson:"size"`
	}

	GeoKml struct {
		ID       bson.ObjectId `bson:"_id" json:"-"`
		MD5      string        `bson:"md5"`
		Size     int           `bson:"size,minsize"`
		Geometry MultiPolygon  `bson:"geometry" json:"-"`
		GzKML    []byte        `bson:"gzKml"`
		Status   string        `bson:"status,omitempty" json:"status,omitempty"`
	}
)
