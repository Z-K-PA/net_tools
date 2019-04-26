package util

import (
	"net/http"
	"net/http/pprof"
)

func PprofServerStart(address string) error{
	serverMux := http.NewServeMux()
	serverMux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	serverMux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	serverMux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	serverMux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	serverMux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	return http.ListenAndServe(address, serverMux)
}
