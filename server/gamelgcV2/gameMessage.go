package gamelgcV2

type GameMessage struct {
	Action gameAction `json:"action"`
	// TODO: seems to not bee needed
	//GameId uuid.UUID  `json:"GameId"`

	UserId string `json:"userId"`
	TileId *int   `json:"tileId,omitempty"` // TileId represents tile when discarding and seizing
}
