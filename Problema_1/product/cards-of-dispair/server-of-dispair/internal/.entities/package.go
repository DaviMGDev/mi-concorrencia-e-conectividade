// TO CHECK

package entities

type Package struct {
	Cards [3]CardInterface
}

type PackageInterface interface {
	GetCards() [3]CardInterface
}

func NewPackage(cards [3]CardInterface) Package {
	return Package{
		Cards: cards,
	}
}

func (p Package) GetCards() [3]CardInterface {
	return p.Cards
}
