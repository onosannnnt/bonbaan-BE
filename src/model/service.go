package model

type PackageInput struct {
	Name          string `json:"name"`
	Item          string `json:"item"`
	Price         int    `json:"price"`
	Description   string `json:"description"`
	PackageTypeID string `json:"package_type_id"`
}

type CreateServiceInput struct {
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Rate          int            `json:"rate"`
	Adress        string         `json:"adress"`
	Categories    []string       `json:"categories"`
	Packages      []PackageInput `json:"packages"`
	Attachments   []string       `json:"attachments,omitempty"`
	CustumPackage bool           `json:"custom_package"`
}

type AttachmentOutput struct {
	URL string `json:"url"`
}

type PackageOutput struct {
	Name          string `json:"name"`
	Item          string `json:"item"`
	Price         int    `json:"price"`
	Description   string `json:"description"`
	PackageTypeID string `json:"package_type_id"`
}

type CategoryOutput struct {
	Name string `json:"name"`
}

type ServiceOutput struct {
	ID            string             `json:"id"`
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	Rate          int                `json:"rate"`
	Adress        string             `json:"adress"`
	Categories    []CategoryOutput   `json:"categories"`
	Packages      []PackageOutput    `json:"packages"`
	Attachments   []AttachmentOutput `json:"attachments"`
	CustomPackage bool               `json:"custom_package"`
}