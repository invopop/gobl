package api

import (
	"net/http"

	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/internal/cli"
)

func handleKeygen(w http.ResponseWriter, _ *http.Request) {
	key := dsig.NewES256Key()
	writeJSON(w, cli.KeygenResponse{
		Private: key,
		Public:  key.Public(),
	})
}
