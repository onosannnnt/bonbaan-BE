package model

type PackageInput struct {
	Name        string `json:"name"`
	Item        string `json:"item"`
	Price       int    `json:"price"`
	Description string `json:"description"`
}

type CreateServiceInput struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Rate        int            `json:"rate"`
	Categories  []string       `json:"categories"`
	Packages    []PackageInput `json:"packages"`
}