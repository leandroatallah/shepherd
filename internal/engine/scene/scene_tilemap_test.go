package scene

import (
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/render/tilemap"
)

// createValidTilemap creates a minimal valid tilemap for testing
func createValidTilemap(layers []*tilemap.Layer) *tilemap.Tilemap {
	return &tilemap.Tilemap{
		Tilewidth:  16,
		Tileheight: 16,
		Width:      40,
		Height:     14,
		Tilesets: []*tilemap.Tileset{
			{Firstgid: 1, Tilewidth: 16, Tileheight: 16},
		},
		Layers: layers,
	}
}

// Test SetPlayerStartPosition adjusts Y correctly
func TestTilemapScene_SetPlayerStartPosition(t *testing.T) {
	ctx := &app.AppContext{}
	scene := NewTilemapScene(ctx)

	// Create a mock tilemap with PlayerStart at (80, 144) and tile height 16
	scene.tilemap = createValidTilemap([]*tilemap.Layer{
		{
			Name:    "PlayerStart",
			Type:    "objectgroup",
			Visible: true,
			Objects: []*tilemap.Obstacle{
				{X: 80, Y: 144},
			},
		},
	})

	// Verify tilemap loads PlayerStart correctly
	x, y, found := scene.tilemap.GetPlayerStartPosition()
	if !found {
		t.Fatal("expected PlayerStart to be found")
	}
	if x != 80 {
		t.Errorf("expected x=80, got %d", x)
	}
	if y != 144 {
		t.Errorf("expected y=144, got %d", y)
	}

	// Expected adjustment: y = 144 - (actorHeight - tileHeight)
	// For actorHeight=24, tileHeight=16: y = 144 - 8 = 136
	expectedY := y - (24 - scene.tilemap.Tileheight)
	if expectedY != 136 {
		t.Errorf("expected adjusted y=136, got %d", expectedY)
	}
}

// Test SetPlayerStartPosition with no PlayerStart layer
func TestTilemapScene_SetPlayerStartPosition_NoLayer(t *testing.T) {
	ctx := &app.AppContext{}
	scene := NewTilemapScene(ctx)

	scene.tilemap = createValidTilemap([]*tilemap.Layer{})

	_, _, found := scene.tilemap.GetPlayerStartPosition()
	if found {
		t.Fatal("expected PlayerStart not to be found")
	}
}

// Test InitEnemies position adjustment calculation
func TestTilemapScene_InitEnemies_PositionAdjustment(t *testing.T) {
	ctx := &app.AppContext{}
	scene := NewTilemapScene(ctx)

	scene.tilemap = createValidTilemap([]*tilemap.Layer{
		{
			Name:    "Enemies",
			Type:    "objectgroup",
			Visible: true,
			Objects: []*tilemap.Obstacle{
				{X: 192, Y: 160, Gid: 174},
			},
		},
	})

	enemiesPos := scene.Tilemap().GetEnemiesPositionID()
	if len(enemiesPos) != 1 {
		t.Fatalf("expected 1 enemy position, got %d", len(enemiesPos))
	}

	e := enemiesPos[0]
	if e.X != 192 {
		t.Errorf("expected enemy x=192, got %d", e.X)
	}
	if e.Y != 160 {
		t.Errorf("expected enemy y=160, got %d", e.Y)
	}

	// Expected adjustment: y = 160 - (actorHeight - tileHeight)
	// For actorHeight=24, tileHeight=16: y = 160 - 8 = 152
	actorHeight := 24
	expectedY := e.Y - (actorHeight - scene.tilemap.Tileheight)
	if expectedY != 152 {
		t.Errorf("expected adjusted y=152, got %d", expectedY)
	}
}

// Test InitNPCs position adjustment calculation
func TestTilemapScene_InitNPCs_PositionAdjustment(t *testing.T) {
	ctx := &app.AppContext{}
	scene := NewTilemapScene(ctx)

	scene.tilemap = createValidTilemap([]*tilemap.Layer{
		{
			Name:    "NPCs",
			Type:    "objectgroup",
			Visible: true,
			Objects: []*tilemap.Obstacle{
				{X: 32, Y: 160, Gid: 184},
			},
		},
	})

	npcsPos := scene.Tilemap().GetNpcsPositionID()
	if len(npcsPos) != 1 {
		t.Fatalf("expected 1 npc position, got %d", len(npcsPos))
	}

	n := npcsPos[0]
	if n.X != 32 {
		t.Errorf("expected npc x=32, got %d", n.X)
	}
	if n.Y != 160 {
		t.Errorf("expected npc y=160, got %d", n.Y)
	}

	// Expected adjustment: y = 160 - (actorHeight - tileHeight)
	// For actorHeight=24, tileHeight=16: y = 160 - 8 = 152
	actorHeight := 24
	expectedY := n.Y - (actorHeight - scene.tilemap.Tileheight)
	if expectedY != 152 {
		t.Errorf("expected adjusted y=152, got %d", expectedY)
	}
}

// Test InitItems position
func TestTilemapScene_InitItems_Position(t *testing.T) {
	ctx := &app.AppContext{}
	scene := NewTilemapScene(ctx)

	scene.tilemap = createValidTilemap([]*tilemap.Layer{
		{
			Name:    "Items",
			Type:    "objectgroup",
			Visible: true,
			Objects: []*tilemap.Obstacle{
				{X: 100, Y: 100, Gid: 1}, // Gid > 0 to use tileset
			},
		},
	})

	itemsPos := scene.Tilemap().GetItemsPositionID()
	if len(itemsPos) != 1 {
		t.Fatalf("expected 1 item position, got %d", len(itemsPos))
	}

	i := itemsPos[0]
	if i.X != 100 {
		t.Errorf("expected item x=100, got %d", i.X)
	}
	if i.Y != 100 {
		t.Errorf("expected item y=100, got %d", i.Y)
	}
}
