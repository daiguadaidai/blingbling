package handle

import (
	"net/http"
	"fmt"
)

func ClientParamsHandle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	requestReviewParam := new(RequestReviewParam)
	fmt.Fprintf(w, requestReviewParam.ClientParams())
	return
}

