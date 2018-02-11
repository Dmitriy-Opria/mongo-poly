package model

import (
	"errors"
)

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

var notValidCenter = errors.New("not valid center field point")

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

func (p *Polygon) Center() (center Point, err error) {
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

			var heierLon, lowestLon, heierLat, lowestLat float64

			heierLon = westLon
			lowestLon = eastLon
			heierLat = northLat
			lowestLat = southLat

			if heierLon < lowestLon {
				heierLon = eastLon
				lowestLon = westLon
			}

			if heierLat < lowestLat {
				heierLat = southLat
				lowestLat = northLat
			}

			var centerLon, centerLat float64

			centerLon = (heierLon-lowestLon)/2 + lowestLon
			centerLat = (heierLat-lowestLat)/2 + lowestLat

			center.Type = "Point"
			center.Coordinates = [2]float64{centerLon, centerLat}

			return center, nil
		}
	}

	if center.Coordinates[0] == 0 && center.Coordinates[1] == 0 {
		return center, notValidCenter
	}
	return
}

func (p *MultiPolygon) Center() (center Point, err error) {

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
				var heierLon, lowestLon, heierLat, lowestLat float64

				heierLon = westLon
				lowestLon = eastLon
				heierLat = northLat
				lowestLat = southLat

				if heierLon < lowestLon {
					heierLon = eastLon
					lowestLon = westLon
				}

				if heierLat < lowestLat {
					heierLat = southLat
					lowestLat = northLat
				}

				var centerLon, centerLat float64

				centerLon = (heierLon-lowestLon)/2 + lowestLon
				centerLat = (heierLat-lowestLat)/2 + lowestLat

				center.Type = "Point"
				center.Coordinates = [2]float64{centerLon, centerLat}

				return center, nil
			}
		}
	}

	if center.Coordinates[0] == 0 && center.Coordinates[1] == 0 {
		return center, notValidCenter
	}
	return
}
