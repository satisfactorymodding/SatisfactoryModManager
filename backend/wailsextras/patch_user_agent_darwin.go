package wailsextras

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework WebKit
#import <Foundation/Foundation.h>
#import "WailsContext.h"

void add_user_agent(void *inctx, char *newUserAgent) {
    WailsContext *context = (__bridge WailsContext*) inctx;
    [context.webview evaluateJavaScript:@"navigator.userAgent" completionHandler:^(NSString *aUserAgent, NSError *aError) {
        NSString *sCustomUserAgent = @(newUserAgent);

        if (aUserAgent.length > 0 && aError == nil) {
            sCustomUserAgent = [NSString stringWithFormat:@"%@ %@", aUserAgent, sCustomUserAgent];
        }

        context.webview.customUserAgent = sCustomUserAgent;
    }];
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
	contextRef := reflect.Indirect(mainWindowRef).FieldByName("context")
	readableContextRef := allowUnexportedFieldAccess(contextRef)
	context := readableContextRef.Interface().(unsafe.Pointer)
	C.add_user_agent(context, C.CString(newUserAgent))
}
