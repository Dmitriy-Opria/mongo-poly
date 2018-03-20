package db

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mongo-poly/model"
	"strconv"
	"time"
)

func insertWeather(weatherList []model.MonthWeather) (ok bool) {

	db, def := getDatabase()
	defer def()

	for _, weather := range weatherList {

		err := db.C(weatherCol).Insert(weather)

		if err != nil {

			fmt.Println(err)
			return
		}
		fmt.Printf("Inserted year: %d, month: %d\n", weather.Month.Year, weather.Month.Month)

	}

	return true
}

func isSavedWeather(codeID string, month model.Month) bool {

	db, def := getDatabase()
	defer def()

	query := bson.M{
		"$or": []bson.M{
			{"codeID": codeID,
				"notAll": bson.M{"$exists": false},
				"month": bson.M{
					"monthIndex": month.Month,
					"yearIndex":  month.Year,
				},
			},
			{"codeID": codeID,
				"notAll": false,
				"month": bson.M{
					"monthIndex": month.Month,
					"yearIndex":  month.Year,
				},
			},
		},
	}
	n, err := db.C(weatherCol).Find(query).Count()

	if err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		if n > 0 {
			return true
		}
	}
	return false
}

func removeNotAll(codeID string, month model.Month) {

	db, def := getDatabase()
	defer def()

	notAllQuery := bson.M{
		"codeID": codeID,
		"notAll": true,
		"month": bson.M{
			"monthIndex": month.Month,
			"yearIndex":  month.Year,
		},
	}

	db.C(weatherCol).RemoveAll(notAllQuery)
}

func GetCodeIDByMD5(md5hash string) (codeID string, err error) {

	db, def := getDatabase()
	defer def()

	var result model.GeoKml

	query := bson.M{
		"md5": md5hash,
	}

	if err = db.C(geoKmlCol).Find(query).One(&result); err != nil {
		fmt.Println(err.Error())
		return
	}

	return result.MeteoCodeID, nil

}

func GetWeather(year, month int) (err error) {

	if year < 2016 || year > time.Now().Year() {
		return invalidYearError
	}

	if month <= 0 || month > 12 {
		return invalidMonthError
	}

	meteoList := GetMeteoList()

	yearStr := strconv.Itoa(year)
	monthStr := strconv.Itoa(month)

	for _, meteo := range meteoList {

		requestUrl, path := GetPath(yearStr, monthStr, meteo.CodeID, BOM)

		err := DownloadFile(path, requestUrl)

		if err != nil {
			fmt.Println(err.Error())
		}
		break
	}

	return
}

func FindFieldWeather(md5hash string, year, month int) (monthWeather *model.MonthWeather) {

	db, def := getDatabase()
	defer def()

	codeID, err := GetCodeIDByMD5(md5hash)

	if err == nil {
		weatherQuery := bson.M{
			"codeID": codeID,
			"month": bson.M{
				"monthIndex": month,
				"yearIndex":  year,
			},
		}

		if err := db.C("weather").Find(weatherQuery).One(&monthWeather); err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	return
}

func GetDayDeg(md5hash string, day time.Time) (dayDeg float64, err error) {

	dayStr := day.Format("2006-01-2")

	db, def := getDatabase()
	defer def()

	codeID, err := GetCodeIDByMD5(md5hash)

	if err == nil {

		monthWeather := model.MonthWeather{}

		weatherQuery := bson.M{
			"codeID": codeID,
			"month": bson.M{
				"monthIndex": int(day.Month()),
				"yearIndex":  day.Year(),
			},
		}

		daySelector := bson.M{
			"days": bson.M{
				"$elemMatch": bson.M{
					"day": dayStr,
				},
			},
		}

		if err = db.C("weather").Find(weatherQuery).Select(daySelector).One(&monthWeather); err != nil {
			fmt.Println(err.Error())
			return
		}

		if len(monthWeather.Days) > 0 {

			dayWeather := monthWeather.Days[0]

			if dayWeather.MinTemperature < 12 {
				dayWeather.MinTemperature = 12
			}

			dayDeg = (dayWeather.MinTemperature - 12 + dayWeather.MaxTemperature - 12) / 2

			if dayDeg < 0 {
				dayDeg = 0
			}

			return
		}
	}

	return 0, emptyWeather
}
