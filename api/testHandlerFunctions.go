package api

import (
	"fmt"
	"net/http"
)

func TestGetHnd(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "{\"method\":\"get\"}")
}

func TestPostHnd(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "{\"method\":\"post\"}")
}

