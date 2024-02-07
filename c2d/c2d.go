package c2d

import "syscall/js"

type C2D js.Value

func (c2d C2D) GetCanvas() js.Value {
	return js.Value(c2d).Get("canvas")
}

func (c2d C2D) SetCanvas(val any) {
	js.Value(c2d).Set("canvas", val)
}

func (c2d C2D) GetGlobalAlpha() float64 {
	return js.Value(c2d).Get("globalAlpha").Float()
}

func (c2d C2D) SetGlobalAlpha(val float64) {
	js.Value(c2d).Set("globalAlpha", val)
}

func (c2d C2D) GetGlobalCompositeOperation() GlobalCompositeOperation {
	return js.Value(c2d).Get("globalCompositeOperation").String()
}

func (c2d C2D) SetGlobalCompositeOperation(val GlobalCompositeOperation) {
	js.Value(c2d).Set("globalCompositeOperation", val)
}

func (c2d C2D) GetFillStyle() js.Value {
	return js.Value(c2d).Get("fillStyle")
}

func (c2d C2D) SetFillStyle(val any) {
	js.Value(c2d).Set("fillStyle", val)
}

func (c2d C2D) GetStrokeStyle() js.Value {
	return js.Value(c2d).Get("strokeStyle")
}

func (c2d C2D) SetStrokeStyle(val any) {
	js.Value(c2d).Set("strokeStyle", val)
}

func (c2d C2D) GetFilter() string {
	return js.Value(c2d).Get("filter").String()
}

func (c2d C2D) SetFilter(val string) {
	js.Value(c2d).Set("filter", val)
}

func (c2d C2D) GetImageSmoothingEnabled() bool {
	return js.Value(c2d).Get("imageSmoothingEnabled").Bool()
}

func (c2d C2D) SetImageSmoothingEnabled(val bool) {
	js.Value(c2d).Set("imageSmoothingEnabled", val)
}

func (c2d C2D) GetImageSmoothingQuality() ImageSmoothingQuality {
	return js.Value(c2d).Get("imageSmoothingQuality").String()
}

func (c2d C2D) SetImageSmoothingQuality(val ImageSmoothingQuality) {
	js.Value(c2d).Set("imageSmoothingQuality", val)
}

func (c2d C2D) GetLineCap() CanvasLineCap {
	return js.Value(c2d).Get("lineCap").String()
}

func (c2d C2D) SetLineCap(val CanvasLineCap) {
	js.Value(c2d).Set("lineCap", val)
}

func (c2d C2D) GetLineDashOffset() float64 {
	return js.Value(c2d).Get("lineDashOffset").Float()
}

func (c2d C2D) SetLineDashOffset(val float64) {
	js.Value(c2d).Set("lineDashOffset", val)
}

func (c2d C2D) GetLineJoin() CanvasLineJoin {
	return js.Value(c2d).Get("lineJoin").String()
}

func (c2d C2D) SetLineJoin(val CanvasLineJoin) {
	js.Value(c2d).Set("lineJoin", val)
}

func (c2d C2D) GetLineWidth() float64 {
	return js.Value(c2d).Get("lineWidth").Float()
}

func (c2d C2D) SetLineWidth(val float64) {
	js.Value(c2d).Set("lineWidth", val)
}

func (c2d C2D) GetMiterLimit() float64 {
	return js.Value(c2d).Get("miterLimit").Float()
}

func (c2d C2D) SetMiterLimit(val float64) {
	js.Value(c2d).Set("miterLimit", val)
}

func (c2d C2D) GetShadowBlur() float64 {
	return js.Value(c2d).Get("shadowBlur").Float()
}

func (c2d C2D) SetShadowBlur(val float64) {
	js.Value(c2d).Set("shadowBlur", val)
}

func (c2d C2D) GetShadowColor() string {
	return js.Value(c2d).Get("shadowColor").String()
}

func (c2d C2D) SetShadowColor(val string) {
	js.Value(c2d).Set("shadowColor", val)
}

func (c2d C2D) GetShadowOffsetX() float64 {
	return js.Value(c2d).Get("shadowOffsetX").Float()
}

