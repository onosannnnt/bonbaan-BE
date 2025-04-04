package model

type PackageInput struct {
	Name        string   `json:"name"`
	Item        []string `json:"item"`
	Price       int      `json:"price"`
	Description string   `json:"description"`
	OrderTypeID string   `json:"order_type_id"`
}

type CreateServiceInput struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Rate        int            `json:"rate"`
	Address     string         `json:"address"`
	Categories  []string       `json:"categories"`
	Packages    []PackageInput `json:"packages"`
	Attachments []string       `json:"attachments,omitempty"`
}
type UpdateServiceInput struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Address     string         `json:"address"`
	Categories  []string       `json:"categories"`
	Packages    []PackageInput `json:"packages"`

	CustomPackage bool `json:"custom_package"`
}

type AttachmentOutput struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type PackageOutput struct {
	Name        string   `json:"name"`
	Item        []string `json:"item"`
	Price       int      `json:"price"`
	Description string   `json:"description"`
	OrderTypeID string   `json:"order_type_id"`
}

type CategoryOutput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ServiceOutput struct {
	ID            string             `json:"id"`
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	Rate          float64            `json:"rate"`
	Address       string             `json:"address"`
	Categories    []CategoryOutput   `json:"categories"`
	Packages      []PackageOutput    `json:"packages"`
	Attachments   []AttachmentOutput `json:"attachments"`
	CustomPackage bool               `json:"custom_package"`
	UpdateAt      string             `json:"update_at"`
	CreateAt      string             `json:"create_at"`
	DeleteAt      string             `json:"delete_at"`
}
