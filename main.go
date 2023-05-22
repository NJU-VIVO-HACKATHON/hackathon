package main

import (
	"fmt"
	"github.com/NJU-VIVO-HACKATHON/hackathon/global"
	"github.com/NJU-VIVO-HACKATHON/hackathon/handler"
)

func main() {
	r := handler.SetupRouter()
	cfg := global.GetConfig()

	if err := r.Run(fmt.Sprintf("%s:%d",
		cfg.Server.Host, cfg.Server.Port,
	)); err != nil {
		panic(err)
	}
}
