package ui

import (
	"sync"

	"github.com/mappu/miqt/qt6"
)

var (
	once        sync.Once
	RegularFont *qt6.QFont
	BoldFont    *qt6.QFont
	HeadingFont *qt6.QFont
)

func LoadFonts() {
	once.Do(func() {
		appFontFamily := qt6.QApplication_Font().Family()

		RegularFont = qt6.NewQFont8(appFontFamily, 12, int(qt6.QFont__Normal), false)
		BoldFont = qt6.NewQFont8(appFontFamily, 12, int(qt6.QFont__Bold), false)
		HeadingFont = qt6.NewQFont8(appFontFamily, 22, int(qt6.QFont__DemiBold), false)
	})
}
