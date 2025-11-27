package adapters

import (
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ encapsula conexión y canal
type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQ() *RabbitMQ {
    var conn *amqp.Connection
    var err error
    
    // Leemos la URL una sola vez
    rabbitURL := os.Getenv("URL_RABBIT")

    // --- BUCLE DE REINTENTOS (La Solución) ---
    maxRetries := 10 // Intentaremos 10 veces
    for i := 0; i < maxRetries; i++ {
        log.Printf("Intentando conectar a RabbitMQ (Intento %d/%d)...", i+1, maxRetries)
        
        conn, err = amqp.Dial(rabbitURL)
        if err == nil {
            // ¡Conexión exitosa! Salimos del bucle
            log.Println("¡Conexión a RabbitMQ exitosa!")
            break
        }

        // Si falló, esperamos 2 segundos antes de volver a intentar
        log.Printf("Fallo al conectar: %v. Esperando 2 segundos...", err)
        time.Sleep(2 * time.Second)
    }
    // -----------------------------------------

    // Si después de 10 intentos (20 segundos) sigue fallando, entonces sí entramos en pánico
    if err != nil {
        log.Fatalf("ERROR FATAL: No se pudo conectar a RabbitMQ después de varios intentos: %v", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("ERROR FATAL: No se pudo abrir el canal: %v", err)
    }

    return &RabbitMQ{conn: conn, ch: ch}
}

func (r *RabbitMQ) Consume(exchange, queue, key string) (<-chan amqp.Delivery, error) {
	err := r.ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	_, err = r.ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	err = r.ch.QueueBind(queue, key, exchange, false, nil)
	if err != nil {
		return nil, err
	}
	return r.ch.Consume(queue, "", true, false, false, false, nil)
}