package chat

import (
	"encoding/hex"
	"fmt"

	"github.com/nem0z/room-chat/src/crypto"
)

func (m *Message) Prettify() string {
	return fmt.Sprintf("From: %v\nPublic key: %v\nWitness %v\nTag: %v\nValid: %v\n => %v",
		crypto.GetAlias(*crypto.PubKeyFromBytes(m.PubKey)),
		hex.EncodeToString(m.PubKey),
		hex.EncodeToString(m.Witness),
		m.Tag,
		crypto.VerifyWitness(m.PubKey, m.Witness, []byte(m.Data)),
		m.Data,
	)
}
