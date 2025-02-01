package storage

type Message struct {
	Alias   string `json:"alias"`
	PubKey  []byte `json:"pub_key"`
	Witness []byte `json:"witness"`
	Tag     string `json:"tag"`
	Data    string `json:"data"`
}

type Storage interface {
	WriteOne(msg Message) error
	ReadAll(tag string) ([]Message, error)
}
