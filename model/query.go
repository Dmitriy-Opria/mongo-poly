package model

type KmlQuery struct {
	Point   [2]float64   `json:"point"`
	Polygon [][2]float64 `json:"polygon"`
}

type KmlAnswer struct {
	MD5   []string `json:"md5"`
	Empty bool     `json:"-"`
}