func (c2d C2D) SetShadowOffsetX(val float64) {
	js.Value(c2d).Set("shadowOffsetX", val)
}

func (c2d C2D) GetShadowOffsetY() float64 {
	return js.Value(c2d).Get("shadowOffsetY").Float()
}

func (c2d C2D) SetShadowOffsetY(val float64) {
	js.Value(c2d).Set("shadowOffsetY", val)
}

func (c2d C2D) GetDirection() CanvasDirection {
	return js.Value(c2d).Get("direction").String()
}

func (c2d C2D) SetDirection(val CanvasDirection) {
	js.Value(c2d).Set("direction", val)
}

func (c2d C2D) GetFont() string {
	return js.Value(c2d).Get("font").String()
}

func (c2d C2D) SetFont(val string) {
	js.Value(c2d).Set("font", val)
}

func (c2d C2D) GetFontKerning() CanvasFontKerning {
	return js.Value(c2d).Get("fontKerning").String()
}

func (c2d C2D) SetFontKerning(val CanvasFontKerning) {
	js.Value(c2d).Set("fontKerning", val)
}

func (c2d C2D) GetTextAlign() CanvasTextAlign {
	return js.Value(c2d).Get("textAlign").String()
}

func (c2d C2D) SetTextAlign(val CanvasTextAlign) {
	js.Value(c2d).Set("textAlign", val)
}

func (c2d C2D) GetTextBaseline() CanvasTextBaseline {
	return js.Value(c2d).Get("textBaseline").String()
}

func (c2d C2D) SetTextBaseline(val CanvasTextBaseline) {
	js.Value(c2d).Set("textBaseline", val)
}

func (c2d C2D) GetContextAttributes() js.Value {
	return js.Value(c2d).Call("getContextAttributes")
}

func (c2d C2D) DrawImage3(image any, dx, dy float64) {
	js.Value(c2d).Call("drawImage", image, dx, dy)
}

func (c2d C2D) DrawImage5(image any, dx, dy, dw, dh float64) {
	js.Value(c2d).Call("drawImage", image, dx, dy, dw, dh)
}

func (c2d C2D) DrawImage9(image any, sx, sy, sw, sh, dx, dy, dw, dh float64) {
	js.Value(c2d).Call("drawImage", image, sx, sy, sw, sh, dx, dy, dw, dh)
}

func (c2d C2D) BeginPath() {
	js.Value(c2d).Call("beginPath")
}

func (c2d C2D) Clip0() {
	js.Value(c2d).Call("clip")
}

func (c2d C2D) Clip1CanvasFillRule(fillRule CanvasFillRule) {
	js.Value(c2d).Call("clip", fillRule)
}

func (c2d C2D) Clip1Path2D(path any) {
	js.Value(c2d).Call("clip", path)
}

func (c2d C2D) Clip2(path any, fillRule CanvasFillRule) {
	js.Value(c2d).Call("clip", path, fillRule)
}

func (c2d C2D) Fill0() {
	js.Value(c2d).Call("fill")
}

func (c2d C2D) Fill1CanvasFillRule(fillRule CanvasFillRule) {
	js.Value(c2d).Call("fill", fillRule)
}

func (c2d C2D) Fill1Path2D(path any) {
	js.Value(c2d).Call("fill", path)
}

func (c2d C2D) Fill2(path any, fillRule CanvasFillRule) {
	js.Value(c2d).Call("fill", path, fillRule)
}

func (c2d C2D) IsPointInPath2(x, y float64) js.Value {
	return js.Value(c2d).Call("isPointInPath", x, y)
}

func (c2d C2D) IsPointInPath3Path2DNumberNumber(path any, x, y float64) js.Value {
	return js.Value(c2d).Call("isPointInPath", path, x, y)
}

func (c2d C2D) IsPointInPath3NumberNumberCanvasFillRule(x, y float64, fillRule CanvasFillRule) js.Value {
	return js.Value(c2d).Call("isPointInPath", x, y, fillRule)
}

func (c2d C2D) IsPointInPath4(path any, x, y float64, fillRule CanvasFillRule) js.Value {
	return js.Value(c2d).Call("isPointInPath", path, x, y, fillRule)
}

