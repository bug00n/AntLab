//

package core

import (
	"antlab/source/entities"
	"antlab/source/geometry_utils"
)

type World struct {
	Width, Height   float64
	Ants            []entities.Ant
	AnthillPosition geometry_utils.Point
}
