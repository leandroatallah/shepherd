package sequences

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/app"
	contractseq "github.com/leandroatallah/firefly/internal/engine/contracts/sequences"
)

type testCommand struct {
	initCalled    bool
	updateCount   int
	completeAfter int
}

func (c *testCommand) Init(appContext any) {
	c.initCalled = true
}

func (c *testCommand) Update() bool {
	c.updateCount++
	return c.updateCount >= c.completeAfter
}

type testSequence struct {
	commands []contractseq.Command
}

func (s *testSequence) Commands() []contractseq.Command {
	return s.commands
}

func TestSequencePlayerPlaysBlockingCommandsToCompletion(t *testing.T) {
	ctx := &app.AppContext{}

	player := NewSequencePlayer(ctx)

	cmd1 := &testCommand{completeAfter: 1}
	cmd2 := &testCommand{completeAfter: 1}

	seq := &testSequence{commands: []contractseq.Command{cmd1, cmd2}}

	if player.IsPlaying() {
		t.Fatalf("expected player to be idle before Play")
	}

	player.Play(seq)

	if !cmd1.initCalled {
		t.Fatalf("expected first command Init to be called on Play")
	}
	if cmd2.initCalled {
		t.Fatalf("expected second command Init not to be called yet")
	}
	if !player.IsPlaying() {
		t.Fatalf("expected player to be playing after Play")
	}

	player.Update()

	if cmd1.updateCount != 1 {
		t.Fatalf("expected first command Update to be called once, got %d", cmd1.updateCount)
	}
	if !cmd2.initCalled {
		t.Fatalf("expected second command Init to be called after first completes")
	}
	if !player.IsPlaying() {
		t.Fatalf("expected player to still be playing while second command runs")
	}

	player.Update()

	if cmd2.updateCount == 0 {
		t.Fatalf("expected second command Update to be called")
	}
	if player.IsPlaying() {
		t.Fatalf("expected player to stop playing after all commands complete")
	}
	if !player.IsOver() {
		t.Fatalf("expected player to report IsOver after completion")
	}
}

func TestCallSequenceCommandLoadsAndPlaysNestedSequence(t *testing.T) {
	// Create a temporary nested sequence file
	tmpDir := t.TempDir()
	nestedSeqPath := filepath.Join(tmpDir, "nested.json")
	nestedSeqContent := `{
		"commands": [
			{ "command": "delay", "frames": 2 }
		],
		"block_player_movement": false
	}`
	if err := os.WriteFile(nestedSeqPath, []byte(nestedSeqContent), 0644); err != nil {
		t.Fatalf("failed to create nested sequence file: %v", err)
	}

	ctx := &app.AppContext{}
	cmd := &CallSequenceCommand{Path: nestedSeqPath}

	cmd.Init(ctx)

	if cmd.isComplete {
		t.Fatalf("expected command not to be complete immediately after Init")
	}
	if cmd.nestedSequence == nil {
		t.Fatalf("expected nested sequence to be loaded")
	}
	if cmd.sequencePlayer == nil {
		t.Fatalf("expected sequence player to be created")
	}
	if !cmd.sequencePlayer.IsPlaying() {
		t.Fatalf("expected sequence player to be playing after Init")
	}

	// Update until complete
	for i := 0; i < 10 && !cmd.isComplete; i++ {
		cmd.Update()
	}

	if !cmd.isComplete {
		t.Fatalf("expected command to complete after nested sequence finishes")
	}
}

func TestCallSequenceCommandHandlesInvalidPath(t *testing.T) {
	ctx := &app.AppContext{}
	cmd := &CallSequenceCommand{Path: "nonexistent.json"}

	cmd.Init(ctx)

	if !cmd.isComplete {
		t.Fatalf("expected command to be complete immediately when path is invalid")
	}
}
