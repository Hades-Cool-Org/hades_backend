package vendors

import (
	"gorm.io/gorm"
	"hades_backend/app/model/vendors"
)

type Vendor struct {
	gorm.Model
	Name     string   `gorm:"type:varchar(255);not null;"` // Name of the vendor
	Email    string   `gorm:"type:varchar(255);not null;unique"`
	Phone    string   `gorm:"type:varchar(255);"`
	Cnpj     string   `gorm:"type:varchar(255);not null;unique"`
	Type     string   `gorm:"type:varchar(255);not null;"`
	Location string   `gorm:"type:varchar(255);not null;"`
	Contact  *Contact `gorm:"embedded;embeddedPrefix:contact_"`
}

type Contact struct {
	Name  string `gorm:"type:varchar(255);"`
	Email string `gorm:"type:varchar(255);"`
	Phone string `gorm:"type:varchar(255);"`
}

func (v *Vendor) ToDTO() *vendors.Vendor {
	return &vendors.Vendor{
		ID:       v.ID,
		Name:     v.Name,
		Email:    v.Email,
		Phone:    v.Phone,
		Cnpj:     v.Cnpj,
		Type:     v.Type,
		Location: v.Location,
		Contact: &vendors.Contact{
			Name:  v.Contact.Name,
			Email: v.Contact.Email,
			Phone: v.Contact.Phone,
		},
	}
}

func ToModel(vendor *vendors.Vendor) *Vendor {
	v := &Vendor{
		Name:     vendor.Name,
		Email:    vendor.Email,
		Phone:    vendor.Phone,
		Cnpj:     vendor.Cnpj,
		Type:     vendor.Type,
		Location: vendor.Location,
		Contact: &Contact{
			Name:  vendor.Contact.Name,
			Email: vendor.Contact.Email,
			Phone: vendor.Contact.Phone,
		},
	}

	if vendor.ID != 0 {
		v.ID = vendor.ID
	}

	return v
}
