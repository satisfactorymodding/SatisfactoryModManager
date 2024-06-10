package wailsextras

import (
	"reflect"
	"strings"

	"github.com/wailsapp/go-webview2/pkg/edge"
)

func addUserAgent(newUserAgent string) {
	frontendRef := getFrontendReflected()
	chromiumRef := reflect.Indirect(frontendRef).FieldByName("chromium")
	callbackRef := reflect.Indirect(chromiumRef).FieldByName("WebResourceRequestedCallback")
	readableCallbackRef := allowUnexportedFieldAccess(callbackRef)

	prevCallback := readableCallbackRef.Interface().(func(req *edge.ICoreWebView2WebResourceRequest, args *edge.ICoreWebView2WebResourceRequestedEventArgs))

	readableCallbackRef.Set(reflect.ValueOf(func(req *edge.ICoreWebView2WebResourceRequest, args *edge.ICoreWebView2WebResourceRequestedEventArgs) {
		// Setting the UserAgent on the CoreWebView2Settings clears the whole default UserAgent of the Edge browser, but
		// we want to just append our ApplicationIdentifier. So we adjust the UserAgent for every request.
		if reqHeaders, err := req.GetHeaders(); err == nil {
			useragent, _ := reqHeaders.GetHeader("User-Agent")
			useragent = strings.Join([]string{useragent, newUserAgent}, " ")
			_ = reqHeaders.SetHeader("User-Agent", useragent)
			_ = reqHeaders.Release()
		}

		prevCallback(req, args)
	}))
}
