package gamelgcV2

import (
	"time"

	"go.uber.org/zap"
)

type handlerFunc func(GameMessage) error

func (g *Game) handleSeize(msg GameMessage) error {
	if g.State.ActiveTile == nil {
		zap.S().Debugf("ActiveTile is nil")
		return nil
	}
	tiles, stype := g.State.PlayerStates[msg.UserId].CheckSeize(g.State.ActiveTile)
	if tiles == nil {
		zap.S().Debugf("No available seize can be preformed")
	}

	switch stype {
	case _SEIZE_CHOW:
		if g.State.seazePower > stype {
			return nil
		}
		// TODO: player next on turn
		if false {
			zap.S().Debugf("Player %s can't preforme a seize, player on turn: %s, turn order: %v", msg.UserId, g.State.PlayerOnTurnId, g.State.PlayerOrderIds)
			return nil
		} else {
			zap.S().Debugf("Player seazing with chow")
		}

	case _SEIZE_PUNG:
		if g.State.seazePower > stype {
			return nil
		}
		zap.S().Debugf("Player %s seazing with pung", msg.UserId)

	case _SEIZE_KUNG:
		if g.State.seazePower > stype {
			return nil
		}
		zap.S().Debugf("Player %s seazing with kung", msg.UserId)
	}
	g.State.discardCallback.Stop()
	g.State.discardCallback = time.AfterFunc(_RESPONSE_TIME, func() {
		if g.State.ActiveTile == nil {
			zap.S().Debugf("ActiveTile is nil")
			return
		}
		for _, t := range tiles {
			tmpTile := g.State.PlayerStates[msg.UserId].discard(t.TileId)
			g.State.PlayerStates[msg.UserId].RevealedHand = append(g.State.PlayerStates[msg.UserId].RevealedHand, *tmpTile)
		}
		g.State.PlayerStates[msg.UserId].RevealedHand = append(g.State.PlayerStates[msg.UserId].RevealedHand, *g.State.ActiveTile)
		g.State.ActiveTile = nil
		g.State.seazePower = 0

		g.Hub.Clients[msg.UserId].Send <- g.State.GetPlayerResp(msg.UserId)
		g.Hub.Broadcast <- g.State.GetPubResp()
	})
	return nil
}

// handleDiscardTile handles game action for discarding tile by moving it to ActiveTile and then after the timer moving it to DiscardPile
func (g *Game) handleDiscardTile(msg GameMessage) error {
	if g.State.ActiveTile != nil {
		if tmp := g.State.discardCallback; tmp != nil {
			zap.S().Debugf("Stoppin old discardCallback")
			tmp.Stop()
		}

		g.State.DiscardPile = append(g.State.DiscardPile, *g.State.ActiveTile)
		zap.S().Debugf("Tile %+v is added to DiscardPile", *g.State.ActiveTile)
		g.State.ActiveTile = nil
	}

	g.State.ActiveTile = g.State.PlayerStates[msg.UserId].discard(*msg.TileId)

	g.State.discardCallback = time.AfterFunc(_RESPONSE_TIME, func() {
		if g.State.ActiveTile != nil {
			g.State.DiscardPile = append(g.State.DiscardPile, *g.State.ActiveTile)
			zap.S().Debugf("Tile %+v is added to DiscardPile", *g.State.ActiveTile)
			g.State.ActiveTile = nil
		} else {
			zap.S().Debugf("ActiveTile is nil")
		}
		g.Hub.Broadcast <- g.State.GetPubResp()
	})
	return nil
}

func (g *Game) handleDrawTile(msg GameMessage) error {
	drawnTiles := g.State.Wall.draw()
	if drawnTiles != nil {
		drawnTile := drawnTiles[0]
		if drawnTile.Category != _CATEGORY_BONUS {
			g.State.PlayerStates[msg.UserId].MainHand = append(g.State.PlayerStates[msg.UserId].MainHand, drawnTile)
		} else {
			zap.S().Debugf("Player %s has drawn a bonus tile adding it to bonus hand", msg.UserId)
			g.State.PlayerStates[msg.UserId].BonusTiles = append(g.State.PlayerStates[msg.UserId].BonusTiles, drawnTile)
			g.handleDrawTile(msg)
		}
	}
	return nil
}

func (g *Game) handlePassTurn(msg GameMessage) error {
	return nil
}
