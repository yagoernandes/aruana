package application

import (
	"aruana/internal/domain"

	"github.com/hajimehoshi/ebiten/v2"
)

// GameService gerencia a lógica do jogo
type GameService struct {
	game *domain.Game
}

// NewGameService cria uma nova instância do GameService
func NewGameService(game *domain.Game) *GameService {
	return &GameService{
		game: game,
	}
}

// Update atualiza o estado do jogo
func (s *GameService) Update() error {
	// Atualiza a posição do mouse
	mouseX, mouseY := ebiten.CursorPosition()
	s.game.SetMousePosition(mouseX, mouseY)

	// Verifica se o botão esquerdo do mouse está pressionado
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !s.game.IsDragging() {
			s.game.SetDragging(true)
			s.game.SetLastMousePosition(mouseX, mouseY)
		} else {
			// Calcula o deslocamento do mouse
			lastX, lastY := s.game.GetLastMousePosition()
			dx := lastX - mouseX
			dy := lastY - mouseY

			// Move a câmera
			cameraX, cameraY := s.game.GetCameraPosition()
			s.game.SetCameraPosition(cameraX+float64(dx), cameraY+float64(dy))

			// Atualiza a última posição do mouse
			s.game.SetLastMousePosition(mouseX, mouseY)
		}
	} else {
		s.game.SetDragging(false)
	}

	// Detecta o scroll do mouse para zoom
	_, dy := ebiten.Wheel()
	if dy != 0 {
		// Calcula o fator de zoom (0.1 é a sensibilidade do zoom)
		zoomFactor := 1.0 + dy*0.1

		// Limita o zoom entre 0.1 e 4.0
		currentZoom := s.game.GetZoom()
		newZoom := currentZoom * zoomFactor
		if newZoom >= 0.1 && newZoom <= 4.0 {
			// Ajusta a posição da câmera para manter o ponto sob o mouse
			cameraX, cameraY := s.game.GetCameraPosition()
			mouseWorldX := (float64(mouseX) + cameraX) / currentZoom
			mouseWorldY := (float64(mouseY) + cameraY) / currentZoom

			s.game.SetZoom(newZoom)

			// Recalcula a posição da câmera para manter o ponto sob o mouse
			s.game.SetCameraPosition(
				mouseWorldX*newZoom-float64(mouseX),
				mouseWorldY*newZoom-float64(mouseY),
			)
		}
	}

	return nil
}

// GetGame retorna a instância do jogo
func (s *GameService) GetGame() *domain.Game {
	return s.game
}
