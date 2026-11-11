package main

import (
	"fmt"
	"bytes"
	_ "embed"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"

	"antlab/source/config"
	"antlab/source/core"
	"antlab/source/entities"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	world     *core.World
	antSprite *ebiten.Image
}

func (g *Game) Update() error {
	for i := range g.world.Ants {
		g.world.Ants[i].Move(g.world.Width, g.world.Height)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 120, 120, 155})

	for _, ant := range g.world.Ants {
		op := &ebiten.DrawImageOptions{}

		w := float64(g.antSprite.Bounds().Dx())
		h := float64(g.antSprite.Bounds().Dy())

		op.GeoM.Translate(-w/2, -h/2)

		op.GeoM.Scale(12/w, 12/h)
		op.GeoM.Rotate(ant.Rotation + math.Pi/2)
		op.GeoM.Translate(ant.Position.X, ant.Position.Y)

		screen.DrawImage(g.antSprite, op)
	}

	ebitenutil.DebugPrintAt(
    screen,
    fmt.Sprintf(
        "FPS: %.1f\nTPS: %.1f",
        ebiten.ActualFPS(),
        ebiten.ActualTPS(),
    ),
    10, 10,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(g.world.Width), int(g.world.Height)
}

//go:embed materials/ant_logo.png
var antLogoPNG []byte

//go:embed materials/ant.png
var antPNG []byte

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatal(err)
	}

	icon, _, err := image.Decode(bytes.NewReader(antLogoPNG))
	if err != nil {
		log.Fatal(err)
	}

	antSprite, _, err := image.Decode(bytes.NewReader(antPNG))
	if err != nil {
		log.Fatal(err)
	}

	world := &core.World{
		Width:  float64(cfg.World.Width),
		Height: float64(cfg.World.Height),
	}

	for i := 0; i < cfg.Ant.Count; i++ {
		ant := entities.Ant{}
		ant.Position.X = world.Width / 2
		ant.Position.Y = world.Height / 2
		ant.Speed = 1
		ant.Rotation = float64(i) * 0.1
		world.Ants = append(world.Ants, ant)
	}

	game := &Game{
		world:     world,
		antSprite: ebiten.NewImageFromImage(antSprite),
	}

	ebiten.SetWindowSize(int(cfg.World.Width), int(cfg.World.Height))
	// ebiten.SetFullscreen(true)

	ebiten.SetWindowTitle("Ant Lab")
	ebiten.SetWindowIcon([]image.Image{icon})
	
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
