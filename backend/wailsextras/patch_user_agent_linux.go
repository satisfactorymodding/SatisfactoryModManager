package wailsextras

/*
#cgo linux pkg-config: gtk+-3.0
#cgo !webkit2_41 pkg-config: webkit2gtk-4.0
#cgo webkit2_41 pkg-config: webkit2gtk-4.1

#include <gtk/gtk.h>
#include <webkit2/webkit2.h>

void add_user_agent(GtkWidget *webview, const char *newUserAgent) {
    WebKitSettings *settings = webkit_web_view_get_settings(WEBKIT_WEB_VIEW(webview));
	const gchar * userAgent = webkit_settings_get_user_agent(settings);
	gchar *newUserAgentWithApp = g_strconcat(userAgent, " ", newUserAgent, NULL);
	webkit_settings_set_user_agent(settings, newUserAgentWithApp);
}
*/
import "C"

import (
	"reflect"
	"unsafe"
)

func addUserAgent(newUserAgent string) {
	frontendRef := getFrontendReflected()
	mainWindowRef := reflect.Indirect(frontendRef).FieldByName("mainWindow")
	webviewRef := reflect.Indirect(mainWindowRef).FieldByName("webview")
	readableWebviewRef := allowUnexportedFieldAccess(webviewRef)
	webview := readableWebviewRef.Interface().(unsafe.Pointer)
	C.add_user_agent((*C.GtkWidget)(webview), C.CString(newUserAgent))
}
