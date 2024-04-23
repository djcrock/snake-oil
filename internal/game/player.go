package game

import (
	"errors"
	"math/rand/v2"
)

func (g *Game) Draw(playerId int) error {
	if g.Phase != PhaseBrew {
		return errors.New("players can only draw during the brew phase")
	}
	if playerId < 0 || playerId >= len(g.Players) {
		return errors.New("player id out of range")
	}
	player := &g.Players[playerId]
	if player.Done {
		return errors.New("player is already done for this round")
	}
	if player.Potion.IsBusted() {
		return errors.New("player is busted")
	}
	if len(player.Options) > 0 {
		return errors.New("player must first resolve their optional tiles")
	}
	if player.Potion.IsFull() {
		return errors.New("player has reached the end of the board")
	}

	drawnIngredient := player.TakeFromBag()

	return player.PlaceIngredient(drawnIngredient)
}

func (p *Player) TakeFromBag() Ingredient {
	drawIndex := rand.IntN(len(p.Bag))
	drawnIngredient := p.Bag[drawIndex]
	// Remove the drawnIngredient from the Bag by replacing it with the last
	// ingredient in the bag and reducing the length of the bag by one.
	p.Bag[drawIndex] = p.Bag[len(p.Bag)-1]
	p.Bag = p.Bag[:len(p.Bag)-1]
	return drawnIngredient
}

func (p *Player) PlaceIngredient(ingredient Ingredient) error {
	nextSpace := p.Potion.GetNextSpace()
	ingredientSpace := (nextSpace - 1) + ingredient.Value
	if ingredient.Type == Red {
		// TODO: handle ingredient variants
		orangeCount := 0
		for _, potionIngredient := range p.Potion.Ingredients {
			if potionIngredient.Type == Orange {
				orangeCount++
			}
		}
		if orangeCount >= 1 && orangeCount < 3 {
			ingredientSpace += 1
		} else if orangeCount >= 3 {
			ingredientSpace += 2
		}
	}
	if ingredient.Type == Blue {
		numOptions := min(ingredient.Value, len(p.Bag))
		p.Options = make([]Ingredient, numOptions)
		for i := 0; i < numOptions; i++ {
			p.Options[i] = p.TakeFromBag()
		}
	}
	if ingredient.Type == Yellow {
		previousIngredient := p.Potion.Ingredients[len(p.Potion.Ingredients)-1]
		if previousIngredient.Type == White {
			p.Bag = append(p.Bag, previousIngredient.Ingredient)
			p.Potion.Ingredients = p.Potion.Ingredients[:len(p.Potion.Ingredients)-1]
		}
	}
	if ingredientSpace >= len(Board) {
		ingredientSpace = len(Board) - 1
	}

	p.Potion.Ingredients = append(p.Potion.Ingredients, IngredientSpace{ingredient, ingredient.Value})
	if p.Potion.IsBusted() || p.Potion.IsFull() {
		p.Done = true
	}
	return nil
}

func (p *Player) RollDice() {
	roll := rand.IntN(6)
	switch roll {
	case 0:
		fallthrough
	case 1:
		p.Score += 1
	case 2:
		p.Score += 2
	case 3:
		p.Buys[0] = append(p.Buys[0], Ingredient{Type: Orange, Value: 1})
	case 4:
		p.Rubies++
	case 5:
		p.PotionLevel++
	}
}

func (p *Potion) IsBusted() bool {
	whiteTotal := 0
	for _, ingredient := range p.Ingredients {
		if ingredient.Type == White {
			whiteTotal += ingredient.Value
		}
	}
	return whiteTotal > BustThreshold
}

func (p *Potion) IsFull() bool {
	return p.GetNextSpace() == len(Board)-1
}

// GetNextSpace returns the index of the next space on the board after the last
// ingredient in the Potion.
func (p *Potion) GetNextSpace() int {
	if len(p.Ingredients) == 0 {
		return p.StartingSpace
	}

	return p.Ingredients[len(p.Ingredients)-1].Space + 1
}
