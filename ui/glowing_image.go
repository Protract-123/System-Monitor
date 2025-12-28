package ui

import (
	"github.com/mappu/miqt/qt6"
)

type GlowImageWidget struct {
	*qt6.QWidget
	pixmap      *qt6.QPixmap
	borderWidth int
}

func NewGlowImageWidget(parent *qt6.QWidget, imagePath string, borderWidth int) *GlowImageWidget {
	w := &GlowImageWidget{}
	w.QWidget = qt6.NewQWidget(parent)
	w.pixmap = qt6.NewQPixmap7(imagePath, "", qt6.AutoColor)
	w.borderWidth = borderWidth

	w.SetContentsMargins(w.borderWidth, w.borderWidth, w.borderWidth, w.borderWidth)
	w.OnPaintEvent(w.PaintEvent)
	w.OnSizeHint(w.SizeHint)

	return w
}

func (w *GlowImageWidget) PaintEvent(super func(e *qt6.QPaintEvent), e *qt6.QPaintEvent) {
	if w.pixmap == nil || w.pixmap.IsNull() {
		return
	}

	// Create painter, which draws UI
	painter := qt6.NewQPainter2(w.QPaintDevice)
	defer painter.Delete()
	painter.SetRenderHint(qt6.QPainter__Antialiasing)

	// ContentsRect is size of content excluding margin
	content := w.ContentsRect()

	// Image rect at native size, centered in content
	imgSize := w.pixmap.Size()
	maxSize := qt6.NewQSize2(256, 256)

	targetSize := content.Size().BoundedTo(maxSize)

	scaledSize := imgSize.Scaled2(targetSize, qt6.KeepAspectRatio)

	imgRect := qt6.NewQRect4(
		0,
		0,
		scaledSize.Width(),
		scaledSize.Height(),
	)

	imgRect.MoveCenter(content.Center())

	outer := imgRect.ToRectF()

	// Create gradient for border rectangle
	gradient := qt6.NewQLinearGradient2(outer.TopLeft(), outer.BottomRight())
	gradient.SetColorAt(0, qt6.NewQColor11(0, 180, 255, 200))
	gradient.SetColorAt(1.0, qt6.NewQColor11(0, 180, 255, 40))

	// Configure pen which defines how objects are drawn
	pen := qt6.NewQPen4(
		qt6.NewQBrush10(gradient.QGradient),
		float64(w.borderWidth),
	)

	pen.SetJoinStyle(qt6.MiterJoin)
	pen.SetCapStyle(qt6.SquareCap)

	painter.SetPenWithPen(pen)
	painter.SetBrush(qt6.NewQBrush())

	// Draw border and then image on top
	painter.DrawRect(outer)
	painter.DrawPixmap10(imgRect, w.pixmap)
}

func (w *GlowImageWidget) SizeHint(super func() *qt6.QSize) *qt6.QSize {
	if w.pixmap == nil || w.pixmap.IsNull() {
		return qt6.NewQSize2(0, 0)
	}

	// Max image size you allow
	maxImage := qt6.NewQSize2(256, 256)

	// Native image size clamped to max
	imageSize := w.pixmap.Size().BoundedTo(maxImage)

	// Add contents margins (border)
	l, t, r, b := w.borderWidth, w.borderWidth, w.borderWidth, w.borderWidth
	imageSize.SetWidth(imageSize.Width() + l + r)
	imageSize.SetHeight(imageSize.Height() + t + b)

	return imageSize
}
