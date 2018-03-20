package db

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/iizotop/baseweb/utils"
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

func GetPath(yearStr, monthStr, codeID string, source uint8) (requestUrl, filePath string) {
	if source == BOM {

		if len(monthStr) == 1 {
			monthStr = "0" + monthStr
		}

		requestUrl = fmt.Sprintf("http://www.bom.gov.au/climate/dwo/%s%s/text/%s.%s%s.csv", yearStr, monthStr, codeID, yearStr, monthStr)
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

	err := db.C(stationCol).Find(nil).All(&meteo)

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

func ReadBOMWeatherFile(filePath string, month model.Month, codeID string) (weatherList []model.MonthWeather) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)

	if err != nil {
		fmt.Printf("Can`t open file %s, error: %s\n", filePath, err.Error())
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

func ReadOZFWeatherFile(filePath string, codeID string) (weatherList []model.MonthWeather) {

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)

	if err != nil {
		fmt.Printf("Can`t open file %s, error: %s\n", filePath, err.Error())
		return
	}
	defer file.Close()

	r := textproto.NewReader(bufio.NewReader(file))

	var count int

	var currentYear, currentMonth int

	var monthWeather = model.MonthWeather{}

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

				monthWeather.Month = model.Month{
					Year:  currentYear,
					Month: currentMonth,
				}

				monthWeather.CodeID = codeID

				tm := time.Now()

				nowMonth := tm.Month()
				nowYear := tm.Year()

				if monthWeather.Month.Month == int(nowMonth) && monthWeather.Month.Year == nowYear {
					monthWeather.NotAll = true
				}
				weatherList = append(weatherList, monthWeather)
				monthWeather = model.MonthWeather{}
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

			monthWeather.Days = append(monthWeather.Days, dayWeather)

		} else {
			fmt.Println(err.Error())
		}
	}
	return
}

func UpdateWeatherData() {

	tm := time.Now()

	from, to := fmt.Sprintf("%d-%d", tm.Year()-2, tm.Month()), fmt.Sprintf("%d-%d", tm.Year(), tm.Month())

	monthList := model.GetMonthList(from, to)

	meteoList := GetMeteoList()

	for _, meteo := range meteoList {

		fmt.Println("get here")

		if meteo.Source == "OZF" {

			fmt.Println("get OZ")

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

			fmt.Println("get BOM")

			source := BOM

			for _, month := range monthList {

				if isSavedWeather(meteo.CodeID, month) {
					fmt.Println("are saved")
					continue
				}

				yearStr := strconv.Itoa(month.Year)
				monthStr := strconv.Itoa(month.Month)

				requestUrl, filePath := GetPath(yearStr, monthStr, meteo.CodeID, source)

				fmt.Println("get path", requestUrl, filePath)
				if err := DownloadFile(filePath, requestUrl); err == nil {

					weather := ReadBOMWeatherFile(filePath, month, meteo.CodeID)

					if ok := insertWeather(weather); ok {
						if err := os.Remove(filePath); err != nil {
							fmt.Println(err)
						}
					} else {
						removeNotAll(meteo.CodeID, month)
					}
				} else {
					fmt.Println(err.Error())
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
