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
	commands      []contractseq.Command
	interruptible bool
	oneTime       bool
	path          string
}

func (s *testSequence) Commands() []contractseq.Command {
	return s.commands
}

func (s *testSequence) Interruptible() bool {
	return s.interruptible
}

func (s *testSequence) OneTime() bool {
	return s.oneTime
}

func (s *testSequence) GetPath() string {
	if s.path != "" {
		return s.path
	}
	return "test_sequence"
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

func TestSequencePlayerBlockedByParentPreventsUnblock(t *testing.T) {
	ctx := &app.AppContext{}

	// Create sequence with block_player_movement: true
	seqContent := `{
		"commands": [
			{ "command": "delay", "frames": 1 }
		],
		"block_player_movement": true
	}`
	tmpDir := t.TempDir()
	seqPath := filepath.Join(tmpDir, "blocking.json")
	if err := os.WriteFile(seqPath, []byte(seqContent), 0644); err != nil {
		t.Fatalf("failed to create sequence file: %v", err)
	}

	seq, err := NewSequenceFromJSON(seqPath)
	if err != nil {
		t.Fatalf("failed to load sequence: %v", err)
	}

	// Create player and mark as blocked by parent
	player := NewSequencePlayer(ctx)
	player.blockedByParent = true

	// Play sequence - should NOT block/unblock because already blocked by parent
	player.Play(seq)

	// Run sequence to completion
	for i := 0; i < 10 && player.IsPlaying(); i++ {
		player.Update()
	}

	// Verify sequence ended
	if player.IsPlaying() {
		t.Fatalf("expected sequence to complete")
	}
	// blockedByParent prevents UnblockMovement from being called
}

func TestSequencePlayerInterruptibleAndOneTime(t *testing.T) {
	ctx := &app.AppContext{}
	player := NewSequencePlayer(ctx)

	// Create interruptible sequence (default)
	interruptibleSeq := &testSequence{
		commands:      []contractseq.Command{&testCommand{completeAfter: 5}},
		interruptible: true,
		path:          "interruptible.json",
	}

	// Create non-interruptible sequence
	nonInterruptibleSeq := &testSequence{
		commands:      []contractseq.Command{&testCommand{completeAfter: 10}},
		interruptible: false,
		path:          "non_interruptible.json",
	}

	// Test 1: Same sequence requested while playing - should not restart
	player.Play(interruptibleSeq)
	if !player.IsPlaying() {
		t.Fatalf("expected player to be playing")
	}
	initialIndex := player.currentCommandIndex
	player.Play(interruptibleSeq)
	if player.currentCommandIndex != initialIndex {
		t.Fatalf("expected same sequence to not restart")
	}

	// Test 2: Different interruptible sequence should interrupt current
	anotherSeq := &testSequence{
		commands:      []contractseq.Command{&testCommand{completeAfter: 2}},
		interruptible: true,
		path:          "another.json",
	}
	player.Play(anotherSeq)
	if !player.IsPlaying() {
		t.Fatalf("expected player to still be playing after interrupt")
	}
	// Let it complete
	for i := 0; i < 10 && player.IsPlaying(); i++ {
		player.Update()
	}

	// Test 3: Non-interruptible sequence should block other sequences
	player.Play(nonInterruptibleSeq)
	if !player.IsPlaying() {
		t.Fatalf("expected player to be playing non-interruptible sequence")
	}

	// Try to interrupt with another sequence - should fail
	interruptSeq := &testSequence{
		commands:      []contractseq.Command{&testCommand{completeAfter: 1}},
		interruptible: true,
		path:          "interrupt.json",
	}
	player.Play(interruptSeq)
	// Should still be running nonInterruptibleSeq, not interruptSeq
	if player.currentSequencePath != "non_interruptible.json" {
		t.Fatalf("expected non-interruptible sequence to continue")
	}

	// Let non-interruptible complete
	for i := 0; i < 15 && player.IsPlaying(); i++ {
		player.Update()
	}
	if player.IsPlaying() {
		t.Fatalf("expected non-interruptible sequence to complete")
	}
}

func TestSequencePlayerOneTimeSequence(t *testing.T) {
	ctx := &app.AppContext{}
	player := NewSequencePlayer(ctx)

	// Create a one-time sequence
	oneTimeSeq := &testSequence{
		commands:      []contractseq.Command{&testCommand{completeAfter: 5}},
		oneTime:       true,
		path:          "one_time.json",
	}

	// Play the sequence first time - should work
	player.Play(oneTimeSeq)
	if !player.IsPlaying() {
		t.Fatalf("expected player to be playing first time")
	}

	// Try to play again BEFORE completion - should be ignored
	player.Play(oneTimeSeq)
	if !player.IsPlaying() {
		t.Fatalf("expected player to still be playing")
	}
	// Verify it didn't restart (command index should not reset)
	initialIndex := player.currentCommandIndex
	player.Play(oneTimeSeq)
	if player.currentCommandIndex != initialIndex {
		t.Fatalf("expected one-time sequence to not restart during playback")
	}

	// Let it complete
	for i := 0; i < 10 && player.IsPlaying(); i++ {
		player.Update()
	}

	if player.IsPlaying() {
		t.Fatalf("expected sequence to complete")
	}

	// Try to play again after completion - should be ignored because it's one-time
	player.Play(oneTimeSeq)
	if player.IsPlaying() {
		t.Fatalf("expected one-time sequence to be ignored on second play after completion")
	}

	// Different sequence should still work
	anotherSeq := &testSequence{
		commands:      []contractseq.Command{&testCommand{completeAfter: 1}},
		oneTime:       false,
		path:          "another.json",
	}
	player.Play(anotherSeq)
	if !player.IsPlaying() {
		t.Fatalf("expected different sequence to play after one-time completed")
	}
}

