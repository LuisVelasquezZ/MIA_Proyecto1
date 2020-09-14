package disco

import (
	"fmt"
	"os"
	"strings"
)

func REP(commandArray []string) {
	nameval := false
	pathval := false
	idval := false
	name := ""
	path := ""
	id := ""

	for i := 1; i < len(commandArray); i++ {
		tmp := strings.Split(commandArray[i], "->")
		if tmp[0] == "-path" {
			pathval = true
			if tmp[1][0] == '"' {
				path = tmp[1] + " " + commandArray[i+1]
				path = strings.Trim(path, "\"")
				i++
			} else {
				path = tmp[1]
			}
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				fmt.Println("Archivo no encontrado")
				return
			}
		} else if tmp[0] == "-name" {
			nameval = true
			name = []byte(tmp[1])
			if len(name) > 16 {
				fmt.Println("Longitud de nombre mayor a 16")
				return
			}
		} else if tmp[0] == "-id" {
			id = tmp[1]
		}
	}
}
