package wkit

type (
	// A special container object that allows objects to be attached and removed.
	Directory interface {
		Item
		Name() string
		Walk([]string) (Item, error)
		AddDirectory(Directory)
		AddItem(Item)
		Children() map[string]Item
		RemoveItem(Item) error
	}
)
