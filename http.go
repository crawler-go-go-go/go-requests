package requests

import (
	"bytes"
	"context"
	"net/http"
)

// DefaultMaxTryTimes 默认情况下的最大重试次数
const DefaultMaxTryTimes = 3

type RequestSetting func(httpRequest *http.Request) error

const DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"

func DefaultUserAgentRequestSetting() RequestSetting {
	return func(httpRequest *http.Request) error {
		httpRequest.Header.Set("User-Agent", DefaultUserAgent)
		return nil
	}
}

// ------------------------------------------------- --------------------------------------------------------------------

func SendRequest[Request any, Response any](ctx context.Context, options *Options[Request, Response]) (Response, error) {

	// TODO set default params

	var zero Response
	var lastErr error
	for tryTimes := 0; tryTimes < options.MaxTryTimes; tryTimes++ {
		var client http.Client
		httpRequest, err := http.NewRequest(options.Method, options.TargetURL, bytes.NewReader(options.Body))
		if err != nil {
			return zero, err
		}

		httpRequest = httpRequest.WithContext(ctx)

		for _, requestSettingFunc := range options.RequestSettingSlice {
			if err := requestSettingFunc(httpRequest); err != nil {
				return zero, err
			}
		}

		httpResponse, err := client.Do(httpRequest)
		if err != nil {
			lastErr = err
			continue
		}
		defer httpResponse.Body.Close()

		response, err := options.ResponseHandler(httpResponse)
		if err != nil {
			lastErr = err
			continue
		}
		return response, nil
	}

	return zero, lastErr
}
