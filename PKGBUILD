# Maintainer: Hesham T. Banafa <hishaminv at gmail dot com>
pkgname=yt-mdl
pkgver=0.1
pkgrel=1
pkgdesc="Concurrent yt-dlp downloads"
arch=('i686' 'pentium4' 'x86_64' 'arm' 'armv7h' 'armv6h' 'aarch64' 'riscv64')
url="https://github.com/HeshamTB/yt-mdl"
options=(!lto)
license=('GPL3')
depends=(
    'yt-dlp'
)

makedepends=('go>=1.19')

source=("yt-mdl::git+https://github.com/HeshamTB/yt-mdl")
sha256sums=("SKIP")

build() {
    export GOPATH="$srcdir"/gopath
    export CGO_CPPFLAGS="${CPPFLAGS}"
    export CGO_CFLAGS="${CFLAGS}"
    export CGO_CXXFLAGS="${CXXFLAGS}"
    export CGO_LDFLAGS="${LDFLAGS}"
    export CGO_ENABLED=1

    cd "$srcdir/$pkgname"
    go mod tidy
    go build yt-mdl.go downloader.go
}

package() {
    cd "${srcdir}/yt-mdl"
    GOBIN="${pkgdir}/usr/bin/" go install
}

