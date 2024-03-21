//go:build js && wasm

package main

import (
	"bytes"
	"fmt"
	"log"
	"syscall/js"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/abstrack"
	"monks.co/piano-alone/game"
	"monks.co/piano-alone/recorder"
)

const (
	fps            = 30
	screenDuration = time.Second * 4
)

var (
	doc       = js.Global().Get("document")
	overlay   = doc.Call("getElementById", "overlay")
	canv      = doc.Call("querySelector", "canvas")
	usesTouch = js.Global().Call("hasOwnProperty", "ontouchstart").Bool()
)

type GameClient struct {
	loopback            chan<- message
	send                chan<- *game.Message
	fingerprint         string
	state               *game.State
	pianoAnimationReady bool
	myScore             *Score
	tutorialScore       *Score
	myRendition         *smf.SMF
	piano               *Piano
	startPlayingAt      time.Time
	startTutorialAt     time.Time
	recorder            *recorder.Recorder
}

func New(fingerprint string, root js.Value) *GameClient {
	return &GameClient{
		state:       game.NewState(),
		fingerprint: fingerprint,
		recorder:    recorder.New(120),
	}
}

func (c *GameClient) Start(send chan<- *game.Message, recv <-chan *game.Message) error {
	c.send = send

	loopback := make(chan message)
	c.loopback = loopback

	c.piano = NewPiano(doc.Call("getElementById", "piano"), c.loopback)

	send <- game.NewMessage(
		game.MessageTypeJoin,
		c.fingerprint,
		[]byte(c.fingerprint),
	)
	go func() { c.loopback <- msgInit{} }()
	for {
		var m message
		var ok bool
		select {
		case m, ok = <-recv:
		case m, ok = <-loopback:
		}
		if !ok {
			break
		}
		log.Printf("incoming message: '%s'", m)
		if err := c.handleMessage(m); err != nil {
			return fmt.Errorf("error handling message: %w", err)
		}
	}
	log.Printf("done")
	return nil
}

type message interface{}

type (
	msgInit                struct{}
	msgPianoAnimationReady struct{}
	msgShowClickthrough    struct{}
	msgStartHeroIntro      struct{}
	msgStartRecording      struct{}
	msgHeroDone            struct{}
	msgKey                 = recorder.Event
)

func (msgInit) String() string                { return "msgInit" }
func (msgPianoAnimationReady) String() string { return "msgPianoAnimationReady" }
func (msgStartHeroIntro) String() string      { return "msgStartHeroIntro" }
func (msgStartRecording) String() string      { return "msgStartRecording" }
func (msgHeroDone) String() string            { return "msgHeroDone" }

func show(el js.Value, duration string) {
	style := el.Get("style")
	style.Call("setProperty", "transition-duration", duration)
	style.Call("setProperty", "opacity", "100%")
}
func hide(el js.Value, duration string) {
	style := el.Get("style")
	style.Call("setProperty", "transition-duration", duration)
	style.Call("setProperty", "opacity", "0%")
}

