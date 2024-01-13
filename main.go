package main

import (
	"fmt"
	"os"
	"time"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("done")
	os.Exit(0)
}

func run() error {
	defer midi.CloseDriver()

	ports := midi.GetOutPorts()
	for _, port := range ports {
		fmt.Printf("%d: %s\n", port.Number(), port.String())
	}

	if len(ports) == 0 {
		return fmt.Errorf("no midi output ports found")
	}

	port, err := midi.FindOutPort(ports[0].String())
	if err != nil {
		return err
	}

	send, err := midi.SendTo(port)
	if err != nil {
		return err
	}

	if err := send(midi.NoteOn(1, 55, 100)); err != nil {
		return err
	}

	time.Sleep(time.Second)

	if err := send(midi.NoteOff(1, 55)); err != nil {
		return err
	}

	return nil
}
