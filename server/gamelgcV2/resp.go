package gamelgcV2

import (
	"encoding/json"

	"go.uber.org/zap"
)

// PublicResponse is public data for each player to see
type PublicResponse struct {
	WallSize        int    `json:"wallSize"`        // WallSize is remaining number of tiles in the wall
	BlockedWallSize int    `json:"blockedWallSize"` // BlockedWallSize is the size of the blocked part of wall from the back
	BonusTile       tile   `json:"bonusTile"`       // BonusTile represents a revealed bonus tile
	DiscardPile     []tile `json:"discardPile"`     // DiscardPile represents tiles that have been discarded from the game
	ActiveTile      *tile  `json:"activeTile"`      // ActiveTile shows discarded tile that players can take action on
	// PlayerState represents publicly available data for each player
	PlayerHands map[string]PublicPlayerStateResponse `json:"playerHands"`
}

type PublicPlayerStateResponse struct {
	MainHandLen int    `json:"mainHandLen"`
	ReveladHand []tile `json:"reveladHand"`
	BonusTiles  []tile `json:"bonusTiles"`
}

func (state *GameState) GetPubResp() []byte {
	resp := PublicResponse{
		WallSize:        len(state.Wall.MainWall),
		BlockedWallSize: len(state.Wall.BackWall),
		BonusTile:       state.Wall.BonusTile,
		DiscardPile:     state.DiscardPile,
		ActiveTile:      state.ActiveTile,
		PlayerHands:     make(map[string]PublicPlayerStateResponse, _PLAYER_COUNT),
	}

	for key, value := range state.PlayerStates {
		resp.PlayerHands[key] = PublicPlayerStateResponse{
			MainHandLen: len(value.MainHand),
			ReveladHand: value.RevealedHand,
			BonusTiles:  value.BonusTiles,
		}
	}

	zap.S().Debugf("Message %+v", resp)
	updatedState, err := json.Marshal(resp)
	if err != nil {
		zap.S().Errorf("Faled to marshal public state, err = %w", err)
		return nil
	}
	return updatedState
}

// PrivatePlayerStateResponse is data only for the player to see
type PrivatePlayerStateResponse PlayerState

func (state *GameState) GetPlayerResp(userid string) []byte {
	updatedState, err := json.Marshal(*state.PlayerStates[userid])
	if err != nil {
		zap.S().Errorf("Faled to marshal public state, err = %w", err)
		return nil
	}
	return updatedState
}
