package target_test

import (
	"fmt"
	"io/ioutil"
	"log"

	"net/http"
)

type Http struct {
}

func NewHttp() Http {
	return Http{}
}

func (r *Http) Test() string {
	host := "luutia.com"

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, host, nil)

	if err != nil {
		fmt.Println("ERROR:", err)
		return ""
	}

	q := req.URL.Query()
	q.Add("foo", "bar")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return err.Error()
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Status)
	fmt.Println(string(responseBody))

	return "ok"
}
