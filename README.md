# semaforoGo
Instalar Go,
get mod tidy,
go run main.go,

Se usa gouritines en controlador.go para el cambio de semaforos, en gestor.go para la generacion de vehiculos, y en main.go para la actualizacion de la interfaz

Resumen de Concurrencia
Componente	Goroutines	Propósito
ControladorSemaforos	1	Cambiar semáforos en intervalos fijos.
GestorTrafico (NS/EW)	2	Generar vehículos en ambas direcciones.
main.go	1	Actualizar estado de la simulación.
Servidor HTTP	1+ (internas)	Manejar peticiones web concurrentes.

