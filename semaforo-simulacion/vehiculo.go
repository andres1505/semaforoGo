package main

import (
	"math"
	"math/rand"
	"time"
)

const (
	DistanciaMinima     = 25
	AnchoInterseccion   = 100
	AltoInterseccion    = 100
	PosicionInterseccion = 200
	MargenFrenado       = 30  // Margen adicional para zona de frenado
)

type Vehiculo struct {
	ID             string
	Tipo           string
	PosX           int
	PosY           int
	Velocidad      int
	LineaDetencion int
	LineaFinal     int
	Detenido       bool
}

func NuevoVehiculo(tipo string, x, y int) *Vehiculo {
	rand.Seed(time.Now().UnixNano())
	lineaDetencion := 180
	if tipo == "ew" {
		lineaDetencion = 180 // Ajustar según necesidad
	}
	
	return &Vehiculo{
		ID:             generarID(),
		Tipo:           tipo,
		PosX:           x,
		PosY:           y,
		Velocidad:      4,
		LineaDetencion: lineaDetencion,
		LineaFinal:     530,
		Detenido:       false,
	}
}

func (v *Vehiculo) EnInterseccion() bool {
	if v.Tipo == "ns" {
		return v.PosY >= PosicionInterseccion && v.PosY <= PosicionInterseccion+AltoInterseccion
	}
	return v.PosX >= PosicionInterseccion && v.PosX <= PosicionInterseccion+AnchoInterseccion
}

func (v *Vehiculo) EnZonaFrenado() bool {
	if v.Tipo == "ns" {
		return v.PosY >= v.LineaDetencion-MargenFrenado && v.PosY <= v.LineaDetencion+MargenFrenado
	}
	return v.PosX >= v.LineaDetencion-MargenFrenado && v.PosX <= v.LineaDetencion+MargenFrenado
}

func (v *Vehiculo) Mover(semaforoVerde bool, vehiculoDelante *Vehiculo, vehiculosContrarios []*Vehiculo) {
	switch v.Tipo {
	case "ns":
		v.moverNS(semaforoVerde, vehiculoDelante, vehiculosContrarios)
	case "ew":
		v.moverEW(semaforoVerde, vehiculoDelante, vehiculosContrarios)
	}
}

func (v *Vehiculo) moverNS(semaforoVerde bool, vehiculoDelante *Vehiculo, vehiculosEW []*Vehiculo) {
	if v.PosY >= v.LineaFinal {
		return
	}

	// Verificar colisión en intersección
	if v.EnInterseccion() {
		for _, ew := range vehiculosEW {
			if ew.EnInterseccion() && distancia(v.PosX, v.PosY, ew.PosX, ew.PosY) < DistanciaMinima {
				v.Detenido = true
				return
			}
		}
	}

	// Verificar si debe detenerse (semáforo o vehículo delante)
	if (!semaforoVerde && v.EnZonaFrenado()) || 
	   (vehiculoDelante != nil && v.PosY >= vehiculoDelante.PosY-DistanciaMinima) {
		v.Detenido = true
		return
	}

	v.Detenido = false
	v.PosY += v.Velocidad
}

func (v *Vehiculo) moverEW(semaforoVerde bool, vehiculoDelante *Vehiculo, vehiculosNS []*Vehiculo) {
	if v.PosX >= v.LineaFinal {
		return
	}

	// Verificar colisión en intersección
	if v.EnInterseccion() {
		for _, ns := range vehiculosNS {
			if ns.EnInterseccion() && distancia(v.PosX, v.PosY, ns.PosX, ns.PosY) < DistanciaMinima {
				v.Detenido = true
				return
			}
		}
	}

	// Verificar si debe detenerse (semáforo o vehículo delante)
	if (!semaforoVerde && v.EnZonaFrenado()) || 
	   (vehiculoDelante != nil && v.PosX >= vehiculoDelante.PosX-DistanciaMinima) {
		v.Detenido = true
		return
	}

	v.Detenido = false
	v.PosX += v.Velocidad
}

func distancia(x1, y1, x2, y2 int) int {
	dx := x1 - x2
	dy := y1 - y2
	return int(math.Sqrt(float64(dx*dx + dy*dy)))
}

func generarID() string {
	return string(rune(65+rand.Intn(26))) + string(rune(65+rand.Intn(26)))
}