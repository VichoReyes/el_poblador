package game

import (
	"el_poblador/board"
	"fmt"
	"strings"
)

// OfferStatus represents the state of a trade offer.
type OfferStatus int
type CanTakeOffer int

const (
	// OfferIsPending is an active offer waiting for responses.
	OfferIsPending OfferStatus = iota
	// OfferIsRetracted means the offerer has withdrawn the offer.
	OfferIsRetracted
	// OfferIsCompleted means the trade has been finalized and executed.
	OfferIsCompleted
)

const (
	CanTakeExcludedPlayer CanTakeOffer = iota
	CanTakeObsolete
	CanTakeNotEnoughResources
	CanTakeIsAmbiguous
	CanTakeTrue
)

// TradeOffer represents a player-to-player trade proposal.
type TradeOffer struct {
	OffererID int

	// The ID of the player the offer is directed to.
	// A value of -1 will represent an offer to anyone
	TargetID int

	// Resources the offerer will give.
	// e.g., map[ResourceType]int{ResourceWood: 1, ResourceBrick: 1}
	// Ambiguity is represented by using board.ResourceInvalid as a key.
	Offering map[board.ResourceType]int

	// Resources the offerer wants to receive.
	Requesting map[board.ResourceType]int

	// The current status of the offer (Pending, Retracted, Completed).
	Status OfferStatus
}

func (t *TradeOffer) String() string {
	var offeringParts []string

	// Check for ambiguous resources first
	if amount, ok := t.Offering[board.ResourceInvalid]; ok && amount > 0 {
		offeringParts = append(offeringParts, fmt.Sprintf("%d of something", amount))
	}

	for _, rt := range board.RESOURCE_TYPES {
		if amount, ok := t.Offering[rt]; ok && amount > 0 {
			offeringParts = append(offeringParts, fmt.Sprintf("%d %s", amount, rt))
		}
	}

	var requestingParts []string

	// Check for ambiguous resources first
	if amount, ok := t.Requesting[board.ResourceInvalid]; ok && amount > 0 {
		requestingParts = append(requestingParts, fmt.Sprintf("%d of something", amount))
	}

	for _, rt := range board.RESOURCE_TYPES {
		if amount, ok := t.Requesting[rt]; ok && amount > 0 {
			requestingParts = append(requestingParts, fmt.Sprintf("%d %s", amount, rt))
		}
	}

	offering := "nothing"
	if len(offeringParts) > 0 {
		offering = strings.Join(offeringParts, ", ")
	}

	requesting := "nothing"
	if len(requestingParts) > 0 {
		requesting = strings.Join(requestingParts, ", ")
	}

	return fmt.Sprintf("%s for %s", offering, requesting)
}

func (t *TradeOffer) canTake(playerId int, g *Game) CanTakeOffer {
	if t.Status != OfferIsPending {
		return CanTakeObsolete
	}
	if t.TargetID >= 0 && t.TargetID != playerId {
		return CanTakeExcludedPlayer
	}

	if t.Offering[board.ResourceInvalid] > 0 {
		return CanTakeIsAmbiguous
	}
	if t.Requesting[board.ResourceInvalid] > 0 {
		return CanTakeIsAmbiguous
	}

	player := &g.Players[playerId]

	for r, v := range t.Requesting {
		if player.Resources[r] < v {
			return CanTakeNotEnoughResources
		}
	}
	return CanTakeTrue
}

// executeTrade performs the actual resource exchange between two players
// and marks the offer as completed. Returns true if successful.
func (t *TradeOffer) executeTrade(acceptorID int, g *Game) bool {
	// Verify the trade can be executed
	if t.canTake(acceptorID, g) != CanTakeTrue {
		return false
	}

	offerer := &g.Players[t.OffererID]
	acceptor := &g.Players[acceptorID]

	// Transfer resources from offerer to acceptor (what offerer is giving)
	for resource, amount := range t.Offering {
		if amount > 0 {
			offerer.Resources[resource] -= amount
			acceptor.Resources[resource] += amount
		}
	}

	// Transfer resources from acceptor to offerer (what offerer is requesting)
	for resource, amount := range t.Requesting {
		if amount > 0 {
			acceptor.Resources[resource] -= amount
			offerer.Resources[resource] += amount
		}
	}

	// Mark the offer as completed
	t.Status = OfferIsCompleted

	return true
}
