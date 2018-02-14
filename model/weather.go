package model

import "time"

type (
	MonthWeather struct {
		Month         Month        `bson:"month"`
		CodeID        string       `bson:"codeID"`
		Description   string       `bson:"description"`
		PreparingInfo string       `bson:"preparingInfo"`
		CopyRight     string       `bson:"copyRight"`
		Observation   string       `bson:"observation"`
		Days          []DayWeather `bson:"days"`
	}

	DayWeather struct {
		Date            time.Time `bson:"dayTime"`
		MinTemperature  float64   `bson:"minTemperature"`
		MaxTemperature  float64   `bson:"maxTemperature"`
		RainFall        float64   `bson:"rainFall"`
		Evaporation     string    `bson:"evaporation"`
		SunShine        string    `bson:"sunShine"`
		WindDirection   string    `bson:"windDirection"`
		WindSpeed       string    `bson:"windSpeed"`
		WindMaxGustTime string    `bson:"windMaxGustTime"`
		NineAM          Period    `bson:"nineAM"`
		TreePM          Period    `bson:"treePM"`
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
)
