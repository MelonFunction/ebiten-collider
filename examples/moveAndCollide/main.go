// Package main 👍
package main

import (
	"errors"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	collider "github.com/melonfunction/ebiten-collider"
)

// vars
var (
	WindowWidth  = 640 * 2
	WindowHeight = 480 * 2

	player *Player
	wall   *collider.SquareShape
	hash   *collider.SpatialHash
	cursor *collider.PointShape

	ErrNormalExit = errors.New("Normal exit")
)

// Player is the moveable shape
type Player struct {
	Bounds *collider.CircleShape
	Speed  float64
}

// Game implements ebiten.Game interface.
type Game struct{}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ErrNormalExit
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyH) {
		player.Bounds.Move(-player.Speed, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyN) {
		player.Bounds.Move(player.Speed, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyC) {
		player.Bounds.Move(0, -player.Speed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyT) {
		player.Bounds.Move(0, player.Speed)
	}

	cx, cy := ebiten.CursorPosition()
	cursor.MoveTo(float64(cx), float64(cy))

	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	hash.Draw(screen)

	red := color.RGBA{255, 0, 0, 128}
	green := color.RGBA{0, 255, 0, 128}
	_ = green

	ebitenutil.DrawCircle(
		screen,
		player.Bounds.Pos.X,
		player.Bounds.Pos.Y,
		player.Bounds.Radius,
		red)

	ebitenutil.DrawRect(
		screen,
		wall.Pos.X-wall.Width/2,
		wall.Pos.Y-wall.Height/2,
		wall.Width,
		wall.Height,
		red)

	ebitenutil.DrawCircle(
		screen,
		cursor.Pos.X,
		cursor.Pos.Y,
		5,
		red)
}

// Layout sets window size
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if WindowWidth != outsideWidth || WindowHeight != outsideHeight {
		log.Println("resize", outsideWidth, outsideHeight)
		WindowWidth = outsideWidth
		WindowHeight = outsideHeight
	}
	return outsideWidth, outsideHeight
}

func main() {
	game := &Game{}
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle("Collisions example")
	ebiten.SetWindowResizable(true)

	hash = collider.NewSpatialHash(128)
	player = &Player{
		Bounds: hash.NewCircleShape(
			float64(WindowWidth)/2-64/2,
			float64(WindowHeight)/2-64/2,
			32),
		Speed: 5,
	}

	wall = hash.NewSquareShape(
		float64(WindowWidth)/2-64/2-96,
		float64(WindowHeight)/2-64/2,
		32,
		320,
	)

	cursor = hash.NewPointShape(0, 0)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
