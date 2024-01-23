package event_consumer

import (
	"log"
	"time"

	"github.com/wlcmtunknwndth/gobot_manager/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize) // нужен ретрай в фетчере(3 попытки) -- экспоненциально или с конст задержкой
		if err != nil {
			// допилить -- проблема с сетью -- мимо данные. Чувствительно, если задержка большая
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}
		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)
			continue
		}
	}
}

/*
1. Потеря событий -- ретраи, возвращение в хранилище, фолбэк, подтверждение для фетчера(не обрабатывает, пока не будет сигнала с инфой об обработке)
2. обработка всей пачки: останавилваемся после одной ошибки или считаем до n ошибок
3. параллельная обработка
*/

func (c *Consumer) handleEvents(events []events.Event) error {
	// waitgroup: sync.WaitGroup{}
	for _, event := range events {
		log.Printf("got event: %s", event.Text)

		if err := c.processor.Process(event); err != nil { //  можно добавить ретрай и/или бэкап
			log.Printf("can't handle event: %s", err.Error())
			continue // network issues -- упали все события
		}

	}
	return nil
}
