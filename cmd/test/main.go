package main

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/pior/runnable"
)

func printf(format string, a ...interface{}) {
	// fmt.Printf(format, a...)
}

type ServerNoShutdown struct{}

func (s *ServerNoShutdown) Run(ctx context.Context) error {
	printf("%T: start\n", s)
	<-make(chan struct{})

	printf("%T: stop\n", s)
	return nil
}

type ServerPanic struct{}

func (s *ServerPanic) Run(ctx context.Context) error {
	printf("%T: start\n", s)

	time.Sleep(time.Second * 1)
	panic("yooooolooooooo")
}

type Server struct {
	deadline time.Duration
}

func (s *Server) Run(ctx context.Context) error {
	printf("%T: start\n", s)

	if s.deadline.Seconds() == 0 {
		s.deadline = time.Second * 10000000
	}

	theEnd := time.After(s.deadline)

	select {
	case <-ctx.Done():
	case <-theEnd:
		printf("%T: sepuku\n", s)
		return errors.New("sepuku")
	}

	printf("%T: stop\n", s)
	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()

	runner := runnable.Signal(runnable.Group(
		&Server{deadline: time.Millisecond * 1500},
		&Server{deadline: time.Millisecond * 2000},
		&Server{deadline: time.Millisecond * 2500},
		&ServerNoShutdown{},
		&ServerPanic{},
	))

	err := runner.Run(ctx)

	if err != nil {
		printf("Error: %s", err)
		os.Exit(1)
	}

	os.Exit(0)
}
