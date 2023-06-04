package rabbit

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"tranc/server/internal/entity"
)

type Rabbit struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func NewRabbit() (*Rabbit, error) {
	conn, err := amqp.Dial("amqp://test:pass@localhost:5672/")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	//err = ch.ExchangeDeclare(
	//	"tranc",
	//	"direct",
	//	true,
	//	false,
	//	false,
	//	false,
	//	nil)

	return &Rabbit{
		Conn: conn,
		Ch:   ch,
	}, nil
}

func (r *Rabbit) SendMessage(id int, data *entity.Tranc) error {
	message, err := json.Marshal(data)

	queue, err := r.Ch.QueueDeclare(
		strconv.Itoa(id),
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	err = r.Ch.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *Rabbit) ReadMessage(id int, complete chan bool) (chan *entity.Tranc, error) {
	msgs, err := r.Ch.Consume(
		strconv.Itoa(id),
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
		return nil, err
	}

	chMsg := make(chan *entity.Tranc)
	go func() error {
		for msg := range msgs {
			q := msg.RoutingKey
			log.Printf("Received transaction: %s from queue %s", msg.Body, q)
			var data *entity.Tranc
			err := json.Unmarshal(msg.Body, &data)
			if err != nil {
				log.Println("error with unmarshal json: ", err)
				return err
			}
			chMsg <- data
			if <-complete {
				err := msg.Ack(false)
				if err != nil {
					log.Println("error with key: ", err)
					return err
				}
			}
		}
		close(chMsg)
		return nil
	}()
	return chMsg, nil
}
