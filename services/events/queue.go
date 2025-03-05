package events

import (
	"log"

	"github.com/realjv3/event-agg/domain"
)

const maxConcurrent = 3
const queueLen = 100

type Queue struct {
	concurrent chan int
	queue      chan *Event
}

type Event struct {
	Event *domain.Event
	Dest  func(event *domain.Event) error
}

func NewQueue() *Queue {
	return &Queue{
		concurrent: make(chan int, maxConcurrent),
		queue:      make(chan *Event, queueLen),
	}
}

func (q *Queue) QueueEvent(event Event) {
	q.queue <- &event
	log.Println("Event added to queue")
}

func (q *Queue) Process() {
	for {
		event := <-q.queue
		q.concurrent <- 1

		log.Println("Processing Event")

		go func(event *Event, sem chan int) {
			log.Println("Sending event to destination")

			err := event.Dest(event.Event)
			if err != nil {
				log.Printf("Error from destination: %v\n", err)
			}

			<-sem
		}(event, q.concurrent)
	}
}
