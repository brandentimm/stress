package main

import (
	`flag`
	`fmt`
	`os`
	`os/signal`
	`runtime`
)

var load int

// On init, bind load to the '-load' flag, then parse the flags
func init() {
	flag.IntVar(&load, `load`, 1, `Amount of load to generate`)
	flag.Parse()
}

func spinUntilCancel(cancel chan bool) {
	var i int64 = 0
	for {
		select {
		case <-cancel:
			return
		default:
			i++
		}
	}
}

func main() {
	// Create a channel to notify on SIGINT or SIGKILL
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	// Create a channel to cancel running goroutines
	cancel := make(chan bool, 1)

	// Set golang's max number of processes to load
	runtime.GOMAXPROCS(load)

	fmt.Printf("Stress running with requested load: %d\n", load)
	fmt.Println("Press ctrl-c at any time to exit.")

	for i := 0; i < load; i++ {
		go spinUntilCancel(cancel)
	}

	select {
	case s := <-sig:
		fmt.Printf("Caught signal %s, exiting.", s)
		cancel <- true
	}
	os.Exit(0)
}
