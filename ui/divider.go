package ui

import "github.com/mappu/miqt/qt6"

type DividerType int

const (
	Vertical DividerType = iota
	Horizontal
)

type Divider struct {
	*qt6.QFrame
}

func NewDivider(lineWidth int, dividerType DividerType) *Divider {
	qFrame := qt6.NewQFrame(nil)
	if dividerType == Horizontal {
		qFrame.SetFrameShape(qt6.QFrame__HLine)
	} else if dividerType == Vertical {
		qFrame.SetFrameShape(qt6.QFrame__VLine)
	}

	qFrame.SetFrameShadow(qt6.QFrame__Plain)
	qFrame.SetLineWidth(lineWidth)

	divider := &Divider{
		QFrame: qFrame,
	}

	return divider
}
