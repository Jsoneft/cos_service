package main

import (
	"context"
	"log"

	"github.com/Jsoneft/cos_service/constant"
	"github.com/Jsoneft/cos_service/service/cos"
)

func init() {
	err := constant.SetupSetting()
	if err != nil {
		panic(err)
	}
}

func main() {

	c := cos.New()
	err := c.Serve(context.Background())
	if err != nil {
		log.Panicln(err)
	}
}
