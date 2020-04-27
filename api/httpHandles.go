package api

import (
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/matjazonline/go-adaptr-httprouter-handlers"
	"net/http"
)

var NAMESPACE_PARAM_NAME = "ns"

func TestGET(lifecycleAdapters *adaptr.RequestLifecycleAdapters) httprouter.Handle {
	return adaptr.ToHttpRouterHandle(
		TestGetHnd,
		lifecycleAdapters,
	)
}

func TestPOST(lifecycleAdapters *adaptr.RequestLifecycleAdapters) httprouter.Handle {
	return adaptr.WrapHandleFuncAdapters(
		[]adaptr.Adapter{
			//adaptr.WriteResponse(`backend wrks :)`),
		},
		TestPostHnd,
		lifecycleAdapters.BeforeHandlerFn,
		lifecycleAdapters.AfterHandlerFn,
	)
}

func AdminJsonAdaptrs(auth *auth.Client, testAuthorization string, testEnv bool) *adaptr.RequestLifecycleAdapters {
	return &adaptr.RequestLifecycleAdapters{
		BeforeHandlerFn: []adaptr.Adapter{
			AuthorizedNs2Ctx(auth, testAuthorization),
			adaptr.Json2Ctx(adaptr.CtxRequestJsonStructKey, false),
		},
		AfterHandlerFn: []adaptr.Adapter{
			adaptr.JsonContentType(),
		},
	}
}

func BasicGETLifecycle() *adaptr.RequestLifecycleAdapters {
	return &adaptr.RequestLifecycleAdapters{
		BeforeHandlerFn: []adaptr.Adapter{getPrintAdapter("BEFR ADAPTR")},
		AfterHandlerFn:  []adaptr.Adapter{getPrintAdapter("AFTR ADAPTR")},
	}
}

func getPrintAdapter(val string) adaptr.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(val)
			h.ServeHTTP(w, r)
		})
	}
}
