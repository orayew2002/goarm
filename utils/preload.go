package utils

import (
	"fmt"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
)

const text = `  ____                            
 / ___| ___   __ _ _ __ _ __ ___  
| |  _ / _ \ / _' | '__| '_ ' _ \ 
| |_| | (_) | (_| | |  | | | | | |
 \____|\___/ \__'_|_|  |_| |_| |_|
`

func ShowPreload() {
	fmt.Println(Green + text + Reset)
}
