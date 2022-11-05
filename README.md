# Async/Await for golang

![coverage-lines](assets/coverage/coverage.svg)


example usage
```go

package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/agc11/gowait"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {

	type apiResponse struct {
		Id        int    `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		FullName  string `json:"fullName"`
		Title     string `json:"title"`
		Family    string `json:"family"`
		Image     string `json:"image"`
		ImageUrl  string `json:"imageUrl"`
	}

	requestApi := func(ctx context.Context) (apiResponse, error) {
		resp, err := http.Get("https://thronesapi.com/api/v2/Characters/0")

		if err != nil {
			return apiResponse{}, errors.New("error request")
		}
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return apiResponse{}, errors.New("error ioutil.ReadAll")
		}

		var res apiResponse
		err = json.Unmarshal(body, &res)

		if err != nil {
			return apiResponse{}, errors.New("error Unmarshal")
		}

		return res, nil
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	future := gowait.NewFuture(context.Background(), requestApi)
	result := future.Await(ctxTimeout)

	if result.Error != nil {
		log.Fatalln(result.Error)
	}

	log.Printf("%#v\n", result.Value)
}

```