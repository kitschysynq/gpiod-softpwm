package main

import (
	"time"

	pwm "github.com/kitschysynq/gpiod-softpwm"
	"github.com/warthog618/gpiod"
)

func main() {
	c, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		panic(err)
	}

	l, err := c.RequestLine(17, gpiod.AsOutput(0))
	if err != nil {
		panic(err)
	}

	p := pwm.New(l)
	for i := 100; i <= 1000; i += 100 {
		p.Set(i)
		<-time.After(5 * time.Second)
	}
	p.Off()
}
