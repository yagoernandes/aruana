package interfaces

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"log"

	"aruana/internal/application"
	"aruana/internal/domain"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
)

const (
	screenWidth  = 240
	screenHeight = 240
	tileSize     = 16
)

// EbitenGame implementa a interface do Ebiten
type EbitenGame struct {
	gameService *application.GameService
	tilesImage  *ebiten.Image
}

// NewEbitenGame cria uma nova instância do EbitenGame
func NewEbitenGame(game *domain.Game, gameService *application.GameService) *EbitenGame {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
	if err != nil {
		log.Fatal(err)
	}
	tilesImage := ebiten.NewImageFromImage(img)

	return &EbitenGame{
		gameService: gameService,
		tilesImage:  tilesImage,
	}
}

// Update implementa a interface do Ebiten
func (g *EbitenGame) Update() error {
	return g.gameService.Update()
}

// Draw implementa a interface do Ebiten
func (g *EbitenGame) Draw(screen *ebiten.Image) {
	w := g.tilesImage.Bounds().Dx()
	tileXCount := w / tileSize

	// Draw each tile with each DrawImage call.
	const xCount = screenWidth / tileSize
	game := g.gameService.GetGame()
	for _, l := range game.GetLayers() {
		for i, t := range l {
			op := &ebiten.DrawImageOptions{}

			// Posiciona o tile primeiro
			tileX := float64((i % xCount) * tileSize)
			tileY := float64((i / xCount) * tileSize)
			op.GeoM.Translate(tileX, tileY)

			// Aplica o zoom a partir do centro da tela
			op.GeoM.Scale(game.GetZoom(), game.GetZoom())

			// Aplica a transformação da câmera por último
			cameraX, cameraY := game.GetCameraPosition()
			op.GeoM.Translate(-cameraX, -cameraY)

			sx := (t % tileXCount) * tileSize
			sy := (t / tileXCount) * tileSize
			screen.DrawImage(g.tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
		}
	}

	// Mostra informações de debug
	mouseX, mouseY := game.GetMousePosition()
	cameraX, cameraY := game.GetCameraPosition()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nCamera: (%.1f, %.1f)\nMouse: (%d, %d)\nZoom: %.2f",
		ebiten.ActualTPS(), cameraX, cameraY, mouseX, mouseY, game.GetZoom()))
}

// Layout implementa a interface do Ebiten
func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
