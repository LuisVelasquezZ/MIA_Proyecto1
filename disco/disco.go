package disco

import (
	"bufio"
	"bytes"
	crypto "crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type MBR struct {
	mbr_tamanio        int64
	mbr_fecha_creacion time.Time
	mbr_disk_signature int64
	particion1         Particion
	particion2         Particion
	particion3         Particion
	particion4         Particion
}

type Particion struct {
	part_status byte
	part_type   byte
	part_fit    []byte
	part_start  int64
	part_size   int64
	part_name   []byte
}

type EBR struct {
	part_status byte
	part_fit    byte
	part_start  int64
	part_size   int64
	part_next   int64
	part_name   []byte
}

func MKDISK(commandArray []string) {
	sizeval := false
	pathval := false
	nameval := false
	var size int64 = 0
	path := ""
	name := ""
	unit := ""
	for i := 1; i < len(commandArray); i++ {
		tmp := strings.Split(commandArray[i], "->")
		if tmp[0] == "-size" {
			sizeval = true
			sizeconv, _ := strconv.ParseInt(tmp[1], 10, 64)
			if sizeconv > 0 {
				size = sizeconv
			} else {
				fmt.Println("Error en paramtetro size")
				return
			}

		} else if tmp[0] == "-path" {
			pathval = true
			if tmp[1][0] == '"' {
				path = tmp[1] + " " + commandArray[i+1]
				path = strings.Trim(path, "\"")
				i++
			} else {
				path = tmp[1]
			}

		} else if tmp[0] == "-name" {
			nameval = true
			name = tmp[1]
			archivoname := strings.Split(name, ".")
			if archivoname[1] == "dsk" {
				fmt.Println("Extensión Incorrecta")
				return
			}
		} else if tmp[0] == "-unit" {
			unit = tmp[1]
		} else {
			fmt.Println("Parametro Incorrecto")
			return
		}
	}

	if sizeval && pathval && nameval {

		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(path, 0777)
		}
		file, err := os.Create(path + name)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		var inicio int8 = 0
		s := &inicio
		var binario bytes.Buffer
		binary.Write(&binario, binary.BigEndian, s)
		escribirBytes(file, binario.Bytes())
		var val int64 = 0
		if unit != "" {
			if unit == "k" {
				val = size*1024 - 1
			} else if unit == "m" {
				val = size*1024*1024 - 1
			} else {
				fmt.Println("Parametro Incorrecto")
				return
			}
		} else {
			val = size*1024 - 1
		}
		file.Seek(val, 0)
		var binario2 bytes.Buffer
		binary.Write(&binario2, binary.BigEndian, s)
		escribirBytes(file, binario2.Bytes())
		file.Seek(0, 0)
		disco := MBR{
			mbr_tamanio:        size + 1,
			mbr_fecha_creacion: time.Now(),
			mbr_disk_signature: newCryptoRand(),
		}
		s1 := &disco
		var binario3 bytes.Buffer
		binary.Write(&binario3, binary.BigEndian, s1)
		escribirBytes(file, binario3.Bytes())

	} else {
		fmt.Println("Faltan Parametros")
		return
	}
}

func RMDISK(commandArray []string) {
	path := ""
	tmp := strings.Split(commandArray[1], "->")
	if tmp[0] == "-path" {
		reader := bufio.NewReader(os.Stdin)
		comando, _ := reader.ReadString('\n')
		path = tmp[1]
		fmt.Println("¿Eliminar disco " + path + "? [S/N]")
		if comando == "S" {
			file := os.Remove(path)
			if file != nil {
				log.Fatal(file)
			}
		} else {
			fmt.Println("Eliminacion Cancelada")
		}

	} else {
		fmt.Println("Faltan parametros")
	}

}

