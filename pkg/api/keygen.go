package api

import (
	"net/http"

	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/internal/ops"
)

func handleKeygen(w http.ResponseWriter, _ *http.Request) {
	key := dsig.NewES256Key()
	WriteJSON(w, ops.KeygenResponse{
		Private: key,
		Public:  key.Public(),
	})
}
