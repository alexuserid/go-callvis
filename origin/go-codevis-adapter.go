package origin

import (
	"log"
	"net/http"
)

func Handler(mainPckgPath string) http.Handler {
	tests := false
	args := []string{mainPckgPath}

	Analysis = new(analysis)
	if err := Analysis.DoAnalysis(CallGraphType(*callgraphAlgo), "", tests, args); err != nil {
		log.Fatal(err)
	}

	return http.HandlerFunc(handler)
}
