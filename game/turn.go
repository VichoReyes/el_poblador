package game

// This file now only contains phase interface definitions and forward references.
// All phase implementations have been moved to separate files:
// - phase_base.go - Base utilities and phaseWithOptions
// - phase_turn.go - phaseDiceRoll, phaseIdle, rollDice function
// - phase_robber.go - phasePlaceRobber, phaseStealCard
// - phase_building.go - phaseBuilding, phaseSettlementPlacement, phaseCityPlacement
// - phase_roads.go - phaseRoadStart, phaseRoadEnd
// - phase_dev_cards.go - phasePlayDevelopmentCard, phaseMonopoly, phaseYearOfPlenty
// - phase_game_end.go - phaseGameEnd

// The Phase interface and related types remain here as they're used by the main Game struct
// All phase implementations have been moved to their respective files