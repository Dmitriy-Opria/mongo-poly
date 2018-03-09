package db

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"golib/utils"
	"gopkg.in/mgo.v2/bson"
	"mongo-poly/model"
	"net/http"
	"strconv"
	"time"
)

var (
	tableHeader = []string{
		"Date",
		"Minimum temperature (°C)",
		"Maximum temperature (°C)",
		"Rainfall (mm)",
		"Evaporation (mm)",
		"Sunshine (hours)",
		"Direction of maximum wind gust",
		"Speed of maximum wind gust (km/h)",
		"Time of maximum wind gust",
		"9am Temperature (°C)",
		"9am relative humidity (%)",
		"9am cloud amount (oktas)",
		"9am wind direction",
		"9am wind speed (km/h)",
		"9am MSL pressure (hPa)",
		"3pm Temperature (°C)",
		"3pm relative humidity (%)",
		"3pm cloud amount (oktas)",
		"3pm wind direction",
		"3pm wind speed (km/h)",
		"3pm MSL pressure (hPa)",
	}
)

func GetWeatherResponse(w http.ResponseWriter, r *http.Request, codeID, format string, monthList []model.Month) {

	weatherList := getPeriodWeather(codeID, monthList)

	tm := time.Now().Unix()

	if len(weatherList) > 0 {

		if format == "csv" {
			if csvCont := WriteCSVWeather(weatherList); csvCont != nil {
				w.Header().Set("Content-Type", "text/csv")
				w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s_%d.csv", codeID, tm))
				w.Write(csvCont.Bytes())
			}
		} else if format == "xlsx" {
			if xlsxCont := WriteXLSXWeather(weatherList); xlsxCont != nil {
				w.Header().Set("Content-Type", "application/vnd.ms-excel")
				w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s_%d.xlsx", codeID, tm))
				w.Write(xlsxCont.Bytes())
			}
		} else if format == "json" {
			if jsonCont := GetJsonWeather(weatherList); jsonCont != nil {
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonCont.Bytes())
			}
		}
	}

}

func getPeriodWeather(codeID string, monthList []model.Month) []model.MonthWeather {

	db, def := getDatabase()
	defer def()

	var weatherList = make([]model.MonthWeather, 0, 12)

	for _, month := range monthList {

		var monthWeather model.MonthWeather

		weatherQuery := bson.M{
			"codeID": codeID,
			"month": bson.M{
				"monthIndex": month.Month,
				"yearIndex":  month.Year,
			},
		}

		if err := db.C("weather").Find(weatherQuery).One(&monthWeather); err == nil {
			weatherList = append(weatherList, monthWeather)
		}

	}

	return weatherList
}

func WriteCSVWeather(weatherList []model.MonthWeather) (csvCont *bytes.Buffer) {

	records := [][]string{
		tableHeader,
	}

	for _, month := range weatherList {

		for _, day := range month.Days {

			rec := []string{
				day.DayIndex,
				utils.FloatToStr(day.MinTemperature),
				utils.FloatToStr(day.MaxTemperature),
				utils.FloatToStr(day.RainFall),
				day.Evaporation,
				day.SunShine,
				day.WindDirection,
				day.WindSpeed,
				day.WindMaxGustTime,
				utils.FloatToStr(day.NineAM.Temperature),
				strconv.Itoa(day.NineAM.Humidity),
				day.NineAM.Cloud,
				day.NineAM.WindDirection,
				day.NineAM.WindSpeed,
				day.NineAM.MslPressure,
				utils.FloatToStr(day.TreePM.Temperature),
				strconv.Itoa(day.TreePM.Humidity),
				day.TreePM.Cloud,
				day.TreePM.WindDirection,
				day.TreePM.WindSpeed,
				day.TreePM.MslPressure,
			}

			records = append(records, rec)
		}
	}

	b := &bytes.Buffer{}
	wr := csv.NewWriter(b)

	err := wr.WriteAll(records)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return b
}

func WriteXLSXWeather(weatherList []model.MonthWeather) (xlsxCont *bytes.Buffer) {

	file := xlsx.NewFile()

	sheet, err := file.AddSheet("Weather")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row := sheet.AddRow()
	for _, field := range tableHeader {
		cell := row.AddCell()
		cell.Value = field
	}
	for _, month := range weatherList {

		for _, field := range month.Days {
			row := sheet.AddRow()
			cell := row.AddCell()
			cell.Value = field.DayIndex
			cell = row.AddCell()
			cell.Value = utils.FloatToStr(field.MinTemperature)
			cell = row.AddCell()
			cell.Value = utils.FloatToStr(field.MaxTemperature)
			cell = row.AddCell()
			cell.Value = utils.FloatToStr(field.RainFall)
			cell = row.AddCell()
			cell.Value = field.Evaporation
			cell = row.AddCell()
			cell.Value = field.SunShine
			cell = row.AddCell()
			cell.Value = field.WindDirection
			cell = row.AddCell()
			cell.Value = field.WindSpeed
			cell = row.AddCell()
			cell.Value = field.WindMaxGustTime
			cell = row.AddCell()
			cell.Value = utils.FloatToStr(field.NineAM.Temperature)
			cell = row.AddCell()
			cell.Value = strconv.Itoa(field.NineAM.Humidity)
			cell = row.AddCell()
			cell.Value = field.NineAM.Cloud
			cell = row.AddCell()
			cell.Value = field.NineAM.WindDirection
			cell = row.AddCell()
			cell.Value = field.NineAM.WindSpeed
			cell = row.AddCell()
			cell.Value = field.NineAM.MslPressure
			cell = row.AddCell()
			cell.Value = utils.FloatToStr(field.TreePM.Temperature)
			cell = row.AddCell()
			cell.Value = strconv.Itoa(field.TreePM.Humidity)
			cell = row.AddCell()
			cell.Value = field.TreePM.Cloud
			cell = row.AddCell()
			cell.Value = field.TreePM.WindDirection
			cell = row.AddCell()
			cell.Value = field.TreePM.WindSpeed
			cell = row.AddCell()
			cell.Value = field.TreePM.MslPressure
		}
	}

	b := &bytes.Buffer{}
	err = file.Write(b)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	return b
}

func getMonthList(from, to time.Time) (monthList []model.Month) {

	monthList = make([]model.Month, 0, 12)

	month := model.Month{}

	for y := from.Year(); y <= to.Year(); y++ {

		fmt.Println(y)

		if y != from.Year() && y != to.Year() {

			for m := 1; m < 12; m++ {
				month.Month = m
				month.Year = y
				monthList = append(monthList, month)
			}

		} else if y == from.Year() && y == to.Year() {
			for m := int(from.Month()); m <= int(to.Month()); m++ {
				month.Month = m
				month.Year = y
				monthList = append(monthList, month)
			}

		} else if y == from.Year() {

			for m := int(from.Month()); m < 12; m++ {
				month.Month = m
				month.Year = from.Year()
				monthList = append(monthList, month)
			}

		} else if y == to.Year() {

			for m := int(1); m < int(to.Month()); m++ {
				month.Month = m
				month.Year = to.Year()
				monthList = append(monthList, month)
			}
		}
	}

	return
}

func GetJsonWeather(weatherList []model.MonthWeather) (jsonCont *bytes.Buffer) {

	jsonCont = new(bytes.Buffer)

	enc := json.NewEncoder(jsonCont)

	enc.Encode(weatherList)
	return
}
