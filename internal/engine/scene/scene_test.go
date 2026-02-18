package scene

import (
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/navigation"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
)

type testScene struct {
	app.AppContextHolder

	drawCount   int
	updateCount int
	startCount  int
	finishCount int
}

func (s *testScene) Draw(screen *ebiten.Image) {
	s.drawCount++
}

func (s *testScene) Update() error {
	s.updateCount++
	return nil
}

func (s *testScene) OnStart() {
	s.startCount++
}

func (s *testScene) OnFinish() {
	s.finishCount++
}

// TODO: AppContextHolder should provide SetAppContext method
func (s *testScene) SetAppContext(ctx any) {
	s.AppContextHolder.SetAppContext(ctx)
}

type testSceneFactory struct {
	scenes map[navigation.SceneType]*testScene
	appCtx any
}

func newTestSceneFactory() *testSceneFactory {
	return &testSceneFactory{
		scenes: make(map[navigation.SceneType]*testScene),
	}
}

func (f *testSceneFactory) Create(sceneType navigation.SceneType, freshInstance bool) (navigation.Scene, error) {
	if !freshInstance {
		if s, ok := f.scenes[sceneType]; ok {
			return s, nil
		}
	}
	s := &testScene{}
	s.SetAppContext(f.appCtx)
	f.scenes[sceneType] = s
	return s, nil
}

// TODO: AppContextHolder should provide SetAppContext method
func (f *testSceneFactory) SetAppContext(ctx any) {
	f.appCtx = ctx
}

type testTransition struct {
	startCalled   bool
	updateCount   int
	drawCount     int
	startCallback func()
}

func (t *testTransition) Update() {
	t.updateCount++
}

func (t *testTransition) Draw(screen *ebiten.Image) {
	t.drawCount++
}

func (t *testTransition) StartTransition(cb func()) {
	t.startCalled = true
	t.startCallback = cb
}

func (t *testTransition) EndTransition(cb func()) {
	cb()
}

func (t *testTransition) Complete() {
	if t.startCallback != nil {
		t.startCallback()
	}
}

func TestSceneManagerNavigateWithTransition(t *testing.T) {
	ctx := &app.AppContext{}
	manager := NewSceneManager()
	manager.SetAppContext(ctx)

	factory := newTestSceneFactory()
	factory.SetAppContext(ctx)
	manager.SetFactory(factory)

	sceneType := navigation.SceneType(1)
	transition := &testTransition{}

	if err := manager.Update(); err != nil {
		t.Fatalf("initial Update error: %v", err)
	}

	manager.NavigateTo(sceneType, transition, false)

	created := factory.scenes[sceneType]
	if created == nil {
		t.Fatalf("expected factory to create scene")
	}
	if !transition.startCalled {
		t.Fatalf("expected transition StartTransition to be called")
	}
	if created.startCount != 0 {
		t.Fatalf("scene should not start before transition completes")
	}

	if err := manager.Update(); err != nil {
		t.Fatalf("Update error while transition active: %v", err)
	}
	if transition.updateCount == 0 {
		t.Fatalf("expected transition Update to be called during manager.Update")
	}

	transition.Complete()

	if err := manager.Update(); err != nil {
		t.Fatalf("Update error after transition completion: %v", err)
	}

	if created.startCount == 0 {
		t.Fatalf("expected scene OnStart to be called after transition completion")
	}

	screen := ebiten.NewImage(1, 1)
	manager.Draw(screen)
	if created.drawCount == 0 {
		t.Fatalf("expected scene Draw to be called after switch")
	}
}

func TestSceneManagerNavigateWithoutTransition(t *testing.T) {
	ctx := &app.AppContext{}
	manager := NewSceneManager()
	manager.SetAppContext(ctx)

	factory := newTestSceneFactory()
	factory.SetAppContext(ctx)
	manager.SetFactory(factory)

	sceneType1 := navigation.SceneType(1)
	sceneType2 := navigation.SceneType(2)

	manager.NavigateTo(sceneType1, nil, false)
	s1 := factory.scenes[sceneType1]
	if s1 == nil {
		t.Fatalf("expected first scene to be created")
	}
	if s1.startCount != 1 {
		t.Fatalf("expected first scene OnStart to be called once, got %d", s1.startCount)
	}

	manager.NavigateTo(sceneType2, nil, false)
	s2 := factory.scenes[sceneType2]
	if s2 == nil {
		t.Fatalf("expected second scene to be created")
	}
	if s1.finishCount != 1 {
		t.Fatalf("expected first scene OnFinish to be called when switching, got %d", s1.finishCount)
	}
	if s2.startCount != 1 {
		t.Fatalf("expected second scene OnStart to be called once, got %d", s2.startCount)
	}
}

func TestSceneManagerNavigateBack(t *testing.T) {
	ctx := &app.AppContext{}
	manager := NewSceneManager()
	manager.SetAppContext(ctx)

	factory := newTestSceneFactory()
	factory.SetAppContext(ctx)
	manager.SetFactory(factory)

	sceneType1 := navigation.SceneType(1)
	sceneType2 := navigation.SceneType(2)

	manager.NavigateTo(sceneType1, nil, false)
	s1 := factory.scenes[sceneType1]

	manager.NavigateTo(sceneType2, nil, false)
	s2 := factory.scenes[sceneType2]

	manager.NavigateBack(nil)

	if s2.finishCount == 0 {
		t.Fatalf("expected second scene OnFinish to be called when navigating back")
	}
	if s1.startCount < 2 {
		t.Fatalf("expected first scene OnStart to be called again when navigated back, got %d", s1.startCount)
	}
}

func TestBaseSceneOnStartClearsPhysicsSpace(t *testing.T) {
	sp := space.NewSpace()
	rect := bodyphysics.NewRect(0, 0, 10, 10)
	obstacle := bodyphysics.NewObstacleRect(rect)
	obstacle.SetID("obstacle")
	sp.AddBody(obstacle)

	ctx := &app.AppContext{
		Space: sp,
	}

	base := NewScene()
	base.SetAppContext(ctx)

	if len(sp.Bodies()) == 0 {
		t.Fatalf("precondition: expected space to contain bodies before OnStart")
	}

	base.OnStart()

	if len(sp.Bodies()) != 0 {
		t.Fatalf("expected BaseScene.OnStart to clear physics space")
	}
}

func TestBaseSceneScheduleExecutesActionAfterDelay(t *testing.T) {
	ctx := &app.AppContext{}
	base := NewScene()
	base.SetAppContext(ctx)

	called := 0
	delay := time.Second

	base.Schedule(delay, func() {
		called++
	})

	frames := timing.FromDuration(delay)
	for i := 0; i < frames; i++ {
		ctx.FrameCount++
		if err := base.Update(); err != nil {
			t.Fatalf("Update error at frame %d: %v", i, err)
		}
	}

	if called != 1 {
		t.Fatalf("expected scheduled action to be called once after delay, got %d", called)
	}

	ctx.FrameCount++
	if err := base.Update(); err != nil {
		t.Fatalf("Update error after delay: %v", err)
	}
	if called != 1 {
		t.Fatalf("expected scheduled action not to be called more than once, got %d", called)
	}
}
