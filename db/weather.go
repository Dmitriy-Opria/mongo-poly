package db

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/iizotop/baseweb/utils"
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

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

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

			yearStr := strconv.Itoa(month.Year)
			monthStr := strconv.Itoa(month.Month)

			requestUrl, filePath := getPath(yearStr, monthStr, meteo.CodeID)

			DownloadFile(filePath, requestUrl)

			monthWeather := ReadWeatherFile(filePath)

			insertWeather(monthWeather)
		}
	}
}

func insertWeather(weather model.MonthWeather) {

	db, def := getDatabase()
	defer def()

	err := db.C("weather").Insert(weather)

	if err != nil {

		fmt.Println(err)
		return
	}
}
