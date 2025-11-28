package game

import "el_poblador/board"

// OfferStatus represents the state of a trade offer.
type OfferStatus int

const (
	// OfferIsPending is an active offer waiting for responses.
	OfferIsPending OfferStatus = iota
	// OfferIsRetracted means the offerer has withdrawn the offer.
	OfferIsRetracted
	// OfferIsCompleted means the trade has been finalized and executed.
	OfferIsCompleted
)

// TradeOffer represents a player-to-player trade proposal.
type TradeOffer struct {
	// A unique identifier for this specific offer.
	ID int

	// The ID of the player making the offer.
	OffererID int

	// The ID of the player the offer is directed to.
	// A value of -1 will represent an offer to all players (broadcast).
	TargetID int

	// Resources the offerer will give.
	// e.g., map[ResourceType]int{ResourceWood: 1, ResourceBrick: 1}
	// Ambiguity is represented by using board.ResourceInvalid as a key.
	Offering map[board.ResourceType]int

	// Resources the offerer wants to receive.
	Requesting map[board.ResourceType]int

	// The current status of the offer (Pending, Retracted, Completed).
	Status OfferStatus

	// For counter-offers or acceptances, this links back to the original offer's ID.
	// A value of 0 indicates a new, unsolicited offer.
	InReplyToID int
}
