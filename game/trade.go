package game

import "el_poblador/board"

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
	return "some trade offer"
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
