//

package core

import (
	"math"
	"math/rand/v2"

	"antlab/source/entities"
	"antlab/source/geometry_utils"
)

type Cell struct {
	FoodSmell          float64
	FoodPheromone      float64
	FearPheromone      float64
	EstimatedPheromone float64
}

type Anthill struct {
	AnthillPosition geometry_utils.Point
	FoodCnt         int
}

type World struct {
	Width, Height float64
	Ants          []entities.Ant
	Foods         []entities.Food
	Anthill       Anthill
	Cells         [][]Cell
	GridSize      int
}

type Smell struct {
	Position    geometry_utils.Point
	Strength    int16
	MaxStrength int16
}

func round(x float64) float64 {
	return math.Round(x*10000) / 10000
}

func (w *World) getGridCoords(x, y float64) (col, row int) {
	if w.GridSize <= 0 {
		return -1, -1
	}
	cellSize := float64(w.GridSize)
	return int(math.Floor(x / cellSize)), int(math.Floor(y / cellSize))
}

func getFoodSmellIntensity(smellIntensity int16, maxIntensity int16) float64 {
	x := 1 - float64(smellIntensity)/float64(maxIntensity)
	return 2 - math.Exp(x) - max(0, x-0.35)/10
}

func HighBiasedIntn(n int) int {
	// Число 4.0 зашито внутрь как коэффициент перекоса (можно поменять под себя)
	biased := 1.0 - math.Pow(1.0-rand.Float64(), 1.0)
	return int(biased * float64(n))
}

func (w *World) smellSpread() {
	if w.GridSize <= 0 || len(w.Cells) == 0 {
		return
	}

	cellSize := float64(w.GridSize)
	q := make([]Smell, 0, len(w.Foods))
	for _, food := range w.Foods {
		x, y := food.Position.X, food.Position.Y
		if food.SmellIntensity == 0 || x < 0 || y < 0 || x >= w.Width || y >= w.Height {
			continue
		}

		if food.Size <= 0 {
			continue
		}

		col, row := w.getGridCoords(x, y)
		if row < 0 || row >= len(w.Cells) || col < 0 || col >= len(w.Cells[row]) {
			continue
		}

		w.Cells[row][col].FoodSmell = 1
		q = append(q, Smell{
			Position:    food.Position,
			Strength:    int16(food.SmellIntensity) - int16(HighBiasedIntn(int(food.SmellIntensity/3))),
			MaxStrength: int16(food.SmellIntensity),
		})
	}

	for head := 0; head < len(q); head++ {
		v := q[head]
		col, row := w.getGridCoords(v.Position.X, v.Position.Y)

		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if dx == 0 && dy == 0 {
					continue
				}

				nextCol, nextRow := col+dx, row+dy
				if nextRow < 0 || nextRow >= len(w.Cells) || nextCol < 0 || nextCol >= len(w.Cells[nextRow]) {
					continue
				}

				newX := float64(nextCol) * cellSize
				newY := float64(nextRow) * cellSize
				if newX >= w.Width || newY >= w.Height {
					continue
				}

				newStrength := int(v.Strength) - 1 - rand.IntN(2)
				if newStrength <= 0 {
					continue
				}

				nextIntensity := getFoodSmellIntensity(int16(newStrength), v.MaxStrength)
				cell := &w.Cells[nextRow][nextCol]
				if nextIntensity > cell.FoodSmell {
					cell.FoodSmell = nextIntensity
					q = append(q, Smell{
						Position:    geometry_utils.Point{X: newX, Y: newY},
						Strength:    int16(newStrength),
						MaxStrength: v.MaxStrength,
					})
				}
			}
		}
	}

	for i := range w.Cells {
		for j := range w.Cells[i] {
			cell := &w.Cells[i][j]
			cell.FoodSmell *= 0.95
			cell.FoodPheromone *= 0.9995
			if cell.FoodPheromone < 0.0001 {
				cell.FoodPheromone = 0
			}
			cell.FearPheromone *= 0.5

			cell.FoodSmell = round(cell.FoodSmell)
			cell.FearPheromone = round(cell.FearPheromone)

			if cell.EstimatedPheromone > 0.3 {
				cell.EstimatedPheromone *= 0.9995
			} else {
				cell.EstimatedPheromone *= 0.999
			}
			cell.EstimatedPheromone = round(cell.EstimatedPheromone)
		}
	}
}

func (w *World) noiseSmell() {

}

