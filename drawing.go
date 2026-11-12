package main

import (
	"antlab/source/entities"
	_ "embed"
	"image/color"
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) drawFoodSmell(screen *ebiten.Image) {
	if g.world.GridSize <= 0 {
		return
	}
	cellSize := float32(g.world.GridSize)

	for row, cells := range g.world.Cells {
		for col, cell := range cells {
			if cell.FoodSmell > 0 {
				strength := math.Min(cell.FoodSmell, 1)
				// Blending blue over the background gives the intensity gradient.
				cellColor := color.NRGBA{B: 255, A: uint8(math.Round(255 * strength))}
				if cellColor.A == 0 {
					continue
				}
				vector.FillRect(screen,
					float32(col)*cellSize, float32(row)*cellSize,
					cellSize, cellSize, cellColor, false,
				)
			}
		}
	}
}

func (g *Game) drawEstimatedPheromone(screen *ebiten.Image) {
	if g.world.GridSize <= 0 {
		return
	}
	cellSize := float32(g.world.GridSize)

	for row, cells := range g.world.Cells {
		for col, cell := range cells {
			if !(cell.EstimatedPheromone > 0) {
				continue
			}
			strength := math.Min(cell.EstimatedPheromone, 1)
			cellColor := color.NRGBA{G: 100, A: uint8(math.Round(255 * strength))}
			if cellColor.A == 0 {
				continue
			}
			vector.FillRect(screen,
				float32(col)*cellSize, float32(row)*cellSize,
				cellSize, cellSize, cellColor, false,
			)
		}
	}
}

func (g *Game) drawFoodPheromone(screen *ebiten.Image) {
	if g.world.GridSize <= 0 {
		return
	}
	cellSize := float32(g.world.GridSize)

	for row, cells := range g.world.Cells {
		for col, cell := range cells {
			if !(cell.FoodPheromone > 0) {
				continue
			}
			strength := math.Min(cell.FoodPheromone, 1)
			cellColor := color.NRGBA{R: 255, G: 165, A: uint8(math.Round(255 * strength))}
			if cellColor.A == 0 {
				continue
			}
			vector.FillRect(screen,
				float32(col)*cellSize, float32(row)*cellSize,
				cellSize, cellSize, cellColor, false,
			)
		}
	}
}

func (g *Game) drawFoods(screen *ebiten.Image) {
	w := float64(g.sprites.foodSprite.Bounds().Dx())
	h := float64(g.sprites.foodSprite.Bounds().Dy())

	for _, food := range g.world.Foods {
		if food.Size <= 0 {
			continue
		}

		var op ebiten.DrawImageOptions
		op.GeoM.Translate(-w/2, -h/2)
		op.GeoM.Scale(30/w, 30/h)
		op.GeoM.Translate(food.Position.X, food.Position.Y)

		screen.DrawImage(g.sprites.foodSprite, &op)
	}
}

func (g *Game) drawAnts(screen *ebiten.Image) {
	for _, ant := range g.world.Ants {
		op := &ebiten.DrawImageOptions{}

		w := float64(g.sprites.antSprite.Bounds().Dx())
		h := float64(g.sprites.antSprite.Bounds().Dy())

		op.GeoM.Translate(-w/2, -h/2)

		if ant.Job == entities.CarryFood {
			op.GeoM.Scale(16/w, 16/h)
		}
		if ant.Job == entities.Search {
			op.GeoM.Scale(12/w, 12/h)
		}
		if ant.Job == entities.Tired {
			op.GeoM.Scale(12/w, 12/h)
		}

		op.GeoM.Rotate(ant.Rotation + math.Pi/2)
		op.GeoM.Translate(ant.Position.X, ant.Position.Y)

		screen.DrawImage(g.sprites.antSprite, op)
	}
}

func (g *Game) drawAnthill(screen *ebiten.Image) {
	w := float64(g.sprites.foodSprite.Bounds().Dx())
	h := float64(g.sprites.foodSprite.Bounds().Dy())

	var op ebiten.DrawImageOptions
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Scale(30/w, 30/h)
	op.GeoM.Translate(g.world.Anthill.AnthillPosition.X, g.world.Anthill.AnthillPosition.Y)

	screen.DrawImage(g.sprites.anthillSprite, &op)
}
