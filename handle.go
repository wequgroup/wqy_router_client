package main

import "encoding/json"

type Shell struct {
	ShellType    int    `json:"shellType"`
	ShellContent string `json:"shellContent"`
}

func ParseShellJSON(Data []byte) *Shell {
	var shell Shell
	_ = json.Unmarshal(Data, &shell)
	return &shell
}
