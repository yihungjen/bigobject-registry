package main

import (
	"github.com/yihungjen/registry-driver"
	"github.com/yihungjen/registry-driver/driver/rabbitmq"
	"github.com/yihungjen/registry-driver/driver/redis"
	"log"
)

func init() {
	go func() {
		queue := make(map[string]chan registry.Serializer)

		// spawn worker to monitor redis keyspace
		regevent := redis.KeySpaceEventLoop()

		for action := range regevent {
			if dest, ok := queue[action.Tier]; !ok {
				queue[action.Tier] = make(chan registry.Serializer, 4)
				queue[action.Tier] <- action
				rabbitmq.NamedSendLoop(action.Tier, queue[action.Tier])
			} else {
				dest <- action
			}
		}
	}()
}

func main() {
	log.Println("begin the redis keyspace tracker to RabbitMQ")

	client := redis.NewClient()
	defer client.Close()

	actions, confirm := rabbitmq.NamedRecvLoop("registry-complete")
	for oneact := range actions {
		var res registry.ResourceEvent
		confirm <- &rabbitmq.ConfirmedEvent{
			DTag: oneact.DTag,
			Ack:  true,
		}
		if err := oneact.Deserialize(&res); err != nil {
			log.Println(err)
			continue
		}
		if err := client.Append(res.VendorName, string(oneact.Val)).Err(); err != nil {
			log.Println(err)
			continue
		}
		log.Println("update:", res)
	}
}
