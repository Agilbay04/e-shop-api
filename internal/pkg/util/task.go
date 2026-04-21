package util

import (
	"log"
)

// SafeGo runs a goroutine and recovers from panics
func SafeGo(task func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Goroutine Panic] Recovered: %v", r)
			}
		}()
		
		task()
	}()
}