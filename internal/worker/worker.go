package worker

import (
	"api_shope/utils/helper"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var ctx = context.Background()

type Worker struct {
	DB      *gorm.DB
	Redis   *redis.Client
	ticker  *time.Ticker
	quit    chan struct{}
	running bool
	mu      sync.Mutex
}

func NewWorker(db *gorm.DB, redis *redis.Client) *Worker {
	return &Worker{
		DB:    db,
		Redis: redis,
	}
}

func (w *Worker) flushPendingItems() {
	keys, _ := w.Redis.Keys(ctx, "behind:pending:*").Result()
	for _, key := range keys {
		data, _ := w.Redis.HGetAll(ctx, key).Result()
		if len(data) == 0 {
			continue
		}

		op := data["op"]
		switch op {
		case "register":
			fmt.Printf("%v masuk ke queue redis register", data["username"])
			helper.SendEmail(data["email"], data["message"])
		case "buy":
			fmt.Printf("%v masuk ke queue redis buy product", data["username"])
			helper.SendEmail(data["email"], data["message"])
		}

		w.Redis.Del(ctx, key)
	}
}

func (w *Worker) StartFlushWorker(interval time.Duration) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.running {
		return
	}

	w.ticker = time.NewTicker(interval)
	w.quit = make(chan struct{})
	w.running = true

	go func() {
		for {
			select {
			case <-w.ticker.C:
				w.flushPendingItems()
			case <-w.quit:
				w.ticker.Stop()
				return
			}
		}
	}()
}

func (w *Worker) StopFlushWorker() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.running {
		return
	}

	close(w.quit)
	w.running = false
}
