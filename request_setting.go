package requests

import (
	"net/http"
	"net/url"
)

type RequestSetting func(client *http.Client, request *http.Request) error

const DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"

// DefaultRequestSetting 每个请求默认添加的设置
func DefaultRequestSetting() RequestSetting {
	return func(client *http.Client, httpRequest *http.Request) error {
		settings := []RequestSetting{
			RequestSettingUserAgent(),
			RequestSettingSkipTlsVerify(),
		}
		for _, setting := range settings {
			err := setting(client, httpRequest)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// RequestSettingUserAgent 设置User-Agent，如果不传递的则使用默认的User-Agent
func RequestSettingUserAgent(userAgent ...string) RequestSetting {
	if len(userAgent) == 0 {
		userAgent = append(userAgent, DefaultUserAgent)
	}
	return func(client *http.Client, httpRequest *http.Request) error {
		httpRequest.Header.Set("User-Agent", userAgent[0])
		return nil
	}
}

// RequestSettingProxy 配置代理IP
func RequestSettingProxy(proxy string) RequestSetting {
	parse, err := url.Parse(proxy)
	return func(client *http.Client, httpRequest *http.Request) error {
		if err != nil {
			return err
		}
		transport, ok := client.Transport.(*http.Transport)
		if !ok {
			transport = &http.Transport{}
		}
		transport.Proxy = http.ProxyURL(parse)
		client.Transport = transport
		return nil
	}
}

// RequestSettingSkipTlsVerify 跳过https证书验证
func RequestSettingSkipTlsVerify() RequestSetting {
	return func(client *http.Client, httpRequest *http.Request) error {
		transport, ok := client.Transport.(*http.Transport)
		if !ok {
			transport = &http.Transport{}
		}
		transport.TLSClientConfig.InsecureSkipVerify = true
		client.Transport = transport
		return nil
	}
}
