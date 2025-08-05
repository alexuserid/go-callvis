package origin

import (
	"fmt"
	"log"
	"net/http"
)

type GoCodevisAdapter struct{}

type Config struct {
	MainPkgPath   string
	CallgraphAlgo CallGraphType
	NoStd         bool
	NoInter       bool
	Group         []string
}

func DefaultConfig() Config {
	return Config{
		MainPkgPath:   "./",
		CallgraphAlgo: CallGraphTypeStatic,
		NoStd:         true,
		NoInter:       true,
		Group:         []string{"pkg,type"},
	}
}

func NewGoCodevisAdapter(cfg Config) *GoCodevisAdapter {
	tests := false
	args := []string{cfg.MainPkgPath}

	Analysis = new(analysis)
	setInitOpts(Analysis, cfg)

	if err := Analysis.DoAnalysis(CallGraphType(*callgraphAlgo), "", tests, args); err != nil {
		log.Fatal(err)
	}

	return &GoCodevisAdapter{}
}

func (a *GoCodevisAdapter) Handler() http.Handler {
	return http.HandlerFunc(goCodevisHandler)
}

func setInitOpts(a *analysis, config Config) {
	Analysis.opts = &renderOpts{
		cacheDir: a.opts.cacheDir,
		focus:    a.opts.focus,
		group:    config.Group,
		ignore:   a.opts.ignore,
		include:  a.opts.include,
		limit:    a.opts.limit,
		nointer:  config.NoInter,
		nostd:    config.NoStd,
	}
}

func goCodevisHandler(w http.ResponseWriter, r *http.Request) {
	logf("----------------------")
	logf(" => handling request:  %v", r.URL)
	logf("----------------------")

	// Allow overriding by HTTP params
	Analysis.OverrideByHTTP(r)

	var img string
	if img = Analysis.FindCachedImg(); img != "" {
		log.Println("serving cached file:", img)
		http.ServeFile(w, r, img)
		return
	}

	// Convert list-style args to []string
	if e := Analysis.ProcessListArgs(); e != nil {
		http.Error(w, "invalid parameters", http.StatusBadRequest)
		return
	}

	output, err := Analysis.Render()
	if err != nil {
		http.Error(w, fmt.Sprintf("rendering failed: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	if r.Form.Get("format") == "dot" {
		log.Println("writing dot output")
		fmt.Fprint(w, string(output))
		return
	}

	log.Printf("converting dot to %s\n", *outputFormat)

	img, err = dotToImage("", *outputFormat, output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Analysis.CacheImg(img)
	if err != nil {
		http.Error(w, "cache img error: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("serving file:", img)
	http.ServeFile(w, r, img)
}
