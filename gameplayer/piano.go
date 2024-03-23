package main

import (
	"fmt"
	"log"
	"syscall/js"
	"time"
)

type Piano struct {
	announceState chan<- message
	state         pianoState
	el            js.Value
	noteset       map[uint8]struct{}
	keys          []*keyEl
	buttons       []*keyEl
	mouseDown     bool
	audio         *Audio
}

//go:generate go run golang.org/x/tools/cmd/stringer -type=pianoState
type pianoState int

const (
	pianoStateAfterInit pianoState = iota
	pianoStateDuringHideInactiveKeys
	pianoStateAfterHideInactiveKeys
	pianoStateDuringMorph
	pianoStateAfterMorph
)

func NewPiano(container js.Value, announceState chan<- message) *Piano {
	piano := &Piano{
		announceState: announceState,
		el:            container,
		keys:          make([]*keyEl, 88),
	}
	for i := 0; i < 88; i++ {
		noteno := uint8(i + 21)
		value := keyOrder[i%12]
		el := doc.Call("createElement", "div")
		el.Get("classList").Call("add", "key")
		el.Get("classList").Call("add", value.color.String())
		el.Get("classList").Call("add", value.name)
		container.Call("appendChild", el)
		piano.keys[i] = &keyEl{
			color:  value.color,
			name:   value.name,
			noteno: noteno,
			el:     el,
		}
	}
	return piano
}

func (p *Piano) HideInactiveKeys(activeKeys []uint8, dur time.Duration) {
	p.changeState(pianoStateAfterInit, pianoStateDuringHideInactiveKeys)
	p.noteset = make(map[uint8]struct{}, len(activeKeys))
	for _, n := range activeKeys {
		p.noteset[n] = struct{}{}
	}
	var secondBatch []*keyEl
	for _, key := range p.keys {
		if _, ok := p.noteset[key.noteno]; ok {
			p.buttons = append(p.buttons, &keyEl{
				color:  key.color,
				name:   key.name,
				noteno: key.noteno,
				el:     key.el,
			})
		} else {
			if key.color == keyColorWhite {
				fadeOut(key.el, dur.Abs().Milliseconds()*2/3)
			} else {
				secondBatch = append(secondBatch, key)
			}
		}
	}
	go func() {
		time.Sleep(dur / 3)
		for _, key := range secondBatch {
			fadeOut(key.el, dur.Abs().Milliseconds()*2/3)
		}
		time.Sleep(dur * 2 / 3)
		p.changeState(pianoStateDuringHideInactiveKeys, pianoStateAfterHideInactiveKeys)
	}()
}

