# Publish a message to create a new hive
rabbitmqadmin publish exchange=amq.default \
    routing_key="hive_queue" \
    properties="{\"delivery_mode\": 2, \"content_type\": \"application/json\"}" \
    payload='{
        "hive_id": 99,
        "sensors": [
            {
                "sensor_id": 10,
                "hive_id": 99,
                "sensor_type": "temperature",
                "last_reading": "",
                "last_reading_time": null
            },
            {
                "sensor_id": 20,
                "hive_id": 99,
                "sensor_type": "humidity",
                "last_reading": "",
                "last_reading_time": null
            },
            {
                "sensor_id": 30,
                "hive_id": 99,
                "sensor_type": "temperature",
                "last_reading": "",
                "last_reading_time": null
            }
        ]
    }'
