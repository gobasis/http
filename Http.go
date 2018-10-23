package http

import (
	"fmt"
	"github.com/gobasis/log"
	"net/http"
)

func Run() {
	addr := fmt.Sprintf("%s:%d", ServerConf.Address, ServerConf.Port)
	log.Info("starting http server", "address", addr)
	Router.GetRawRouter().NotFound = http.FileServer(http.Dir(ServerConf.Public))
	err := http.ListenAndServe(addr, Router.GetRawRouter())
	if err != nil {
		log.Fatal("http server start failed", "error", err)
	}
}