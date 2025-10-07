package game

import (
	"el_poblador/board"
	"fmt"
	"math/rand/v2"
	"strings"
)

type phasePlaceRobber struct {
	game         *Game
	tileCoord    board.TileCoord
	continuation Phase
	invalid      string
}

func PhasePlaceRobber(game *Game, continuation Phase) Phase {
	return &phasePlaceRobber{
		game:         game,
		tileCoord:    game.Board.ValidTileCoord(),
		continuation: continuation,
	}
}

func (p *phasePlaceRobber) BoardCursor() interface{} {
	return p.tileCoord
}

func (p *phasePlaceRobber) MoveCursor(direction string) {
	dest, ok := moveTileCursor(p.tileCoord, direction)
	if !ok {
		return
	}
	p.tileCoord = dest
}

func (p *phasePlaceRobber) HelpText() string {
	if p.invalid != "" {
		return p.invalid
	}
	return "Select a tile to place the robber on"
}

func (p *phasePlaceRobber) Confirm() Phase {
	// Check if trying to place robber on the same tile it's already on
	if p.tileCoord == p.game.Board.GetRobber() {
		p.invalid = "Robber cannot be moved to the same tile it's already on"
		return p
	}

	playerIds := p.game.Board.PlaceRobber(p.tileCoord)

	currentPlayer := &p.game.Players[p.game.PlayerTurn]
	p.game.LogAction(fmt.Sprintf("%s moved the robber", currentPlayer.RenderName()))

	var stealablePlayers []Player
	for _, playerId := range playerIds {
		p := p.game.Players[playerId]
		if p.TotalResources() > 0 {
			stealablePlayers = append(stealablePlayers, p)
		}
	}
	if len(stealablePlayers) == 0 { // no one to steal from? skip
		return p.continuation
	}
	return &phaseStealCard{
		game:             p.game,
		continuation:     p.continuation,
		stealablePlayers: stealablePlayers,
	}
}

type phaseStealCard struct {
	game             *Game
	continuation     Phase
	stealablePlayers []Player
	selected         int
}

func (p *phaseStealCard) BoardCursor() interface{} {
	return nil
}

func (p *phaseStealCard) MoveCursor(direction string) {
	switch direction {
	case "up":
		p.selected--
	case "down":
		p.selected++
	}
	p.selected = (p.selected + len(p.stealablePlayers)) % len(p.stealablePlayers)
}

func (p *phaseStealCard) HelpText() string {
	return "Select a player to steal from"
}

func (p *phaseStealCard) Confirm() Phase {
	player := p.stealablePlayers[p.selected]
	var resourcePool []board.ResourceType
	for resType, count := range player.Resources {
		for i := 0; i < count; i++ {
			resourcePool = append(resourcePool, resType)
		}
	}
	if len(resourcePool) > 0 {
		selectedResource := resourcePool[rand.IntN(len(resourcePool))]
		player.Resources[selectedResource] -= 1
		p.game.Players[p.game.PlayerTurn].AddResource(selectedResource)

		currentPlayer := &p.game.Players[p.game.PlayerTurn]
		p.game.LogAction(fmt.Sprintf("%s stole a card from %s", currentPlayer.RenderName(), player.RenderName()))
	}
	return p.continuation
}

func (p *phaseStealCard) Menu() string {
	var paddedOptions []string
	for i, player := range p.stealablePlayers {
		if i == p.selected {
			paddedOptions = append(paddedOptions, "> "+player.Render(player.Name))
		} else {
			paddedOptions = append(paddedOptions, player.Render(" "+player.Name))
		}
	}
	return strings.Join(paddedOptions, "\n")
}