package ui

import "github.com/mappu/miqt/qt6"

type BorderContainer struct {
	*qt6.QWidget

	borderWidth int
	padding     int
	borderColor *qt6.QColor
}

func NewBorderContainer(
	parent *qt6.QWidget,
	borderWidth int,
	padding int,
	borderColor *qt6.QColor,
) *BorderContainer {

	w := &BorderContainer{
		QWidget:     qt6.NewQWidget(parent),
		borderWidth: borderWidth,
		borderColor: borderColor,
		padding:     padding,
	}

	w.OnPaintEvent(w.PaintEvent)

	return w
}

func (w *BorderContainer) PaintEvent(
	super func(e *qt6.QPaintEvent),
	e *qt6.QPaintEvent,
) {
	super(e)

	painter := qt6.NewQPainter2(w.QPaintDevice)
	defer painter.Delete()

	painter.SetRenderHint(qt6.QPainter__Antialiasing)

	r := w.Rect()

	half := float64(w.borderWidth) / 2.0
	borderRect := qt6.NewQRectF4(
		float64(r.X())+half,
		float64(r.Y())+half,
		float64(r.Width())-float64(w.borderWidth),
		float64(r.Height())-float64(w.borderWidth),
	)

	pen := qt6.NewQPen4(
		qt6.NewQBrush3(w.borderColor),
		float64(w.borderWidth),
	)
	pen.SetJoinStyle(qt6.MiterJoin)
	pen.SetCapStyle(qt6.SquareCap)

	painter.SetPenWithPen(pen)
	painter.SetBrush(qt6.NewQBrush())

	painter.DrawRect(borderRect)
}

func (w *BorderContainer) SetLayout(layout *qt6.QLayout) {
	contentsMargins := w.borderWidth + w.padding
	layout.SetContentsMargins(
		contentsMargins,
		contentsMargins,
		contentsMargins,
		contentsMargins,
	)
	layout.SetSpacing(layout.Spacing())

	w.QWidget.SetLayout(layout)
}
