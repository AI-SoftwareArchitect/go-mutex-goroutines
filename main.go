package main

import (
	"fmt"
	"sync"
	"time"
)

var mRW = sync.RWMutex{}
var wg = sync.WaitGroup{}
var cursor int64 = 0
var data = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

// Artış sonrası değeri iletmek için kanal ekliyoruz.
var ch = make(chan int64, len(data))

func main() { 
	t0 := time.Now()

	for i := 0; i < len(data); i++ {
		wg.Add(2)
		go openWorker()
		go read()
	}
	
	wg.Wait()
	fmt.Println("Time: ", time.Since(t0))
}

func openWorker() {
	// Her çağrıda yalnızca bir artış yapılacak.
	if get() < int64(len(data)) {
		save()
	}
}

func save() {
	mRW.Lock()
	cursor++
	// Güncellenen değeri kanala gönderiyoruz.
	ch <- cursor
	mRW.Unlock()
	wg.Done()
}

func read() {
	// Kanal üzerinden sırayla gönderilen değeri alıyoruz.
	val := <- ch
	fmt.Println("Cursor: ", val)
	wg.Done()
}

func get() int64 {
	mRW.RLock()
	defer mRW.RUnlock()
	return cursor
}
