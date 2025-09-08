#!/bin/bash


/opt/bitnami/scripts/kafka/setup.sh 
/run.sh &
KAFKA_PID=$!

KAFKA_BIN_DIR="/opt/bitnami/kafka/bin"
KAFKA_API_VERSIONS="${KAFKA_BIN_DIR}/kafka-broker-api-versions.sh"
KAFKA_TOPICS="${KAFKA_BIN_DIR}/kafka-topics.sh"
KAFKA_SCRIPTS_DIR="/opt/bitnami/scripts/kafka/"


MAX_ATTEMPTS=10
ATTEMPT=1
echo "Waiting for Kafka to be ready..."
until ${KAFKA_API_VERSIONS} --bootstrap-server :9092; do
    if [ ${ATTEMPT} -ge ${MAX_ATTEMPTS} ]; then
        echo "Kafka not ready after ${MAX_ATTEMPTS} attempts, exiting."
        exit 1
    fi
    echo "Kafka not ready, attempt ${ATTEMPT}/${MAX_ATTEMPTS}, retrying in 5 seconds..."
    sleep 5
    ATTEMPT=$((ATTEMPT+1))
done

echo "Topics creation -------------------"
$KAFKA_TOPICS --create --if-not-exists --bootstrap-server :9092 --topic new-announcement --partitions 5 
$KAFKA_TOPICS --create --if-not-exists --bootstrap-server :9092 --topic warning-logs --partitions 5
$KAFKA_TOPICS --create --if-not-exists --bootstrap-server :9092 --topic info-logs --partitions 5
$KAFKA_TOPICS --create --if-not-exists --bootstrap-server :9092 --topic error-logs --partitions 5
$KAFKA_TOPICS --create --if-not-exists --bootstrap-server :9092 --topic update-announcement --partitions 5
$KAFKA_TOPICS --create --if-not-exists --bootstrap-server :9092 --topic delete-announcement --partitions 5

echo "All topics created successfully!"

wait $KAFKA_PID