package ui

import (
	"github.com/mappu/miqt/qt6"
)

type GlowImageWidget struct {
	*qt6.QWidget
	pixmap *qt6.QPixmap
}

func NewGlowImageWidget(parent *qt6.QWidget, imagePath string) *GlowImageWidget {
	w := &GlowImageWidget{}
	w.QWidget = qt6.NewQWidget(parent)
	w.pixmap = qt6.NewQPixmap7(imagePath, "", qt6.AutoColor)

	w.SetMinimumSize2(200, 200)
	w.OnPaintEvent(w.PaintEvent)

	return w
}

func (w *GlowImageWidget) PaintEvent(super func(e *qt6.QPaintEvent), e *qt6.QPaintEvent) {
	super(e)

	if w.pixmap == nil || w.pixmap.IsNull() {
		return
	}

	// Create painter, which draws UI
	painter := qt6.NewQPainter2(w.QPaintDevice)
	defer painter.Delete()
	painter.SetRenderHint(qt6.QPainter__Antialiasing)

	// Create centered container for image
	imgRect := w.pixmap.Rect()
	imgRect.MoveCenter(w.Rect().Center())

	borderWidth := 2

	// Create border rectangle
	outer := imgRect.Adjusted(-borderWidth, -borderWidth, borderWidth, borderWidth)
	outer = outer.Adjusted(0, 0, -1, -1)

	// Create gradient for border rectangle
	gradient := qt6.NewQLinearGradient2(outer.TopLeft().ToPointF(), outer.BottomRight().ToPointF())
	gradient.SetColorAt(0, qt6.NewQColor11(0, 180, 255, 200))
	gradient.SetColorAt(1.0, qt6.NewQColor11(0, 180, 255, 40))

	// Configure pen which defines how objects are drawn
	pen := qt6.NewQPen4(
		qt6.NewQBrush10(gradient.QGradient),
		float64(borderWidth),
	)

	pen.SetJoinStyle(qt6.MiterJoin)
	pen.SetCapStyle(qt6.SquareCap)

	painter.SetPenWithPen(pen)
	painter.SetBrush(qt6.NewQBrush())

	// Draw border and then image on top
	painter.DrawRect(outer.ToRectF())
	painter.DrawPixmap10(imgRect, w.pixmap)
}
