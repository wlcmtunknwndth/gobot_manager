package event_consumer

import (
	"context"
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

// Start —— starts bot event handler
func (c Consumer) Start(ctx context.Context) error {
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
		if err := c.handleEvents(ctx, gotEvents); err != nil {
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

func (c *Consumer) handleEvents(ctx context.Context, events []events.Event) error {
	// waitgroup: sync.WaitGroup{}
	for _, event := range events {
		log.Printf("got event: %s", event.Text)

		if err := c.processor.Process(ctx, event); err != nil { //  можно добавить ретрай и/или бэкап
			log.Printf("can't handle event: %s", err.Error())
			continue // network issues -- упали все события
		}

	}
	return nil
}
