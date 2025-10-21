// Package gamelgc provides all game logic and handling
package gamelgcV2

import (
	"encoding/json"
	"slices"
	"sync"
	"template/model"
	"template/util/ws"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Game holds information about the game and clients
type Game struct {
	GameId     uuid.UUID `json:"gameId"` // GameId is a uuid representing a unique game
	State      GameState
	Mutex      sync.RWMutex // GameMutex is used to lock game action until it finishes but Spectators can still get current game status
	Hub        ws.Hub
	Spectators []*model.User // Spectators represents users spectating the game
	Players    []*model.User // Players represents players and clients clients in the game key is player.ID
	handlers   map[gameAction]handlerFunc
}

// HandleMsg implements ws.MsgHandler.
func (g *Game) HandleMsg(data []byte) {
	var msg GameMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		zap.S().Errorf("error unmarshalling message: %v", err)
		return
	}
	zap.S().Debugf("Message: %+v", msg)

	if g.State.PlayerOnTurnId != msg.UserId {
		zap.S().Debugf("Player id: %s trying to take action not on his turn, turn of player %s", msg.UserId, g.State.PlayerOnTurnId)
	}

	g.Mutex.Lock()
	if fn, ok := g.handlers[msg.Action]; ok {
		err := fn(msg)
		if err != nil {
			zap.S().Errorf("Error while handling action %s, error: %w", msg.Action, err)
			// TODO: check for some error response
			return
		}
	} else {
		zap.S().Errorf("Unknown gameAction %s", msg.Action)
		return
	}
	g.Mutex.Unlock()

	if resp := g.State.GetPlayerResp(msg.UserId); resp != nil {
		g.Hub.Clients[msg.UserId].Send <- resp
	}

	if resp := g.State.GetPubResp(); resp != nil {
		g.Hub.Broadcast <- resp
	}

}

// Update send all available information to the client
// implements ws.MsgHandler.
func (g *Game) Update(client *ws.Client) {
	// Broadcast the message to client
	isPlayer := slices.ContainsFunc(g.Players, func(p *model.User) bool {
		return p.ID == client.UserId
	})
	// TODO: Work for basic update on join

	if isPlayer {
		zap.S().Debugf("Client clinet is a player userId: %s", client.UserId)
		client.Send <- g.State.GetPlayerResp(client.UserId)
	}
	zap.S().Debugf("Updating client userId: %s", client.UserId)
	client.Send <- g.State.GetPubResp()
}

func NewGameV2() *Game {
	// NOTE: somehow transfer players for the lobby to the game

	var game = Game{
		GameId:   uuid.New(),
		Hub:      ws.NewHub(),
		Players:  make([]*model.User, _PLAYER_COUNT),
		handlers: make(map[gameAction]handlerFunc),
	}

	game.Players[0] = &model.User{ID: "279700487249592320", Username: ".f.c."}
	game.Players[1] = &model.User{ID: "453991398648315914", Username: "liimun"}
	game.Players[2] = &model.User{ID: "289425583232909313", Username: "Savoric"}
	game.Players[3] = &model.User{ID: "168360634843398145", Username: "divald"}

	game.State = newGameStateV2(game.Players)

	game.Hub.Handler = &game

	// Register handlers
	game.handlers[_ACTION_DRAW_TILE] = game.handleDrawTile
	game.handlers[_ACTION_DISCARD_TILE] = game.handleDiscardTile
	game.handlers[_ACTION_SEIZE] = game.handleSeize
	game.handlers[_ACTION_PASS_TURN] = game.handlePassTurn

	return &game
}
