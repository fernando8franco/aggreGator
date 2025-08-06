package main

import (
	"fmt"

	"github.com/fernando8franco/aggreGator/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	err = conf.SetUser("franco")
	if err != nil {
		fmt.Println(err)
	}
	conf, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(conf)
}
