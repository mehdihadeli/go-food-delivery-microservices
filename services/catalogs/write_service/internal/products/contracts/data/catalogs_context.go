package data

// CatalogContext provides access to datastores that can be
// used inside a Unit-of-Work. All data changes done through
// them will be executed atomically (inside a DB transaction).
type CatalogContext interface {
	Products() ProductRepository
}
