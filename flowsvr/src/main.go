package main

import (
	"fmt"
	"github.com/niuniumart/asyncflow/flowsvr/src/config"
	"github.com/niuniumart/asyncflow/flowsvr/src/initialize"
	"github.com/niuniumart/gosdk/gin"

	"github.com/niuniumart/gosdk/martlog"
)

func main() {
	config.Init()
	err := initialize.InitResource()
	if err != nil {
		fmt.Printf("initialize.InitResource err %s", err.Error())
		martlog.Errorf("initialize.InitResource err %s", err.Error())
		return
	}
	router := gin.CreateGin()
	initialize.RegisterRouter(router)
	fmt.Println("before router run")
	err = gin.RunByPort(router, config.Conf.Common.Port)
	fmt.Println(err)
}