func (w *World) Step(dt float64) {
	if dt <= 0 {
		return
	}

	for i := range w.Ants {
		w.Ants[i].Move(w.Width, w.Height, dt, w.foodSmellAt, w.getFood, w.homeSmellAt, w.anthillScan, w.foodPheromoneAt)
		w.Ants[i].Update(w.Width, w.Height, dt, w.foodSmellAt, w.getFood, w.homeSmellAt, w.anthillScan, w.foodPheromoneAt, w.putFood)

		col, row := w.getGridCoords(w.Ants[i].Position.X, w.Ants[i].Position.Y)

		if row >= 0 && row < len(w.Cells) && col >= 0 && col < len(w.Cells[row]) {
			cell := &w.Cells[row][col]

			if w.Ants[i].Job == entities.Search {
				cell.EstimatedPheromone = max(w.Ants[i].Stamina, cell.EstimatedPheromone)
				if w.Ants[i].Stamina < 0.01 {
					w.Ants[i].Stamina = 0.005
				} else {
					w.Ants[i].Stamina *= 0.997
					w.Ants[i].Stamina = round(w.Ants[i].Stamina)
				}
			}
			if w.Ants[i].Job == entities.CarryFood {
				cell.FoodPheromone = max(w.Ants[i].FoodPheromoneCapacity, cell.FoodPheromone)
				w.Ants[i].FoodPheromoneCapacity *= 0.995
			}
		}
	}

	w.smellSpread()
	w.noiseSmell()
}

func (w *World) foodSmellAt(x, y float64) float64 {
	if w.GridSize <= 0 ||
		x < 0 || y < 0 || x >= w.Width || y >= w.Height {
		return 0
	}

	col, row := w.getGridCoords(x, y)

	if row < 0 || row >= len(w.Cells) {
		return 0
	}
	if col < 0 || col >= len(w.Cells[row]) {
		return 0
	}

	return w.Cells[row][col].FoodSmell
}

func (w *World) getFood(x, y float64, maxCnt int) int {
	cellX, cellY := w.getGridCoords(x, y)

	for id, f := range w.Foods {
		foodX, foodY := w.getGridCoords(f.Position.X, f.Position.Y)

		if (cellX-foodX)*(cellX-foodX)+(cellY-foodY)*(cellY-foodY) <= 10 {
			// В будущем тут нужен мьютекс
			diff := min(w.Foods[id].Size, maxCnt)
			w.Foods[id].Size -= diff
			return diff
		}
	}

	return 0
}

func (w *World) homeSmellAt(x, y float64) (float64, float64) {
	if w.GridSize <= 0 ||
		x < 0 || y < 0 || x >= w.Width || y >= w.Height {
		return -1, -1
	}

	col, row := w.getGridCoords(x, y)

	if row < 0 || row >= len(w.Cells) {
		return -1, -1
	}
	if col < 0 || col >= len(w.Cells[row]) {
		return -1, -1
	}

	bestCol, bestRow := col, row

	for dx := -2; dx <= 2; dx++ {
		for dy := -2; dy <= 2; dy++ {
			nextCol, nextRow := col+dx, row+dy

			if 0 <= nextRow && nextRow < len(w.Cells) &&
				0 <= nextCol && nextCol < len(w.Cells[nextRow]) {
				if w.Cells[bestRow][bestCol].EstimatedPheromone < w.Cells[nextRow][nextCol].EstimatedPheromone {
					bestCol, bestRow = nextCol, nextRow
				}
			}
		}
	}

	if bestCol == col && bestRow == row {
		return x, y
	}

	cellSize := float64(w.GridSize)
	return (float64(bestCol) + 0.5) * cellSize, (float64(bestRow) + 0.5) * cellSize
}

func (w *World) foodPheromoneAt(x, y float64) (float64, float64) {
	if w.GridSize <= 0 || x < 0 || y < 0 || x >= w.Width || y >= w.Height {
		return -1, -1
	}

	col, row := w.getGridCoords(x, y)
	if row < 0 || row >= len(w.Cells) || col < 0 || col >= len(w.Cells[row]) {
		return -1, -1
	}

	const radius = 5
	bestCol, bestRow := -1, -1
	// A destination must be stronger than the current cell, as with the home trail.
	bestStrength := max(0.005, w.Cells[row][col].FoodPheromone)

	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nextCol, nextRow := col+dx, row+dy
			if nextRow < 0 || nextRow >= len(w.Cells) || nextCol < 0 || nextCol >= len(w.Cells[nextRow]) {
				continue
			}

			strength := w.Cells[nextRow][nextCol].FoodPheromone
			if strength > bestStrength {
				bestCol, bestRow = nextCol, nextRow
				bestStrength = strength
			}
		}
	}

	if bestCol < 0 {
		return -1, -1
	}
	cellSize := float64(w.GridSize)
	return (float64(bestCol) + 0.5) * cellSize, (float64(bestRow) + 0.5) * cellSize
}

func (w *World) anthillScan(x, y float64) bool {
	cellX, cellY := w.getGridCoords(x, y)

	hillX, hillY := w.getGridCoords(w.Anthill.AnthillPosition.X, w.Anthill.AnthillPosition.Y)

	if (cellX-hillX)*(cellX-hillX)+(cellY-hillY)*(cellY-hillY) <= 10 {
		return true
	}

	return false
}

func (w *World) putFood(cnt int) {
	w.Anthill.FoodCnt += cnt
}
