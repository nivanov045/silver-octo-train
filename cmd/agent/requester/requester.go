package requester

import (
	"net/http"
)

type requester struct{}

func (*requester) Send(a string) error {
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, a, nil)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "text/plain")
	_, err = client.Do(request)
	if err != nil {
		return err
	}
	return nil
}

func New() *requester {
	return &requester{}
}
