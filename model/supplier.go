package model

type Supplier struct{
	ID int32 `json:"id"`
	Name string `json:"name"`
	IconURL string `json:"icon_url"`
	Description string `json:"description"`
	AddressName string `json:"address_name"`
	
}
