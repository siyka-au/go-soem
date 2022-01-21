package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/siyka-au/go-soem/soem"
)

/*
 * Slave layout for testing
 * EK1100 > EL1008 > EL1004 > EL2008 > EL2004
 */
func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()
	go func() {
		select {
		case <-signalChan: // first signal, cancel context
			cancel()
		case <-ctx.Done():
		}
		<-signalChan // second signal, hard exit
		os.Exit(10)
	}()
	if err := run(ctx, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(11)
	}
}

func run(ctx context.Context, args []string) error {

	master, err := soem.NewSOEMMaster(args[1])
	if err != nil {
		return err
	}
	defer master.Close()

	master.ConfigInit()
	fmt.Printf("Found %d attached slaves\n", master.SlaveCount)

	// if master.ConfigDC() && master.Slaves[2].HasDC {
	// 	master.DCSync0(2, 200*time.Millisecond, 0)
	// }

	master.ConfigMap(1024)

	if err := stateCheck(master, soem.EC_STATE_SAFE_OP); err != nil {
		fmt.Println(err)
	}

	printSlaveDetails(master)

	// send one valid process data to make outputs of the slaves happy
	master.SendProcessData()
	wkc := master.ReceiveProcessData(soem.EC_TIMEOUTRET)
	fmt.Printf("WKC: %d after ReceiveProcessData\n", wkc)

	if wkc, err := master.SetState(soem.EC_STATE_OPERATIONAL); err != nil {
		return err
	} else {
		fmt.Printf("WKC: %d after SetState()\n", wkc)
	}

	if err := stateCheck(master, soem.EC_STATE_OPERATIONAL); err != nil {
		fmt.Println(err)
	}

	ticker := time.NewTicker(200 * time.Millisecond)

	go func() {
		const light0Max uint8 = 1 << 7
		const light1Max uint8 = 1 << 3
		// light0DirUp := true
		// light1DirUp := true
		// lights0 := uint8(1)
		// lights1 := uint8(1)

		// calcDir := func(dir bool, max, min, val uint8) bool {
		// 	return (!(val == max) && (val == min)) || (dir && !(val == max))
		// }

		// stepLight := func(dir bool, val uint8) uint8 {
		// 	if dir {
		// 		return val << 1
		// 	} else {
		// 		return val >> 1
		// 	}
		// }

		/*
		 * Main PDO loop
		 * Print inputs as binary strings
		 * Bit shift a bit through the output range
		 */
		for {
			select {
			case <-ticker.C:
				fmt.Println("Processing I/O")
				master.SendProcessData()
				master.ReceiveProcessData(soem.EC_TIMEOUTRET)

				// fmt.Printf("Inputs: %08b %08b\r", master.Slaves[1].Read()[0], master.Slaves[2].Read()[0])

				// light0DirUp = calcDir(light0DirUp, light0Max, 1, lights0)
				// light1DirUp = calcDir(light1DirUp, light1Max, 1, lights1)

				// lights0 = stepLight(light0DirUp, lights0)
				// lights1 = stepLight(light1DirUp, lights1)

				// master.Slaves[3].Write([]byte{lights0})
				// master.Slaves[4].Write([]byte{lights1})

				master.Slaves[0].Write([]byte{1, 2, 3, 4, 5, 6, 7, 8})
			case <-ctx.Done():
				ticker.Stop()
				return
			}

		}
	}()

	<-ctx.Done()
	if _, err := master.SetState(soem.EC_STATE_INIT); err != nil {
		return err
	}
	if err := stateCheck(master, soem.EC_STATE_INIT); err != nil {
		return err
	}
	return nil

}

func stateCheck(master *soem.Master, state soem.EtherCATState) error {
	actState, err := master.CheckState(0, state, soem.EC_TIMEOUTSTATE)
	if err != nil {
		master.ReadState()
		for i, slave := range master.Slaves {
			if slave.State != state {
				return fmt.Errorf("slave %d has state %s, expected %s", i, actState, state)
			}
		}
	}

	fmt.Printf("Current State: %s\n", state)
	return nil
}

func printSlaveDetails(master *soem.Master) {
	for i, slave := range master.Slaves {
		fmt.Printf("Slave %d\n%s\n", i+1, "  "+strings.ReplaceAll(slave.String(), "\n", "\n  "))
	}
}
