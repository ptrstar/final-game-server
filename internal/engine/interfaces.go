package engine

type GameItf interface {
	Run()
	Broadcast([]byte)
	AddPlayer(Client) error
	RemovePlayer(Client)
	HandleInput(Client, []byte)
	Update()
	SerializeState()
}
