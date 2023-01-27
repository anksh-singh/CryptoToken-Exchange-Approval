package utils

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type IHttpRequest interface {
	GetRequest(url string) ([]byte, error)
	PostRequest(url string, body io.Reader) ([]byte, error)
	MakeGqlRequest(graphqlRequest *graphql.Request, url string) (map[string]interface{}, error)
	GetRequestWithHeaders(url string, headerKey string, headerValue string) ([]byte, error)
	PostRequestWithHeaders(url string, body string, headerKey string, headerValue string) ([]byte, error)
	GetRequestWithErrorResponse(url string) ([]byte, []byte, error)
	PutRequest(url string, body string) ([]byte, error)
}

type HttpRequest struct {
	logger       *zap.SugaredLogger
	successCodes []int
}

func NewHttpRequest(logger *zap.SugaredLogger) *HttpRequest {
	return &HttpRequest{
		logger:       logger,
		successCodes: []int{200, 201},
	}
}

func (h *HttpRequest) PostRequestWithHeaders(url string, body string, headerKey string, headerValue string) ([]byte, error) {
	client := &http.Client{}
	var data = strings.NewReader(body)
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		h.logger.Errorf(" Error in Post Request : %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set(headerKey, headerValue)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorf(" Post Request  Logging Error  is : %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "Http Client Request Error")
	}
	if contains(h.successCodes, resp.StatusCode) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				h.logger.Errorf("Error in Post Request  : %v", err.Error())
			}
		}(resp.Body)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "EOF")
		}
		return body, nil
	} else {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				h.logger.Errorf("Error in Post Request  : %v", err.Error())
			}
		}(resp.Body)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			h.logger.Errorf("Get Request  LoggingError  is : %v", err.Error())
		}
		return body, errors.New(string(body))
	}

}

func (h *HttpRequest) GetRequestWithHeaders(url string, headerKey string, headerValue string) ([]byte, error) {
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(headerKey, headerValue)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 200 {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				h.logger.Errorf("Error in GET Request  : %v", err.Error())
			}
		}(res.Body)

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "EOF")
		}
		return body, nil
	} else {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				h.logger.Errorf("Error in GET Request   : %v", err.Error())
			}
		}(res.Body)

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "EOF")
		}
		return body, status.Errorf(codes.Internal, string(body), "Error in http response")
	}

}

func (h *HttpRequest) GetRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		h.logger.Errorf(" Error in GET Request  : %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "Http Request Error")
	}
	if resp.StatusCode == 200 {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				h.logger.Errorf("Error in GET Request   : %v", err.Error())
			}
		}(resp.Body)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "EOF")
		}
		return body, nil
	} else {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				h.logger.Errorf("Error in GET Request  : %v", err.Error())
			}
		}(resp.Body)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			h.logger.Errorf("Error in GET Request   : %v", err.Error())
		}
		return body, status.Errorf(codes.Internal, string(body), "Error in http response")
	}

}

func (h *HttpRequest) PutRequest(url string, body string) ([]byte, error) {
	client := &http.Client{}
	var data = strings.NewReader(body)
	req, err := http.NewRequest("PUT", url, data)
	if err != nil {
		h.logger.Errorf(" Error in Put Request : %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	req.Header.Set("accept", "application/json")
	//req.Header.Set(headerKey, headerValue)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorf(" Put Request  Logging Error  is : %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if contains(h.successCodes, resp.StatusCode) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				h.logger.Errorf("Error in Put Request  : %v", err.Error())
			}
		}(resp.Body)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		return body, nil
	} else {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				h.logger.Errorf("Error in Put Request  : %v", err.Error())
			}
		}(resp.Body)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			h.logger.Errorf("Put Request  LoggingError  is : %v", err.Error())
		}
		return body, status.Errorf(codes.Internal, string(body))
	}
}

func (h *HttpRequest) PostRequest(url string, req io.Reader) ([]byte, error) {
	resp, err := http.Post(url, "application/json", req)
	if err != nil {
		h.logger.Errorf(" Error in POST Request : %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error(), "Error in Post Request")
	}
	if contains(h.successCodes, resp.StatusCode) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				h.logger.Errorf("Error in POST Request : %v", err.Error())
			}
		}(resp.Body)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error(), "EOF")
		}
		return body, nil
	} else {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				h.logger.Errorf("Error in POST Request : %v", err.Error())
			}
		}(resp.Body)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			h.logger.Errorf("Error in POST Request  : %v", err.Error())
		}
		return body, status.Errorf(codes.Internal, string(body), "Error in http response")
	}
}

func (h *HttpRequest) MakeGqlRequest(graphqlRequest *graphql.Request, url string) (map[string]interface{}, error) {
	ctx := context.Background()
	h.logger.Infof("making gql request")
	var responseData map[string]interface{}
	graphqlRequest.Header.Set("Cache-Control", "no-cache")
	graphqlClient := graphql.NewClient(url)
	if err := graphqlClient.Run(ctx, graphqlRequest, &responseData); err != nil {
		h.logger.Errorf(" graph ql request error with error message  is : %v", err.Error())
		return nil, err
	}
	jsonString, err := json.Marshal(responseData)
	if err != nil {
		h.logger.Errorf(" graph ql request error with error message  is : %v", err.Error())
		return nil, err
	}
	err = json.Unmarshal(jsonString, &responseData)
	if err != nil {
		h.logger.Errorf(" graph ql request error with error message  is : %v", err.Error())
		return nil, err
	}
	return responseData, err
}

type MockHttpRequest struct {
	mock.Mock
}

func (m *MockHttpRequest) GetRequest(url string) ([]byte, error) {
	args := m.Called()
	result := args.Get(0)
	return result.([]byte), args.Error(1)
}

func (m *MockHttpRequest) MakeGqlRequest(graphqlRequest *graphql.Request, url string) (map[string]interface{}, error) {
	args := m.Called()
	result := args.Get(0)
	return result.(map[string]interface{}), args.Error(1)
}

func (h *HttpRequest) GetRequestWithErrorResponse(url string) ([]byte, []byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		h.logger.Errorf(" Get Request  Logging Error  is : %v", err.Error())
		return nil, nil, status.Errorf(codes.Internal, err.Error(), "Error in Get Request")
	}
	if resp.StatusCode == 200 {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				h.logger.Errorf("Get Request  LoggingError  is : %v", err.Error())
			}
		}(resp.Body)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, status.Errorf(codes.Internal, err.Error(), "EOF")
		}
		return body, body, nil
	} else {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				h.logger.Errorf("Get Request  LoggingError  is : %v", err.Error())
			}
		}(resp.Body)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			h.logger.Errorf("Get Request  LoggingError  is : %v", err.Error())
		}
		return nil, body, status.Errorf(codes.Internal, string(body), "Error in http response")
	}

}
func contains(slice []int, element int) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}