func FDISK(commandArray []string) {
	sizeval := false
	pathval := false
	nameval := false
	var size int64 = 0
	path := ""
	name := []byte("")
	unit := ""
	tipo := ""
	fit := ""
	delete := ""
	var add int64 = 0
	//leyendo atributos del comando
	for i := 1; i < len(commandArray); i++ {
		tmp := strings.Split(commandArray[i], "->")
		if tmp[0] == "-size" {
			sizeval = true
			sizeconv, _ := strconv.ParseInt(tmp[1], 10, 64)
			if sizeconv > 0 {
				size = sizeconv
			} else {
				fmt.Println("Error en paramtetro size")
				return
			}

		} else if tmp[0] == "-path" {
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
		} else if tmp[0] == "-unit" {
			unit = tmp[1]
		} else if tmp[0] == "-type" {
			tipo = tmp[1]
			if tipo != "p" && tipo != "e" && tipo != "l" {
				fmt.Println("Error en parametro type")
				return
			}
		} else if tmp[0] == "-fit" {
			fit = tmp[1]
			if fit != "bf" && fit != "ff" && fit != "wf" {
				fmt.Println("Error en el parametro fit")
				return
			}
		} else if tmp[0] == "-delete" {
			delete = tmp[1]
		} else if tmp[0] == "-add" {
			addconv, _ := strconv.ParseInt(tmp[1], 10, 64)
			if addconv != 0 {
				add = addconv
			} else {
				fmt.Println("Error en el parametro add")
				return
			}
		}
	}

	//verificando comandos obligatoiros
	if sizeval && pathval && nameval {

		//calculando size
		var val int64 = 0
		if unit != "" {
			if size != 0 {
				if unit == "k" {
					val = size * 1024
				} else if unit == "m" {
					val = size * 1024 * 1024
				} else if unit == "b" {
					val = size
				} else {
					fmt.Println("Parametro Incorrecto")
					return
				}
			} else {
				if unit == "k" {
					val = add * 1024
				} else if unit == "m" {
					val = add * 1024 * 1024
				} else if unit == "b" {
					val = add
				} else {
					fmt.Println("Parametro Incorrecto")
					return
				}
			}

		} else {
			if size != 0 {
				val = size * 1024
			} else {
				val = add * 1024
			}

		}
		//verificando disco
		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		//cargando MBR
		m := MBR{}
		var sizeMBR int = int(unsafe.Sizeof(m))
		data := leerBytes(file, sizeMBR)
		buffer := bytes.NewBuffer(data)

		err = binary.Read(buffer, binary.BigEndian, &m)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		//verificando cantidad de particiones y espacio
		cantParitcionesPrimarias := 0
		cantParticionExtendida := 0
		particionvacia := 0
		espacioDisponible := 0 //m.mbr_tamanio
		part1size := 0
		part2size := 0
		part3size := 0
		//part4size := 0
		patrExtendida := 0

		if &m.particion4 != nil {
			if m.particion4.part_type == 'p' {
				cantParitcionesPrimarias++
			} else if m.particion4.part_type == 'e' {
				cantParticionExtendida++
				patrExtendida = 4
			}
			//part4size = int(m.particion4.part_size)
			//espacioDisponible -= m.particion4.part_size
			if bytes.Compare(m.particion4.part_name, name) == 0 {
				fmt.Println("Nombre Repetido")
				return
			}
		} else {
			particionvacia = 4
			if &m.particion3 != nil {
				espacioDisponible = int(m.mbr_tamanio - (m.particion3.part_start + m.particion3.part_size))
			}
		}

		if &m.particion3 != nil {
			if m.particion3.part_type == 'p' {
				cantParitcionesPrimarias++
			} else if m.particion3.part_type == 'e' {
				cantParticionExtendida++
				patrExtendida = 3
			}
			part3size = int(m.particion3.part_size)
			//espacioDisponible -= m.particion3.part_size
			if bytes.Compare(m.particion3.part_name, name) == 0 {
				fmt.Println("Nombre Repetido")
				return
			}
		} else {
			particionvacia = 3
			if &m.particion2 != nil {
				if &m.particion4 != nil {
					espacioDisponible = int(m.particion4.part_start - (m.particion2.part_start + m.particion2.part_size))
				} else {
					espacioDisponible = int(m.mbr_tamanio - (m.particion2.part_start + m.particion2.part_size))
				}

			}
		}

		if &m.particion2 != nil {
			if m.particion2.part_type == 'p' {
				cantParitcionesPrimarias++
			} else if m.particion2.part_type == 'e' {
				cantParticionExtendida++
				patrExtendida = 2
			}
			part2size = int(m.particion2.part_size)
			//espacioDisponible -= m.particion2.part_size
			if bytes.Compare(m.particion2.part_name, name) == 0 {
				fmt.Println("Nombre Repetido")
				return
			}
		} else {
			particionvacia = 2
			if &m.particion1 != nil {
				if &m.particion3 != nil {
					espacioDisponible = int(m.particion3.part_start - (m.particion1.part_start + m.particion1.part_size))
					if espacioDisponible < int(val) {
						if &m.particion4 == nil {
							if (m.mbr_tamanio - (m.particion3.part_start + m.particion3.part_size)) >= val {
								particionvacia = 4
								espacioDisponible = int(m.mbr_tamanio - (m.particion3.part_start + m.particion3.part_size))
							}
						}
					}
				} else {
					if &m.particion4 != nil {
						espacioDisponible = int(m.particion4.part_start - (m.particion1.part_start + m.particion1.part_size))
					} else {
						espacioDisponible = int(m.mbr_tamanio - (m.particion1.part_start + m.particion1.part_size))
					}

				}

			}
		}

		if &m.particion1 != nil {
			if m.particion1.part_type == 'p' {
				cantParitcionesPrimarias++
			} else if m.particion1.part_type == 'e' {
				cantParticionExtendida++
				patrExtendida = 1
			}
			part1size = int(m.particion1.part_size)
			//espacioDisponible -= m.particion1.part_size
			if bytes.Compare(m.particion1.part_name, name) == 0 {
				fmt.Println("Nombre Repetido")
				return
			}
		} else {
			particionvacia = 1
			if &m.particion2 != nil {
				espacioDisponible = int(m.particion2.part_start) - sizeMBR
				if espacioDisponible < int(val) {
					if &m.particion3 == nil {
						if &m.particion4 != nil {
							if (m.particion4.part_start - (m.particion2.part_start + m.particion2.part_size)) >= val {
								particionvacia = 3
								espacioDisponible = int(m.particion4.part_start - (m.particion2.part_start + m.particion2.part_size))
							}
						} else {
							if (m.mbr_tamanio - (m.particion2.part_start + m.particion2.part_size)) >= val {
								particionvacia = 3
								espacioDisponible = int(m.mbr_tamanio) - int(m.particion2.part_start+m.particion2.part_size)
							}
						}

					} else if &m.particion4 == nil {
						if (m.mbr_tamanio - (m.particion3.part_start + m.particion3.part_size)) >= val {
							particionvacia = 4
							espacioDisponible = int(m.mbr_tamanio) - int(m.particion3.part_start+m.particion3.part_size)
						}
					}
				}
			} else {
				if &m.particion3 != nil {
					espacioDisponible = int(m.particion3.part_start) - sizeMBR
				} else {
					if &m.particion4 != nil {
						espacioDisponible = int(m.particion4.part_start) - sizeMBR
					} else {
						espacioDisponible = int(m.mbr_tamanio) - sizeMBR
					}
				}

			}

		}
		//creacion de particion
		if tipo == "p" {
			if espacioDisponible >= int(val) {
				if cantParticionExtendida+cantParitcionesPrimarias < 4 {
					switch particionvacia {
					case 1:
						var start int = int(unsafe.Sizeof(m))
						part := Particion{
							part_status: '1',
							part_type:   'p',
							part_fit:    []byte(fit),
							part_start:  int64(start),
							part_size:   val,
							part_name:   name,
						}
						m.particion1 = part

					case 2:
						var start int = int(unsafe.Sizeof(m)) + part1size
						part := Particion{
							part_status: '1',
							part_type:   'p',
							part_fit:    []byte(fit),
							part_start:  int64(start),
							part_size:   val,
							part_name:   name,
						}
						m.particion2 = part
					case 3:
						var start int = int(unsafe.Sizeof(m)) + part1size + part2size
						part := Particion{
							part_status: '1',
							part_type:   'p',
							part_fit:    []byte(fit),
							part_start:  int64(start),
							part_size:   val,
							part_name:   name,
						}
						m.particion3 = part
					case 4:
						var start int = int(unsafe.Sizeof(m)) + part1size + part2size + part3size
						part := Particion{
							part_status: '1',
							part_type:   'p',
							part_fit:    []byte(fit),
							part_start:  int64(start),
							part_size:   val,
							part_name:   name,
						}
						m.particion4 = part
					default:
					}
				} else {
					fmt.Println("No se puede crear mas particiones primarias")
					return
				}
			} else {
				fmt.Println("Espacio insuficiente en disco")
				return
			}

		} else if tipo == "e" {
			if espacioDisponible >= int(val) {
				if cantParticionExtendida < 1 {
					switch particionvacia {
					case 1:
						var start int = int(unsafe.Sizeof(m))
						part := Particion{
							part_status: '1',
							part_type:   'e',
							part_fit:    []byte(fit),
							part_start:  int64(start),
							part_size:   val,
							part_name:   name,
						}
						m.particion1 = part
						file.Seek(m.particion1.part_start, 0)
					case 2:
						var start int = int(unsafe.Sizeof(m)) + part1size
						part := Particion{
							part_status: '1',
							part_type:   'e',
							part_fit:    []byte(fit),
							part_start:  int64(start),
							part_size:   val,
							part_name:   name,
						}
						m.particion2 = part
						file.Seek(m.particion2.part_start, 0)
					case 3:
						var start int = int(unsafe.Sizeof(m)) + part1size + part2size
						part := Particion{
							part_status: '1',
							part_type:   'e',
							part_fit:    []byte(fit),
							part_start:  int64(start),
							part_size:   val,
							part_name:   name,
						}
						m.particion3 = part
						file.Seek(m.particion3.part_start, 0)
					case 4:
						var start int = int(unsafe.Sizeof(m)) + part1size + part2size + part3size
						part := Particion{
							part_status: '1',
							part_type:   'e',
							part_fit:    []byte(fit),
							part_start:  int64(start),
							part_size:   val,
							part_name:   name,
						}
						m.particion4 = part
						file.Seek(m.particion4.part_start, 0)
					default:
					}
					primerEbr := EBR{
						part_status: '0',
						part_fit:    ' ',
						part_start:  0,
						part_size:   0,
						part_name:   []byte(" "),
						part_next:   -1,
					}
					valoresEBR := &primerEbr
					var binarioEBR bytes.Buffer
					binary.Write(&binarioEBR, binary.BigEndian, valoresEBR)
					escribirBytes(file, binarioEBR.Bytes())
				} else {
					fmt.Println("No se puede crear mas particiones extendias")
				}
			} else {
				fmt.Println("Espacio insuficienta en disco")
				return
			}
		} else if tipo == "l" {
			if cantParticionExtendida > 0 {
				e := EBR{}
				//start donde inicia el primer EBR
				var start int64 = 0
				var sizeextend int64 = 0
				var espacioDisponibleExtended int64 = 0
				//sizeextend
				switch patrExtendida {
				case 1:
					start = int64(unsafe.Sizeof(m))
					sizeextend = m.particion1.part_size
				case 2:
					start = int64(unsafe.Sizeof(m)) + int64(part1size)
					sizeextend = m.particion2.part_size
				case 3:
					start = int64(unsafe.Sizeof(m)) + int64(part1size) + int64(part2size)
					sizeextend = m.particion3.part_size
				case 4:
					start = int64(unsafe.Sizeof(m)) + int64(part1size) + int64(part2size) + int64(part3size)
					sizeextend = m.particion4.part_size
				default:
				}

				//cargando primer EBR
				file.Seek(start, 0)
				var sizeEBR int = int(unsafe.Sizeof(e))
				dataEBR := leerBytes(file, sizeEBR)
				bufferEBR := bytes.NewBuffer(dataEBR)
				err = binary.Read(bufferEBR, binary.BigEndian, &e)
				espacioDisponibleExtended = sizeextend
				if err != nil {
					log.Fatal("binary.Read failed", err)
				}
				for {
					if e.part_next != -1 {
						espacioDisponibleExtended -= e.part_size
						file.Seek(e.part_next, 0)
						data = leerBytes(file, sizeEBR)
						err = binary.Read(bufferEBR, binary.BigEndian, &e)
					} else {
						break
					}

				}
				if espacioDisponibleExtended >= val {
					file.Seek(e.part_start, 0)
					valoresPartLogicaanterior := &e
					var binarioLogicaA bytes.Buffer
					binary.Write(&binarioLogicaA, binary.BigEndian, valoresPartLogicaanterior)

					start = e.part_start + e.part_size
					fitEbr := ' '
					if fit == "bf" {
						fitEbr = 'b'
					} else if fit == "ff" {
						fitEbr = 'f'
					} else if fit == "wf" {
						fitEbr = 'w'
					}
					part := EBR{
						part_status: '1',
						part_fit:    byte(fitEbr),
						part_start:  start,
						part_size:   val,
						part_name:   name,
						part_next:   -1,
					}
					file.Seek(start, 0)
					valoresPartLogicanueva := &part
					var binarioLogicaN bytes.Buffer
					binary.Write(&binarioLogicaN, binary.BigEndian, valoresPartLogicanueva)

				} else {
					fmt.Println("Espacio insuficiente para crear particion")
					return
				}

			} else {
				fmt.Println("No hay particion extendida para crear logica")
				return
			}
		}

		//eliminacion de particion
		if delete != "" {
			encontrada := false
			if &m.particion1 != nil {
				if bytes.Compare(name, m.particion1.part_name) == 0 {
					if delete == "fast" {
						m.particion1 = Particion{}
					} else {
						var vacio int8 = 0
						vaciando := &vacio
						var binarioVaciando bytes.Buffer
						for i := m.particion1.part_start; i < m.particion1.part_size; i++ {
							file.Seek(i, 0)
							binary.Write(&binarioVaciando, binary.BigEndian, vaciando)
						}
					}
					encontrada = true
				}
			}
			if &m.particion2 != nil {
				if bytes.Compare(name, m.particion2.part_name) == 0 {
					if delete == "fast" {
						m.particion2 = Particion{}
					} else {
						var vacio int8 = 0
						vaciando := &vacio
						var binarioVaciando bytes.Buffer
						for i := m.particion2.part_start; i < m.particion2.part_size; i++ {
							file.Seek(i, 0)
							binary.Write(&binarioVaciando, binary.BigEndian, vaciando)
						}
					}
					encontrada = true
				}

			}
			if &m.particion3 != nil {
				if bytes.Compare(name, m.particion3.part_name) == 0 {
					if delete == "fast" {
						m.particion3 = Particion{}
					} else {
						var vacio int8 = 0
						vaciando := &vacio
						var binarioVaciando bytes.Buffer
						for i := m.particion3.part_start; i < m.particion3.part_size; i++ {
							file.Seek(i, 0)
							binary.Write(&binarioVaciando, binary.BigEndian, vaciando)
						}
					}
					encontrada = true
				}
			}
			if &m.particion4 != nil {
				if bytes.Compare(name, m.particion4.part_name) == 0 {
					if delete == "fast" {
						m.particion4 = Particion{}
					} else {
						var vacio int8 = 0
						vaciando := &vacio
						var binarioVaciando bytes.Buffer
						for i := m.particion4.part_start; i < m.particion4.part_size; i++ {
							file.Seek(i, 0)
							binary.Write(&binarioVaciando, binary.BigEndian, vaciando)
						}
					}
					encontrada = true
				}
			}
			if !encontrada {
				if cantParticionExtendida > 0 {
					e := EBR{}
					ebrAnterior := EBR{}
					//start donde inicia el primer EBR
					var start int64 = 0
					var sizeextend int64 = 0
					//sizeextend
					switch patrExtendida {
					case 1:
						start = int64(unsafe.Sizeof(m))
						sizeextend = m.particion1.part_size
					case 2:
						start = int64(unsafe.Sizeof(m)) + int64(part1size)
						sizeextend = m.particion2.part_size
					case 3:
						start = int64(unsafe.Sizeof(m)) + int64(part1size) + int64(part2size)
						sizeextend = m.particion3.part_size
					case 4:
						start = int64(unsafe.Sizeof(m)) + int64(part1size) + int64(part2size) + int64(part3size)
						sizeextend = m.particion4.part_size
					default:
					}

					//Buscando particion logica
					file.Seek(start, 0)
					var sizeEBR int = int(unsafe.Sizeof(e))
					dataEBR := leerBytes(file, sizeEBR)
					bufferEBR := bytes.NewBuffer(dataEBR)
					err = binary.Read(bufferEBR, binary.BigEndian, &e)
					if err != nil {
						log.Fatal("binary.Read failed", err)
					}
					for {
						if bytes.Compare(e.part_name, name) != 0 {
							ebrAnterior = e
							if e.part_next != -1 {
								file.Seek(e.part_next, 0)
								data = leerBytes(file, sizeEBR)
								err = binary.Read(bufferEBR, binary.BigEndian, &e)
							} else {
								fmt.Println("Particion no encontrada")
								break
								return
							}
						} else {
							break
						}

					}

					file.Seek(e.part_start, 0)
					valoresPartLogicaanterior := &e
					var binarioLogicaA bytes.Buffer
					binary.Write(&binarioLogicaA, binary.BigEndian, valoresPartLogicaanterior)

					start = e.part_start + e.part_size
					fitEbr := ' '
					if fit == "bf" {
						fitEbr = 'b'
					} else if fit == "ff" {
						fitEbr = 'f'
					} else if fit == "wf" {
						fitEbr = 'w'
					}
					part := EBR{
						part_status: '1',
						part_fit:    byte(fitEbr),
						part_start:  start,
						part_size:   val,
						part_name:   name,
						part_next:   -1,
					}
					file.Seek(start, 0)
					valoresPartLogicanueva := &part
					var binarioLogicaN bytes.Buffer
					binary.Write(&binarioLogicaN, binary.BigEndian, valoresPartLogicanueva)

				} else {
					fmt.Println("No se encontro particion")
					return
				}
			}

		}
		//reducir o aumentar particion
		//Guardando MBR Modificado
		file.Seek(0, 0)
		s1 := &m
		var binario bytes.Buffer
		binary.Write(&binario, binary.BigEndian, s1)
		escribirBytes(file, binario.Bytes())
	} else {
		fmt.Println("Faltan Parametros")
	}

}

func leerBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number) //array de bytes

	_, err := file.Read(bytes) // Leido -> bytes
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func escribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

func newCryptoRand() int64 {
	safeNum, err := crypto.Int(crypto.Reader, big.NewInt(100234))
	if err != nil {
		panic(err)
	}
	return safeNum.Int64()
}

func graficar() {
	//comando := exec.Command("dot ejemplo1.dot -o ejemplo1.png -Tpng -Gcharset=utf8")
}
