package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	path := os.Args[1]
	args := os.Args[2:]
	flagSet := flag.NewFlagSet("flag", flag.ExitOnError)
	bulk := flagSet.Bool("bulk", false, "Crea un lote de archivos")
	simu := flagSet.Bool("simu", false, "Simula el SCADA")
	flagSet.Parse(args)
	contenido, err := ioutil.ReadFile("plantilla.txt")
	if err != nil {
		log.Fatal(err)
	}
	if *bulk {
		i := 111
		for i < 1000 {
			name := strconv.Itoa(i)
			nuevoContenido := strings.Replace(string(contenido), "XXXXX", name, 1)
			err = ioutil.WriteFile(filepath.Join(path, name+".txt"), []byte(nuevoContenido), 0666)
			if err != nil {
				log.Fatal(err)
			}
			i++
		}
	}
	if *simu {
		i := 1000
		for {
			name := strconv.Itoa(i)
			nuevoContenido := strings.Replace(string(contenido), "XXXXX", name, 1)
			err = ioutil.WriteFile(filepath.Join(path, name+".txt"), []byte(nuevoContenido), 0666)
			if err != nil {
				log.Fatal(err)
			}
			i++
			time.Sleep(30 * time.Second)
		}
	}
}
