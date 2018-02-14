package db

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/iizotop/baseweb/utils"
	"gopkg.in/mgo.v2/bson"
	"io"
	"mongo_kml/model"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	domainName = "http://www.bom.gov.au/climate/dwo/"
	textField  = "/text/"
	pointPart  = "."
	csvPart    = ".csv"

	invalidYearError  = errors.New("invalid year value")
	invalidMonthError = errors.New("invalid month value")
	badStatusCode     = errors.New("bad status code")
)

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

		requestUrl, path := getPath(yearStr, monthStr, meteo.CodeID)

		err := DownloadFile(path, requestUrl)

		if err != nil {
			fmt.Println(err.Error())
		}
		break
	}

	return
}

func getPath(yearStr, monthStr, codeID string) (requestUrl, filePath string) {

	if len(monthStr) == 1 {
		monthStr = "0" + monthStr
	}

	requestUrl = domainName + yearStr + monthStr + textField + codeID + pointPart + yearStr + monthStr + csvPart
	filePath = codeID + pointPart + yearStr + monthStr + csvPart

	return
}

func GetMeteoList() (meteo []model.MeteoUnit) {

	db, def := getDatabase()
	defer def()

	err := db.C("meteoStations").Find(nil).All(&meteo)

	if err != nil {
		fmt.Println(err)
	}
	return
}

func DownloadFile(filepath string, url string) error {

	resp, err := http.Get(url)

	fmt.Println(url)
	fmt.Println(resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return badStatusCode
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Downloaded:", filepath)
	return nil
}

func ReadWeatherFile(filepath string) (weather model.MonthWeather) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0644)

	if err != nil {
		fmt.Printf("Can`t open file %s, error: %s\n", filepath, err.Error())
		return
	}
	defer file.Close()

	r := csv.NewReader(file)

	weather.Days = make([]model.DayWeather, 0, 31)

	var count int

	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}
		if err != nil {

			if !strings.Contains(err.Error(), csv.ErrFieldCount.Error()) {
				fmt.Println(err.Error())
				continue
			}
		}
		if count < 5 {

			if len(record) > 0 {

				switch count {
				case 0:
					weather.Description = record[0]
				case 1:
					weather.PreparingInfo = record[0]
				case 2:
					weather.CopyRight = record[0]
				case 3:
					weather.Observation = record[0]
				}
			}

			count++
			continue
		}

		if len(record) < 22 {
			continue
		}

		dayWeather := model.DayWeather{}

		layout := "2006-01-2"

		tm, _ := time.Parse(layout, record[1])

		dayWeather.Date = tm
		dayWeather.MinTemperature = utils.ToFloat64(record[2])
		dayWeather.MaxTemperature = utils.ToFloat64(record[3])
		dayWeather.RainFall = utils.ToFloat64(record[4])
		dayWeather.Evaporation = record[5]
		dayWeather.SunShine = record[6]
		dayWeather.WindDirection = record[7]
		dayWeather.WindSpeed = record[8]
		dayWeather.WindMaxGustTime = record[9]
		dayWeather.NineAM.Temperature = utils.ToFloat64(record[10])
		dayWeather.NineAM.Humidity = utils.ToInt(record[11])
		dayWeather.NineAM.Cloud = record[12]
		dayWeather.NineAM.WindDirection = record[13]
		dayWeather.NineAM.WindSpeed = record[14]
		dayWeather.NineAM.MslPressure = record[15]
		dayWeather.TreePM.Temperature = utils.ToFloat64(record[16])
		dayWeather.TreePM.Humidity = utils.ToInt(record[17])
		dayWeather.TreePM.Cloud = record[18]
		dayWeather.TreePM.WindDirection = record[19]
		dayWeather.TreePM.WindSpeed = record[20]
		dayWeather.TreePM.MslPressure = record[21]

		weather.Days = append(weather.Days, dayWeather)
	}
	return
}

func SaveRangeWeather(monthList []model.Month) {

	meteoList := GetMeteoList()

	for _, meteo := range meteoList {

		for _, month := range monthList {

			if !isSavedWeather(meteo.CodeID, month) {

				yearStr := strconv.Itoa(month.Year)
				monthStr := strconv.Itoa(month.Month)

				requestUrl, filePath := getPath(yearStr, monthStr, meteo.CodeID)

				if err := DownloadFile(filePath, requestUrl); err == nil {

					monthWeather := ReadWeatherFile(filePath)

					monthWeather.Month = month

					monthWeather.CodeID = meteo.CodeID

					if ok := insertWeather(monthWeather); ok {
						if err := os.Remove(filePath); err != nil {
							fmt.Println(err)
						}
					}
				}

				time.Sleep(3 * time.Second)
			}

		}
	}
}

func insertWeather(weather model.MonthWeather) (ok bool) {

	db, def := getDatabase()
	defer def()

	err := db.C("weather").Insert(weather)

	if err != nil {

		fmt.Println(err)
		return
	}
	fmt.Printf("Inserted year: %d, month: %d\n", weather.Month.Year, weather.Month.Month)

	return true
}

func isSavedWeather(codeID string, month model.Month) bool {

	db, def := getDatabase()
	defer def()

	query := bson.M{
		"codeID": codeID,
		"month": bson.M{
			"monthIndex": month.Month,
			"yearIndex":  month.Year,
		},
	}
	n, err := db.C("weather").Find(query).Count()

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

func FindFieldWeather(md5hash string, year, month int) (monthWeather *model.MonthWeather) {

	db, def := getDatabase()
	defer def()

	var result model.GeoKml

	query := bson.M{
		"md5": md5hash,
	}

	if err := db.C("geoKml").Find(query).One(&result); err != nil {
		fmt.Println(err.Error())
		return
	}

	weatherQuery := bson.M{
		"codeID": result.MeteoCodeID,
		"month": bson.M{
			"monthIndex": month,
			"yearIndex":  year,
		},
	}

	if err := db.C("weather").Find(weatherQuery).One(&monthWeather); err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}
