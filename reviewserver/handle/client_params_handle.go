package handle

import (
	"net/http"
	"fmt"
)

func ClientParamsHandle(w http.ResponseWriter, r *http.Request) {
	requestReviewParam := new(RequestReviewParam)
	fmt.Fprintf(w, requestReviewParam.ClientParams())
	return
}

