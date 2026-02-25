package sequences

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/sequences"
)

// SequencePlayer manages the execution of a sequence.
type SequencePlayer struct {
	app.AppContextHolder

	currentSequence     sequences.Sequence
	currentCommandIndex int
	isPlaying           bool

	hasActiveCommands bool
	lastBlockingIndex int
	blockingEnded     bool
	blockedByParent   bool

	backgroundCommands []sequences.Command
}

// NewSequencePlayer creates a new player.
func NewSequencePlayer(appContext *app.AppContext) *SequencePlayer {
	ctx := app.AppContextHolder{}
	ctx.SetAppContext(appContext)
	return &SequencePlayer{AppContextHolder: ctx}
}

// Play starts executing a sequence.
func (p *SequencePlayer) Play(sequence sequences.Sequence) {
	if p.hasActiveCommands {
		return
	}

	p.currentSequence = sequence
	p.currentCommandIndex = -1
	p.hasActiveCommands = true
	p.blockingEnded = false
	p.lastBlockingIndex = -1
	p.backgroundCommands = nil

	if seq, ok := sequence.(*Sequence); ok && len(seq.blockSequenceFlags) == len(seq.commands) {
		for i, flag := range seq.blockSequenceFlags {
			if flag {
				p.lastBlockingIndex = i
			}
		}
	} else {
		p.lastBlockingIndex = len(sequence.Commands()) - 1
	}

	p.isPlaying = p.lastBlockingIndex >= 0

	if seq, ok := sequence.(*Sequence); ok && seq.BlockPlayerMovement && p.lastBlockingIndex >= 0 && !p.blockedByParent {
		if player, found := p.AppContext().ActorManager.GetPlayer(); found {
			player.BlockMovement()
		}
	}

	p.advanceToNextCommand()
}

// IsPlaying returns true if a sequence is currently being played.
func (p *SequencePlayer) IsPlaying() bool {
	return p.isPlaying
}

func (p *SequencePlayer) IsOver() bool {
	return p.currentCommandIndex >= len(p.currentSequence.Commands())
}

// Update should be called every frame. It updates the current command.
func (p *SequencePlayer) Update() {
	if !p.hasActiveCommands {
		return
	}

	// Update the currently blocking command (if any)
	if p.currentCommandIndex < len(p.currentSequence.Commands()) {
		currentCommand := p.currentSequence.Commands()[p.currentCommandIndex]
		if currentCommand.Update() {
			p.advanceToNextCommand()
		}
	}

	// Update background (non-blocking) commands
	if len(p.backgroundCommands) > 0 {
		remaining := p.backgroundCommands[:0]
		for _, cmd := range p.backgroundCommands {
			if !cmd.Update() {
				remaining = append(remaining, cmd)
			}
		}
		p.backgroundCommands = remaining
	}

	// End entire sequence only when no more commands and no more background commands
	if p.currentCommandIndex >= len(p.currentSequence.Commands()) && len(p.backgroundCommands) == 0 {
		p.endSequence()
		return
	}
}

// advanceToNextCommand moves to the next command in the queue and initializes it.
func (p *SequencePlayer) advanceToNextCommand() {
	// Keep advancing until we hit a blocking command or the end
	for {
		p.currentCommandIndex++
		if p.currentCommandIndex >= len(p.currentSequence.Commands()) {
			// No more commands in pipeline; background may still be running
			if !p.blockingEnded {
				p.endBlockingPhase()
			}
			return
		}

		if !p.blockingEnded && p.currentCommandIndex > p.lastBlockingIndex {
			p.endBlockingPhase()
		}

		nextCommand := p.currentSequence.Commands()[p.currentCommandIndex]

		// Check if this command is blocking or not
		isBlocking := true
		if seq, ok := p.currentSequence.(*Sequence); ok && p.currentCommandIndex < len(seq.blockSequenceFlags) {
			isBlocking = seq.blockSequenceFlags[p.currentCommandIndex]
		}

		nextCommand.Init(p.AppContext())
		if isBlocking {
			// Wait on this command in Update()
			return
		}

		// Non-blocking: run in background and immediately advance to try next
		p.backgroundCommands = append(p.backgroundCommands, nextCommand)
		// Loop to look for the next blocking command (or finish)
	}
}

func (p *SequencePlayer) endSequence() {
	p.hasActiveCommands = false
	if !p.blockingEnded {
		p.endBlockingPhase()
	}
}

func (p *SequencePlayer) endBlockingPhase() {
	p.blockingEnded = true
	p.isPlaying = false

	if seq, ok := p.currentSequence.(*Sequence); ok && seq.BlockPlayerMovement && !p.blockedByParent {
		if player, found := p.AppContext().ActorManager.GetPlayer(); found {
			player.UnblockMovement()
		}
	}
}
