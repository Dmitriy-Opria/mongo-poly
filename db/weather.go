package db

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/iizotop/baseweb/utils"
	"gopkg.in/mgo.v2/bson"
	"io"
	"mongo-poly/model"
	"net/http"
	"net/textproto"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	invalidYearError  = errors.New("invalid year value")
	invalidMonthError = errors.New("invalid month value")
	badStatusCode     = errors.New("bad status code")
	emptyWeather      = errors.New("empty field weather")
)

const (
	OZF = uint8(1)
	BOM = uint8(2)
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

		requestUrl, path := GetPath(yearStr, monthStr, meteo.CodeID, BOM)

		err := DownloadFile(path, requestUrl)

		if err != nil {
			fmt.Println(err.Error())
		}
		break
	}

	return
}

func GetPath(yearStr, monthStr, codeID string, source uint8) (requestUrl, filePath string) {
	if source == BOM {

		if len(monthStr) == 1 {
			monthStr = "0" + monthStr
		}

		requestUrl = fmt.Sprintf("http://www.bom.gov.au/climate/dwo/%s/%s/text/%s.%s%s.csv", yearStr, monthStr, codeID, yearStr, monthStr)
		filePath = fmt.Sprintf("./%s.%s%s.csv", codeID, yearStr, monthStr)

	} else if source == OZF {

		reqUrl := fmt.Sprintf("http://ozforecast.com.au/cgi-bin/aws_export.cgi?pagetype=csv&aws=%s&year=%s", codeID, yearStr)

		requestUrl = "http://ozforecast.com.au"

		if doc, err := goquery.NewDocument(reqUrl); err == nil {

			selection := doc.Find("table")
			if selection != nil {
				selection.Contents().Find("tr").Each(func(i int, tr *goquery.Selection) {
					a := tr.Find("a")
					if href, ok := a.Attr("href"); ok {

						if strings.Contains(href, ".csv") {
							requestUrl += href
						}
					}
				})
			}
		}

		filePath = fmt.Sprintf("./%s.%s.csv", codeID, yearStr)

	}

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

	if err != nil {
		return err
	}

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

func ReadBOMWeatherFile(filepath string, month model.Month, codeID string) (weatherList []model.MonthWeather) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0644)

	if err != nil {
		fmt.Printf("Can`t open file %s, error: %s\n", filepath, err.Error())
		return
	}
	defer file.Close()

	r := csv.NewReader(file)

	weather := model.MonthWeather{}

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

		dayWeather.DayIndex = record[1]
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

	weather.Month = month

	weather.CodeID = codeID

	tm := time.Now()

	nowMonth := tm.Month()

	nowYear := tm.Year()

	if month.Month == int(nowMonth) && month.Year == nowYear {
		weather.NotAll = true
	}

	weatherList = append(weatherList, weather)
	return
}

