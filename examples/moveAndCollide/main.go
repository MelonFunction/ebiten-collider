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
	wall   *collider.RectangleShape
	wall2  *collider.RectangleShape
	obs    *collider.CircleShape
	obs2   *collider.CircleShape
	hash   *collider.SpatialHash
	cursor *collider.PointShape

	ErrNormalExit = errors.New("Normal exit")
)

// Player is the moveable shape
type Player struct {
	Bounds *collider.CircleShape
	// Bounds *collider.RectangleShape
	Speed float64
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

	collisions := hash.CheckCollisions(player.Bounds)
	for _, collision := range collisions {
		sep := collision.SeparatingVector
		player.Bounds.Move(sep.X, sep.Y)
		// player.Bounds.Move(sep.X/2, sep.Y/2)
		// collision.Other.Move(-sep.X/2, -sep.Y/2)
	}

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
	// ebitenutil.DrawRect(
	// 	screen,
	// 	player.Bounds.Pos.X-player.Bounds.Width/2,
	// 	player.Bounds.Pos.Y-player.Bounds.Height/2,
	// 	player.Bounds.Width,
	// 	player.Bounds.Height,
	// 	red)

	ebitenutil.DrawCircle(
		screen,
		obs.Pos.X,
		obs.Pos.Y,
		obs.Radius,
		red)
	ebitenutil.DrawCircle(
		screen,
		obs2.Pos.X,
		obs2.Pos.Y,
		obs2.Radius,
		red)

	ebitenutil.DrawRect(
		screen,
		wall.Pos.X-wall.Width/2,
		wall.Pos.Y-wall.Height/2,
		wall.Width,
		wall.Height,
		red)
	ebitenutil.DrawRect(
		screen,
		wall2.Pos.X-wall2.Width/2,
		wall2.Pos.Y-wall2.Height/2,
		wall2.Width,
		wall2.Height,
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

	x := float64(WindowWidth)/2 - 64/2
	y := float64(WindowHeight)/2 - 64/2

	hash = collider.NewSpatialHash(128)
	player = &Player{
		Bounds: hash.NewCircleShape(x, y, 32),
		// Bounds: hash.NewRectangleShape(x, y, 64, 64),
		Speed: 1,
	}

	wall = hash.NewRectangleShape(
		x,
		y-16,
		128,
		128,
	)
	wall2 = hash.NewRectangleShape(
		x+128,
		y-16,
		128,
		128*2,
	)

	obs = hash.NewCircleShape(
		x+96,
		y+256+64,
		32)
	obs2 = hash.NewCircleShape(
		x+128,
		y+256,
		64)

	cursor = hash.NewPointShape(0, 0)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
