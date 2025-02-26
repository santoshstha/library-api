package logger

import (
	"log"
	"os"
	"sync"
)

var Logger *AsyncLogger // Export Logger globally

type AsyncLogger struct {
	logger *log.Logger
	ch     chan string
	wg     sync.WaitGroup
}

func InitLogger() { // Initialize the global Logger
	if Logger == nil {
		Logger = &AsyncLogger{
			logger: log.New(os.Stdout, "", log.LstdFlags),
			ch:     make(chan string, 100),
		}
		Logger.wg.Add(1)
		go func() {
			defer Logger.wg.Done()
			for msg := range Logger.ch {
				Logger.logger.Println(msg)
			}
		}()
	}
}

func (l *AsyncLogger) Log(msg string) {
	select {
	case l.ch <- msg:
	default:
		l.logger.Println(msg)
	}
}

func (l *AsyncLogger) Close() {
	close(l.ch)
	l.wg.Wait()
}