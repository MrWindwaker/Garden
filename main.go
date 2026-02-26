package main

import (
	"encoding/json"
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	GridSize = 3
	TileSize = 100
)

var NEWGAME Game = Game{
	Day:   0,
	Money: 1,
	Tiles: make([]Tile, GridSize*GridSize),
	PlantType: []PlantType{
		{"Carrot", 3, 1, 2, rl.Green},
		{"Pumpkin", 4, 2, 6, rl.Orange},
		{"Potato", 8, 6, 12, rl.Magenta},
	},
	Selected: 0,
	GameOver: false,
}

type Game struct {
	Money     int
	Day       int
	Tiles     []Tile
	PlantType []PlantType
	Selected  int
	GameOver  bool
}

type Tile struct {
	Plant   *Plant
	Watered bool
}

type PlantType struct {
	Name     string
	GrowDays int
	Cost     int
	Reward   int
	Color    rl.Color
}

type Plant struct {
	Growth    int
	TypeIndex int
}

func main() {
	fmt.Println("Welcome")

	rl.SetConfigFlags(rl.FlagWindowHighdpi)

	rl.InitWindow(400, 400, "Garden")
	defer rl.CloseWindow()

	shouldClose := rl.WindowShouldClose()

	g := Game{}
	loadGame(&g)

	if g.PlantType == nil || g.GameOver {
		g = NEWGAME
	}

	gridWidth := GridSize * TileSize
	gridHeight := GridSize * TileSize

	offsetX := (rl.GetScreenWidth() - gridWidth) / 2
	offsetY := (rl.GetScreenHeight() - gridHeight) / 2

	for !shouldClose {

		shouldClose = rl.WindowShouldClose()

		if g.GameOver {
			if rl.IsKeyPressed(rl.KeyR) {
				g = NEWGAME
			}
		}

		if !g.GameOver {
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

			if rl.IsKeyPressed(rl.KeyTab) {
				g.Selected = (g.Selected + 1) % len(g.PlantType)
			}
		}

		if g.IsGameOver() {
			g.GameOver = true
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.GetColor(0x222222FF))

		if !g.GameOver {

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
						pt := g.PlantType[tile.Plant.TypeIndex]

						growthRatio := float32(tile.Plant.Growth) / float32(pt.GrowDays)
						ratio := int32(10 + 20*growthRatio)

						rl.DrawCircle(
							x+TileSize/2,
							y+TileSize/2,
							float32(ratio),
							pt.Color,
						)

						if tile.Plant.Growth >= pt.GrowDays {
							rl.DrawCircleLines(
								x+TileSize/2,
								y+TileSize/2,
								float32(ratio+1),
								rl.White,
							)
						}
					}
				}
			}

			rl.DrawText(
				fmt.Sprintf("DAY: %d | Money: %d | Selected: %s", g.Day, g.Money, g.PlantType[g.Selected].Name),
				5,
				10,
				20,
				rl.White,
			)
		} else {
			rl.DrawText("GAME OVER Press r to restart", 10, 10, 20, rl.RayWhite)
		}

		rl.EndDrawing()

	}

	saveGame(&g)

}

func handleTileClick(g *Game, index int) {
	tile := &g.Tiles[index]
	pt := g.PlantType[g.Selected]

	if tile.Plant == nil && g.Money >= pt.Cost {
		tile.Plant = &Plant{
			Growth:    0,
			TypeIndex: g.Selected,
		}
		g.Money -= pt.Cost
		return
	}

	if tile.Plant != nil {
		pt = g.PlantType[tile.Plant.TypeIndex]

		if tile.Plant.Growth >= pt.GrowDays {
			g.Money += pt.Reward
			tile.Plant = nil
			return
		}

		if !tile.Watered {
			tile.Watered = true
		}
	}
}

func saveGame(g *Game) {
	data, _ := json.MarshalIndent(g, "", " ")
	os.WriteFile("save.json", data, 0644)
}

func loadGame(g *Game) {
	data, err := os.ReadFile("save.json")
	if err != nil {
		return
	}
	json.Unmarshal(data, g)
}

func (g *Game) HasGrowingPlants() bool {
	for i := range g.Tiles {
		tile := &g.Tiles[i]

		if tile.Plant != nil {
			pt := g.PlantType[tile.Plant.TypeIndex]

			if tile.Plant.Growth > pt.GrowDays {
				return true
			}

			return true
		}
	}
	return false
}

func (g *Game) IsGameOver() bool {
	minCost := g.PlantType[0].Cost

	for _, pt := range g.PlantType {
		if pt.Cost < minCost {
			minCost = pt.Cost
		}

		if g.Money >= minCost {
			return false
		}

		if g.HasGrowingPlants() {
			return false
		}
	}

	return true
}

func (g *Game) AdvanceDay() {
	g.Day++

	for i := range g.Tiles {
		tile := &g.Tiles[i]

		if tile.Plant != nil {
			pt := g.PlantType[tile.Plant.TypeIndex]

			if tile.Watered {
				tile.Plant.Growth++
			}

			if !tile.Watered && tile.Plant.Growth < pt.GrowDays {
				tile.Plant.Growth--

				if tile.Plant.Growth < 0 {
					tile.Plant = nil
				}
			}
		}

		tile.Watered = false
	}
}
