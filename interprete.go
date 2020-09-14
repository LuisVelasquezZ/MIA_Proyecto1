package main

import (
	"MIA_Proyecto1/disco"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Println("Ingresar Comandos: ('Enter' para finalizar)")
	interpretar("")
}

func exec(ruta string) {
	ruta = strings.Trim(ruta, "\n")
	if ruta[len(ruta)-3] == 'm' && ruta[len(ruta)-2] == 'i' && ruta[len(ruta)-1] == 'a' {
		file, err := os.Open(ruta)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		anterior := ""
		comando := ""
		for scanner.Scan() {
			comando = scanner.Text()
			comando = strings.Trim(comando, " ")
			if comando == "pause" {
				reader := bufio.NewReader(os.Stdin)
				fmt.Println(comando)
				reader.ReadRune()
			} else {
				if anterior == "" {
					if comando[len(comando)-2] == 92 && comando[len(comando)-1] == 42 {
						comand := []rune(comando)
						comandListo := string(comand[0 : len(comando)-2])
						comandListo = strings.Trim(comandListo, " ")
						comandListo += " "
						anterior = comandListo
					} else {
						fmt.Println(comando)
						lineaComando(comando)
					}
				} else {
					if comando[len(comando)-2] == 92 && comando[len(comando)-1] == 42 {
						comand := []rune(comando)
						comandListo := string(comand[0 : len(comando)-2])
						comandListo = strings.Trim(comandListo, " ")
						comandListo += " "
						anterior += comandListo
					} else {
						anterior += comando
						fmt.Println(anterior)
						lineaComando(anterior)
						anterior = ""
					}
				}
			}

		}
	} else {
		fmt.Println("Tipo de archivo incorrecto")
	}

}

func interpretar(anterior string) {
	finalizar := 0

	//Leyendo en consola
	reader := bufio.NewReader(os.Stdin)
	comando, _ := reader.ReadString('\n')

	//fmt.Println(comando)

	if comando == "\n" {
		finalizar = 1
	} else {
		comando = strings.Trim(comando, "\n")
		comando = strings.Trim(comando, " ")
		if comando != "" {
			if comando[0] != '#' {

				if comando[len(comando)-2] == 92 && comando[len(comando)-1] == 42 {
					comand := []rune(comando)
					comandListo := string(comand[0 : len(comando)-2])
					comandListo += " "
					if anterior != "" {
						interpretar(anterior + comandListo)
						anterior = ""
						return

					} else {
						interpretar(comandListo)
						return
					}

				} else {
					if anterior != "" {
						lineaComando(anterior + comando)
						anterior = ""
					} else {
						lineaComando(comando)
					}

				}
			}

		}
	}

	if finalizar != 1 {
		interpretar("")
		return
	}
}

func lineaComando(comando string) {
	if comando[0] == '#' {
		fmt.Println(comando)
	} else {
		//fmt.Println(comando)
		var commandArray []string
		commandArray = strings.Split(comando, " ")
		//fmt.Println(commandArray[1])
		ejecutarComando(commandArray) //Ejecutamos el comando.
	}

}

func ejecutarComando(commandArray []string) {
	particionesMontadas := []disco.MOUNTPart{}
	data := strings.ToLower(commandArray[0])
	if data == "mkdisk" {
		disco.MKDISK(commandArray)
	} else if data == "exec" {
		var ruta []string
		ruta = strings.Split(commandArray[1], "->")
		if ruta[0] == "-path" {
			exec(ruta[1])
		} else {
			fmt.Println("comando incorrecto")
		}

		//fmt.Println("Otro Comando")
	} else if data == "rmdisk" {
		disco.RMDISK(commandArray)

	} else if data == "fdisk" {
		disco.FDISK(commandArray)
	} else if data == "mount" {
		particionesMontadas = disco.MOUNT(commandArray, particionesMontadas)
	} else if data == "unmount" {

	} else if data == "mkfs" {

	} else if data == "login" {

	} else if data == "logout" {

	} else if data == "mkgpr" {

	} else if data == "rmgrp" {

	} else if data == "mkusr" {

	} else if data == "rmusr" {

	} else if data == "chmod" {

	} else if data == "mkfile" {

	} else if data == "cat" {

	} else if data == "rm" {

	} else if data == "edit" {

	} else if data == "REN" {

	} else if data == "mkdir" {

	} else if data == "cp" {

	} else if data == "mv" {

	} else if data == "find" {

	} else if data == "chown" {

	} else if data == "chgrp" {

	} else if data == "rep" {

	}
}
