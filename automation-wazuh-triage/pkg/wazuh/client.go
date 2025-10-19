package wazuh

import (
	"os"

	"github.com/go-resty/resty/v2"
)

type Wazuh struct {
	Client *resty.Client
}

func NewWazuh() *Wazuh {
	rest := resty.New()

	rest.BaseURL = os.Getenv("WAZUH_URL")
	rest.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	})

	return &Wazuh{
		Client: rest,
	}
}

func (w *Wazuh) authenticate() error {
	username := os.Getenv("WAZUH_USERNAME")
	password := os.Getenv("WAZUH_PASSWORD")

	w.Client.SetBasicAuth(username, password)

	res, err := w.Client.R().Get("/security/user/authenticate?raw=true")
	if err != nil {
		return err
	}

	w.Client.SetAuthToken(res.String())

	return nil
}
