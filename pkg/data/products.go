package data

type ProductData struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var ListProduct = []ProductData{
	{ID: 1, Name: "Product 1", Description: "This is product 1"},
	{ID: 2, Name: "Product 2", Description: "This is product 2"},
	{ID: 3, Name: "Product 3", Description: "This is product 3"},
}
