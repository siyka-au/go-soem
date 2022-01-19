package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/siyka-au/go-soem/soem"
)

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

	master.ConfigMap(1024)

	if err := stateCheck(master, soem.EC_STATE_SAFE_OP); err != nil {
		fmt.Println(err)
	}

	// printSlaveDetails(master)

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

	go func() {
		slave := master.Slaves[1]

		for {
			fmt.Print("Inputs: ")
			fmt.Println(slave.Read())
			time.Sleep(1 * time.Second)
		}

		//

		// tmp = bytearray([0 for i in range(output_len)])

		// toggle = True
		// try:
		//     while 1:
		//         if toggle:
		//             tmp[0] = 0x00
		//         else:
		//             tmp[0] = 0x02
		//         self._master.slaves[2].output = bytes(tmp)

		//         toggle ^= True

		//         time.sleep(1)
	}()

	for {
		select {
		case <-ctx.Done():
			if _, err := master.SetState(soem.EC_STATE_INIT); err != nil {
				return err
			}
			if err := stateCheck(master, soem.EC_STATE_INIT); err != nil {
				return err
			}
			return nil
		default:
			// Process data
			master.SendProcessData()
			master.ReceiveProcessData(soem.EC_TIMEOUTRET)
			time.Sleep(10 * time.Millisecond)
		}
	}
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
		fmt.Printf(
			"Slave %d Name %s\n"+
				"  Vendor ID 0x%08x\n"+
				"  Product Code 0x%08x\n"+
				"  Revision 0x%08x\n"+
				"  Configured Address 0x%04x\n"+
				"  Alias Address 0x%04x\n"+
				"  Input Bits %d\n"+
				"  Input Bytes %d\n"+
				"  Output Bits %d\n"+
				"  Output Bytes %d\n"+
				"  Configured Address 0x\n",
			i, slave.Name, slave.VendorID, slave.ProductCode, slave.Revision,
			slave.ConfiguredAddress, slave.AliasAddress,
			slave.InputBits, slave.InputBytes,
			slave.OutputBits, slave.OutputBytes)
	}
}
