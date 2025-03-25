document.addEventListener('DOMContentLoaded', function() {
    const cityContainer = document.getElementById('city-container');
    const nsStatus = document.getElementById('ns-status');
    const ewStatus = document.getElementById('ew-status');
    const startBtn = document.getElementById('start-btn');
    const stopBtn = document.getElementById('stop-btn');
    
    let simulationRunning = true;
    let vehicles = {};
    
    // Función para actualizar semáforos
    function updateLights() {
        fetch('/estado')
            .then(response => response.text())
            .then(data => {
                const [ns, ew] = data.split(',');
                const nsColor = ns.split(':')[1];
                const ewColor = ew.split(':')[1];
                
                nsStatus.textContent = nsColor === 'green' ? 'verde' : 'rojo';
                ewStatus.textContent = ewColor === 'green' ? 'verde' : 'rojo';
                
                nsStatus.style.color = nsColor;
                ewStatus.style.color = ewColor;
                
                document.querySelector('.ns-light').style.backgroundColor = nsColor;
                document.querySelector('.ew-light').style.backgroundColor = ewColor;
            })
            .catch(error => console.error('Error actualizando semáforos:', error));
    }
    
    // Función para actualizar vehículos
    function updateVehicles() {
        fetch('/vehiculos')
            .then(response => response.json())
            .then(data => {
                // Marcar todos los vehículos como no actualizados
                Object.keys(vehicles).forEach(id => {
                    vehicles[id].updated = false;
                });
                
                // Procesar los vehículos recibidos
                data.forEach(vehicle => {
                    if (!vehicles[vehicle.id]) {
                        // Crear nuevo vehículo
                        const vehicleElement = document.createElement('div');
                        vehicleElement.className = `vehicle ${vehicle.type}`;
                        vehicleElement.id = vehicle.id;
                        
                        // Color basado en tipo y velocidad
                        if(vehicle.type === 'ns') {
                            vehicleElement.style.backgroundColor = `rgb(0, 0, ${150 + (vehicle.vel * 30)})`;
                        } else {
                            vehicleElement.style.backgroundColor = `rgb(${150 + (vehicle.vel * 30)}, 0, 0)`;
                        }
                        
                        cityContainer.appendChild(vehicleElement);
                        vehicles[vehicle.id] = {
                            element: vehicleElement,
                            updated: true
                        };
                    }
                    
                    // Actualizar posición
                    const vehicleElement = vehicles[vehicle.id].element;
                    vehicleElement.style.left = `${vehicle.x}px`;
                    vehicleElement.style.top = `${vehicle.y}px`;
                    vehicles[vehicle.id].updated = true;
                });
                
                // Eliminar vehículos que no fueron actualizados
                Object.keys(vehicles).forEach(id => {
                    if (!vehicles[id].updated) {
                        vehicles[id].element.remove();
                        delete vehicles[id];
                    }
                });
            })
            .catch(error => console.error('Error actualizando vehículos:', error));
    }
    
    // Función principal de actualización
    function updateSimulation() {
        if (simulationRunning) {
            updateLights();
            updateVehicles();
        }
    }
    
    // Configurar botones
    startBtn.addEventListener('click', function() {
        simulationRunning = true;
        startBtn.disabled = true;
        stopBtn.disabled = false;
    });
    
    stopBtn.addEventListener('click', function() {
        simulationRunning = false;
        startBtn.disabled = false;
        stopBtn.disabled = true;
    });
    
    // Iniciar simulación automáticamente
    simulationRunning = true;
    startBtn.disabled = true;
    stopBtn.disabled = false;
    
    // Configurar intervalo de actualización
    setInterval(updateSimulation, 50);
});