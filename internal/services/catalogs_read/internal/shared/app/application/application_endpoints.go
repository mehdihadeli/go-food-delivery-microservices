package application

func (a *CatalogReadApplication) MapEndpoints() {
	a.ResolveFunc(func() error {
		return nil
	})
}
