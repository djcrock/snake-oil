package game

import (
	"errors"
	"fmt"
)

const NumRounds = 9
const MaxPlayers = 4

var GameBoard = []Space{
	{Value: 0, VictoryPoints: 0},
	{Value: 1, VictoryPoints: 0},
	{Value: 2, VictoryPoints: 0},
	{Value: 3, VictoryPoints: 0},
	{Value: 4, VictoryPoints: 0},
	{Value: 5, VictoryPoints: 0, Ruby: true},
	{Value: 6, VictoryPoints: 1},
	{Value: 7, VictoryPoints: 1},
	{Value: 8, VictoryPoints: 1},
	{Value: 9, VictoryPoints: 1, Ruby: true},
	{Value: 10, VictoryPoints: 2},
	{Value: 11, VictoryPoints: 2},
	{Value: 12, VictoryPoints: 2},
	{Value: 13, VictoryPoints: 2, Ruby: true},
	{Value: 14, VictoryPoints: 3},
	{Value: 15, VictoryPoints: 3},
	{Value: 15, VictoryPoints: 3, Ruby: true},
	{Value: 16, VictoryPoints: 3},
	{Value: 16, VictoryPoints: 4},
	{Value: 17, VictoryPoints: 4},
	{Value: 17, VictoryPoints: 4, Ruby: true},
	{Value: 18, VictoryPoints: 4},
	{Value: 18, VictoryPoints: 5},
	{Value: 19, VictoryPoints: 5},
	{Value: 19, VictoryPoints: 5, Ruby: true},
	{Value: 20, VictoryPoints: 5},
	{Value: 20, VictoryPoints: 6},
	{Value: 21, VictoryPoints: 6},
	{Value: 21, VictoryPoints: 6, Ruby: true},
	{Value: 22, VictoryPoints: 7},
	{Value: 22, VictoryPoints: 7, Ruby: true},
	{Value: 23, VictoryPoints: 7},
	{Value: 23, VictoryPoints: 8},
	{Value: 24, VictoryPoints: 8},
	{Value: 24, VictoryPoints: 8, Ruby: true},
	{Value: 25, VictoryPoints: 9},
	{Value: 25, VictoryPoints: 9, Ruby: true},
	{Value: 26, VictoryPoints: 9},
	{Value: 27, VictoryPoints: 10},
	{Value: 27, VictoryPoints: 10, Ruby: true},
	{Value: 28, VictoryPoints: 11},
	{Value: 28, VictoryPoints: 11, Ruby: true},
	{Value: 29, VictoryPoints: 11},
	{Value: 29, VictoryPoints: 12},
	{Value: 30, VictoryPoints: 12},
	{Value: 30, VictoryPoints: 12, Ruby: true},
	{Value: 31, VictoryPoints: 12},
	{Value: 31, VictoryPoints: 13},
	{Value: 32, VictoryPoints: 13},
	{Value: 32, VictoryPoints: 13, Ruby: true},
	{Value: 33, VictoryPoints: 14},
	{Value: 33, VictoryPoints: 14, Ruby: true},
	{Value: 35, VictoryPoints: 15},
}

var Prices = []IngredientPrice{
	// Orange
	{Ingredient: Ingredient{Type: Orange, Value: 1}, Price: 3},
	// Green
	{Ingredient: Ingredient{Type: Green, Value: 1}, Price: 4},
	{Ingredient: Ingredient{Type: Green, Value: 2}, Price: 8},
	{Ingredient: Ingredient{Type: Green, Value: 4}, Price: 14},
	// Red
	{Ingredient: Ingredient{Type: Red, Value: 1}, Price: 6},
	{Ingredient: Ingredient{Type: Red, Value: 2}, Price: 10},
	{Ingredient: Ingredient{Type: Red, Value: 4}, Price: 16},
	// Blue
	{Ingredient: Ingredient{Type: Blue, Value: 1}, Price: 5},
	{Ingredient: Ingredient{Type: Blue, Value: 2}, Price: 10},
	{Ingredient: Ingredient{Type: Blue, Value: 4}, Price: 19},
	// Black
	{Ingredient: Ingredient{Type: Black, Value: 1}, Price: 10},
	// Yellow
	{Ingredient: Ingredient{Type: Yellow, Value: 1}, Price: 8},
	{Ingredient: Ingredient{Type: Yellow, Value: 2}, Price: 12},
	{Ingredient: Ingredient{Type: Yellow, Value: 4}, Price: 18},
	// Purple
	{Ingredient: Ingredient{Type: Purple, Value: 1}, Price: 9},
}

// RatTailLocations lists all victory point values that are preceded by a rat tail.
var RatTailLocations = []int{2, 5, 8, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29, 31, 33, 35, 37, 39, 41, 43, 45, 47, 49}

func (g *Game) IsFull() bool {
	return len(g.Players) >= MaxPlayers
}

func (g *Game) AddPlayer(name string) error {
	if g.Phase != PhaseLobby {
		return fmt.Errorf("players can only be added during the lobby phase")
	}
	if g.IsFull() {
		return errors.New("game is full")
	}

	player := Player{
		Name:  name,
		Bag:   []Ingredient{},
		Buys:  make([][]Ingredient, NumRounds),
		Flask: true,
	}
	player.Buys[0] = []Ingredient{
		{White, 1},
		{White, 1},
		{White, 1},
		{White, 1},
		{White, 2},
		{White, 2},
		{White, 3},
		{Orange, 1},
		{Green, 1},
	}
	g.Players = append(g.Players, player)

	return nil
}

func (g *Game) Start() error {
	if g.Phase != PhaseLobby {
		return errors.New("game cannot be started unless in the lobby phase")
	}
	if len(g.Players) < 2 {
		return errors.New("not enough players")
	}

	g.Round = 1
	return g.StartNormalBrew()
}

func (g *Game) StartNormalBrew() error {
	g.Phase = PhaseBrew
	for i := range g.Players {
		g.Players[i].Bag = nil
		for _, buys := range g.Players[i].Buys {
			g.Players[i].Bag = append(g.Players[i].Bag, buys...)
		}
	}

	return nil
}

func (g *Game) GetRatTails() []int {
	ratTails := make([]int, len(g.Players))
	maxScore := 0
	for i := range g.Players {
		if g.Players[i].Score > maxScore {
			maxScore = g.Players[i].Score
		}
	}
	for i := range g.Players {
		if g.Players[i].Score < maxScore {
			for _, ratTailLocation := range RatTailLocations {
				if maxScore < ratTailLocation {
					break
				}
				if g.Players[i].Score < ratTailLocation {
					ratTails[i]++
				}
			}
		}
	}
	return ratTails
}
