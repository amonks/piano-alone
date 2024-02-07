//go:build js && wasm

package vdom

import (
	"monks.co/piano-alone/c2d"
)

func Fill(color string, children ...SceneNode) SceneNode {
	return SceneNodeFunc(func(c2d c2d.C2D, bounds Bounds) {
		c2d.SetFillStyle(color)
		for _, c := range children {
			c.Draw(c2d, bounds)
		}
		c2d.Fill0()
	})
}

func Stroke(color string, children ...SceneNode) SceneNode {
	return SceneNodeFunc(func(c2d c2d.C2D, bounds Bounds) {
		c2d.SetStrokeStyle(color)
		for _, c := range children {
			c.Draw(c2d, bounds)
		}
		c2d.Stroke0()
	})
}

func Box(b Bounds) SceneNode {
	return SceneNodeFunc(func(c2d c2d.C2D, bounds Bounds) {
		c2d.SetLineWidth(1)
		rect := b.Within(bounds)
		c2d.StrokeRect(rect.X, rect.Y, rect.Width, rect.Height)
	})
}
