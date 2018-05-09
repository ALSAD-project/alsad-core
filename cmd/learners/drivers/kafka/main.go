// Kafka Driver
//
// Example use:
// $ export DISPATCHER_LISTEN_URL=":9999"
// $ export KAFKA_BROKER_URL=":9093"
// $ ./kafka

package main

import (
    "bufio"
    "fmt"
    "strings"
    "log"
    "os"

    "github.com/kelseyhightower/envconfig"
    "github.com/Shopify/sarama"
)

func publish(message string, producer sarama.AsyncProducer) {

    msg := &sarama.ProducerMessage {
        Topic: "streamin",
        Value: sarama.StringEncoder(message),
    }

    producer.Input() <- msg
    select {
    case <-producer.Successes():
        fmt.Printf("[Produced] topic: %s, value: %s\n", 
                    msg.Topic, msg.Value)
    case err := <-producer.Errors():
        log.Fatalln(err)
    }

}

func main() {

    // Load ENV configurations
    driverConfig := config{}
    if err := envconfig.Process("driver", &driverConfig); err != nil {
        log.Fatalf("Error on processing configuration: %s", err.Error())
        return
    }

    // Kafka settings
    brokers := []string{driverConfig.KafkaBrokerURL}
    config := sarama.NewConfig()
    config.Producer.Return.Errors = true
    config.Producer.Return.Successes = true
    config.Producer.Retry.Max = 3
    config.Consumer.Return.Errors = true

    // Consumer routine
    consumer, err := sarama.NewConsumer(brokers, config)
    if err != nil {
        panic(err)
    }

    defer func() {
        if err := consumer.Close(); err != nil {
            panic(err)
        }
    }()

    partition, err := consumer.ConsumePartition("streamout", 0, sarama.OffsetNewest)
    if err != nil {
        panic(err)
    }

    go func() {
        for {
            msg := <-partition.Messages()
            fmt.Printf("[Consumed] topic: %s, offset: %d, key: %s, value: %s\n", 
                    msg.Topic, msg.Offset, msg.Key, msg.Value)
        }
    }()

    // Producer routine
    producer, err := sarama.NewAsyncProducer(brokers, config)
    if err != nil {
        panic(err)
    }

    defer producer.AsyncClose()

    reader := bufio.NewReader(os.Stdin)
    for {
        msg, _ := reader.ReadString('\n')
        msg = strings.TrimSuffix(msg, "\n")

        go publish(msg, producer)
    }

}