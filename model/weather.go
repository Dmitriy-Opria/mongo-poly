package model

/*
import (
	"gopkg.in/mgo.v2/bson"
)*/

type (
	MonthWeather struct {
		//ID            bson.ObjectId `bson:"_id" json:"-"`
		Month         Month        `bson:"month" json:"month"`
		CodeID        string       `bson:"codeID" json:"codeID"`
		Description   string       `bson:"description,omitempty" json:"-"`
		PreparingInfo string       `bson:"preparingInfo,omitempty" json:"-"`
		CopyRight     string       `bson:"copyRight,omitempty" json:"-"`
		Observation   string       `bson:"observation,omitempty" json:"-"`
		NotAll        bool         `bson:"notAll"json:"-"`
		Days          []DayWeather `bson:"days,omitempty" json:"days"`
	}

	DayWeather struct {
		DayIndex        string  `bson:"day,omitempty" json:"day,omitempty"`
		MinTemperature  float64 `bson:"minTemperature,omitempty" json:"minTemperature,omitempty"`
		MaxTemperature  float64 `bson:"maxTemperature,omitempty" json:"maxTemperature,omitempty"`
		RainFall        float64 `bson:"rainFall,omitempty" json:"rainFall,omitempty"`
		Evaporation     string  `bson:"evaporation,omitempty" json:"evaporation,omitempty"`
		SunShine        string  `bson:"sunShine,omitempty" json:"sunShine,omitempty"`
		WindDirection   string  `bson:"windDirection,omitempty" json:"winDirection,omitempty"`
		WindSpeed       string  `bson:"windSpeed,omitempty" json:"windSpeed,omitempty"`
		WindMaxGustTime string  `bson:"windMaxGustTime,omitempty" json:"windMaxGustTime,omitempty"`
		MinRH           int     `bson:"minRH,omitempty,omitempty" json:"minRH,omitempty,omitempty"`
		MaxRH           int     `bson:"maxRH,omitempty,omitempty" json:"maxRH,omitempty,omitempty"`
		NineAM          Period  `bson:"nineAM,omitempty" json:"nineAM,omitempty"`
		TreePM          Period  `bson:"treePM,omitempty" json:"treePM,omitempty"`
	}

	Period struct {
		Temperature   float64 `bson:"temperature,omitempty" json:"temperature,omitempty"`
		Humidity      int     `bson:"humidity,omitempty" json:"humidity,omitempty"`
		Cloud         string  `bson:"cloud,omitempty" json:"cloud,omitempty"`
		WindDirection string  `bson:"windDirection,omitempty" json:"windDirection,omitempty"`
		WindSpeed     string  `bson:"windSpeed,omitempty" json:"windSpeed,omitempty"`
		MslPressure   string  `bson:"mslPressure,omitempty" json:"mslPressure,omitempty"`
	}

	Month struct {
		Month int `bson:"monthIndex,omitempty" json:"month,omitempty"`
		Year  int `bson:"yearIndex,omitempty" json:"year,omitempty"`
	}

	Years []string

	DayDegree struct {
		Date      string  `json:"date,omitempty"`
		DayDegree float64 `json:"day_degree,omitempty"`
	}
)
