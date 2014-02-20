package interfaces

type World interface {
	Start() error

	AddClient(c Client) error
	RemoveClient(c Client) error
}
