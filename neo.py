import pika

# 1) Conexión a RabbitMQ
connection = pika.BlockingConnection(
    pika.ConnectionParameters(host='34.235.24.218')
)
channel = connection.channel()

# 2) Asegúrate de que el exchange existe
channel.exchange_declare(
    exchange='amq.topic',
    exchange_type='topic',
    durable=True
)

# 3) Declara la cola (durable y compartida)
result = channel.queue_declare(
    queue='sensor_NEO',
    durable=True,
    exclusive=False
)
queue_name = result.method.queue

# 4) Enlaza la cola al topic MQTT que te interesa
channel.queue_bind(
    exchange='amq.topic',
    queue=queue_name,
    routing_key='neo'  # o '#' para todo
)

print(' [*] Esperando mensajes en amq.topic. Para salir CTRL+C')

# 5) Callback para procesar mensajes
def callback(ch, method, properties, body):
    print(f" [x] Recibido {body}")

channel.basic_consume(
    queue=queue_name,
    on_message_callback=callback,
    auto_ack=True
)

# 6) Comienza a consumir
channel.start_consuming()
