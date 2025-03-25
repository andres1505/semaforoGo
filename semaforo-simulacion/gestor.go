package main

import "time"

type GestorTrafico struct {
	VehiculosNS          []*Vehiculo
	VehiculosEW          []*Vehiculo
	IntervaloGeneracionNS time.Duration
	IntervaloGeneracionEW time.Duration
	DetenerNS             chan bool
	DetenerEW             chan bool
}

func (g *GestorTrafico) GenerarVehiculoNS() {
	espacioInicial := 50 // Aumentado para mejor espaciado

	if len(g.VehiculosNS) > 1 {
		ultimo := g.VehiculosNS[len(g.VehiculosNS)-1]
		if ultimo.PosY > -espacioInicial {
			return
		}
	}

	g.VehiculosNS = append(g.VehiculosNS, NuevoVehiculo("ns", 240, -espacioInicial))
}

func (g *GestorTrafico) GenerarVehiculoEW() {
	espacioInicial := 50 // Aumentado para mejor espaciado

	if len(g.VehiculosEW) > 1 {
		ultimo := g.VehiculosEW[len(g.VehiculosEW)-1]
		if ultimo.PosX > -espacioInicial {
			return
		}
	}

	g.VehiculosEW = append(g.VehiculosEW, NuevoVehiculo("ew", -espacioInicial, 245))
}

func (g *GestorTrafico) IniciarGeneracionNS() {
	ticker := time.NewTicker(g.IntervaloGeneracionNS)
	go func() {
		for {
			select {
			case <-ticker.C:
				g.GenerarVehiculoNS()
			case <-g.DetenerNS:
				ticker.Stop()
				return
			}
		}
	}()
}

func (g *GestorTrafico) IniciarGeneracionEW() {
	ticker := time.NewTicker(g.IntervaloGeneracionEW)
	go func() {
		for {
			select {
			case <-ticker.C:
				g.GenerarVehiculoEW()
			case <-g.DetenerEW:
				ticker.Stop()
				return
			}
		}
	}()
}

func (g *GestorTrafico) ActualizarVehiculos(semaforoNSVerde, semaforoEWVerde bool) {
	// Mover vehículos NS (de atrás hacia adelante)
	for i := len(g.VehiculosNS) - 1; i >= 0; i-- {
		var vehiculoDelante *Vehiculo
		if i > 0 {
			vehiculoDelante = g.VehiculosNS[i-1]
		}
		g.VehiculosNS[i].Mover(semaforoNSVerde, vehiculoDelante, g.VehiculosEW)
	}

	// Mover vehículos EW (de atrás hacia adelante)
	for i := len(g.VehiculosEW) - 1; i >= 0; i-- {
		var vehiculoDelante *Vehiculo
		if i > 0 {
			vehiculoDelante = g.VehiculosEW[i-1]
		}
		g.VehiculosEW[i].Mover(semaforoEWVerde, vehiculoDelante, g.VehiculosNS)
	}
}

func (g *GestorTrafico) EliminarVehiculosFuera() {
	// Eliminar vehículos NS que han salido del área
	var nuevosNS []*Vehiculo
	for _, v := range g.VehiculosNS {
		if v.PosY < v.LineaFinal {
			nuevosNS = append(nuevosNS, v)
		}
	}
	g.VehiculosNS = nuevosNS

	// Eliminar vehículos EW que han salido del área
	var nuevosEW []*Vehiculo
	for _, v := range g.VehiculosEW {
		if v.PosX < v.LineaFinal {
			nuevosEW = append(nuevosEW, v)
		}
	}
	g.VehiculosEW = nuevosEW
}