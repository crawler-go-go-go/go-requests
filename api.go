package requests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// ------------------------------------------------- --------------------------------------------------------------------

// GetYaml 响应内容是YAML格式的
func GetYaml[Response any](ctx context.Context, targetUrl string, options ...*Options[any, Response]) (Response, error) {

	if len(options) == 0 {
		options = append(options, NewOptions[any, Response](targetUrl, YamlResponseHandler[Response]()))
	}

	options[0] = options[0].WithTargetURL(targetUrl).WithYamlResponseHandler()

	return SendRequest[any, Response](ctx, options[0])
}

// ------------------------------------------------- --------------------------------------------------------------------

// GetJson 响应内容是JSON格式的
func GetJson[Response any](ctx context.Context, targetUrl string, options ...*Options[any, Response]) (Response, error) {

	if len(options) == 0 {
		options = append(options, NewOptions[any, Response](targetUrl, JsonResponseHandler[Response]()))
	}

	options[0] = options[0].WithTargetURL(targetUrl).
		WithJsonResponseHandler().
		AppendRequestSetting(func(client *http.Client, request *http.Request) error {
			request.Header.Set("Content-Type", "application/json")
			return nil
		})

	return SendRequest[any, Response](ctx, options[0])
}

func PostJson[Request any, Response any](ctx context.Context, targetUrl string, request Request, options ...*Options[Request, Response]) (Response, error) {

	if len(options) == 0 {
		options = append(options, NewOptions[Request, Response](targetUrl, JsonResponseHandler[Response]()))
	}

	marshal, err := json.Marshal(request)
	if err != nil {
		var zero Response
		return zero, fmt.Errorf("PostJson json marshal request error: %s, typer = %s", err.Error(), reflect.TypeOf(request).String())
	}

	options[0] = options[0].WithTargetURL(targetUrl).
		WithJsonResponseHandler().
		AppendRequestSetting(func(client *http.Client, request *http.Request) error {
			request.Header.Set("Content-Type", "application/json")
			return nil
		}).
		WithBody(marshal)

	return SendRequest[Request, Response](ctx, options[0])
}

// ------------------------------------------------- --------------------------------------------------------------------

func GetBytes(ctx context.Context, targetUrl string, options ...*Options[any, []byte]) ([]byte, error) {

	if len(options) == 0 {
		options = append(options, NewOptions[any, []byte](targetUrl, BytesResponseHandler()))
	}

	options[0] = options[0].WithTargetURL(targetUrl).
		WithResponseHandler(BytesResponseHandler())

	return SendRequest[any, []byte](ctx, options[0])
}

// ------------------------------------------------- --------------------------------------------------------------------

func GetString(ctx context.Context, targetUrl string, options ...*Options[any, []byte]) (string, error) {
	responseBytes, err := GetBytes(ctx, targetUrl, options...)
	if err != nil {
		return "", err
	}
	return string(responseBytes), nil
}

// ------------------------------------------------- --------------------------------------------------------------------

