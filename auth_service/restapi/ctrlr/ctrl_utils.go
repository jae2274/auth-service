package ctrlr

import (
	"context"
	"net/http"

	"github.com/jae2274/goutils/llog"
)

func errorHandler(ctx context.Context, w http.ResponseWriter, err error) bool {
	if err != nil {
		llog.LogErr(ctx, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return true
	}
	return false
}
