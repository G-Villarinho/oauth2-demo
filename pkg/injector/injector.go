package injector

import (
	"fmt"
	"log"

	"go.uber.org/dig"
)

func Provide(container *dig.Container, constructor any) {
	if err := container.Provide(constructor); err != nil {
		log.Fatalf("[dig] failed to Provide: %T: %v", constructor, err)
	}
}

func Resolve[T any](c *dig.Container) T {
	var out T
	if err := c.Invoke(func(dep T) {
		out = dep
	}); err != nil {
		panic(fmt.Sprintf("resolve failed: %v", err))
	}
	return out
}
