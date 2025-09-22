package game

import (
	"el_poblador/board"
)

type phaseBuilding struct {
	phaseWithOptions
	previousPhase Phase
}

func PhaseBuilding(game *Game, previousPhase Phase) Phase {
	player := &game.players[game.playerTurn]

	// Build the list of available building options
	var options []string
	strikethrough := strikethroughStyle()

	if player.CanBuildRoad() {
		options = append(options, "Road")
	} else {
		options = append(options, strikethrough.Render("Road"))
	}

	if player.CanBuildSettlement() {
		options = append(options, "Settlement")
	} else {
		options = append(options, strikethrough.Render("Settlement"))
	}

	if player.CanBuildCity() {
		options = append(options, "City")
	} else {
		options = append(options, strikethrough.Render("City"))
	}

	if player.CanBuyDevelopmentCard() {
		options = append(options, "Development Card")
	} else {
		options = append(options, strikethrough.Render("Development Card"))
	}

	options = append(options, "Cancel (or 'esc')")

	return &phaseBuilding{
		phaseWithOptions: phaseWithOptions{
			game:    game,
			options: options,
		},
		previousPhase: previousPhase,
	}
}

func (p *phaseBuilding) Confirm() Phase {
	player := &p.game.players[p.game.playerTurn]

	switch p.selected {
	case 0: // Road
		if player.CanBuildRoad() {
			return PhaseRoadStart(p.game, p)
		}
		return p
	case 1: // Settlement
		if player.CanBuildSettlement() {
			return PhaseSettlementPlacement(p.game, p)
		}
		return p
	case 2: // City
		if player.CanBuildCity() {
			return PhaseCityPlacement(p.game, p)
		}
		return p
	case 3: // Development Card
		if player.CanBuyDevelopmentCard() {
			if player.BuyDevelopmentCard() {
				if card := p.game.DrawDevelopmentCard(); card != nil {
					player.hiddenDevCards = append(player.hiddenDevCards, *card)
					
					// Check for game end after buying development card (in case it's a victory point card)
					if winner := p.game.CheckGameEnd(); winner != nil {
						return PhaseGameEnd(p.game, winner)
					}
					
					return p.previousPhase
				}
			}
		}
		return p
	case 4: // Cancel
		return p.previousPhase
	default:
		panic("Invalid option selected")
	}
}

func (p *phaseBuilding) Cancel() Phase {
	return p.previousPhase
}

func (p *phaseBuilding) HelpText() string {
	return "Choose what to build"
}

type phaseSettlementPlacement struct {
	game          *Game
	cursorCross   board.CrossCoord
	previousPhase Phase
	invalid       string
}

func PhaseSettlementPlacement(game *Game, previousPhase Phase) Phase {
	cursorCross := game.board.ValidCrossCoord()
	return &phaseSettlementPlacement{
		game:          game,
		cursorCross:   cursorCross,
		previousPhase: previousPhase,
	}
}

func (p *phaseSettlementPlacement) Confirm() Phase {
	player := &p.game.players[p.game.playerTurn]
	playerId := p.game.playerTurn

	if !p.game.board.CanPlaceSettlementForPlayer(p.cursorCross, playerId) {
		p.invalid = "Can't build settlement here"
		return p
	}

	if !player.BuildSettlement() {
		p.invalid = "Not enough resources"
		return p
	}

	p.game.board.SetSettlement(p.cursorCross, playerId)

	// Check for game end after building settlement
	if winner := p.game.CheckGameEnd(); winner != nil {
		return PhaseGameEnd(p.game, winner)
	}

	return PhaseIdleWithNotification(p.game, "Settlement built!")
}

func (p *phaseSettlementPlacement) Cancel() Phase {
	return p.previousPhase
}

func (p *phaseSettlementPlacement) HelpText() string {
	if p.invalid != "" {
		return p.invalid
	}
	return "Select where to place your settlement"
}

func (p *phaseSettlementPlacement) BoardCursor() interface{} {
	return p.cursorCross
}

func (p *phaseSettlementPlacement) MoveCursor(direction string) {
	dest, ok := moveCrossCursor(p.cursorCross, direction)
	if !ok {
		return
	}
	p.cursorCross = dest
}

type phaseCityPlacement struct {
	game          *Game
	cursorCross   board.CrossCoord
	previousPhase Phase
	invalid       string
}

func PhaseCityPlacement(game *Game, previousPhase Phase) Phase {
	cursorCross := game.board.ValidCrossCoord()
	return &phaseCityPlacement{
		game:          game,
		cursorCross:   cursorCross,
		previousPhase: previousPhase,
	}
}

func (p *phaseCityPlacement) Confirm() Phase {
	player := &p.game.players[p.game.playerTurn]
	playerId := p.game.playerTurn

	if !p.game.board.CanUpgradeToCity(p.cursorCross, playerId) {
		p.invalid = "Can't upgrade to city here"
		return p
	}

	if !player.BuildCity() {
		p.invalid = "Not enough resources"
		return p
	}

	p.game.board.UpgradeToCity(p.cursorCross, playerId)

	// Check for game end after building city
	if winner := p.game.CheckGameEnd(); winner != nil {
		return PhaseGameEnd(p.game, winner)
	}

	return PhaseIdleWithNotification(p.game, "City built!")
}

func (p *phaseCityPlacement) Cancel() Phase {
	return p.previousPhase
}

func (p *phaseCityPlacement) HelpText() string {
	if p.invalid != "" {
		return p.invalid
	}
	return "Select a settlement to upgrade to a city"
}

func (p *phaseCityPlacement) BoardCursor() interface{} {
	return p.cursorCross
}

func (p *phaseCityPlacement) MoveCursor(direction string) {
	dest, ok := moveCrossCursor(p.cursorCross, direction)
	if !ok {
		return
	}
	p.cursorCross = dest
}