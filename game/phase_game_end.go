package game

import (
	"fmt"
	"strings"
)

type phaseGameEnd struct {
	game   *Game
	winner *Player
}

func PhaseGameEnd(game *Game, winner *Player) Phase {
	return &phaseGameEnd{
		game:   game,
		winner: winner,
	}
}

func (p *phaseGameEnd) Confirm() Phase {
	// Game is over, stay in this phase
	return p
}

func (p *phaseGameEnd) MoveCursor(direction string) {
	// No cursor movement in game end phase
}

func (p *phaseGameEnd) BoardCursor() interface{} {
	return nil
}

func (p *phaseGameEnd) Menu() string {
	lines := []string{
		fmt.Sprintf("🎉 GAME OVER! 🎉"),
		"",
		fmt.Sprintf("%s WINS!", p.winner.Render(p.winner.Name)),
		"",
		"Final Scores:",
	}

	for _, player := range p.game.Players {
		points := player.VictoryPoints(p.game)
		marker := "  "
		if &player == p.winner {
			marker = "🏆"
		}
		lines = append(lines, fmt.Sprintf("%s %s: %d points", marker, player.Render(player.Name), points))
	}

	return strings.Join(lines, "\n")
}

func (p *phaseGameEnd) HelpText() string {
	return fmt.Sprintf("Game complete! %s reached 10 victory points!", p.winner.Render(p.winner.Name))
}