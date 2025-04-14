package domain

// Game representa a entidade principal do jogo
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
}

// NewGame cria uma nova instância do jogo
func NewGame(layers [][]int) *Game {
	return &Game{
		layers: layers,
		zoom:   1.0,
	}
}

// GetLayers retorna as camadas do jogo
func (g *Game) GetLayers() [][]int {
	return g.layers
}

// GetCameraPosition retorna a posição da câmera
func (g *Game) GetCameraPosition() (float64, float64) {
	return g.cameraX, g.cameraY
}

// GetZoom retorna o nível de zoom atual
func (g *Game) GetZoom() float64 {
	return g.zoom
}

// SetCameraPosition define a posição da câmera
func (g *Game) SetCameraPosition(x, y float64) {
	g.cameraX = x
	g.cameraY = y
}

// SetZoom define o nível de zoom
func (g *Game) SetZoom(zoom float64) {
	g.zoom = zoom
}

// SetMousePosition define a posição do mouse
func (g *Game) SetMousePosition(x, y int) {
	g.mouseX = x
	g.mouseY = y
}

// GetMousePosition retorna a posição atual do mouse
func (g *Game) GetMousePosition() (int, int) {
	return g.mouseX, g.mouseY
}

// IsDragging retorna se o jogo está em estado de arrasto
func (g *Game) IsDragging() bool {
	return g.isDragging
}

// SetDragging define o estado de arrasto
func (g *Game) SetDragging(dragging bool) {
	g.isDragging = dragging
}

// GetLastMousePosition retorna a última posição do mouse
func (g *Game) GetLastMousePosition() (int, int) {
	return g.lastMouseX, g.lastMouseY
}

// SetLastMousePosition define a última posição do mouse
func (g *Game) SetLastMousePosition(x, y int) {
	g.lastMouseX = x
	g.lastMouseY = y
}
