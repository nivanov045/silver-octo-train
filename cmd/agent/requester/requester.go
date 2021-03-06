package requester

import (
	"fmt"
	"net/http"
)

type requester struct{}

func (*requester) Send(a string) error {
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/"+a, nil)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "text/plain")
	response, err := client.Do(request)
	fmt.Println(request.URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

func New() *requester {
	return &requester{}
}
