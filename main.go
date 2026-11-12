package main

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"

	"antlab/source/config"
	"antlab/source/core"
	"antlab/source/entities"
	"antlab/source/geometry_utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	world   *core.World
	dt      float64
	sprites Sprites
}

func (g *Game) Update() error {
	g.world.Step(g.dt)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 120, 120, 155}) // Set background color

	g.drawFoodSmell(screen)
	g.drawEstimatedPheromone(screen)
	g.drawFoodPheromone(screen)
	g.drawFoods(screen)
	g.drawAnthill(screen)

	g.drawAnts(screen)

	metrics := fmt.Sprintf("FPS: %.1f\nTPS: %.1f", ebiten.ActualFPS(), ebiten.ActualTPS())
	ebitenutil.DebugPrintAt(screen, metrics, 10, 10)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(g.world.Width), int(g.world.Height)
}

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatal(err)
	}

	sprites := SpritesImport()

	sieve := make([][]core.Cell, cfg.World.Height/cfg.World.GridSize)
	for i := range sieve {
		sieve[i] = make([]core.Cell, cfg.World.Width/cfg.World.GridSize)
	}

	food := entities.Food{
		Position: geometry_utils.Point{
			X: 1000,
			Y: 650,
		},
		Size:           1000,
		SmellIntensity: 50,
	}

	anthill := core.Anthill{
		AnthillPosition: geometry_utils.Point{
			X: float64(cfg.World.Width) / 2,
			Y: float64(cfg.World.Height) / 2,
		},
		FoodCnt: 0,
	}

	world := &core.World{
		Width:  float64(cfg.World.Width),
		Height: float64(cfg.World.Height),
		Anthill: anthill,
		Cells:    sieve,
		Foods:    []entities.Food{food},
		GridSize: cfg.World.GridSize,
	}

	for i := 0; i < cfg.Ant.Count; i++ {
		ant := entities.Ant{}
		ant.Position.X = world.Width / 2
		ant.Position.Y = world.Height / 2
		ant.Speed = 35
		ant.Rotation = float64(i) * 0.1
		ant.Job = entities.Search
		ant.Capacity = 5
		ant.Stamina = 1.0
		ant.FoodPheromoneCapacity = 1.0
		world.Ants = append(world.Ants, ant)
	}

	game := &Game{
		world:   world,
		dt:      1.0 / float64(cfg.TPS),
		sprites: sprites,
	}

	ebiten.SetWindowSize(int(cfg.World.Width), int(cfg.World.Height))
	// ebiten.SetFullscreen(true)
	ebiten.SetTPS(cfg.TPS)
	ebiten.SetVsyncEnabled(true)

	ebiten.SetWindowTitle("Ant Lab")
	ebiten.SetWindowIcon([]image.Image{sprites.iconSprite})

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
