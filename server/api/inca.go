package api

import (
	"net/http"

	"github.com/dragonrider23/inca/poller"
)

func incaAPI(r *http.Request, urlPieces []string) (interface{}, *apiError) {
	switch urlPieces[0] {
	case "hb":
		return incaHeartbeat(r)
	default:
		return "", newError("Endpoint /inca/"+urlPieces[0]+" not found", 1)
	}

	return "Invalid job", nil
}

func incaHeartbeat(r *http.Request) (string, *apiError) {
	c, _, err := poller.Process(poller.Job{
		Cmd: "echo",
		Data: map[string]interface{}{
			"echo": "heartbeat",
		},
	})

	if err != nil {
		return "", newError(err.Error(), 1)
	}

	res := <-c
	return res.Data.(string), nil
}
