package main

import (
	"log"
	"sync"
	"syscall/js"

	"monks.co/piano-alone/notefreqs"
)

type Audio struct {
	mu          sync.Mutex
	notes       []uint8
	context     js.Value
	oscillators map[uint8]*oscillator
}

type oscillator struct {
	handle      js.Value
	isConnected bool
}

func NewAudio() *Audio {
	a := &Audio{}
	a.context = js.Global().Get("AudioContext").New()
	a.oscillators = map[uint8]*oscillator{}
	return a
}

func (a *Audio) NoteOn(note uint8) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, hasOsc := a.oscillators[note]; !hasOsc {
		osc := a.context.Call("createOscillator")
		osc.Set("type", "triangle")
		osc.Get("frequency").Set("value", notefreqs.NoteFreqs[note])
		osc.Call("start")
		a.oscillators[note] = &oscillator{handle: osc}
	}

	osc := a.oscillators[note]
	if !osc.isConnected {
		log.Printf("connect %d", note)
		osc.isConnected = true
		osc.handle.Call("connect", a.context.Get("destination"))
	}
}

func (a *Audio) NoteOff(note uint8) {
	a.mu.Lock()
	defer a.mu.Unlock()

	osc, hasOsc := a.oscillators[note]
	if !hasOsc || !osc.isConnected {
		return
	}

	log.Printf("disconnect %d", note)
	osc.isConnected = false
	osc.handle.Call("disconnect", a.context.Get("destination"))
}
