package controller

import (
	"errors"
	"net/http"
	"template/app"
	"template/gamelgcV2"
	"template/lobbylgc"
	"template/util/ws"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const _MAX_GAME_CNT = 10
const _MAX_LOBBY_CNT = 4

// GameCnt provides functions for WebSocket communication.
type GameCnt struct {
	logger  *zap.SugaredLogger
	games   map[string]*gamelgcV2.Game // games map of games with id discordSdkId
	lobbies map[string]*lobbylgc.Lobby // lobbies is a map of lobbies with id discordSdkId
}

// NewGameCnt creates a new controller for the game.
func NewGameCnt() app.Controller {
	var controller *GameCnt

	app.Invoke(func(logger *zap.SugaredLogger) {
		controller = &GameCnt{
			logger:  logger,
			games:   make(map[string]*gamelgcV2.Game, _MAX_GAME_CNT),
			lobbies: make(map[string]*lobbylgc.Lobby, _MAX_LOBBY_CNT),
		}
	})

	return controller
}

// RegisterEndpoints registers the WebSocket endpoint.
func (cnt *GameCnt) RegisterEndpoints(router *gin.RouterGroup) {
	router.GET("/ws/lobby/:userId/:clientId", cnt.serveLobbyWs)
	router.GET("/ws/game/:userId/:clientId", cnt.serveGameWs)
}

// serveGameWs godoc
//
//	@Summary		web socket for lobby
//	@Description	Web socket for lobby providing data for players
//	@Tags			lobby
//	@Failure		500
//	@Router			/ws/game/{userId}/{clientId} [get]
func (cnt *GameCnt) serveGameWs(c *gin.Context) {
	gameId := c.Param("clientId")
	if gameId == "" {
		cnt.logger.Errorf("clientId parameter is empty")
		c.AbortWithError(http.StatusBadRequest, errors.New("id parameter is required"))
		return
	}

	userId := c.Param("userId")
	if userId == "" {
		cnt.logger.Errorf("clientId parameter is empty")
		c.AbortWithError(http.StatusBadRequest, errors.New("id parameter is required"))
		return
	}

	conn, err := ws.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		cnt.logger.Errorf("Failed to upgrade connection: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// crate new game from lobby
	if cnt.games[gameId] == nil {
		zap.S().Debugf("Creating a new game with id %s", gameId)
		cnt.games[gameId] = gamelgcV2.NewGameV2()
		go cnt.games[gameId].Hub.Run()
	}

	// if lobby exists add new connection
	ws.NewClient(&cnt.games[gameId].Hub, conn, userId, func() {
		zap.S().Debugf("Empty game with id %s is removed", cnt.games[gameId].GameId)
		delete(cnt.lobbies, gameId)
		zap.S().Infof("Lobby count is %d", len(cnt.lobbies))
	})
	zap.S().Infof("Game count is %d", len(cnt.games))
}

// serveLobbyWs godoc
//
//	@Summary		web socket for lobby
//	@Description	Web socket for lobby providing data for players
//	@Tags			lobby
//	@Produce		json
//	@Success		200	{object}	lobbylgc.LobbyState
//	@Failure		500
//	@Router			/ws/lobby/{userId}/{clientId} [get]
func (cnt *GameCnt) serveLobbyWs(c *gin.Context) {

	// TODO: implement checking for lobby size

	clientId := c.Param("clientId")
	if clientId == "" {
		cnt.logger.Errorf("clientId parameter is empty")
		c.AbortWithError(http.StatusBadRequest, errors.New("id parameter is required"))
		return
	}

	userId := c.Param("userId")
	if userId == "" {
		cnt.logger.Errorf("userId parameter is empty")
		c.AbortWithError(http.StatusBadRequest, errors.New("id parameter is required"))
		return
	}

	conn, err := ws.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		cnt.logger.Errorf("Failed to upgrade connection: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// crate new lobby
	if cnt.lobbies[clientId] == nil {
		cnt.lobbies[clientId] = lobbylgc.NewLobby(clientId)
	}

	// if lobby exists add new connection
	ws.NewClient(&cnt.lobbies[clientId].Hub, conn, userId, func() {
		zap.S().Debugf("Empty lobby with id %s is removed", clientId)
		delete(cnt.lobbies, clientId)
		zap.S().Infof("Lobby count is %d", len(cnt.lobbies))
	})
	zap.S().Infof("Lobby count is %d", len(cnt.lobbies))
}
