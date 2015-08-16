package api

import (
	"net/http"

	"github.com/dragonrider23/inca/configs"
)

func incaAPI(r *http.Request, urlPieces []string) (interface{}, *apiError) {
	switch urlPieces[0] {
	case "hb":
		return incaHeartbeat(r)
	default:
		return "", newError("Endpoint /cp/"+urlPieces[0]+" not found", 1)
	}

	return "Invalid job", nil
}

func incaHeartbeat(r *http.Request) (string, *apiError) {
	s, e := configs.HeartBeat()
	if e != nil {
		return "", newError(e.Error(), 1)
	}
	return s, nil
}