func (c2d C2D) IsPointInStroke2(x, y float64) js.Value {
	return js.Value(c2d).Call("isPointInStroke", x, y)
}

func (c2d C2D) IsPointInStroke3(path any, x, y float64) js.Value {
	return js.Value(c2d).Call("isPointInStroke", path, x, y)
}

func (c2d C2D) Stroke0() {
	js.Value(c2d).Call("stroke")
}

func (c2d C2D) Stroke1(path any) {
	js.Value(c2d).Call("stroke", path)
}

func (c2d C2D) CreateConicGradient(startAngle, x, y float64) js.Value {
	return js.Value(c2d).Call("createConicGradient", startAngle, x, y)
}

func (c2d C2D) CreateLinearGradient(x0, y0, x1, y1 float64) js.Value {
	return js.Value(c2d).Call("createLinearGradient", x0, y0, x1, y1)
}

func (c2d C2D) CreatePattern(image any, repetition string) js.Value {
	return js.Value(c2d).Call("createPattern", image, repetition)
}

func (c2d C2D) CreateRadialGradient(x0, y0, r0, x1, y1, r1 float64) js.Value {
	return js.Value(c2d).Call("createRadialGradient", x0, y0, r0, x1, y1, r1)
}

func (c2d C2D) CreateImageData1(imagedata any) js.Value {
	return js.Value(c2d).Call("createImageData", imagedata)
}

func (c2d C2D) CreateImageData2(sw, sh float64) js.Value {
	return js.Value(c2d).Call("createImageData", sw, sh)
}

func (c2d C2D) CreateImageData3(sw, sh float64, settings any) js.Value {
	return js.Value(c2d).Call("createImageData", sw, sh, settings)
}

func (c2d C2D) GetImageData4(sx, sy, sw, sh float64) js.Value {
	return js.Value(c2d).Call("getImageData", sx, sy, sw, sh)
}

func (c2d C2D) GetImageData5(sx, sy, sw, sh float64, settings any) js.Value {
	return js.Value(c2d).Call("getImageData", sx, sy, sw, sh, settings)
}

func (c2d C2D) PutImageData3(imagedata any, dx, dy float64) {
	js.Value(c2d).Call("putImageData", imagedata, dx, dy)
}

func (c2d C2D) PutImageData7(imagedata any, dx, dy, dirtyX, dirtyY, dirtyWidth, dirtyHeight float64) {
	js.Value(c2d).Call("putImageData", imagedata, dx, dy, dirtyX, dirtyY, dirtyWidth, dirtyHeight)
}

func (c2d C2D) Arc5(x, y, radius, startAngle, endAngle float64) {
	js.Value(c2d).Call("arc", x, y, radius, startAngle, endAngle)
}

func (c2d C2D) Arc6(x, y, radius, startAngle, endAngle float64, counterclockwise bool) {
	js.Value(c2d).Call("arc", x, y, radius, startAngle, endAngle, counterclockwise)
}

func (c2d C2D) ArcTo(x1, y1, x2, y2, radius float64) {
	js.Value(c2d).Call("arcTo", x1, y1, x2, y2, radius)
}

func (c2d C2D) BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	js.Value(c2d).Call("bezierCurveTo", cp1x, cp1y, cp2x, cp2y, x, y)
}

func (c2d C2D) ClosePath() {
	js.Value(c2d).Call("closePath")
}

func (c2d C2D) Ellipse7(x, y, radiusX, radiusY, rotation, startAngle, endAngle float64) {
	js.Value(c2d).Call("ellipse", x, y, radiusX, radiusY, rotation, startAngle, endAngle)
}

func (c2d C2D) Ellipse8(x, y, radiusX, radiusY, rotation, startAngle, endAngle float64, counterclockwise bool) {
	js.Value(c2d).Call("ellipse", x, y, radiusX, radiusY, rotation, startAngle, endAngle, counterclockwise)
}

func (c2d C2D) LineTo(x, y float64) {
	js.Value(c2d).Call("lineTo", x, y)
}

func (c2d C2D) MoveTo(x, y float64) {
	js.Value(c2d).Call("moveTo", x, y)
}

