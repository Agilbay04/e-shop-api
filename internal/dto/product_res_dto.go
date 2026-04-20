package dto

type ProductResponse struct {
	ID    		string  `json:"id"`
	Name  		string  `json:"name"`
	Description string  `json:"description"`
	Price 		int     `json:"price"`
	Stock 		int     `json:"stock"`
	Unit  		string  `json:"unit"`
	IsActive 	bool 	`json:"is_active"`
	CreatedAt 	string 	`json:"created_at"`
	CreatedBy 	string 	`json:"created_by"`
	UpdatedAt 	string 	`json:"updated_at"`
	UpdatedBy 	string 	`json:"updated_by"`
	DeletedAt 	string 	`json:"deleted_at"`
	StoreID 	string 	`json:"store_id"`
	StoreName 	string 	`json:"store_name"`
}