package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Dato struct {
	Parametro string `json:"parametro"`
	Valor     string `json:"valor"`
}

func borrarTabla(db *sql.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS Archivos")
	if err != nil {
		log.Fatal(err)
	}
}
func crearTabla(db *sql.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS Archivos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL,
		fechaCreacion DATE,
		data TEXT
	);`)
	if err != nil {
		log.Fatal(err)
	}
}
func insertarDatos(db *sql.DB, parametros []string, valores []string) {
	var datos []string
	for i := 0; i < len(parametros); i++ {
		dato, _ := json.Marshal(Dato{Parametro: parametros[i], Valor: valores[i]})
		datos = append(datos, string(dato))
	}
	_, err := db.Exec("INSERT INTO Archivos (nombre, fechaCreacion, data) VALUES (?, ?, ?)",
		valores[0], time.Now(), strings.Join(datos, ";"))
	if err != nil {
		log.Fatal(err)
	}
}
func mostrarBaseDatos(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM Archivos")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var nombre string
		var fechaCreacion string
		var data string
		err = rows.Scan(&id, &nombre, &fechaCreacion, &data)
		if err != nil {
			log.Fatal(err)
		}
		datosExtraidos := strings.Split(data, ";")
		var datos []Dato
		for _, datoJSON := range datosExtraidos {
			var dato Dato
			if err := json.Unmarshal([]byte(datoJSON), &dato); err != nil {
				log.Println("Error al decodificar dato JSON:", err)
				continue
			}
			datos = append(datos, dato)
		}
		fmt.Printf("ID: %d ,Nombre: %s, FechaCreación: %s, Data: %#v\n", id, nombre, fechaCreacion, data)
	}
}
func conectar() *sql.DB {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
func lectura(path string) ([]string, []string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	i := 1
	parametros := []string{}
	valores := []string{}
	for scanner.Scan() {
		if i%2 != 0 {
			parametros = append(parametros, scanner.Text())
		} else {
			valores = append(valores, scanner.Text())
		}
		i++
	}
	return parametros, valores
}
func main() {
	path := os.Args[1]
	args := os.Args[2:]
	flagSet := flag.NewFlagSet("flag", flag.ExitOnError)
	create := flagSet.Bool("create", false, "Crea la tabla")
	borrar := flagSet.Bool("borrar", false, "Borra la tabla")
	insert := flagSet.Bool("insert", false, "Inserta registros en la tabla")
	show := flagSet.Bool("show", false, "Muestra todos los registros de la tabla")
	flagSet.Parse(args)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	db := conectar()
	defer db.Close()
	if *borrar {
		borrarTabla(db)
	}
	if *create {
		crearTabla(db)
	}
	if *insert {
		for _, file := range files {
			parametros, valores := lectura(filepath.Join(path, file.Name()))
			insertarDatos(db, parametros, valores)
		}
	}
	if *show {
		mostrarBaseDatos(db)
	}
}
