package main

import (
	"github.com/stratosnet/sds/pp/peers"
	"github.com/stratosnet/sds/pp/setting"
)

func main() {
	setting.LoadConfig("./configs/config1.yaml")
	peers.StartPP()

	for {

	}
}