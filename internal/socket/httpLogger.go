package socket

import (
	"fmt"
	"net/http"
)

func httpLogger(request *http.Request) {
	fmt.Print(request.Method + " ")
	fmt.Print(request.URL)

	fmt.Println()
}
