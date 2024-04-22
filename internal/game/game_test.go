package game

import (
	"strconv"
	"testing"
)

func TestGame_AddPlayer(t *testing.T) {
	t.Run("can only add players during lobby phase", func(t *testing.T) {
		g := Game{
			Phase: PhaseLobby,
		}
		err := g.AddPlayer("player1")
		if err != nil {
			t.Fatal(err)
		}
		err = g.AddPlayer("player2")
		if err != nil {
			t.Fatal(err)
		}
		err = g.Start()
		if err != nil {
			t.Fatal(err)
		}
		err = g.AddPlayer("player3")
		if err == nil {
			t.Fatal("expected error when adding player to already started game")
		}
	})
	t.Run("can add players up to maximum", func(t *testing.T) {
		g := Game{
			Phase: PhaseLobby,
		}
		for i := range MaxPlayers {
			if g.IsFull() {
				t.Fatalf("game should not be full")
			}
			err := g.AddPlayer("player" + strconv.Itoa(i))
			if err != nil {
				t.Fatalf("player%d not added: %v", i, err)
			}
		}
		if !g.IsFull() {
			t.Fatalf("game should be full")
		}
		err := g.AddPlayer("playerN")
		if err == nil {
			t.Fatal("expected error when adding to full game")
		}
	})
}

func TestGame_Start(t *testing.T) {
	g := Game{
		Phase: PhaseLobby,
	}
	err := g.Start()
	if err == nil {
		t.Fatal("expected error when starting empty game")
	}

	err = g.AddPlayer("player1")
	if err != nil {
		t.Fatal(err)
	}
	err = g.Start()
	if err == nil {
		t.Fatal("expected error when starting game with one player")
	}

	err = g.AddPlayer("player2")
	if err != nil {
		t.Fatal(err)
	}
	err = g.Start()
	if err != nil {
		t.Fatal(err)
	}

	err = g.Start()
	if err == nil {
		t.Fatal("expected error when starting already started game")
	}
}

func TestGame_StartNormalBrew(t *testing.T) {
	t.Run("rat tails are updated", func(t *testing.T) {
		g := Game{Phase: PhaseLobby, Round: 1, Players: []Player{
			{Score: 40},
			{Score: 4},
		}}

		err := g.StartNormalBrew()
		if err != nil {
			t.Fatal(err)
		}
		if g.Players[0].RatTails != 0 {
			t.Fatalf("expected 0 ratTails for player 1, got %v", g.Players[0].RatTails)
		}
		if g.Players[1].RatTails != 17 {
			t.Fatalf("expected 17 ratTails for player 2, got %v", g.Players[1].RatTails)
		}
	})
}

func TestGame_Draw(t *testing.T) {
	t.Run("ingredients are moved from Bag to Potion", func(t *testing.T) {
		g := Game{Phase: PhaseLobby}
		err := g.AddPlayer("player1")
		if err != nil {
			t.Fatal(err)
		}
		err = g.AddPlayer("player2")
		if err != nil {
			t.Fatal(err)
		}
		err = g.Start()
		if err != nil {
			t.Fatal(err)
		}
		if len(g.Players[0].Potion.Ingredients) != 0 {
			t.Fatalf("expected potion to start with no ingredients, got %d", g.Players[0].Potion)
		}
		startBagLen := len(g.Players[0].Bag)
		err = g.Draw(0)
		if err != nil {
			t.Fatal(err)
		}
		if len(g.Players[0].Potion.Ingredients) != 1 {
			t.Fatalf("expected 1 potion ingredient for player1, got %v", g.Players[0].Potion)
		}
		if len(g.Players[0].Bag) != startBagLen-1 {
			t.Fatalf("expected potion bag for player1 to contain %d ingredients, got %d", startBagLen-1, len(g.Players[0].Bag))
		}
	})
	tests := []struct {
		name        string
		potion      Potion
		bag         []Ingredient
		expectError bool
	}{
		{
			"cannot draw if busted",
			Potion{Ingredients: []IngredientSpace{
				{Ingredient: Ingredient{Type: White, Value: 3}},
				{Ingredient: Ingredient{Type: White, Value: 5}},
			}},
			[]Ingredient{{Type: Orange, Value: 1}},
			true,
		},
		{
			"can draw if not busted",
			Potion{Ingredients: []IngredientSpace{
				{Ingredient: Ingredient{Type: White, Value: 3}},
				{Ingredient: Ingredient{Type: Green, Value: 5}},
			}},
			[]Ingredient{{Type: Orange, Value: 1}},
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := Game{Phase: PhaseLobby}
			err := g.AddPlayer("player1")
			if err != nil {
				t.Fatal(err)
			}
			err = g.AddPlayer("player2")
			if err != nil {
				t.Fatal(err)
			}
			err = g.Start()
			if err != nil {
				t.Fatal(err)
			}
			g.Players[0].Potion = test.potion
			err = g.Draw(0)
			if test.expectError && err == nil {
				t.Fatal("expected error when trying to draw, but got none")
			} else if !test.expectError && err != nil {
				t.Fatal(err)
			}
		})
	}
}
