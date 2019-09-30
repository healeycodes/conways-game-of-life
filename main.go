package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Game struct {
	board [][]int
}

var (
	game *Game
)

func newGame() *Game {
	board := make([][]int, 480)
	for i := 0; i < 480; i++ {
		board[i] = make([]int, 480)
	}
	return &Game{board: board}
}

func draw(game *Game, background *ebiten.Image) {
	for x := 0; x < 480; x++ {
		for y := 0; y < 480; y++ {
			if game.board[x][y] == 1 {
				ebitenutil.DrawRect(background, float64(x), float64(y), 1, 1, color.White)
			}
		}
	}
}

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	screen.Fill(color.RGBA{0, 0, 0, 0xff})
	background, _ := ebiten.NewImage(480, 480, ebiten.FilterDefault)
	draw(game, background)
	screen.DrawImage(background, &ebiten.DrawImageOptions{})
	return nil
}

func main() {
	game = newGame()
	if err := ebiten.Run(update, 480, 480, 2, "Hello, World!"); err != nil {
		log.Fatal(err)
	}
}
