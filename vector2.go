package space_traders_api

import (
	"math"
)

type Vector2 struct {
	x int
	y int
}

func (self *Vector2) Distance(other Vector2) float64 {
	d := Vector2{0, 0}
	d.x = self.x - other.x
	d.y = self.y - other.y

	if d.x < 0 {
		d.x *= -1
	}
	if d.y < 0 {
		d.y *= -1
	}

	return math.Sqrt(float64(d.x * d.x  + d.y * d.y))
}

// TODO: Implement
func (self *Vector2) FindNearestMarket() Vector2 {
	return Vector2{0, 0}
}

