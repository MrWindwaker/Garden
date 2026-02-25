package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	GridSize = 3
	TileSize = 100
)

type Game struct {
	Money int
	Day   int
	Tiles []Tile
}

type Tile struct {
	Plant   *Plant
	Watered bool
}

type Plant struct {
	Growth     int
	MaxGrowth  int
	PlantedDay int
}

func (g *Game) AdvanceDay() {
	g.Day++

	for i := range g.Tiles {
		tile := &g.Tiles[i]

		if tile.Plant != nil && tile.Watered {
			tile.Plant.Growth++
		}
		tile.Watered = false
	}
}

func main() {
	fmt.Println("Welcome")

	rl.SetConfigFlags(rl.FlagWindowHighdpi)

	rl.InitWindow(400, 400, "Garden")
	defer rl.CloseWindow()

	g := Game{
		Day:   0,
		Money: 1,
		Tiles: make([]Tile, GridSize*GridSize),
	}

	gridWidht := GridSize * TileSize
	gridHeight := GridSize * TileSize

	offsetX := (rl.GetScreenWidth() - gridWidht) / 2
	offsetY := (rl.GetScreenHeight() - gridHeight) / 2

	for !rl.WindowShouldClose() {

		mouse := rl.GetMousePosition()

		scaleX := float32(rl.GetRenderWidth()) / float32(rl.GetScreenWidth())
		scaleY := float32(rl.GetRenderHeight()) / float32(rl.GetScreenHeight())

		mouse.X *= scaleX
		mouse.Y *= scaleY

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {

			for row := range GridSize {
				for col := range GridSize {
					index := row*GridSize + col
					x := int32(offsetX + col*TileSize)
					y := int32(offsetY + row*TileSize)

					rect := rl.Rectangle{
						X:      float32(x),
						Y:      float32(y),
						Width:  TileSize,
						Height: TileSize,
					}

					if rl.CheckCollisionPointRec(mouse, rect) {
						rl.DrawRectangleLinesEx(rect, 3, rl.Red)
						handleTileClick(&g, index)
					}
				}
			}
		}

		if rl.IsKeyPressed(rl.KeySpace) {
			g.AdvanceDay()
		}

		if rl.IsKeyPressed(rl.KeyA) {
			fmt.Println("Mouse:", rl.GetMousePosition())
			fmt.Println("Screen:", rl.GetScreenWidth(), rl.GetScreenHeight())
			fmt.Println("Render:", rl.GetRenderWidth(), rl.GetRenderHeight())
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.GetColor(0x222222FF))

		for row := range GridSize {
			for col := range GridSize {
				index := row*GridSize + col
				tile := g.Tiles[index]

				x := int32(offsetX + col*TileSize)
				y := int32(offsetY + row*TileSize)

				rl.DrawRectangle(x, y, TileSize, TileSize, rl.Brown)
				rl.DrawRectangleLines(x, y, TileSize, TileSize, rl.Black)

				if tile.Watered {
					rl.DrawRectangle(x, y, TileSize, TileSize, rl.Fade(rl.Blue, 0.3))
				}

				if tile.Plant != nil {
					ratio := float32(tile.Plant.Growth) / float32(tile.Plant.MaxGrowth)
					green := uint8(100 + 155*ratio)

					rl.DrawCircle(
						x+TileSize/2,
						y+TileSize/2,
						30,
						rl.NewColor(0, green, 0, 255),
					)
				}
			}
		}

		rl.DrawText(
			fmt.Sprintf("DAY: %d | Money: %d", g.Day, g.Money),
			10,
			10,
			20,
			rl.White,
		)

		rl.EndDrawing()

	}

}

func handleTileClick(g *Game, index int) {
	tile := &g.Tiles[index]

	if tile.Plant == nil && g.Money >= 1 {
		tile.Plant = &Plant{
			Growth:     0,
			MaxGrowth:  3,
			PlantedDay: g.Day,
		}
		g.Money--
		return
	}

	if tile.Plant != nil {
		if tile.Plant.Growth >= tile.Plant.MaxGrowth {
			g.Money += 2
			tile.Plant = nil
			return
		}

		if !tile.Watered {
			tile.Watered = true
		}
	}
}
