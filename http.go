package requests

import (
	"bytes"
	"context"
	"crypto/tls"
	"net/http"
)

// DefaultMaxTryTimes 默认情况下的最大重试次数
const DefaultMaxTryTimes = 3

// SendRequest 底层API，不建议直接调用
func SendRequest[Request any, Response any](ctx context.Context, options *Options[Request, Response]) (Response, error) {

	// TODO set default params

	var zero Response
	var lastErr error
	for tryTimes := 0; tryTimes < options.MaxTryTimes; tryTimes++ {
		var client http.Client

		// 忽略证书验证
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		httpRequest, err := http.NewRequest(options.Method, options.TargetURL, bytes.NewReader(options.Body))
		if err != nil {
			return zero, err
		}

		httpRequest = httpRequest.WithContext(ctx)

		for _, requestSettingFunc := range options.RequestSettingSlice {
			if err := requestSettingFunc(&client, httpRequest); err != nil {
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
