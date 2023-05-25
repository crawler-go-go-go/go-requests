package requests

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
	"reflect"
)

type ResponseHandler[Response any] func(httpResponse *http.Response) (Response, error)

func BytesResponseHandler(readResponseOnStatusCodeIn ...int) ResponseHandler[[]byte] {

	// By default, the response body is read only when the status code is 200
	if len(readResponseOnStatusCodeIn) == 0 {
		readResponseOnStatusCodeIn = append(readResponseOnStatusCodeIn, http.StatusOK, http.StatusNotFound)
	}

	return func(httpResponse *http.Response) ([]byte, error) {
		for _, status := range readResponseOnStatusCodeIn {
			if status == httpResponse.StatusCode {
				responseBodyBytes, err := io.ReadAll(httpResponse.Body)
				if err != nil {
					return nil, fmt.Errorf("response statuc code: %d, read body error: %s", httpResponse.StatusCode, err.Error())
				}
				return responseBodyBytes, nil
			}
		}
		return nil, fmt.Errorf("response status code: %d", httpResponse.StatusCode)
	}
}

func StringResponseHandler(readResponseOnStatusCodeIn ...int) ResponseHandler[string] {
	return func(httpResponse *http.Response) (string, error) {
		responseBytes, err := BytesResponseHandler(readResponseOnStatusCodeIn...)(httpResponse)
		if err != nil {
			return "", err
		}
		return string(responseBytes), nil
	}
}

func YamlResponseHandler[Response any](readResponseOnStatusCodeIn ...int) ResponseHandler[Response] {
	return func(httpResponse *http.Response) (Response, error) {

		var r Response

		responseBytes, err := BytesResponseHandler(readResponseOnStatusCodeIn...)(httpResponse)
		if err != nil {
			return r, err
		}

		err = yaml.Unmarshal(responseBytes, &r)
		if err != nil {
			return r, fmt.Errorf("response body yaml unmarshal error: %s, type: %s, response body: %s", err.Error(), reflect.TypeOf(r).String(), string(responseBytes))
		}
		return r, nil
	}
}

func JsonResponseHandler[Response any](readResponseOnStatusCodeIn ...int) ResponseHandler[Response] {
	return func(httpResponse *http.Response) (Response, error) {

		var r Response

		responseBytes, err := BytesResponseHandler(readResponseOnStatusCodeIn...)(httpResponse)
		if err != nil {
			return r, err
		}

		err = json.Unmarshal(responseBytes, &r)
		if err != nil {
			return r, fmt.Errorf("response body json unmarshal error: %s, type: %s, response body: %s", err.Error(), reflect.TypeOf(r).String(), string(responseBytes))
		}
		return r, nil
	}
}
