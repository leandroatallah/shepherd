package tilemaplayer

type TilemapDimensionsProvider interface {
	GetTilemapWidth() int
	GetTilemapHeight() int
}
