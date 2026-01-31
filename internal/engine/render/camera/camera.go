package camera

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	"github.com/setanarut/kamera/v2"
)

type Controller struct {
	cam              *kamera.Camera
	target           body.Collidable
	followTarget     body.Body
	DeadZoneRadius   float64
	SmoothingFactor  float64
	isFollowing      bool
	centerX, centerY float64
	screenWidth      float64
	screenHeight     float64
	bounds           *image.Rectangle
}

func NewController(x, y float64) *Controller {
	cfg := config.Get()
	cam := kamera.NewCamera(x, y, float64(cfg.ScreenWidth), float64(cfg.ScreenHeight))
	cam.SmoothType = kamera.SmoothDamp
	cam.ShakeEnabled = true

	// Create a body to be the camera's direct target
	// targetBody := physics.NewPhysicsBody(physics.NewRect(0, 0, 1, 1))
	targetBody := bodyphysics.NewCollidableBodyFromRect(bodyphysics.NewRect(0, 0, 1, 1))

	return &Controller{
		cam:          cam,
		target:       targetBody,
		isFollowing:  false,
		centerX:      x,
		centerY:      y,
		screenWidth:  float64(cfg.ScreenWidth),
		screenHeight: float64(cfg.ScreenHeight),
	}
}

func NewCamera(x, y int) *kamera.Camera {
	cfg := config.Get()
	c := kamera.NewCamera(
		float64(x),
		float64(x),
		float64(cfg.ScreenWidth),
		float64(cfg.ScreenHeight),
	)
	c.SmoothType = kamera.SmoothDamp
	c.ShakeEnabled = true
	return c
}

func (c *Controller) SetFollowing(following bool) {
	c.isFollowing = following
}

// SetBounds restricts the camera movement to the specified rectangle.
func (c *Controller) SetBounds(bounds *image.Rectangle) {
	c.bounds = bounds
}

func (c *Controller) SetCenter(x, y float64) {
	c.centerX = x
	c.centerY = y
	c.Kamera().SetCenter(x, y)
}

func (c *Controller) SetFollowTarget(b body.Body) {
	c.followTarget = b
	x, y := b.GetPositionMin()
	c.Kamera().SetCenter(float64(x), float64(y))
}

func (c *Controller) Update() {
	var targetX, targetY float64
	if c.isFollowing && c.followTarget != nil {
		x, y := c.followTarget.GetPositionMin()
		// Use center of target for following
		w, h := c.followTarget.GetShape().Width(), c.followTarget.GetShape().Height()
		targetX = float64(x) + float64(w)/2
		targetY = float64(y) + float64(h)/2
	} else {
		targetX = c.centerX
		targetY = c.centerY
	}

	if c.bounds != nil {
		// Calculate viewport half-dimensions
		halfW := c.screenWidth / 2
		halfH := c.screenHeight / 2

		// Calculate min and max center positions
		minX := float64(c.bounds.Min.X) + halfW
		maxX := float64(c.bounds.Max.X) - halfW
		minY := float64(c.bounds.Min.Y) + halfH
		maxY := float64(c.bounds.Max.Y) - halfH

		// Clamp targetX
		if targetX < minX {
			targetX = minX
		}
		if targetX > maxX {
			targetX = maxX
		}

		// Clamp targetY
		if targetY < minY {
			targetY = minY
		}
		if targetY > maxY {
			targetY = maxY
		}
	}

	c.cam.LookAt(targetX, targetY)
}

func (c *Controller) Draw(
	src *ebiten.Image, options *ebiten.DrawImageOptions, dst *ebiten.Image,
) {
	c.cam.Draw(src, options, dst)
}

// Useful for debugging
func (c *Controller) Kamera() *kamera.Camera {
	return c.cam
}

func (c *Controller) Position() image.Rectangle {
	// return c.target.Position()
	return c.followTarget.Position()
}

func (c *Controller) Target() body.Body {
	// return c.target
	return c.followTarget
}
