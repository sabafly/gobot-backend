package cli

import (
	"github.com/ikafly144/gobot-backend/pkg/mc"
	"github.com/ikafly144/gobot-backend/pkg/worker"
)

func Run() {
	go mc.Start()
	worker.StartServer()
}
