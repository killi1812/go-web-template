package gamelgcV2

import (
	"slices"
	"template/model"
	"time"

	"go.uber.org/zap"
)

func newGameStateV2(players []*model.User) GameState {
	state := GameState{
		// TODO: initialize  PlayerOrderIds
		PlayerStates:   make(map[string]*PlayerState, 4),
		DiscardPile:    make([]tile, 0, _WALL_SIZE),
		PlayerOrderIds: make([]string, 0, _PLAYER_COUNT),
		ActiveTile:     nil,
	}
	for _, p := range players {
		state.PlayerOrderIds = append(state.PlayerOrderIds, p.ID)
		state.PlayerStates[p.ID] = &PlayerState{
			user:         p,
			MainHand:     make(hand, 0, _HAND_SIZE),
			BonusTiles:   make([]tile, 0),
			RevealedHand: make([]tile, 0),
			// TODO: init wind value
		}
	}

	// Start of the game procedure
	// 1. create a new wall
	// 2. shuffel the wall
	// 3. break the wall and revial the bonus tile
	state.Wall = breakWall(shuffe(newWall()))

	state.deal()
	state.handleBonusTiles()

	return state
}

// GameState holds information about the state of the game
type GameState struct {
	//DealerId       string                  // DealerId is an id of the Player who is a dealer for this game
	BreakPlayerId  string
	PlayerOnTurnId string                  // PlayerOnTurnId is an id of the Player whose turn it is
	PlayerOrderIds []string                // PlayerOrderId is order of players
	PlayerStates   map[string]*PlayerState // PlayerStates is a storage for each players state

	Wall        wall   `json:"wall"`        // Wall represents a drawing pile an array of tiles
	DiscardPile []tile `json:"discardPile"` // DiscardPile represents tiles that have been discarded from the game
	ActiveTile  *tile  `json:"activeTile"`  // ActiveTile shows discarded tile that players can take action on

	discardCallback *time.Timer // Callback to remove tile  form ActiveTile after some time
	seazePower      seizeType   // seazePower is the power of active seaze
}

// deal deals player hands wall need to bee inicialized in this state of the game
func (state *GameState) deal() {
	// TODO: shift based on dealer

	// 1. deal player hands 4 rounds of 4 tiles each for 12 piece size hand
	for range 3 {
		for _, p := range state.PlayerOrderIds {
			if state.PlayerStates[p] == nil {
				zap.S().Errorf("Player wiht id: %s is missing", p)
				continue
			}
			state.PlayerStates[p].MainHand = append(state.PlayerStates[p].MainHand, state.Wall.draw(4)...)
		}
	}

	// 2. then each player draws an additional tile to complete a 13 piece starting hand
	for _, p := range state.PlayerOrderIds {
		if state.PlayerStates[p] == nil {
			zap.S().Errorf("Player wiht id: %s is missing", p)
			continue
		}
		state.PlayerStates[p].MainHand = append(state.PlayerStates[p].MainHand, state.Wall.draw()...)
	}
}

func (state *GameState) handleBonusTiles() {
	for _, p := range state.PlayerOrderIds {
		for {
			bonusIndex := slices.IndexFunc(state.PlayerStates[p].MainHand, func(t tile) bool {
				return t.Category == _CATEGORY_BONUS
			})
			if bonusIndex == -1 {
				break
			}
			state.PlayerStates[p].BonusTiles = append(state.PlayerStates[p].BonusTiles, *state.PlayerStates[p].discard(state.PlayerStates[p].MainHand[bonusIndex].TileId))
			state.PlayerStates[p].MainHand = append(state.PlayerStates[p].MainHand, state.Wall.drawBack())
		}
	}
}