func (c *GameClient) handleMessage(m message) error {
	switch m := m.(type) {

	case msgInit:
		go func() {
			show(overlay, "1s")
			time.Sleep(time.Second / 2)
			show(doc.Call("getElementById", "piano"), "1s")
			time.Sleep(time.Second * 2)
			c.loopback <- msgPianoAnimationReady{}
		}()

	case msgPianoAnimationReady:
		c.pianoAnimationReady = true
		if c.myScore != nil {
			go func() { c.loopback <- msgShowClickthrough{} }()
		}

	case msgShowClickthrough:
		go func() {
			hide(overlay, "0.5s")
			time.Sleep(time.Second)

			var handler js.Func
			handler = js.FuncOf(func(_ js.Value, args []js.Value) any {
				args[0].Call("preventDefault")
				c.piano.audio = NewAudio()
				if usesTouch {
					canv.Call("removeEventListener", "touchend", handler)
				} else {
					canv.Call("removeEventListener", "click", handler)
				}
				c.loopback <- msgStartHeroIntro{}
				return nil
			})
			if usesTouch {
				canv.Call("addEventListener", "touchend", handler)
				overlay.Set("innerText", "Tap the screen to continue.")
			} else {
				canv.Call("addEventListener", "click", handler)
				overlay.Set("innerText", "Click anywhere to continue.")
			}
			show(overlay, "0.5s")
			time.Sleep(time.Second * 2)

		}()

	case msgStartHeroIntro:
		keys := c.state.Players[c.fingerprint].Notes
		go func() {
			hide(overlay, "0.5s")
			time.Sleep(time.Second)

			overlay.Set("innerText", "We’re going to play the piano together.")
			show(overlay, "0.5s")
			time.Sleep(time.Second * 2)
			hide(overlay, "0.5s")
			time.Sleep(time.Second)

			overlay.Set("innerText", "Playing all 88 keys is hard.")
			show(overlay, "0.5s")
			time.Sleep(time.Second * 2)
			hide(overlay, "0.5s")
			time.Sleep(time.Second)

			overlay.Set("innerText", "But If we each handle just a few…")
			show(overlay, "0.5s")
			time.Sleep(time.Second / 2)
			c.piano.HideInactiveKeys(keys, time.Second*4)
			time.Sleep(time.Second * 4)
			hide(overlay, "0.5s")
			time.Sleep(time.Second)

			hide(overlay, "0.5s")
			c.piano.Morph(time.Second*2, func(on bool, noteno uint8) {
				now := time.Now()
				var msg midi.Message
				if on {
					msg = midi.NoteOn(1, noteno, 100)
				} else {
					msg = midi.NoteOff(1, noteno)
				}
				c.loopback <- msgKey{Timestamp: now, Message: msg}
			})
			time.Sleep(time.Second * 4)

			overlay.Set("innerText", "Let’s handle these keys here.")
			show(overlay, "0.5s")
			time.Sleep(time.Second * 2)
			hide(overlay, "0.5s")
			time.Sleep(time.Second)

			overlay.Set("innerText", "Notes will fall down towards the keys.")
			show(overlay, "0.5s")
			time.Sleep(time.Second * 2)
			hide(overlay, "0.5s")
			time.Sleep(time.Second)

			overlay.Set("innerText", "When the note reaches the key, press it.")
			show(overlay, "0.5s")
			time.Sleep(time.Second * 2)
			hide(overlay, "0.5s")
			time.Sleep(time.Second)

			c.startTutorialAt = time.Now()
			tutorialDone := c.animate(c.tutorialSceneGraph)
			time.Sleep(screenDuration)
			for i, _ := range keys {
				buttonEl := c.piano.buttons[i].el
				style := buttonEl.Get("style")
				style.Call("setProperty", "border", "solid 15px red")
				time.Sleep(time.Second)
				style.Call("removeProperty", "border")
			}
			time.Sleep(time.Second)
			for i, _ := range keys {
				buttonEl := c.piano.buttons[i].el
				style := buttonEl.Get("style")
				style.Call("setProperty", "border", "solid 15px red")
			}
			time.Sleep(time.Second)
			for i, _ := range keys {
				buttonEl := c.piano.buttons[i].el
				style := buttonEl.Get("style")
				style.Call("removeProperty", "border")
			}

			<-tutorialDone

			overlay.Set("innerText", "I hope you’ll get the hang of it soon.")
			show(overlay, "0.5s")
			time.Sleep(time.Second * 2)
			hide(overlay, "0.5s")
			time.Sleep(time.Second)

			overlay.Set("innerHTML", "Our score is Rachmaninoff’s <em>Prelude in C♯ Minor</em>.")
			show(overlay, "0.5s")
			time.Sleep(time.Second * 2)
			hide(overlay, "0.5s")
			time.Sleep(time.Second)

			overlay.Set("innerHTML", "Let’s go!")
			show(overlay, "0.5s")
			time.Sleep(time.Second * 2)
			hide(overlay, "0.5s")
			time.Sleep(time.Second)

			c.loopback <- msgStartRecording{}

		}()

	case msgStartRecording:
		c.startPlayingAt = time.Now()
		c.animate(c.performanceSceneGraph)

		go func() {
			time.Sleep(c.myScore.NoteTracks[0].Track.Dur() + screenDuration + time.Second)
			c.loopback <- msgHeroDone{}
		}()

	case msgHeroDone:
		c.recorder.Close()
		bs, err := c.recorder.Bytes()
		if err != nil {
			panic(err)
		}
		c.send <- game.NewMessage(game.MessageTypeSubmitPartialTrack, "", bs)

	case msgKey:
		c.recorder.Record(m)

	case *game.Message:
		switch m.Type {

		case game.MessageTypeInitialState:
			c.state = game.StateFromBytes(m.Data)

		case game.MessageTypeBroadcastPhase:
			msg := game.PhaseFromBytes(m.Data)
			c.state.Phase = msg

		case game.MessageTypeBroadcastConnectedPlayer:
			player := game.PlayerFromBytes(m.Data)
			c.state.Players[player.Fingerprint] = player

		case game.MessageTypeBroadcastDisconnectedPlayer:
			c.state.Players[string(m.Data)].ConnectionState = game.ConnectionStateDisconnected

		case game.MessageTypeAssignment:
			me := c.state.Players[c.fingerprint]
			me.Notes = m.Data
			c.myScore = NewScore(abstrack.FromSMF(c.state.Score, 0).Select(me.Notes))
			c.tutorialScore = BuildTutorialScore(me.Notes)
			if c.pianoAnimationReady {
				go func() { c.loopback <- msgShowClickthrough{} }()
			}

		case game.MessageTypeBroadcastSubmittedTrack:
			fingerprint := string(m.Data)
			if _, ok := c.state.Players[fingerprint]; !ok {
				c.state.Players[fingerprint] = &game.Player{}
			}
			c.state.Players[fingerprint].HasSubmitted = true

		case game.MessageTypeBroadcastCombinedTrack:
			r := bytes.NewReader(m.Data)
			rendition, err := smf.ReadFrom(r)
			if err != nil {
				return fmt.Errorf("error reading combined smf: %w", err)
			}
			c.state.Rendition = rendition

		default:
			log.Printf("not handling message (type: '%s')", m.Type.String())
		}
	}
	return nil
}
