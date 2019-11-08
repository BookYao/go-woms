package mqmodule

import (
	"../../github.com/streadway/amqp"
	"fmt"
)

const (
	url = "amqp://baseuser:basepasswd@192.168.100.186:5800/vhost_basemq"
	exchange = "exc_basemq"
	exchangeType = "direct"
)

func MqPublish(user string, msg []byte) bool {
	connect, err := amqp.Dial(url)
	if err != nil {
		fmt.Println("AMQP Dial error:", err.Error())
		return false
	}
	defer connect.Close()

	channel, err := connect.Channel()
	if err != nil {
		fmt.Println("AMQP Channel error. ", err.Error())
		return false
	}
	defer channel.Close()

	err = channel.ExchangeDeclarePassive(exchange, exchangeType, true, false, false, false, nil)
	if err != nil {
		err = channel.ExchangeDeclare(exchange, exchangeType, true, false, false, false, nil)
		if err != nil {
			fmt.Println("ExchangeDeclare failed.", err.Error())
			return false
		}
	}

	queueName := "q_" + user
	fmt.Println("QueueName:", queueName)
	_, err = channel.QueueDeclarePassive(queueName,true, false, false, false, nil)
	if err != nil {
		_, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
		if err != nil {
			fmt.Println("Declare queue failed.", err.Error())
			return false
		}

		fmt.Println("Queue Declare is inPassive...")
	}

	fmt.Println("Queue Is Exist...")
	err = channel.QueueBind(queueName, user, exchange, false, nil)
	if err != nil {
		fmt.Println("QueueBind failed.", err.Error())
		return false
	}

	err = channel.Publish(exchange, user, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: msg,
		})

	fmt.Println("Publish Msg: ", string(msg))
	return true
}
