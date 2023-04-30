package main_test

import (
	"dev07"
	"fmt"
	"testing"
	"time"
)

func TestOr(t *testing.T) {

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	out := main.Or(
		// sig(2*time.Hour),
		// sig(5*time.Minute),
		sig(1*time.Second),
		sig(2*time.Second),
		// sig(1*time.Hour),
		// sig(1*time.Minute),
	)
	fmt.Println(<-out)

	took := time.Since(start)
	fmt.Printf("done after %v\n", took.Seconds())

	if took < time.Duration(1*time.Second) || took > time.Duration(2*time.Second) {
		t.Error("\nГоруитны завершились слишком быстро или слишком поздно, прошло - ", took)
	}
}
