package api

import (
	"firebase.google.com/go/auth"
	"github.com/matjazonline/go-adaptr-httprouter-handlers"
	"log"
	"net/http"
	"strings"
)

func AuthorizedNs2Ctx(firebaseAuth *auth.Client, verifyAuthorizationRoleId string) adaptr.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			token, err := firebaseAuth.VerifyIDToken(r.Context(), idToken)
			//token, err := firebaseAuth.VerifyIDTokenAndCheckRevoked(r.Context(), idToken)
			if err == nil {
				namespaceAuthorizations := token.Claims["nsAuth"]

				if (namespaceAuthorizations != nil) {
					nsRequestVal := r.Header.Get(NAMESPACE_PARAM_NAME)
					if (nsRequestVal == "") {
						nsRequestVal, _ = adaptr.GetHttprouterQueryParameterValue(r, NAMESPACE_PARAM_NAME)
					}

					//ns, nsOk := r.Context().Value(ctxNamespaceKey).(string)
					if (nsRequestVal != "") {
						tokenNamespacesAuthorizationsMap, ok := namespaceAuthorizations.(map[string]interface{})
						if ok {
							tokenAuthorizationsForRequestedNamespace, strArrOk := tokenNamespacesAuthorizationsMap[nsRequestVal].([]interface{})
							if (strArrOk) {
								for _, authIdent := range tokenAuthorizationsForRequestedNamespace {
									if (strings.TrimSpace(authIdent.(string)) == strings.TrimSpace(verifyAuthorizationRoleId)) {

										//TODO ...
										r = adaptr.SetCtxValue(r, "ns", nsRequestVal)
										h.ServeHTTP(w, r)
										return
									}
								}
							}
						}
					}
				}

			} else {
				log.Printf("error verifying ID token: %v\n", idToken)
			}
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		})
	}
}

func NsParam2Ctx(ctxNamespaceKey string) adaptr.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO get from host addr

			nsRequestVal, _ := adaptr.GetHttprouterQueryParameterValue(r, NAMESPACE_PARAM_NAME)
			if (len(nsRequestVal) < 1) {
				nsRequestVal = r.Header.Get(NAMESPACE_PARAM_NAME)
			}

			if (len(nsRequestVal) > 0) {
				r = adaptr.SetCtxValue(r, ctxNamespaceKey, nsRequestVal)
				h.ServeHTTP(w, r)
				return
			}

			http.Error(w, http.StatusText(http.StatusPartialContent), http.StatusPartialContent)
		})
	}
}

