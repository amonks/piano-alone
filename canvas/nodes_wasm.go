//go:build js && wasm

package canvas

import (
	"math"

	"monks.co/piano-alone/c2d"
)

func Fill(color string, children ...SceneNode) SceneNodeFunc {
	return SceneNodeFunc(func(c2d c2d.C2D, bounds Bounds) {
		c2d.BeginPath()
		c2d.SetFillStyle(color)
		for _, c := range children {
			c.Draw(c2d, bounds)
		}
		c2d.ClosePath()
		c2d.Fill0()
	})
}

func Stroke(color string, children ...SceneNode) SceneNodeFunc {
	return SceneNodeFunc(func(c2d c2d.C2D, bounds Bounds) {
		c2d.BeginPath()
		c2d.SetStrokeStyle(color)
		for _, c := range children {
			c.Draw(c2d, bounds)
		}
		c2d.ClosePath()
		c2d.Stroke0()
	})
}

func Text(font string, x, y float64, str string) SceneNodeFunc {
	return SceneNodeFunc(func(c2d c2d.C2D, bounds Bounds) {
		position := bounds.Proportional(Bounds{x, y, 0, 0})
		c2d.SetFont(font)
		c2d.FillText3(str, position.X, position.Y)
	})
}

func Dot(x, y float64, radius int) SceneNodeFunc {
	return SceneNodeFunc(func(c2d c2d.C2D, bounds Bounds) {
		position := bounds.Proportional(Bounds{x, y, 0, 0})
		c2d.Ellipse7(position.X, position.Y, 10, 10, 0, 0, 2*math.Pi)
		c2d.Stroke0()
	})
}

func Rect(b Bounds) SceneNodeFunc {
	return SceneNodeFunc(func(c2d c2d.C2D, bounds Bounds) {
		c2d.SetLineWidth(1)
		rect := bounds.Proportional(b)
		c2d.FillRect(rect.X, rect.Y, rect.Width, rect.Height)
	})
}

func Box(b Bounds) SceneNodeFunc {
	return SceneNodeFunc(func(c2d c2d.C2D, bounds Bounds) {
		c2d.SetLineWidth(1)
		rect := bounds.Proportional(b)
		c2d.StrokeRect(rect.X, rect.Y, rect.Width, rect.Height)
	})
}

type Group struct {
	bounds   Bounds
	children []SceneNode
}

func G(b Bounds, children ...SceneNode) *Group {
	return &Group{bounds: b, children: children}
}

func (c *Group) Draw(c2d c2d.C2D, bounds Bounds) {
	for _, child := range c.children {
		child.Draw(c2d, bounds.Proportional(c.bounds))
	}
}
