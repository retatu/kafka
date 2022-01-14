package main

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	deliveryChan := make(chan kafka.Event)
	producer := NewKafkaProducer()
	Publish("transferiu", "teste", producer, []byte("Transferência"), deliveryChan) // com a key (3 parâmetro) a msg vai cair sempre na msm partição
	go DeliveryReport(deliveryChan)
	time.Sleep(2 * time.Second)

	// //sincrono
	// e := <-deliveryChan
	// msg := e.(*kafka.Message)

	// if msg.TopicPartition.Error != nil {
	// 	fmt.Println("Erro ao enviar")
	// } else {
	// 	fmt.Println("Mensagem enviada: ", msg.TopicPartition)
	// }

	//assincrono
}

func NewKafkaProducer() *kafka.Producer {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":   "host.docker.internal:9094",
		"delivery.timeout.ms": "0",    //tempo máximo de entrega de uma msg, 0 é infinito
		"acks":                "all",  //acknowledge: 0, 1 e all. 0 sem retorno, 1 retorno qd o lider recebe, all qd todas as replicações recebem
		"enable.idempotence":  "true", //producor indepotente ou não. Se utilizar como true, o acks precisa ser ALL
	}
	p, err := kafka.NewProducer(configMap)
	if err != nil {
		log.Println(err.Error())
	}
	return p
}

func Publish(msg string, topic string, producer *kafka.Producer, key []byte, deliveryChan chan kafka.Event) error {
	message := &kafka.Message{
		Value:          []byte(msg),
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
	}
	// o resultado da publicação vai para O deliveryChan
	err := producer.Produce(message, deliveryChan)

	if err != nil {
		return err
	}
	return nil
}

func DeliveryReport(deliveryChan chan kafka.Event) {
	for e := range deliveryChan {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Println("Erro ao enviar")
			} else {
				fmt.Println("Mensagem enviada: ", ev.TopicPartition)
			}
		}
	}
}
