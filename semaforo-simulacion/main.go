package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	// Inicializar componentes
	interfaz := NuevaInterfaz()
	if err := interfaz.CargarTemplates(); err != nil {
		log.Fatalf("Error inicializando interfaz: %v", err)
	}

	semaforoNS := &SemaforoNS{EsVerde: false}
	semaforoEW := &SemaforoEW{EsVerde: true}

	controlador := &ControladorSemaforos{
		SemaforoNS: semaforoNS,
		SemaforoEW: semaforoEW,
		Intervalo:  10 * time.Second,
		Detener:    make(chan bool),
	}

	gestor := &GestorTrafico{
		VehiculosNS:          make([]*Vehiculo, 0),
		VehiculosEW:          make([]*Vehiculo, 0),
		IntervaloGeneracionNS: 1500 * time.Millisecond,
		IntervaloGeneracionEW: 1500 * time.Millisecond,
		DetenerNS:            make(chan bool),
		DetenerEW:            make(chan bool),
	}

	// Iniciar procesos
	controlador.Iniciar()
	gestor.IniciarGeneracionNS()
	gestor.IniciarGeneracionEW()

	// Canal para enviar datos de vehículos a la interfaz
	vehiculosChan := make(chan []VehiculoData, 10)

	// Actualizar interfaz periódicamente
	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		for range ticker.C {
			// Actualizar semáforos
			estadoNS := semaforoNS.Actualizar()
			estadoEW := semaforoEW.Actualizar()
			interfaz.ActualizarEstado(estadoNS, estadoEW)
			
			// Mover vehículos
			gestor.ActualizarVehiculos(estadoNS == "green", estadoEW == "green")
			gestor.EliminarVehiculosFuera()
			
			// Preparar datos de vehículos para la interfaz
			var vehiculosData []VehiculoData
			for _, v := range gestor.VehiculosNS {
				vehiculosData = append(vehiculosData, VehiculoData{
					ID:   v.ID,
					Type: v.Tipo,
					X:    v.PosX,
					Y:    v.PosY,
					Vel:  v.Velocidad,
				})
			}
			for _, v := range gestor.VehiculosEW {
				vehiculosData = append(vehiculosData, VehiculoData{
					ID:   v.ID,
					Type: v.Tipo,
					X:    v.PosX,
					Y:    v.PosY,
					Vel:  v.Velocidad,
				})
			}
			
			// Enviar datos a la interfaz
			select {
			case vehiculosChan <- vehiculosData:
			default:
			}
		}
	}()

	// Manejador para datos de vehículos
	http.HandleFunc("/vehiculos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		select {
		case vehiculos := <-vehiculosChan:
			jsonData, _ := json.Marshal(vehiculos)
			w.Write(jsonData)
		case <-time.After(50 * time.Millisecond):
			w.Write([]byte("[]"))
		}
	})

	// Iniciar servidor web
	interfaz.ServidorWeb()
}

type VehiculoData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Vel  int    `json:"vel"`
}