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

func TestGame_GetRatTails(t *testing.T) {
	g := Game{Phase: PhaseBrew, Round: 1, Players: []Player{
		{Score: 40},
		{Score: 4},
	}}

	ratTails := g.GetRatTails()
	if len(ratTails) != 2 {
		t.Fatalf("expected 2 ratTails, got %v", len(ratTails))
	}
	if ratTails[0] != 0 {
		t.Fatalf("expected 0 ratTails for player 1, got %v", ratTails[0])
	}
	if ratTails[1] != 17 {
		t.Fatalf("expected 17 ratTails for player 2, got %v", ratTails[1])
	}
}
