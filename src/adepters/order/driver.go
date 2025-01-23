package orderAdepter

import (
	"fmt"

	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
	orderUsecase "github.com/onosannnnt/bonbaan-BE/src/usecases/order"
	"gorm.io/gorm"
)

type OrderDriver struct {
	db *gorm.DB
}

func NewOrderDriver(db *gorm.DB) orderUsecase.OrderRepository {
	return &OrderDriver{
		db: db,
	}
}

func (d *OrderDriver) Insert(order *Entities.Order) error {
	if err := d.db.Create(order).Error; err != nil {
		return err
	}
	return nil
}

func (d *OrderDriver) FindAll(page *int, count *int) ([]*Entities.Order, error) {
	var selectOrder []*Entities.Order
	if err := d.db.Preload("Status").Preload("User").Order("created_at desc").Limit(*count).Find(&selectOrder).Offset(*page).Error; err != nil {
		return nil, err
	}
	return selectOrder, nil
}

func (d *OrderDriver) FindOne(id *string) (*Entities.Order, error) {
	var selectOrder Entities.Order
	if err := d.db.Preload("Status").Preload("User").Where("id = ?", id).First(&selectOrder).Error; err != nil {
		return nil, err
	}
	return &selectOrder, nil
}

func (d *OrderDriver) Update(id *string, order *Entities.Order) error {
	if err := d.db.Model(order).Updates(order).Error; err != nil {
		return err
	}
	return nil
}

func (d *OrderDriver) Delete(id *string) error {
	fmt.Println(*id)
	if err := d.db.Where("id = ?", *id).Delete(&Entities.Order{}).Error; err != nil {
		return err
	}
	return nil
}
