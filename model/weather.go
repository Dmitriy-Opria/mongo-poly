package model
/*
import (
	"gopkg.in/mgo.v2/bson"
)*/

type (
	MonthWeather struct {
		//ID            bson.ObjectId `bson:"_id" json:"-"`
		Month         Month         `bson:"month"`
		CodeID        string        `bson:"codeID"`
		Description   string        `bson:"description"`
		PreparingInfo string        `bson:"preparingInfo"`
		CopyRight     string        `bson:"copyRight"`
		Observation   string        `bson:"observation"`
		NotAll        bool          `bson:"notAll"`
		Days          []DayWeather  `bson:"days"`
	}

	DayWeather struct {
		DayIndex        string  `bson:"day"`
		MinTemperature  float64 `bson:"minTemperature"`
		MaxTemperature  float64 `bson:"maxTemperature"`
		RainFall        float64 `bson:"rainFall"`
		Evaporation     string  `bson:"evaporation"`
		SunShine        string  `bson:"sunShine"`
		WindDirection   string  `bson:"windDirection"`
		WindSpeed       string  `bson:"windSpeed"`
		WindMaxGustTime string  `bson:"windMaxGustTime"`
		NineAM          Period  `bson:"nineAM"`
		TreePM          Period  `bson:"treePM"`
	}

	Period struct {
		Temperature   float64 `bson:"temperature"`
		Humidity      int     `bson:"humidity"`
		Cloud         string  `bson:"cloud"`
		WindDirection string  `bson:"windDirection"`
		WindSpeed     string  `bson:"windSpeed"`
		MslPressure   string  `bson:"mslPressure"`
	}

	Month struct {
		Month int `bson:"monthIndex"`
		Year  int `bson:"yearIndex"`
	}

	DayDegree struct {
		Date string `json:"date"`
		DayDegree float64 `json:"day_degree"`
	}
)
