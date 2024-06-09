# README

## 1. Install all dependencies

Use command `go mod download` to install all the dependencies.

## 2. Create .env file

Include the following:

```
KAFKA_BROKERS=localhost:9092
KAFKA_CRYPTO_TOPIC=crypto_data

BINANCE_SPOT_WS_URL=wss://stream.binance.com:443/ws/!ticker@arr
BINANCE_USD_M_FUT_WS_URL=wss://fstream.binance.com/ws/!ticker@arr
BINANCE_COIN_M_FUT_WS_URL=wss://dstream.binance.com/ws/!ticker@arr

WEBSOCKET_FAIL_THRESHOLD=10
MESSAGE_FAIL_THRESHOLD=10
KAFKA_CONNECTION_FAIL_THRESHOLD = 5
```

## 3. Start service locally using docker compose
It starts the kafka, zookeeper & crypto-pulse service locally and expose kafka on port `9092`
<br>
<br>
`NOTE : No need to expose the crypto-pulse service on any port as we only want the data being streamed to the kafka to be consumed by the consumer`
```bash
docker compose up -d
```

## 4 . Run Kafka Console Consumer
To see if the messages are being consumed by consumer perfectly, run kafka console consumer.
```
docker ps

docker exec -it kafka /bin/bash

kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic crypto_data
```

![alt text](screenshots/image.png)