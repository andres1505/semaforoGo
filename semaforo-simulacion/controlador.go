package main

import "time"

// ControladorSemaforos gestiona la sincronización de los semáforos
type ControladorSemaforos struct {
	SemaforoNS    *SemaforoNS
	SemaforoEW    *SemaforoEW
	Intervalo     time.Duration
	Detener       chan bool
}

// Iniciar comienza el ciclo de cambios de semáforos
func (c *ControladorSemaforos) Iniciar() {
	ticker := time.NewTicker(c.Intervalo)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.ActualizarSemaforos()
			case <-c.Detener:
				ticker.Stop()
				return
			}
		}
	}()
}

// ActualizarSemaforos alterna los estados de los semáforos
func (c *ControladorSemaforos) ActualizarSemaforos() {
	c.SemaforoNS.mu.Lock()
	c.SemaforoEW.mu.Lock()
	
	c.SemaforoNS.EsVerde = !c.SemaforoNS.EsVerde
	c.SemaforoEW.EsVerde = !c.SemaforoEW.EsVerde
	
	c.SemaforoEW.mu.Unlock()
	c.SemaforoNS.mu.Unlock()
}