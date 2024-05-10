package wailsextras

import (
	"reflect"
	"strings"
	"unsafe"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/common"
)

func AddUserAgent(userAgentName, userAgentVersion string) {
	fullUserAgent := strings.Join([]string{userAgentName, userAgentVersion}, "/")
	fullUserAgent = strings.TrimSuffix(fullUserAgent, "/") // in case no version is provided
	addUserAgent(fullUserAgent)
}

func getFrontendReflected() reflect.Value {
	frontend := common.AppContext.Value("frontend")
	return getInnermostFrontend(reflect.ValueOf(frontend))
}

func getInnermostFrontend(frontend reflect.Value) reflect.Value {
	for i := 0; i < reflect.Indirect(frontend).NumField(); i++ {
		if reflect.Indirect(frontend).Field(i).Type().String() == "frontend.Frontend" {
			return getInnermostFrontend(reflect.Indirect(frontend).Field(i).Elem())
		}
	}
	return frontend
}

func allowUnexportedFieldAccess(field reflect.Value) reflect.Value {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
}
