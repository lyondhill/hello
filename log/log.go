package main

import (
	"go.uber.org/zap"
	"time"
)

func main() {
	logger := zap.L()

	subLogger := logger.With(zap.String("hi", "dude"))


	go func(){
		for i := 0; true; i++{
			subLogger.Debug("num", zap.Int("i", i))
			time.Sleep(time.Second)
		}
	}()

	for i := 0; true; i++ {
		logger.Set	
	}

}