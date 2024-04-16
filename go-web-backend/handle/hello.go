package handle

import (
	"fmt"
	"go-web-backend/servercontext"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	ctx := servercontext.FromRequest(r)
	ctx.Logger.Info("logging!")
	fmt.Fprintln(w, "Hello!")
}
