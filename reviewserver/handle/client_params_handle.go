package handle

import (
	"fmt"
	"net/http"
)

func ClientParamsHandle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	requestReviewParam := new(RequestReviewParam)
	fmt.Fprintf(w, requestReviewParam.ClientParams())
	return
}
