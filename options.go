package requests

import "net/http"

const DefaultMethod = http.MethodGet

type Options[Request any, Response any] struct {
	MaxTryTimes         int
	TargetURL           string
	Method              string
	Body                []byte
	RequestSettingSlice []RequestSetting
	ResponseHandler     ResponseHandler[Response]
}

func NewOptions[Request any, Response any](targetUrl string, responseHandler ResponseHandler[Response]) *Options[Request, Response] {
	return &Options[Request, Response]{
		MaxTryTimes:         DefaultMaxTryTimes,
		TargetURL:           targetUrl,
		Method:              DefaultMethod,
		Body:                []byte{},
		RequestSettingSlice: []RequestSetting{DefaultRequestSetting()},
		ResponseHandler:     responseHandler,
	}
}

func (x *Options[Request, Response]) WithMaxTryTimes(maxTryTimes int) *Options[Request, Response] {
	x.MaxTryTimes = maxTryTimes
	return x
}

func (x *Options[Request, Response]) WithTargetURL(targetURL string) *Options[Request, Response] {
	x.TargetURL = targetURL
	return x
}

func (x *Options[Request, Response]) WithMethod(method string) *Options[Request, Response] {
	x.Method = method
	return x
}

func (x *Options[Request, Response]) WithBody(body []byte) *Options[Request, Response] {
	if x.Method == DefaultMethod {
		x.Method = http.MethodPost
	}
	x.Body = body
	return x
}

func (x *Options[Request, Response]) WithRequestSettingSlice(requestSettingSlice []RequestSetting) *Options[Request, Response] {
	x.RequestSettingSlice = requestSettingSlice
	return x
}

func (x *Options[Request, Response]) AppendRequestSetting(requestSetting RequestSetting) *Options[Request, Response] {
	x.RequestSettingSlice = append(x.RequestSettingSlice, requestSetting)
	return x
}

func (x *Options[Request, Response]) WithResponseHandler(responseHandler ResponseHandler[Response]) *Options[Request, Response] {
	x.ResponseHandler = responseHandler
	return x
}

func (x *Options[Request, Response]) WithYamlResponseHandler() *Options[Request, Response] {
	x.ResponseHandler = YamlResponseHandler[Response]()
	return x
}

func (x *Options[Request, Response]) WithJsonResponseHandler() *Options[Request, Response] {
	x.ResponseHandler = JsonResponseHandler[Response]()
	return x
}