func ReadOZFWeatherFile(filepath string, codeID string) (weatherList []model.MonthWeather) {

	file, err := os.OpenFile(filepath, os.O_RDONLY, 0644)

	if err != nil {
		fmt.Printf("Can`t open file %s, error: %s\n", filepath, err.Error())
		return
	}
	defer file.Close()

	r := textproto.NewReader(bufio.NewReader(file))

	var count int

	var currentYear, currentMonth int

	var monthWeaher = model.MonthWeather{}

	var firstCheck = true

	for {
		line, err := r.ReadLine()
		record := strings.Split(line, ", ")

		if err == io.EOF {
			break
		}
		if count < 4 {
			count++
			continue
		}

		if len(record) < 11 {
			continue
		}

		if tm, err := time.Parse("2006-01-02", record[0]); err == nil {

			year, month := tm.Year(), tm.Month()

			if firstCheck {
				currentYear = year
				currentMonth = int(month)
				firstCheck = false
			}

			if currentYear != year || currentMonth != int(month) {

				monthWeaher.Month = model.Month{
					Year:  currentYear,
					Month: currentMonth,
				}

				monthWeaher.CodeID = codeID

				tm := time.Now()

				nowMonth := tm.Month()
				nowYear := tm.Year()

				if monthWeaher.Month.Month == int(nowMonth) && monthWeaher.Month.Year == nowYear {
					monthWeaher.NotAll = true
				}
				weatherList = append(weatherList, monthWeaher)
				monthWeaher = model.MonthWeather{}
				currentYear = year
				currentMonth = int(month)
			}

			dayWeather := model.DayWeather{}

			sunShine := record[7]
			evaporation := record[8]

			if sunShine == "?" {
				sunShine = ""
			}

			if evaporation == "?" {
				evaporation = ""
			}

			dayWeather.DayIndex = record[0]
			dayWeather.MinTemperature = utils.ToFloat64(record[1])
			dayWeather.MaxTemperature = utils.ToFloat64(record[2])
			dayWeather.WindSpeed = record[3]
			dayWeather.NineAM.WindSpeed = record[4]
			dayWeather.NineAM.WindDirection = record[5]
			dayWeather.RainFall = utils.ToFloat64(record[6])
			dayWeather.SunShine = sunShine
			dayWeather.Evaporation = evaporation
			dayWeather.MinRH = utils.ToInt(record[9])
			dayWeather.MaxRH = utils.ToInt(record[10])

			monthWeaher.Days = append(monthWeaher.Days, dayWeather)

		} else {
			fmt.Println(err.Error())
		}
	}
	return
}

func SaveRangeWeather(monthList []model.Month) {

	meteoList := GetMeteoList()

	for _, meteo := range meteoList {

		if meteo.Source == "OZF" {

			source := OZF

			for _, year := range getYearList(monthList) {

				for _, month := range monthList {

					if isSavedWeather(meteo.CodeID, month) {
						continue
					}

					requestUrl, filePath := GetPath(year, "", meteo.CodeID, source)

					if err := DownloadFile(filePath, requestUrl); err == nil {

						weather := ReadOZFWeatherFile(filePath, meteo.CodeID)

						if ok := insertWeather(weather); ok {
							if err := os.Remove(filePath); err != nil {
								fmt.Println(err)
							}
						} else {
							removeNotAll(meteo.CodeID, month)
						}
					}
				}
			}

		} else {

			source := BOM

			for _, month := range monthList {

				if !isSavedWeather(meteo.CodeID, month) {

					yearStr := strconv.Itoa(month.Year)
					monthStr := strconv.Itoa(month.Month)

					requestUrl, filePath := GetPath(yearStr, monthStr, meteo.CodeID, source)

					if err := DownloadFile(filePath, requestUrl); err == nil {

						weather := ReadBOMWeatherFile(filePath, month, meteo.CodeID)

						if ok := insertWeather(weather); ok {
							if err := os.Remove(filePath); err != nil {
								fmt.Println(err)
							}
						} else {
							removeNotAll(meteo.CodeID, month)
						}
					}
				}
			}
		}
	}
}

func getYearList(monthList []model.Month) (years model.Years) {

	currentYear := 0

	for i, month := range monthList {
		if i == 0 {
			currentYear = month.Year
		}
		if currentYear != month.Year {
			currentYear = month.Year
			years = append(years, strconv.Itoa(currentYear))
		}
	}
	return
}

func insertWeather(weatherList []model.MonthWeather) (ok bool) {

	db, def := getDatabase()
	defer def()

	for _, weather := range weatherList {

		err := db.C("weather").Insert(weather)

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

	db.C("weather").RemoveAll(notAllQuery)
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

func GetMonthDayDegree(codeID string, monthList []model.Month) []model.DayDegree {

	db, def := getDatabase()
	defer def()

	var weatherList = make([]model.MonthWeather, 0, 8)

	for _, month := range monthList {

		monthWeather := model.MonthWeather{}

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

func GetCodeIDByMD5(md5hash string) (codeID string, err error) {

	return "IDCJDW3001", nil

	db, def := getDatabase()
	defer def()

	var result model.GeoKml

	query := bson.M{
		"md5": md5hash,
	}

	if err = db.C("geoKml").Find(query).One(&result); err != nil {
		fmt.Println(err.Error())
		return
	}

	return result.MeteoCodeID, nil

}
