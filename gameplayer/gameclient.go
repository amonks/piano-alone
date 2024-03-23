//go:build js && wasm

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"syscall/js"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
	"monks.co/piano-alone/abstrack"
	"monks.co/piano-alone/baseurl"
	"monks.co/piano-alone/data"
	"monks.co/piano-alone/game"
	"monks.co/piano-alone/recorder"
)

const (
	fps            = 30
	screenDuration = time.Second * 4
)

var (
	doc = js.Global().Get("document")

	html = doc.Call("querySelector", "html")

	page         = doc.Call("getElementById", "page")
	performances = doc.Call("getElementById", "performances")

	app     = doc.Call("getElementById", "app")
	overlay = doc.Call("getElementById", "overlay")
	canv    = doc.Call("querySelector", "canvas")

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
	}
}

func (c *GameClient) Start(send chan<- *game.Message, recv <-chan *game.Message) error {
	c.send = send

	loopback := make(chan message)
	c.loopback = loopback

	c.piano = NewPiano(doc.Call("getElementById", "piano"), c.loopback)

	noteCapacity := 1
	if usesTouch {
		noteCapacity = 3
	}
	send <- game.NewMessage(
		game.MessageTypeJoin,
		c.fingerprint,
		game.JoinMessage{
			Fingerprint:  c.fingerprint,
			NoteCapacity: noteCapacity,
		}.Bytes(),
	)
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
	msgLookAtDisklavier    struct{}
	msgPerformanceIsOver   struct{}
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
			hide(page, "1s")
			time.Sleep(time.Second)
			html.Get("classList").Call("add", "app")

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
				overlay.Set("innerText", "Tap the screen to begin.")
			} else {
				canv.Call("addEventListener", "click", handler)
				overlay.Set("innerText", "Click anywhere to begin.")
			}
			show(overlay, "0.5s")
			time.Sleep(time.Second * 2)

		}()

	case msgStartHeroIntro:
		keys := c.state.Players[c.fingerprint].AssignedNotes
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

			overlay.Set("innerText", "But if we each handle just a few…")
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

			if len(c.state.Players[c.fingerprint].AssignedNotes) > 1 {
				overlay.Set("innerText", "Focus on these keys here.")
			} else {
				overlay.Set("innerText", "Focus on this key here.")
			}
			show(overlay, "0.5s")
			time.Sleep(time.Second * 2)
			hide(overlay, "0.5s")
			time.Sleep(time.Second)

			if len(c.state.Players[c.fingerprint].AssignedNotes) > 1 {
				overlay.Set("innerText", "Notes will fall down towards the keys.")
			} else {
				overlay.Set("innerText", "Notes will fall down towards the key.")
			}
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
			for i := range keys {
				buttonEl := c.piano.buttons[i].el
				style := buttonEl.Get("style")
				style.Call("setProperty", "border", "solid 15px red")
				time.Sleep(time.Second)
				style.Call("removeProperty", "border")
			}
			time.Sleep(time.Second)
			for i := range keys {
				buttonEl := c.piano.buttons[i].el
				style := buttonEl.Get("style")
				style.Call("setProperty", "border", "solid 15px red")
			}
			time.Sleep(time.Second)
			for i := range keys {
				buttonEl := c.piano.buttons[i].el
				style := buttonEl.Get("style")
				style.Call("removeProperty", "border")
			}

			<-tutorialDone
			time.Sleep(time.Second)

			overlay.Set("innerText", "That’s all there is to it!")
			show(overlay, "0.5s")
			time.Sleep(time.Second * 2)
			hide(overlay, "0.5s")
			time.Sleep(time.Second)

			overlay.Set("innerHTML", "Our score is Rachmaninoff’s <em>Prelude in C♯ Minor</em>.")
			show(overlay, "0.5s")
			time.Sleep(time.Second * 2)
			hide(overlay, "0.5s")
			time.Sleep(time.Second)

			overlay.Set("innerHTML", "We’ll each record our parts alone.")
			show(overlay, "0.5s")
			time.Sleep(time.Second * 2)
			hide(overlay, "0.5s")
			time.Sleep(time.Second)

			overlay.Set("innerHTML", "Then, we’ll listen to all our parts together on the Disklavier.")
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
		c.recorder = recorder.New(120)
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
		c.send <- game.NewMessage(game.MessageTypeSubmitPartialTrack, c.fingerprint, bs)
		go func() { c.loopback <- msgLookAtDisklavier{} }()

	case msgLookAtDisklavier:
		overlay.Set("innerHTML", "Done!<br />When everyone else is finished, we’ll hear our performance on the disklavier.")
		show(overlay, "0.5s")
		time.Sleep(time.Second * 2)

	case msgPerformanceIsOver:
		go func() {
			if resp, err := http.Get(baseurl.NoHost.Rest(data.PathFeaturedPerformances)); err == nil {
				defer resp.Body.Close()
				if html, err := io.ReadAll(resp.Body); err == nil {
					fmt.Println(string(html))
					performances.Set("outerHTML", string(html))
				}
			}
			hide(app, "1s")
			time.Sleep(time.Second)

			html.Get("classList").Call("remove", "app")
			show(page, "1s")
		}()

	case msgKey:
		if c.recorder != nil {
			c.recorder.Record(m)
		}

	case *game.Message:
		switch m.Type {

		case game.MessageTypeState:
			c.state = game.StateFromBytes(m.Data)
			if c.state.Phase != game.GamePhaseUninitialized &&
				c.state.Phase != game.GamePhaseDone {
				go func() { c.loopback <- msgInit{} }()
			}

		case game.MessageTypeBroadcastPhase:
			before := c.state.Phase
			after := game.PhaseFromBytes(m.Data)
			c.state.Phase = after
			if before == game.GamePhaseUninitialized && after != game.GamePhaseUninitialized {
				go func() { c.loopback <- msgInit{} }()
			} else {
				switch after {
				case game.GamePhaseProcessing:
					go func() { c.loopback <- msgLookAtDisklavier{} }()
				case game.GamePhaseDone:
					go func() { c.loopback <- msgPerformanceIsOver{} }()
				}
			}

		case game.MessageTypeBroadcastConnectedPlayer:
			player := game.PlayerFromBytes(m.Data)
			c.state.Players[player.Fingerprint] = player

		case game.MessageTypeBroadcastDisconnectedPlayer:
			c.state.Players[string(m.Data)].ConnectionState = game.ConnectionStateDisconnected

		case game.MessageTypeAssignment:
			me := c.state.Players[c.fingerprint]
			me.AssignedNotes = m.Data
			r := bytes.NewReader(c.state.Configuration.Score)
			smf, err := smf.ReadFrom(r)
			if err != nil {
				return err
			}
			c.myScore = NewScore(abstrack.FromSMF(smf, 0).Select(me.AssignedNotes))
			c.tutorialScore = BuildTutorialScore(me.AssignedNotes)
			if c.pianoAnimationReady {
				go func() { c.loopback <- msgShowClickthrough{} }()
			}

		case game.MessageTypeBroadcastSubmittedTrack:
			fingerprint := string(m.Data)
			if _, ok := c.state.Players[fingerprint]; !ok {
				c.state.Players[fingerprint] = &game.Player{}
			}
			c.state.Players[fingerprint].HasSubmitted = true

		default:
			log.Printf("not handling message (type: '%s')", m.Type.String())
		}
	}
	return nil
}
