package db

import (
	"gopkg.in/mgo.v2/bson"
	"mongo-poly/model"
)

func GetMonthDayDegree(query model.WeatherQuery) []model.DayDegree {

	db, def := getDatabase()
	defer def()

	var weatherList = make([]model.MonthWeather, 0, 8)

	for _, month := range query.MonthList {

		monthWeather := model.MonthWeather{}

		weatherQuery := bson.M{
			"codeID": query.CodeID,
			"month": bson.M{
				"monthIndex": month.Month,
				"yearIndex":  month.Year,
			},
		}

		if err := db.C(weatherCol).Find(weatherQuery).One(&monthWeather); err == nil {

			weatherList = append(weatherList, monthWeather)

		}
	}

	var dayDegList = make([]model.DayDegree, 0, 31*len(weatherList))

	for _, month := range weatherList {

		for _, dayWeather := range month.Days {

			if dayWeather.MinTemperature < 12 {
				dayWeather.MinTemperature = 12
			}

			dayDeg := (dayWeather.MinTemperature - 12 + dayWeather.MaxTemperature - 12) / 2

			if dayDeg < 0 {
				dayDeg = 0
			}

			deg := model.DayDegree{
				Date:      dayWeather.DayIndex,
				DayDegree: dayDeg,
			}

			dayDegList = append(dayDegList, deg)
		}
	}

	return dayDegList
}
