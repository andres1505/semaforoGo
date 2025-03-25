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
//Se usa una goroutine para manejar el ciclo automático de los semáforos sin bloquear el hilo principal:
//Se crea un ticker que se ejecuta cada Intervalo de tiempo y se ejecuta un ciclo select que espera a que el ticker envíe una señal o que se envíe una señal para detener el ciclo.
//Cuando se recibe una señal de detener, se detiene el ticker y la goroutine termina.
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