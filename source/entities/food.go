package entities

import "antlab/source/geometry_utils"

type Food struct {
	Position geometry_utils.Point
	Size     		  int
	SmellIntensity    uint16
}

