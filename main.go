package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 240
	screenHeight = 240
)

const (
	tileSize = 16
)

var (
	tilesImage *ebiten.Image
)

func init() {
	// Abre o arquivo tiles.png local
	file, err := os.Open("assets/tiles.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Decodifica a imagem
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	tilesImage = ebiten.NewImageFromImage(img)
}

type Character struct {
	x        float64
	y        float64
	targetX  float64
	targetY  float64
	speed    float64
	isMoving bool
}

type Game struct {
	layers     [][]int
	cameraX    float64
	cameraY    float64
	mouseX     int
	mouseY     int
	isDragging bool
	lastMouseX int
	lastMouseY int
	zoom       float64
	character  Character
}

func (g *Game) Update() error {
	// Atualiza a posição do mouse
	g.mouseX, g.mouseY = ebiten.CursorPosition()

	// Verifica se o botão esquerdo do mouse está pressionado
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !g.isDragging {
			// Converte a posição do mouse para coordenadas do mundo
			worldX := (float64(g.mouseX) + g.cameraX) / g.zoom
			worldY := (float64(g.mouseY) + g.cameraY) / g.zoom

			// Define o alvo do personagem
			g.character.targetX = worldX
			g.character.targetY = worldY
			g.character.isMoving = true
		}
	}

	// Atualiza a posição do personagem
	if g.character.isMoving {
		dx := g.character.targetX - g.character.x
		dy := g.character.targetY - g.character.y
		distance := math.Sqrt(dx*dx + dy*dy)

		if distance < 1.0 {
			g.character.x = g.character.targetX
			g.character.y = g.character.targetY
			g.character.isMoving = false
		} else {
			g.character.x += dx * 0.05
			g.character.y += dy * 0.05
		}
	}

	// Detecta o scroll do mouse para zoom
	_, dy := ebiten.Wheel()
	if dy != 0 {
		// Calcula o fator de zoom (0.1 é a sensibilidade do zoom)
		zoomFactor := 1.0 + dy*0.1

		// Limita o zoom entre 0.1 e 4.0
		newZoom := g.zoom * zoomFactor
		if newZoom >= 0.1 && newZoom <= 4.0 {
			// Ajusta a posição da câmera para manter o ponto sob o mouse
			mouseWorldX := (float64(g.mouseX) + g.cameraX) / g.zoom
			mouseWorldY := (float64(g.mouseY) + g.cameraY) / g.zoom

			g.zoom = newZoom

			// Recalcula a posição da câmera para manter o ponto sob o mouse
			g.cameraX = mouseWorldX*g.zoom - float64(g.mouseX)
			g.cameraY = mouseWorldY*g.zoom - float64(g.mouseY)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	w := tilesImage.Bounds().Dx()
	tileXCount := w / tileSize

	// Draw each tile with each DrawImage call.
	// As the source images of all DrawImage calls are always same,
	// this rendering is done very efficiently.
	// For more detail, see https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage
	const xCount = screenWidth / tileSize
	for _, l := range g.layers {
		for i, t := range l {
			op := &ebiten.DrawImageOptions{}

			// Posiciona o tile primeiro
			tileX := float64((i % xCount) * tileSize)
			tileY := float64((i / xCount) * tileSize)
			op.GeoM.Translate(tileX, tileY)

			// Aplica o zoom a partir do centro da tela
			op.GeoM.Scale(g.zoom, g.zoom)

			// Aplica a transformação da câmera por último
			op.GeoM.Translate(-g.cameraX, -g.cameraY)

			sx := (t % tileXCount) * tileSize
			sy := (t / tileXCount) * tileSize
			screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
		}
	}

	// Desenha o personagem
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.character.x, g.character.y)
	op.GeoM.Scale(g.zoom, g.zoom)
	op.GeoM.Translate(-g.cameraX, -g.cameraY)

	// Usa o tile 26 para representar o personagem (um personagem mais visível)
	sx := (26 % tileXCount) * tileSize
	sy := (26 / tileXCount) * tileSize
	screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)

	// Mostra informações de debug
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nCamera: (%.1f, %.1f)\nMouse: (%d, %d)\nZoom: %.2f\nCharacter: (%.1f, %.1f)",
		ebiten.ActualTPS(), g.cameraX, g.cameraY, g.mouseX, g.mouseY, g.zoom, g.character.x, g.character.y))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		layers: [][]int{
			{
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 218, 243, 244, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 244, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 219, 243, 243, 243, 219, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 244, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			},
			{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 26, 27, 28, 29, 30, 31, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 51, 52, 53, 54, 55, 56, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 76, 77, 78, 79, 80, 81, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 101, 102, 103, 104, 105, 106, 0, 0, 0, 0,

				0, 0, 0, 0, 0, 126, 127, 128, 129, 130, 131, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 303, 303, 245, 242, 303, 303, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,

				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
			},
		},
		zoom: 1.0,
		character: Character{
			x:        100,
			y:        100,
			speed:    2.0,
			isMoving: false,
		},
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Tiles (Ebitengine Demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
