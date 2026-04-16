package dto

type CreateProductResponse struct {
	ID    		string  `json:"id"`
	Name  		string  `json:"name"`
	Price 		int     `json:"price"`
	Stock 		int     `json:"stock"`
	Unit  		string  `json:"unit"`
	StoreID 	string 	`json:"store_id"`
	StoreName 	string 	`json:"store_name"`
}