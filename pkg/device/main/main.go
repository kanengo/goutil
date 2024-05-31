package main

import (
	"github.com/kanengo/goutil/pkg/device"
)

func main() {
	_, err := device.GetBaseBoardID()
	if err != nil {
		panic(err)
	}
}
