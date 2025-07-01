package main

import (
	"os"

	"dns-update/internal/service"

	"github.com/alibabacloud-go/tea/tea"
)

func main() {
	err := service.Run(tea.StringSlice(os.Args[1:]))
	if err != nil {
		panic(err)
	}
}