func (p *Piano) Morph(dur time.Duration, f func(bool, uint8)) {
	p.changeState(pianoStateAfterHideInactiveKeys, pianoStateDuringMorph)
	p.state = pianoStateDuringMorph

	// add buttons
	for _, key := range p.buttons {
		rect := key.el.Call("getBoundingClientRect")

		key.el = doc.Call("createElement", "div")

		key.el.Get("classList").Call("add", "button")
		style := key.el.Get("style")
		style.Call("setProperty", "left", fmt.Sprintf("%fpx", rect.Get("left").Float()))
		style.Call("setProperty", "width", fmt.Sprintf("%fpx", rect.Get("width").Float()))
		style.Call("setProperty", "height", fmt.Sprintf("%fpx", rect.Get("height").Float()))
		style.Call("setProperty", "background-color", key.color.String())
		style.Call("setProperty", "transition-duration", fmt.Sprintf("%dms", dur.Abs().Milliseconds()))
		style.Call("setProperty", "transition-property", "left, width, height, background-color")
		p.el.Call("appendChild", key.el)
	}

	// remove keys
	for _, key := range p.keys {
		p.el.Call("removeChild", key.el)
	}

	// Wait one second for some reason?
	time.Sleep(time.Second)

	// move buttons
	for i, key := range p.buttons {
		style := key.el.Get("style")
		style.Call("setProperty", "background-color", "white")
		style.Call("setProperty", "left", fmt.Sprintf("calc(%d * 100%% / %d)", i, len(p.buttons)))
		style.Call("setProperty", "width", fmt.Sprintf("calc(100%% / %d)", len(p.buttons)))
		style.Call("setProperty", "height", "100%")
	}
	go func() {
		time.Sleep(dur)
		for _, key := range p.buttons {
			style := key.el.Get("style")
			style.Call("setProperty", "transition-property", "none")
		}
		handler := js.FuncOf(func(_ js.Value, args []js.Value) any {
			var (
				ev          = args[0]
				evType      = ev.Get("type").String()
				screenWidth = p.el.Get("clientWidth").Float()
				keyCount    = float64(len(p.buttons))
				keyWidth    = screenWidth / keyCount
				touched     = map[uint8]struct{}{}
			)
			ev.Call("preventDefault")
			switch evType {
			case "mouseup":
				p.mouseDown = false
			case "mousedown":
				p.mouseDown = true
				fallthrough
			case "mousemove":
				if p.mouseDown {
					x := ev.Get("clientX").Int()
					k := x / int(keyWidth)
					touched[p.buttons[k].noteno] = struct{}{}
				}
			case "touchstart", "touchmove", "touchend":
				touches := ev.Get("touches")
				for i := 0; i < touches.Get("length").Int(); i++ {
					touch := touches.Call("item", i)
					x := touch.Get("clientX").Int()
					k := x / int(keyWidth)
					touched[p.buttons[k].noteno] = struct{}{}
				}
			}
			for _, k := range p.buttons {
				_, isTouched := touched[k.noteno]
				if k.isOn && !isTouched {
					k.el.Get("style").Call("setProperty", "background-color", "white")
					k.isOn = false
					if p.audio != nil {
						p.audio.NoteOff(k.noteno)
					}
					f(false, k.noteno)
				} else if !k.isOn && isTouched {
					k.el.Get("style").Call("setProperty", "background-color", "#dddddd")
					k.isOn = true
					if p.audio != nil {
						p.audio.NoteOn(k.noteno)
					}
					f(true, k.noteno)
				}
			}
			return nil
		})
		if usesTouch {
			log.Println("using touch")
			p.el.Call("addEventListener", "touchstart", handler)
			p.el.Call("addEventListener", "touchend", handler)
			p.el.Call("addEventListener", "touchmove", handler)
		} else {
			log.Println("using mouse")
			p.el.Call("addEventListener", "mousedown", handler)
			p.el.Call("addEventListener", "mouseup", handler)
			p.el.Call("addEventListener", "mousemove", handler)
		}
		p.changeState(pianoStateDuringMorph, pianoStateAfterMorph)
	}()
}

func (p *Piano) changeState(before, after pianoState) {
	if p.state != before {
		panic(fmt.Sprintf("invalid piano state; in %s, expected %s->%s", p.state, before, after))
	}
	p.state = after
	go func() { p.announceState <- after }()
}

type keyEl struct {
	color  keyColor
	name   string
	noteno uint8
	el     js.Value
	isOn   bool
}

func fadeIn(el js.Value, durationMS int64) {
	el.Get("style").Call("setProperty", "transition-duration", fmt.Sprintf("%dms", durationMS))
	el.Get("style").Call("setProperty", "opacity", "100%")
}

func fadeOut(el js.Value, durationMS int64) {
	el.Get("style").Call("setProperty", "transition-duration", fmt.Sprintf("%dms", durationMS))
	el.Get("style").Call("setProperty", "opacity", "0")
}

type keyColor int

func (c keyColor) String() string {
	switch c {
	case keyColorWhite:
		return "white"
	case keyColorBlack:
		return "black"
	default:
		return "pink"
	}
}

const (
	keyColorWhite keyColor = iota
	keyColorBlack
)

var keyOrder = []struct {
	color keyColor
	name  string
}{
	{keyColorWhite, "A"},
	{keyColorBlack, "As"},
	{keyColorWhite, "B"},
	{keyColorWhite, "C"},
	{keyColorBlack, "Cs"},
	{keyColorWhite, "D"},
	{keyColorBlack, "Ds"},
	{keyColorWhite, "E"},
	{keyColorWhite, "F"},
	{keyColorBlack, "Fs"},
	{keyColorWhite, "G"},
	{keyColorBlack, "Gs"},
}
