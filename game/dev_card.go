package game

import "math/rand/v2"

type DevCard string

const (
	DevCardKnight       DevCard = "Knight"
	DevCardRoadBuilding DevCard = "Road Building"
	DevCardMonopoly     DevCard = "Monopoly"
	DevCardYearOfPlenty DevCard = "Year of Plenty"
	DevCardVictoryPoint DevCard = "Victory Point"
)

func shuffleDevCards() []DevCard {
	unshuffled := []DevCard{}

	for i := 0; i < 14; i++ {
		unshuffled = append(unshuffled, DevCardKnight)
	}
	for i := 0; i < 2; i++ {
		unshuffled = append(unshuffled, DevCardRoadBuilding)
	}
	for i := 0; i < 2; i++ {
		unshuffled = append(unshuffled, DevCardMonopoly)
	}
	for i := 0; i < 2; i++ {
		unshuffled = append(unshuffled, DevCardYearOfPlenty)
	}
	for i := 0; i < 5; i++ {
		unshuffled = append(unshuffled, DevCardVictoryPoint)
	}

	rand.Shuffle(len(unshuffled), func(i, j int) {
		unshuffled[i], unshuffled[j] = unshuffled[j], unshuffled[i]
	})

	return unshuffled
}
