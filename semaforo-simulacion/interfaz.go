package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

// Interfaz maneja la representación visual de la simulación
type Interfaz struct {
	CityContainer  string
	Intersection   string
	templates      *template.Template
	mu             sync.Mutex
	EstadoSemaforo struct {
		NS string
		EW string
	}
}

// NuevaInterfaz crea una nueva instancia de Interfaz
func NuevaInterfaz() *Interfaz {
	return &Interfaz{
		CityContainer: "city-container",
		Intersection:  "intersection",
		EstadoSemaforo: struct {
			NS string
			EW string
		}{
			NS: "red",
			EW: "green",
		},
	}
}

// CargarTemplates carga las plantillas HTML
func (i *Interfaz) CargarTemplates() error {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		return fmt.Errorf("error cargando plantillas: %v", err)
	}
	i.templates = tmpl
	return nil
}

// ActualizarEstado actualiza el estado visual de los semáforos
func (i *Interfaz) ActualizarEstado(ns, ew string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.EstadoSemaforo.NS = ns
	i.EstadoSemaforo.EW = ew
}

// ServidorWeb inicia el servidor HTTP
func (i *Interfaz) ServidorWeb() {
	http.HandleFunc("/", i.ManejadorPrincipal)
	http.HandleFunc("/estado", i.ManejadorEstado)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	
	fmt.Println("Servidor iniciado en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// ManejadorPrincipal maneja las solicitudes a la página principal
func (i *Interfaz) ManejadorPrincipal(w http.ResponseWriter, r *http.Request) {
	i.templates.ExecuteTemplate(w, "index.html", nil)
}

// ManejadorEstado devuelve el estado actual de los semáforos
func (i *Interfaz) ManejadorEstado(w http.ResponseWriter, r *http.Request) {
	i.mu.Lock()
	defer i.mu.Unlock()
	fmt.Fprintf(w, "NS:%s,EW:%s", i.EstadoSemaforo.NS, i.EstadoSemaforo.EW)
}