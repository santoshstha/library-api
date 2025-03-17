package logger

import (
	"log"
	"os"
	"sync"
)

type AsyncLogger struct {
	logger *log.Logger
	ch     chan string
	wg     sync.WaitGroup
}

func NewAsyncLogger() *AsyncLogger {
	l := &AsyncLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		ch:     make(chan string, 100),
	}
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		for msg := range l.ch {
			l.logger.Println(msg)
		}
	}()
	return l
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