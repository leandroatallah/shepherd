package tilemap

import (
	"fmt"
	"log"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
)

type LayerNameID int

const (
	ObstaclesLayer LayerNameID = iota
	EnemiesLayer
	ItemsLayer
	PlayerStartLayer
	EndpointLayer
)

var LayerNameMap = map[string]LayerNameID{
	"Obstacles":   ObstaclesLayer,
	"Enemies":     EnemiesLayer,
	"Items":       ItemsLayer,
	"PlayerStart": PlayerStartLayer,
	"Endpoint":    EndpointLayer,
}

func (t *Tilemap) CreateCollisionBodies(space body.BodiesSpace, triggerEndpoint body.Touchable) {
	foundEndpoint := false
	foundObstacles := false

	for _, layer := range t.Layers {
		if !layer.Visible {
			continue
		}

		if layer.Name == "Endpoint" {
			foundEndpoint = true
			if layer.Type == "tilelayer" {
				for i, tileID := range layer.Data {
					if tileID == 0 {
						continue
					}

					// Calculate the x and y coordinates of the tile
					x := (i % layer.Width) * t.Tilewidth
					y := (i / layer.Width) * t.Tileheight

					rect := bodyphysics.NewRect(x, y, t.Tilewidth, t.Tileheight)
					obstacle := bodyphysics.NewObstacleRect(rect)
					obstacle.SetPosition(x, y)
					// Generate a unique ID for the obstacle based on its position
					obstacle.SetID(fmt.Sprintf("ENDPOINT_%d_%d", x, y))
					obstacle.AddCollisionBodies()
					obstacle.SetIsObstructive(false)
					obstacle.SetTouchable(triggerEndpoint)
					space.AddBody(obstacle)
				}
			} else {
				for _, obj := range layer.Objects {
					obstacle := t.NewObstacleRect(obj, "Endpoint", false)
					obstacle.SetTouchable(triggerEndpoint)
					space.AddBody(obstacle)
				}
			}
		}

		if layer.Name == "Obstacles" {
			foundObstacles = true
			if layer.Type == "tilelayer" {
				for i, tileID := range layer.Data {
					if tileID == 0 {
						continue
					}

					// Calculate the x and y coordinates of the tile
					x := (i % layer.Width) * t.Tilewidth
					y := (i / layer.Width) * t.Tileheight

					rect := bodyphysics.NewRect(x, y, t.Tilewidth, t.Tileheight)
					obstacle := bodyphysics.NewObstacleRect(rect)
					obstacle.SetPosition(x, y)
					// Generate a unique ID for the obstacle based on its position
					obstacle.SetID(fmt.Sprintf("OBSTACLE_%d_%d", x, y))
					obstacle.AddCollisionBodies()
					obstacle.SetIsObstructive(true)
					space.AddBody(obstacle)
				}
			} else {
				for _, obj := range layer.Objects {
					obstacle := t.NewObstacleRect(obj, "OBSTACLE", true)
					space.AddBody(obstacle)
				}
			}
		}
	}

	if !foundEndpoint {
		log.Printf("Endpoint layer not found in tilemap")
	}

	if !foundObstacles {
		log.Printf("Obstacles layer not found in tilemap")
	}
}

func (t *Tilemap) NewObstacleRect(obj *Obstacle, prefix string, isObstructive bool) *bodyphysics.ObstacleRect {
	y := int(obj.Y)

	rect := bodyphysics.NewRect(int(obj.X), y, int(obj.Width), int(obj.Height))
	o := bodyphysics.NewObstacleRect(rect)
	o.SetPosition(int(obj.X), y)
	var id string
	for _, p := range obj.Properties {
		if p.Name == "body_id" {
			id = p.Value
			break
		}
	}
	o.SetID(fmt.Sprintf("%v_%v", prefix, id))
	o.AddCollisionBodies()
	o.SetIsObstructive(isObstructive)
	return o
}
