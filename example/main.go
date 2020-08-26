package main

import (
	"fmt"

	"github.com/shideqin/golog"
)

func main() {
	log := golog.NewLogger(map[string]interface{}{"dirname": "./log/", "output": "file", "format": "yyyyMMddHH"})
	fmt.Println(log)
}