func (c2d C2D) QuadraticCurveTo(cpx, cpy, x, y float64) {
	js.Value(c2d).Call("quadraticCurveTo", cpx, cpy, x, y)
}

func (c2d C2D) Rect(x, y, w, h float64) {
	js.Value(c2d).Call("rect", x, y, w, h)
}

func (c2d C2D) RoundRect4(x, y, w, h float64) {
	js.Value(c2d).Call("roundRect", x, y, w, h)
}

func (c2d C2D) RoundRect5(x, y, w, h float64, radii any) {
	js.Value(c2d).Call("roundRect", x, y, w, h, radii)
}

func (c2d C2D) GetLineDash() js.Value {
	return js.Value(c2d).Call("getLineDash")
}

func (c2d C2D) SetLineDash(segments any) {
	js.Value(c2d).Call("setLineDash", segments)
}

func (c2d C2D) ClearRect(x, y, w, h float64) {
	js.Value(c2d).Call("clearRect", x, y, w, h)
}

func (c2d C2D) FillRect(x, y, w, h float64) {
	js.Value(c2d).Call("fillRect", x, y, w, h)
}

func (c2d C2D) StrokeRect(x, y, w, h float64) {
	js.Value(c2d).Call("strokeRect", x, y, w, h)
}

func (c2d C2D) Reset() {
	js.Value(c2d).Call("reset")
}

func (c2d C2D) Restore() {
	js.Value(c2d).Call("restore")
}

func (c2d C2D) Save() {
	js.Value(c2d).Call("save")
}

func (c2d C2D) FillText3(text string, x, y float64) {
	js.Value(c2d).Call("fillText", text, x, y)
}

func (c2d C2D) FillText4(text string, x, y, maxWidth float64) {
	js.Value(c2d).Call("fillText", text, x, y, maxWidth)
}

func (c2d C2D) MeasureText(text string) js.Value {
	return js.Value(c2d).Call("measureText", text)
}

func (c2d C2D) StrokeText3(text string, x, y float64) {
	js.Value(c2d).Call("strokeText", text, x, y)
}

func (c2d C2D) StrokeText4(text string, x, y, maxWidth float64) {
	js.Value(c2d).Call("strokeText", text, x, y, maxWidth)
}

func (c2d C2D) GetTransform() js.Value {
	return js.Value(c2d).Call("getTransform")
}

func (c2d C2D) ResetTransform() {
	js.Value(c2d).Call("resetTransform")
}

func (c2d C2D) Rotate(angle float64) {
	js.Value(c2d).Call("rotate", angle)
}

func (c2d C2D) Scale(x, y float64) {
	js.Value(c2d).Call("scale", x, y)
}

func (c2d C2D) SetTransform0() {
	js.Value(c2d).Call("setTransform")
}

func (c2d C2D) SetTransform1(transform any) {
	js.Value(c2d).Call("setTransform", transform)
}

func (c2d C2D) SetTransform6(a, b, c, d, e, f float64) {
	js.Value(c2d).Call("setTransform", a, b, c, d, e, f)
}

func (c2d C2D) Transform(a, b, c, d, e, f float64) {
	js.Value(c2d).Call("transform", a, b, c, d, e, f)
}

func (c2d C2D) Translate(x, y float64) {
	js.Value(c2d).Call("translate", x, y)
}

func (c2d C2D) DrawFocusIfNeeded1(element any) {
	js.Value(c2d).Call("drawFocusIfNeeded", element)
}

func (c2d C2D) DrawFocusIfNeeded2(path any, element any) {
	js.Value(c2d).Call("drawFocusIfNeeded", path, element)
}

type GlobalCompositeOperation = string

