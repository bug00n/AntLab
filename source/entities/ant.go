package entities

import (
	"antlab/source/geometry_utils"
	"math"
	"math/rand/v2"
)

type Ant struct {
	Position geometry_utils.Point
	Speed	 float64
	Rotation float64
}

func getRotation(rotation float64) float64 {
	if rand.Float64() < 0.08 {
		rotation += (rand.Float64() - 0.5) * math.Pi / 4
	}
	return rotation
}

func (ant *Ant) Move(Width, Height float64) {
	ant.Position.X += ant.Speed * float64(math.Cos(ant.Rotation))
	ant.Position.Y += ant.Speed * float64(math.Sin(ant.Rotation))

	x_padding, y_padding := Width*0.01, Height*0.01

	if (ant.Position.X < x_padding) {
		ant.Rotation = math.Pi - ant.Rotation
	}
	if (ant.Position.Y < y_padding) {
		ant.Rotation = -ant.Rotation
	}

	if (ant.Position.X > Width-x_padding) {
		ant.Rotation = math.Pi - ant.Rotation
	}
	if (ant.Position.Y > Height-y_padding) {
		ant.Rotation = -ant.Rotation
	}

	ant.Rotation = getRotation(ant.Rotation)
}
