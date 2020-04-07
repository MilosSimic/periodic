# periodic
Simple task scheduling lib


# Usage:
```go

import (
"fmt"
"time"
)

func main() {
	p := New(10)

	timer1 := time.NewTimer(20 * time.Second)

	p.Repeat(5*time.Second, func(d ...interface{}, 1, 2) {
		fmt.Println("Repeat 5 ", d)
	})

	p.Repeat(10*time.Second, func(d ...interface{}) {
		fmt.Println("Repeat 10")
	})

	p.Once(10*time.Second, func(d ...interface{}) {
		fmt.Println("Once 10")
	})

	<-timer1.C
}

```
