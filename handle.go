package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)

type Content struct {
	Data []byte
}

type Shell struct {
	ShellType    int    `json:"shellType"`
	ShellContent string `json:"shellContent"`
}

func Handle(data []byte) Content {
	return Content{data}
}

func (d Content) ParseShellJSON() Shell {
	var shell Shell
	_ = json.Unmarshal(d.Data, &shell)
	return shell
}

func (s Shell) Run() {
	switch s.ShellType {
	case 0:
		switch runtime.GOOS {
		case "linux":
			switch runtime.GOARCH {
			case "386", "amd64", "arm", "arm64":
				cmd := exec.Command("/bin/bash", "-c", s.ShellContent)
				_, err := cmd.StdoutPipe()
				if err != nil {
					fmt.Println(err)
				}
				_ = cmd.Start()
			case "mips", "mipsle":
				cmd := exec.Command("/bin/ash", "-c", s.ShellContent)
				_, err := cmd.StdoutPipe()
				if err != nil {
					fmt.Println(err)
				}
				_ = cmd.Start()
			}
		case "darwin":
			cmd := exec.Command("/bin/zsh", "-c", s.ShellContent)
			_, err := cmd.StdoutPipe()
			if err != nil {
				fmt.Println(err)
			}
			_ = cmd.Start()
		default:
			fmt.Println("OS dose not support")
		}
	case 1:
		fmt.Println("Command dose not support")
		fmt.Println("Please use GUI program")
	default:
		fmt.Println("Command Error!")
	}
}
