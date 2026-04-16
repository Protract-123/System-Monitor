# Licenses

This directory contains the licenses of all third party libraries/software used by the project. A list of all these
libraries/software along with their license type can be found below.

| Dependency                                           | Version | License      | License File                                                         |
|------------------------------------------------------|---------|--------------|----------------------------------------------------------------------|
| [qt/qtbase](https://github.com/qt/qtbase)            | v6.11.0 | GPL-3.0      | [GPL-3.0-qt-qtbase.txt](./GPL-3.0-qt-qtbase.txt)                     |
| [qt/qtcharts](https://github.com/qt/qtcharts)        | v6.11.0 | GPL-3.0      | [GPL-3.0-qt-qtcharts.txt](./GPL-3.0-qt-qtcharts.txt)                 |
| [mappu/miqt](https://github.com/mappu/miqt)          | v0.13.0 | MIT          | [MIT-mappu-miqt.txt](./MIT-mappu-miqt.txt)                           |
| [golang.org/sys](https://go.googlesource.com/sys/)   | v0.43.0 | BSD-3-Clause | [BSD-3-Clause-golang.org-sys.txt](./BSD-3-Clause-golang.org-sys.txt) |
| [golang/go](https://go.googlesource.com/go)          | v1.25   | BSD-3-Clause | [BSD-3-Clause-golang.txt](./BSD-3-Clause-golang.txt)                 |

## Qt-bundled third-party libraries

Qt itself bundles additional third-party libraries which are included inside the Qt frameworks
shipped in the macOS `.app` bundle. The license texts for these libraries are maintained by the Qt Project and are listed at:

<https://doc.qt.io/qt-6/licenses-used-in-qt.html>

They are not redistributed alongside this project's binaries; refer to the link above for the current set of
licenses applicable to the Qt version this project is built against.
