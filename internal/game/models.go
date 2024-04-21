package game

const (
	PhaseLobby = "lobby"
	PhaseBrew  = "brew"
)

type Game struct {
	Players []Player
	Round   int
	Phase   string
}

type Player struct {
	Name   string
	Potion []Ingredient
	Bag    []Ingredient
	Buys   [][]Ingredient
	Score  int
	Rubies int
	Flask  bool
}

type Potion struct {
	StartingSpace int
	Ingredients   []Ingredient
}

const (
	White  = 'w'
	Orange = 'o'
	Green  = 'g'
	Red    = 'r'
	Blue   = 'u'
	Black  = 'b'
	Yellow = 'y'
	Purple = 'p'
)

type Ingredient struct {
	Type  byte
	Value int
}

type IngredientPrice struct {
	Ingredient
	Price int
}

type Space = struct {
	Value         int
	VictoryPoints int
	Ruby          bool
}
