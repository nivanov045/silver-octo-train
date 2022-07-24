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
	request.Header.Set("Content-Type", "application/json")
	request.Close = true
	fmt.Println(request.URL)
	fmt.Println(string(a))
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("err != nil 2 with ", err)
		return err
	}
	defer response.Body.Close()
	return nil
}

func New() *requester {
	return &requester{}
}
