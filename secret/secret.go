package secret

import "github.com/sirupsen/logrus"

type Secret interface {
	Get(key string) (secret string, err error)
}

var conn Secret

func SetEngine(engine string) {
	var err error
	switch engine {
	case "conjur":
	case "":
		conn, err = NewConjurClient()
		if err != nil {
			logrus.Fatal(err)
		}
	default:
		conn, err = NewConjurClient()
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func Get(key string) (string, error) {
	if conn == nil {
		SetEngine("conjur")
	}

	return conn.Get(key)
}
