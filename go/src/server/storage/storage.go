package storage

type Message struct {
	SenderPubKey []byte `json:"sender"`
	Witness      []byte `json:"witness"`
	Tag          string `json:"tag"`
	Data         string `json:"data"`
}

type Storage interface {
	WriteOne(msg Message) error
	ReadAll(tag string) ([]Message, error)
}
