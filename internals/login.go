package internals

import (
	"time"
)

type Login struct {
	Login               string
	PrivateKeyEncrypted []byte
	CreateTime          time.Time
}

func AddLogin(name, password string) {

}
