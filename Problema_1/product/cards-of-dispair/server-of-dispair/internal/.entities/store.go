// TO CHECK
package entities

type Store struct {
	Package chan PackageInterface
}

func NewStore() *Store {
	return &Store{
		Package: make(chan PackageInterface, 100),
	}
}

type StoreInterface interface {
	AddPackage(pkg PackageInterface) error
	GetPackage() (PackageInterface, error)
}

func (s *Store) AddPackage(pkg PackageInterface) error {
	s.Package <- pkg
	return nil
}

func (s *Store) GetPackage() (PackageInterface, error) {
	pkg := <-s.Package
	return pkg, nil
}
