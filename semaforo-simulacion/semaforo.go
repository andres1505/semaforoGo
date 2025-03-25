package main

import "sync"

// SemaforoNS representa el semáforo en dirección Norte-Sur
type SemaforoNS struct {
	ElementoHTMLElement string
	EsVerde            bool
	mu                 sync.Mutex
}

// Actualizar actualiza el estado del semáforo NS
func (s *SemaforoNS) Actualizar() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.EsVerde {
		return "green"
	}
	return "red"
}

// SemaforoEW representa el semáforo en dirección Este-Oeste
type SemaforoEW struct {
	ElementoHTMLElement string
	EsVerde            bool
	mu                 sync.Mutex
}

// Actualizar actualiza el estado del semáforo EW
func (s *SemaforoEW) Actualizar() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.EsVerde {
		return "green"
	}
	return "red"
}