package model

import "time"

type MonthWeather struct {
	Description   string
	PreparingInfo string
	CopyRight     string
	Observation   string
	Days          []DayWeather
}

type DayWeather struct {
	Date                time.Time
	MinTemperature      float64
	MaxTemperature      float64
	RainFall            float64
	Evaporation         string
	SunShine            string
	WindDirection       string
	WindSpeed           string
	WindMaxGustTime     string
	NineAmTemperature   float64
	NineAmHumidity      int
	NineAmCloud         string
	NineAmWindDirection string
	NineAmWindSpeed     string
	NineAmMslPressure   string
	TreePmTemperature   float64
	TreePmHumidity      int
	TreePmCloud         string
	TreePmWindDirection string
	TreePmWindSpeed     string
	TreePmMslPressure   string
}
