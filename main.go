package main

import (
	"IM/router"
	"IM/utils"
)

func main() {

	r := router.Router()
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	r.Run(":8081")
}
