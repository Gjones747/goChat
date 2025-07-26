package middleWare

import (
	"log"
	"net/http"
)

func LogMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Printf("%s %s", r.Method, r.URL)

		next.ServeHTTP(w, r)
	})

}