const (
	GlobalCompositeOperationColor           GlobalCompositeOperation = "color"
	GlobalCompositeOperationColorBurn       GlobalCompositeOperation = "color-burn"
	GlobalCompositeOperationColorDodge      GlobalCompositeOperation = "color-dodge"
	GlobalCompositeOperationCopy            GlobalCompositeOperation = "copy"
	GlobalCompositeOperationDarken          GlobalCompositeOperation = "darken"
	GlobalCompositeOperationDestinationAtop GlobalCompositeOperation = "destination-atop"
	GlobalCompositeOperationDestinationIn   GlobalCompositeOperation = "destination-in"
	GlobalCompositeOperationDestinationOut  GlobalCompositeOperation = "destination-out"
	GlobalCompositeOperationDestinationOver GlobalCompositeOperation = "destination-over"
	GlobalCompositeOperationDifference      GlobalCompositeOperation = "difference"
	GlobalCompositeOperationExclusion       GlobalCompositeOperation = "exclusion"
	GlobalCompositeOperationHardLight       GlobalCompositeOperation = "hard-light"
	GlobalCompositeOperationHue             GlobalCompositeOperation = "hue"
	GlobalCompositeOperationLighten         GlobalCompositeOperation = "lighten"
	GlobalCompositeOperationLighter         GlobalCompositeOperation = "lighter"
	GlobalCompositeOperationLuminosity      GlobalCompositeOperation = "luminosity"
	GlobalCompositeOperationMultiply        GlobalCompositeOperation = "multiply"
	GlobalCompositeOperationOverlay         GlobalCompositeOperation = "overlay"
	GlobalCompositeOperationSaturation      GlobalCompositeOperation = "saturation"
	GlobalCompositeOperationScreen          GlobalCompositeOperation = "screen"
	GlobalCompositeOperationSoftLight       GlobalCompositeOperation = "soft-light"
	GlobalCompositeOperationSourceAtop      GlobalCompositeOperation = "source-atop"
	GlobalCompositeOperationSourceIn        GlobalCompositeOperation = "source-in"
	GlobalCompositeOperationSourceOut       GlobalCompositeOperation = "source-out"
	GlobalCompositeOperationSourceOver      GlobalCompositeOperation = "source-over"
	GlobalCompositeOperationXor             GlobalCompositeOperation = "xor"
)

type CanvasFillRule = string

const (
	CanvasFillRuleEvenodd CanvasFillRule = "evenodd"
	CanvasFillRuleNonzero CanvasFillRule = "nonzero"
)

type ImageSmoothingQuality = string

const (
	ImageSmoothingQualityHigh   ImageSmoothingQuality = "high"
	ImageSmoothingQualityLow    ImageSmoothingQuality = "low"
	ImageSmoothingQualityMedium ImageSmoothingQuality = "medium"
)

type CanvasLineCap = string

const (
	CanvasLineCapButt   CanvasLineCap = "butt"
	CanvasLineCapRound  CanvasLineCap = "round"
	CanvasLineCapSquare CanvasLineCap = "square"
)

type CanvasLineJoin = string

const (
	CanvasLineJoinRound CanvasLineJoin = "round"
	CanvasLineJoinBevel CanvasLineJoin = "bevel"
	CanvasLineJoinMiter CanvasLineJoin = "miter"
)

type CanvasDirection = string

const (
	CanvasDirectionInherit CanvasDirection = "inherit"
	CanvasDirectionLtr     CanvasDirection = "ltr"
	CanvasDirectionRtl     CanvasDirection = "rtl"
)

type CanvasFontKerning = string

const (
	CanvasFontKerningAuto   CanvasFontKerning = "auto"
	CanvasFontKerningNone   CanvasFontKerning = "none"
	CanvasFontKerningNormal CanvasFontKerning = "normal"
)

type CanvasTextAlign = string

const (
	CanvasTextAlignCenter CanvasTextAlign = "center"
	CanvasTextAlignEnd    CanvasTextAlign = "end"
	CanvasTextAlignLeft   CanvasTextAlign = "left"
	CanvasTextAlignRight  CanvasTextAlign = "right"
	CanvasTextAlignStart  CanvasTextAlign = "start"
)

type CanvasTextBaseline = string

const (
	CanvasTextBaselineAlphabetic  CanvasTextBaseline = "alphabetic"
	CanvasTextBaselineBottom      CanvasTextBaseline = "bottom"
	CanvasTextBaselineHanging     CanvasTextBaseline = "hanging"
	CanvasTextBaselineIdeographic CanvasTextBaseline = "ideographic"
	CanvasTextBaselineMiddle      CanvasTextBaseline = "middle"
	CanvasTextBaselineTop         CanvasTextBaseline = "top"
)
