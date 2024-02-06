package internal

// Warehouse is an struct that represents a warehouse
type Warehouse struct {
	// ID is the unique identifier of the warehouse
	Id int
	// Name is the name of the warehouse
	Name string
	// Address is the address of the warehouse
	Address string
	//Telephone is the telephone of the warehouse
	Telephone string
	// Capacity is the capacity of the warehouse
	Capacity int
}

type ReportProduct struct {
	// Name is the name of the warehouse
	Name string
	// ProductCount is the number of products in the warehouse
	ProductCount int
}
