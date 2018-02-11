package model

type MeteoUnit struct {
	Location string  `bson:"location"`
	Station  string  `bson:"station"`
	CodeID   string  `bson:"codeId"`
	Site     string  `bson:"site"`
	Dist     string  `bson:"dist"`
	Start    int     `bson:"start"`
	End      int     `bson:"end"`
	Point    Point   `bson:"point"`
	Source   string  `bson:"source"`
	STA      string  `bson:"sta"`
	Height   float64 `bson:"height"`
	Bar      float64 `bson:"bar"`
	WMO      string  `bson:"wmo"`
}
