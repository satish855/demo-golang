package domain

import (
	"database/sql"
	"encoding/json"
	"github.com/byteintellect/go_commons/entity"
	"github.com/byteintellect/protos_go/users/v1"
)

type Address struct {
	entity.BaseDomain
	Line1     string `gorm:"column:line_1"`
	Line2     string `gorm:"column:line_2"`
	City      string
	Province  string
	Country   string
	Latitude  float64
	Longitude float64
	ZipCode   string `gorm:"column:zipcode"`
	UserID    string
}

func (a *Address) GetTable() entity.DomainName {
	return "addresses"
}

func (a *Address) ToDto() interface{} {
	return &usersv1.AddressDto{
		Line_1:     a.Line1,
		Line_2:     a.Line2,
		City:       a.City,
		Province:   a.Province,
		Country:    a.Country,
		Zipcode:    a.ZipCode,
		Latitude:   a.Latitude,
		Longitude:  a.Longitude,
		ExternalId: a.ExternalId,
	}
}

func (a *Address) FromDto(dto interface{}) (entity.Base, error) {
	aDto := dto.(*usersv1.AddressDto)
	a.Line1 = aDto.Line_1
	a.Line2 = aDto.Line_2
	a.City = aDto.City
	a.Province = aDto.Province
	a.Country = aDto.Country
	a.ZipCode = aDto.Zipcode
	a.Latitude = aDto.Latitude
	a.Longitude = aDto.Longitude
	return a, nil
}

func (a *Address) Merge(other interface{}) {
	aDto := other.(*Address)
	if aDto.Line1 != "" {
		a.Line1 = aDto.Line1
	}
	if aDto.Line2 != "" {
		a.Line2 = aDto.Line2
	}
	if aDto.City != "" {
		a.City = aDto.City
	}
	if aDto.Province != "" {
		a.Province = aDto.Province
	}
	if aDto.Country != "" {
		a.Country = aDto.Country
	}
	if aDto.ZipCode != "" {
		a.ZipCode = aDto.ZipCode
	}
	if aDto.Latitude != 0 {
		a.Latitude = aDto.Latitude
	}
	if aDto.Longitude != 0 {
		a.Longitude = aDto.Longitude
	}
}

func (a *Address) FromSqlRow(rows *sql.Rows) (entity.Base, error) {
	var err error
	for rows.Next() {
		err = rows.Scan(&a.Id, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt, &a.Status, &a.Line1, &a.Line2,
			&a.City, &a.Province, &a.Country, &a.ZipCode, &a.Longitude, &a.Longitude)
	}
	return a, err
}

func (a *Address) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

func (a *Address) UnmarshalBinary(buffer []byte) error {
	return json.Unmarshal(buffer, a)
}

func NewAddress(dto *usersv1.AddressDto) *Address {
	address := Address{}
	address.FromDto(dto)
	return &address
}
