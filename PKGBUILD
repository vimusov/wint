pkgname=wint
pkgver=1.0
pkgrel=1
pkgdesc='An utility which waits for Internet connection being established'
arch=('x86_64')
url='https://git.vimusov.space/me/wint'
license=('GPL')
makedepends=('go' 'just')
source=(go.mod go.sum justfile service wint.go)
b2sums=('SKIP' 'SKIP' 'SKIP' 'SKIP' 'SKIP')

build()
{
    just
}

package()
{
    just install "$pkgdir"
    install -D --mode=0644 "$srcdir"/service "$pkgdir"/usr/lib/systemd/system/${pkgname}.service
}
