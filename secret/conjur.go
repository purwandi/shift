package secret

import (
	"log"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/cyberark/conjur-api-go/conjurapi"
)

type ConjurClient struct {
	client *conjurapi.Client
}

func NewConjurClient() (*ConjurClient, error) {
	config, err := conjurapi.LoadConfig()
	if err != nil {
		return nil, err
	}

	conjur, err := conjurapi.NewClientFromEnvironment(config)
	if err != nil {
		return nil, err
	}

	return &ConjurClient{
		client: conjur,
	}, nil
}

func (c *ConjurClient) Get(key string) (secret string, err error) {
	var secretValue []byte

	connect := func() error {
		secretValue, err = c.client.RetrieveSecret(key)
		return err
	}
	notify := func(err error, t time.Duration) {
		log.Println("[config]", err.Error(), t)
	}

	bcf := backoff.NewExponentialBackOff()
	bcf.MaxElapsedTime = 5 * time.Minute

	cerr := backoff.RetryNotify(connect, bcf, notify)
	if cerr != nil {
		log.Fatal("[config] giving up connecting to retrieve secret config ")
	}

	res, err := strconv.Unquote("\"" + string(secretValue) + "\"")
	if err != nil {
		return
	}

	return res, nil
}
