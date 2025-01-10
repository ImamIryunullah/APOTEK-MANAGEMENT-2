package models

import "gorm.io/gorm"
import "github.com/go-playground/validator/v10"

type Stok struct {
    gorm.Model
    NamaObat string  `json:"nama_obat" validate:"required,min=3,max=100"`
    Jumlah   int     `json:"jumlah" validate:"required,min=1"`
    Harga    float64 `json:"harga" validate:"required,gt=0"`
}

var validate *validator.Validate

func init() {
    validate = validator.New()
}

func (s *Stok) Validate() error {
    return validate.Struct(s)
}
