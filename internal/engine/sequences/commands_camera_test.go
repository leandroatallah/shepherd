package sequences

import (
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/mocks"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
	"github.com/leandroatallah/firefly/internal/engine/render/camera"
	"github.com/leandroatallah/firefly/internal/engine/scene"
)

type mockSceneWithCamera struct {
	mocks.MockScene
	cam *camera.Controller
}

func (m *mockSceneWithCamera) Camera() *camera.Controller {
	return m.cam
}

func TestCameraSetTargetCommand_SmoothTransition(t *testing.T) {
	// 1. Setup
	appContext := &app.AppContext{}
	sceneManager := scene.NewSceneManager()
	sceneManager.SetAppContext(appContext)
	appContext.SceneManager = sceneManager
	appContext.Space = space.NewSpace()

	cam := camera.NewController(0, 0)
	mockScene := &mockSceneWithCamera{cam: cam}
	mockScene.SetAppContext(appContext)
	appContext.SceneManager.SwitchTo(mockScene)

	target := &mocks.MockActor{Id: "test_target"}
	target.SetPosition(100, 100)
	appContext.Space.AddBody(target)

	cmd := &CameraSetTargetCommand{
		TargetID: "test_target",
		Duration: 60,
	}

	// 2. Init
	cmd.Init(appContext)

	if cmd.camera == nil {
		t.Fatal("camera should not be nil after Init")
	}
	if cmd.target == nil {
		t.Fatal("target should not be nil after Init")
	}
	if cmd.camera.IsFollowing() {
		t.Fatal("camera should not be following during smooth transition")
	}

	startX, startY := cmd.camera.Kamera().Center()
	if startX != 0 || startY != 0 {
		t.Fatalf("expected camera to start at (0,0), got (%f, %f)", startX, startY)
	}

	// 3. Update (simulate frames)
	finished := false
	for i := 0; i < 60; i++ {
		finished = cmd.Update()
		if finished {
			break
		}
	}

	if !finished {
		t.Fatal("command should have finished after 60 frames")
	}

	// 4. Verify final state
	if !cmd.camera.IsFollowing() {
		t.Fatal("camera should be following after transition is complete")
	}

	finalX, finalY := cmd.camera.Kamera().Center()
	targetX, targetY := target.GetPositionMin()
	w, h := target.Width(), target.Height()
	expectedX := float64(targetX) + float64(w)/2
	expectedY := float64(targetY) + float64(h)/2

	// Use a small tolerance for floating point comparison
	if finalX < expectedX-0.1 || finalX > expectedX+0.1 || finalY < expectedY-0.1 || finalY > expectedY+0.1 {
		t.Fatalf("expected camera to be centered on target (%f, %f), got (%f, %f)", expectedX, expectedY, finalX, finalY)
	}
}

func TestCameraSetTargetCommand_InstantTransition(t *testing.T) {
	// 1. Setup
	appContext := &app.AppContext{}
	sceneManager := scene.NewSceneManager()
	sceneManager.SetAppContext(appContext)
	appContext.SceneManager = sceneManager
	appContext.Space = space.NewSpace()

	cam := camera.NewController(0, 0)
	mockScene := &mockSceneWithCamera{cam: cam}
	mockScene.SetAppContext(appContext)
	appContext.SceneManager.SwitchTo(mockScene)

	target := &mocks.MockActor{Id: "test_target"}
	target.SetPosition(100, 100)
	appContext.Space.AddBody(target)

	cmd := &CameraSetTargetCommand{
		TargetID: "test_target",
		Duration: 0, // Instant
	}

	// 2. Init
	cmd.Init(appContext)

	if !cmd.camera.IsFollowing() {
		t.Fatal("camera should be following immediately on instant transition")
	}

	// 3. Update
	finished := cmd.Update()
	if !finished {
		t.Fatal("instant command should finish in one update")
	}

	// 4. Verify final state
	finalX, finalY := cmd.camera.Kamera().Center()
	targetX, targetY := target.GetPositionMin()
	w, h := target.Width(), target.Height()
	expectedX := float64(targetX) + float64(w)/2
	expectedY := float64(targetY) + float64(h)/2

	if finalX != expectedX || finalY != expectedY {
		t.Fatalf("expected camera to be centered on target (%f, %f), got (%f, %f)", expectedX, expectedY, finalX, finalY)
	}
}
