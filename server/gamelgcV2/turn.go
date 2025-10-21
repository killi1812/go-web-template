package gamelgcV2

import "template/model"

type turn struct {
	playerOrderIds   []string
	startingPlayerId string
	playerOnTurnId   string
	players          []*model.User
}

// NextPlayer will return a pointer to the next player
func (t *turn) NextPlayer() *model.User {
	panic("unimplemented")
	return nil
}

// PassTurn will pass the current players turn to the next player
func (t *turn) PassTurn(currentPlayerId string) *model.User {
	panic("unimplemented")
	return nil
}
