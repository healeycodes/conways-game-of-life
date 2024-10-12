package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	RES       int = 20
	CELL_SIZE int = 10 // Size of each cell in pixels
)

type Game struct {
	generation  int
	board       [][]int
	lastUpdate  time.Time
	updateDelay time.Duration
}

var (
	rng *rand.Rand
)

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func emptyGeneration() *Game {
	board := make([][]int, RES)
	for i := 0; i < RES; i++ {
		board[i] = make([]int, RES)
	}
	return &Game{
		board:       board,
		generation:  1,
		lastUpdate:  time.Now(),
		updateDelay: 500 * time.Millisecond, // Update every 500ms
	}
}

func giveState(g *Game) {
	for x := 0; x < RES; x++ {
		for y := 0; y < RES; y++ {
			if rng.Intn(15) == 1 {
				g.board[x][y] = 1
			}
		}
	}
}

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		interaction(x, y, g)
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.updateDelay = max(g.updateDelay-50*time.Millisecond, 50*time.Millisecond)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.updateDelay = min(g.updateDelay+50*time.Millisecond, 2*time.Second)
	}

	if time.Since(g.lastUpdate) >= g.updateDelay {
		g.board = logic(g.board)
		g.generation++
		g.lastUpdate = time.Now()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	draw(g, screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Generation:%v\nUpdate Delay: %v", g.generation, g.updateDelay))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return RES * CELL_SIZE, RES * CELL_SIZE
}

func logic(board [][]int) [][]int {
	newBoard := make([][]int, RES)
	for i := 0; i < RES; i++ {
		newBoard[i] = make([]int, RES)
	}

	for x := 0; x < RES; x++ {
		for y := 0; y < RES; y++ {
			neighbors := checkNeighbors(x, y, board)
			live := board[x][y] == 1

			if live && neighbors < 2 {
				newBoard[x][y] = 0
			} else if live && (neighbors == 2 || neighbors == 3) {
				newBoard[x][y] = 1
			} else if live && neighbors > 3 {
				newBoard[x][y] = 0
			} else if !live && neighbors == 3 {
				newBoard[x][y] = 1
			} else {
				newBoard[x][y] = board[x][y]
			}
		}
	}
	return newBoard
}

func checkNeighbors(x int, y int, board [][]int) int {
	neighbors := 0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < RES && ny >= 0 && ny < RES && board[nx][ny] == 1 {
				neighbors++
			}
		}
	}
	return neighbors
}

func draw(g *Game, screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x00, 0x00, 0x00, 0xff})
	for x := 0; x < RES; x++ {
		for y := 0; y < RES; y++ {
			if g.board[x][y] == 1 {
				ebitenutil.DrawRect(screen, float64(x*CELL_SIZE), float64(y*CELL_SIZE), float64(CELL_SIZE), float64(CELL_SIZE), color.White)
			}
			if x < RES-1 {
				ebitenutil.DrawLine(screen, float64((x+1)*CELL_SIZE), float64(y*CELL_SIZE), float64((x+1)*CELL_SIZE), float64((y+1)*CELL_SIZE), color.RGBA{0x33, 0x33, 0x33, 0xff})
			}
			if y < RES-1 {
				ebitenutil.DrawLine(screen, float64(x*CELL_SIZE), float64((y+1)*CELL_SIZE), float64((x+1)*CELL_SIZE), float64((y+1)*CELL_SIZE), color.RGBA{0x33, 0x33, 0x33, 0xff})
			}
		}
	}
}

func interaction(x int, y int, g *Game) {
	x = clamp(x/CELL_SIZE, 0, RES-1)
	y = clamp(y/CELL_SIZE, 0, RES-1)

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			nx, ny := clamp(x+dx, 0, RES-1), clamp(y+dy, 0, RES-1)
			g.board[nx][ny] = 1
		}
	}
}

func main() {
	g := emptyGeneration()
	giveState(g)

	ebiten.SetWindowSize(RES*CELL_SIZE*2, RES*CELL_SIZE*2)
	ebiten.SetWindowTitle("Conway's Game of Life")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func clamp(x int, min int, max int) int {
	if x < min {
		return min
	} else if x > max {
		return max
	}
	return x
}

func max(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}

func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
