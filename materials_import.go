// Load pics from '/materials'

package main

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

//go:embed materials/ant_logo.png
var antLogoPNG []byte

//go:embed materials/ant.png
var antPNG []byte

//go:embed materials/food.png
var foodPNG []byte

//go:embed materials/anthill.png
var anthillPNG []byte

type Sprites struct {
	iconSprite		*ebiten.Image
	antSprite  		*ebiten.Image
	foodSprite 		*ebiten.Image
	anthillSprite 	*ebiten.Image
}

func SpritesImport() Sprites {
	icon, _, err := image.Decode(bytes.NewReader(antLogoPNG))
	if err != nil {
		log.Fatal(err)
	}

	antSprite, _, err := image.Decode(bytes.NewReader(antPNG))
	if err != nil {
		log.Fatal(err)
	}

	foodSprite, _, err := image.Decode(bytes.NewReader(foodPNG))
	if err != nil {
		log.Fatal(err)
	}

	anthillSprite, _, err := image.Decode(bytes.NewReader(anthillPNG))
	if err != nil {
		log.Fatal(err)
	}

	return Sprites{
		iconSprite:		ebiten.NewImageFromImage(icon),
		antSprite:		ebiten.NewImageFromImage(antSprite),
		foodSprite: 	ebiten.NewImageFromImage(foodSprite),
		anthillSprite:	ebiten.NewImageFromImage(anthillSprite),
	}
}
