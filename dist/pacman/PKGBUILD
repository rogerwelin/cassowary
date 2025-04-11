# Maintainer: Canmi21 <9997200@qq.com>
# Contributor: Canmi (Canmi21)

pkgname=cassowary-git
pkgver="0.18.0"
pkgrel=1
pkgdesc="🚀 Modern cross-platform HTTP load-testing tool written in Go"
arch=('x86_64')
url="https://github.com/rogerwelin/cassowary"
license=('MIT')
depends=()
makedepends=('go' 'git')
provides=('cassowary')
conflicts=('cassowary-bin')
source=("cassowary::git+https://github.com/rogerwelin/cassowary.git")
sha256sums=('SKIP')

pkgver() {
  cd "$srcdir/cassowary"
  echo "r$(git rev-list --count HEAD).$(git rev-parse --short HEAD)"
}

build() {
  cd "$srcdir/cassowary/cmd/cassowary"
  export CGO_ENABLED=0
  go build -o cassowary -trimpath -ldflags="-s -w" .
}

package() {
  install -Dm755 "$srcdir/cassowary/cmd/cassowary/cassowary" "$pkgdir/usr/bin/cassowary"
}