package gowait

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"
)

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

var daenerys = apiResponse{
	Id:        0,
	FirstName: "Daenerys",
	LastName:  "Targaryen",
	FullName:  "Daenerys Targaryen",
	Title:     "Mother of Dragons",
	Family:    "House Targaryen",
	Image:     "daenerys.jpg",
	ImageUrl:  "https://thronesapi.com/assets/images/daenerys.jpg",
}

func TestNewFuture(t *testing.T) {
	type args struct {
		withCancelTime *time.Duration
		callback       func(ctx context.Context) (any, error)
	}

	tests := []struct {
		name string
		args args
		want Result[any]
	}{
		{
			name: "should cancel if context deadline exceeded",
			args: args{
				withCancelTime: pointer(time.Millisecond * 1),
				callback: func(ctx context.Context) (any, error) {
					time.Sleep(time.Second * 10)
					return 1, nil
				},
			},
			want: Result[any]{
				Value: nil,
				Error: context.DeadlineExceeded,
			},
		},
		{
			name: "should return value",
			args: args{
				callback: func(ctx context.Context) (any, error) {
					return 1, nil
				},
			},
			want: Result[any]{
				Value: 1,
				Error: nil,
			},
		},
		{
			name: "should return error",
			args: args{
				callback: func(ctx context.Context) (any, error) {
					return nil, errors.New("")
				},
			},
			want: Result[any]{
				Value: nil,
				Error: errors.New(""),
			},
		},
		{
			name: "should return result from api",
			args: args{
				callback: func(ctx context.Context) (any, error) {
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
				},
			},
			want: Result[any]{
				Value: daenerys,
				Error: nil,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			if tt.args.withCancelTime != nil {
				ctxTimeout, cancel := context.WithTimeout(ctx, *tt.args.withCancelTime)
				ctx = ctxTimeout
				defer cancel()
			}

			future := NewFuture(context.Background(), tt.args.callback)
			result := future.Await(ctx)
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("NewFuture() = %v, want %v", result, tt.want)
			}
		})
	}
}
