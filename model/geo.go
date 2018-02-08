package model

type Polygon struct {
	Type        string         `bson:"type"`
	Coordinates [][][2]float64 `bson:"coordinates"`
}

type MultiPoint struct {
	Type        string       `bson:"type"`
	Coordinates [][2]float64 `bson:"coordinates"`
}

type MultiPolygon struct {
	Type        string           `bson:"type"`
	Coordinates [][][][2]float64 `bson:"coordinates"`
}

type Point struct {
	Type        string     `bson:"type"`
	Coordinates [2]float64 `bson:"coordinates"`
}

func (p *Polygon) Bounds() [][2]float64 {
	if len(p.Coordinates) > 0 {
		if len(p.Coordinates[0]) > 0 {

			var westLon, eastLon, northLat, southLat float64 // lon x1, x2; lat y1, y2

			v := p.Coordinates[0][0]
			lon, lat := v[0], v[1]
			westLon, northLat = lon, lat
			eastLon, southLat = lon, lat

			for i := range p.Coordinates {
				for _, v = range p.Coordinates[i] {
					lon, lat = v[0], v[1]
					if lon < westLon {
						westLon = lon
					}
					if lon > eastLon {
						eastLon = lon
					}
					if lat < northLat {
						northLat = lat
					}
					if lat > southLat {
						southLat = lat
					}
				}
			}

			return [][2]float64{
				{westLon, northLat},
				{eastLon, northLat},
				{eastLon, southLat},
				{westLon, southLat},
				{westLon, northLat},
			}
		}
	}
	return nil
}

func (p *MultiPolygon) Bounds() [][2]float64 {

	if len(p.Coordinates) > 0 {
		if len(p.Coordinates[0]) > 0 {
			if len(p.Coordinates[0][0]) > 0 {

				var westLon, eastLon, northLat, southLat float64 // lon x1, x2; lat y1, y2

				v := p.Coordinates[0][0][0]
				lon, lat := v[0], v[1]
				westLon, northLat = lon, lat
				eastLon, southLat = lon, lat

				for i := range p.Coordinates {
					for j := range p.Coordinates[i] {
						for _, v = range p.Coordinates[i][j] {
							lon, lat = v[0], v[1]
							if lon < westLon {
								westLon = lon
							}
							if lon > eastLon {
								eastLon = lon
							}
							if lat < northLat {
								northLat = lat
							}
							if lat > southLat {
								southLat = lat
							}
						}
					}
				}

				return [][2]float64{
					{westLon, northLat},
					{eastLon, northLat},
					{eastLon, southLat},
					{westLon, southLat},
					{westLon, northLat},
				}
			}
		}
	}
	return nil
}
