# Maintainer: PolJak <polesnik.jaka@gmail.com>
pkgname='curd-polland-bin'
pkgver=1.1.4
pkgrel=1
pkgdesc="Watch anime in CLI with AniList Tracking, Discord RPC, Intro/Outro/Filler/Recap Skipping, etc. Uses PolLand visual configs for rofi."
arch=('x86_64')
url="https://github.com/Pol-Jak-295/curd-polland"
license=('GPL')
depends=('mpv' 'rofi' 'ueberzugpp')
provides=('curd')
conflicts=('curd')
source=("https://github.com/Pol-Jak-295/curd-polland/releases/download/v${pkgver}/curd-linux-x86_64")
sha256sums=('SKIP')

package() {
  install -Dm755 "$srcdir/curd-linux-x86_64" "$pkgdir/usr/bin/curd"
}
