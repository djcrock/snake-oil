package game

const (
	PhaseLobby      = "lobby"
	PhaseBrew       = "brew"
	PhaseEvaluation = "evaluation"
)

type Game struct {
	Players []Player
	Round   int
	Phase   string
}

type Player struct {
	Name        string
	Potion      Potion
	Bag         []Ingredient
	Options     []Ingredient
	Buys        [][]Ingredient
	Score       int
	Rubies      int
	PotionLevel int
	RatTails    int
	Flask       bool
	Done        bool
}

type Potion struct {
	StartingSpace int
	Ingredients   []IngredientSpace
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

type IngredientSpace struct {
	Ingredient
	Space int
}

type Space = struct {
	Value         int
	VictoryPoints int
	Ruby          bool
}
