package config

import (
	"flag"
	"fmt"
	"os"
)

var (
	Name string
	Key  string
)

var CliName = flag.String("username", "", "Input Your Username")
var CliKey = flag.String("passwd", "", "Input Your Device Password")

func init() {
	flag.Parse()
	if *CliName == "" || *CliKey == "" {
		fmt.Println("Param error")
		fmt.Println("Please try again")
		fmt.Println("Use 'wqy -h' for help")
		os.Exit(-1)
	} else {
		Name = *CliName
		Key = *CliKey
	}
}
