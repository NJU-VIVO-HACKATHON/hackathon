package main

import "github.com/NJU-VIVO-HACKATHON/hackathon/router"

func main() {
	r := router.SetupRouter()
	r.Run("0.0.0.0:8974")
}
