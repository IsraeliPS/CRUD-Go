package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

type Empleados struct {
	Id     int
	Nombre string
	Correo string
}

func conexionDB() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contrasenia := ""
	Nombre := "sistema"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}

	return conexion
}

func Inicio(w http.ResponseWriter, r *http.Request) {
	fmt.Println("iniciando la API")

	conexionEstablecida := conexionDB()

	registros, err := conexionEstablecida.Query("SELECT * FROM empleados")

	if err != nil {
		panic(err.Error())
	}

	empleado := Empleados{}
	arregloEmpleado := []Empleados{}

	for registros.Next() {
		var id int
		var nombre, correo string
		err = registros.Scan(&id, &nombre, &correo)

		if err != nil {
			panic(err.Error())
		}

		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		arregloEmpleado = append(arregloEmpleado, empleado)
	}
	plantillas.ExecuteTemplate(w, "inicio", arregloEmpleado)
}

func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionDB()
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados(nombre, correo)VALUES(?,?)")

		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(nombre, correo)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func Editar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	conexionEstablecida := conexionDB()
	editaRegistro, err := conexionEstablecida.Query("SELECT *FROM empleados WHERE id=?", idEmpleado)
	empleado := Empleados{}

	if err != nil {
		panic(err.Error())
	}

	for editaRegistro.Next() {
		var id int
		var nombre, correo string
		err = editaRegistro.Scan(&id, &nombre, &correo)

		if err != nil {
			panic(err.Error())
		}

		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
	}
	plantillas.ExecuteTemplate(w, "editar", empleado)
}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionDB()
		modificarRegistros, err := conexionEstablecida.Prepare("UPDATE empleados SET nombre=?, correo=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}

		modificarRegistros.Exec(nombre, correo, id)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	conexionEstablecida := conexionDB()
	borrarRegistro, err := conexionEstablecida.Prepare("DELETE FROM empleados WHERE id=?")

	if err != nil {
		panic(err.Error())
	}
	borrarRegistro.Exec(idEmpleado)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", Inicio)
	router.HandleFunc("/crear", Crear)
	router.HandleFunc("/insertar", Insertar)
	router.HandleFunc("/borrar", Borrar)
	router.HandleFunc("/editar", Editar)
	router.HandleFunc("/actualizar", Actualizar)

	fmt.Print("Servidor corriendo...")
	log.Fatal(http.ListenAndServe(":3000", router))
}
