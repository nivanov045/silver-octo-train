package requester

import (
	"bytes"
	"fmt"
	"net/http"
)

type requester struct{}

func (*requester) Send(a []byte) error {
	fmt.Println("Snd")
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/update/", bytes.NewBuffer(a))
	if err != nil {
		fmt.Println("err != nil")
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Close = true
	response, err := client.Do(request)
	fmt.Println(request.URL)
	fmt.Println(string(a))
	if err != nil {
		fmt.Println("err != nil 2")
		return err
	}
	defer response.Body.Close()
	return nil
}

func New() *requester {
	return &requester{}
}
