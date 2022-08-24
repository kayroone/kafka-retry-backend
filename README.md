# kafka-retry-backend

A simple retry mechanism for failed Kafka messages written in Golang

# retry-app

Messages that could not be processed by components, for example, because a further processing component or the database
is not available, are "parked" in a so-called retry topic. The Retry app runs in a interval as K8's CronJob, consumes
the messages from the retry topic and forwards them to the input topic for reprocessing.

## Usage

The Retry app is configured with the following parameters via environment variables, CLI parameters or config file
entries
configure. The app checks for the presence of the parameters according to the following priority:

```shell
ENV > CLI > Config
```

Folgende Parameter stehen zur Verfügung:

    kafka_retry_topic --> Retry topic from which to be consumed
    kafka_eingangs_topic --> Input topic to be published to for reprocessing
    kafka_deadletter_topic --> DeadLetter topic to be swapped out to the messages after X attempted processing (retries)
    kafka_brokers --> Comma separated string of Kafka Broker addresses
    kafka_user --> Kafka Username 
    kafka_password --> Kafka Passwort
    kafka_consumer_group_id --> Consumer Group ID that is set on the consumer. Necessary if several instances are operated in parallel
    kafka_retrys_header_name --> Name of the Kafka Message Header in which the number of retries are recorded
    kafka_failed_consumer_group_id_header_name --> Name of the Kafka Message Header in which the Consumer Group ID is stored in the message previously failed
    kafka_batch_size --> How many messages should be collected in a batch before it is published?
    kafka_retry_limit --> After how many retries should a message be moved out to the DeadLetter Topic?
    kafka_auth_enabled --> Activate Kafka Auth?

Beispiel CLI Aufruf:

    RETRY --kafka_retry_topic=fehlerdienst.retry --kafka_eingangs_topic=fehlerdienst.fehlereingang --kafka_consumer_group_id=0 --kafka_brokers=10.0.0.0:9092,11.0.0.0:9092,12.0.0.0:9092 --kafka_auth_enabled=true --kafka_user=user --kafka_password=password --kafka_batch_size=5 --kafka_deadletter_topic=fehlerdienst.deadletter --kafka_retrys_header_name=retrys --kafka_failed_consumer_group_id_header_name=failed.consumer.group.id 

Example ENV call (ENV variables ALWAYS start with the prefix **RETRY_**). Here is an excerpt from the cron.yaml for K8s:

    containers:
            - name: retry-app
              image: retry-app:latest
              imagePullPolicy: Always
              env:
                - name: RETRY_KAFKA_BROKERS
                  valueFrom:
                    configMapKeyRef:
                      key: kafka_brokers
                      name: kafka-brokers
                - name: RETRY_KAFKA_RETRY_TOPIC
                  valueFrom:
                    configMapKeyRef:
                      key: topic.fehlerdienst.retry
                      name: kafka-topics
                - name: RETRY_KAFKA_EINGANGS_TOPIC
                  valueFrom:
                    configMapKeyRef:
                      key: topic.fehlereingang
                      name: kafka-topics
                - name: RETRY_KAFKA_DEADLETTER_TOPIC
                  valueFrom:
                    configMapKeyRef:
                      key: topic.fehlerdienst.deadletter
                      name: kafka-topics
                - name: RETRY_KAFKA_AUTH_ENABLED
                  value: "false"
                - name: RETRY_KAFKA_CONSUMER_GROUP_ID
                  value: fehlerdienst
                - name: RETRY_KAFKA_BATCH_SIZE
                  value: "5"
                - name: RETRY_KAFKA_FAILED_CONSUMER_GROUP_ID_HEADER_NAME
                  value: "failed.consumer.group.id"
                - name: RETRY_KAFKA_RETRYS_HEADER_NAME
                  value: "retrys"
