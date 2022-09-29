package domain

import (
	"database/sql"
	"encoding/json"
	"github.com/byteintellect/go_commons/entity"
	"github.com/byteintellect/protos_go/commons/v1"
	"github.com/byteintellect/protos_go/users/v1"
)

type Identity struct {
	entity.BaseDomain
	IdentityValue string
	IdentityType  commonsv1.IdentityType
	UserID        string
}

func (i *Identity) GetTable() entity.DomainName {
	return "identities"
}

func (i *Identity) ToDto() interface{} {
	return &usersv1.IdentityDto{
		Type:       i.IdentityType,
		Value:      i.IdentityValue,
		ExternalId: i.ExternalId,
	}
}

func (i *Identity) FromDto(dto interface{}) (entity.Base, error) {
	iDto := dto.(*usersv1.IdentityDto)
	i.IdentityType = iDto.Type
	i.IdentityValue = iDto.Value
	return i, nil
}

func (i *Identity) Merge(other interface{}) {
	iDto := other.(*Identity)
	if iDto.IdentityType != commonsv1.IdentityType_IDENTITY_TYPE_INVALID {
		i.IdentityType = iDto.IdentityType
	}
	if iDto.IdentityValue != "" {
		i.IdentityValue = iDto.IdentityValue
	}
}

func (i *Identity) FromSqlRow(rows *sql.Rows) (entity.Base, error) {
	var err error
	for rows.Next() {
		err = rows.Scan(&i.Id, &i.CreatedAt, &i.UpdatedAt, &i.DeletedAt, &i.Status, &i.IdentityValue, &i.IdentityType)
	}
	return i, err
}

func (i *Identity) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func (i *Identity) UnmarshalBinary(buffer []byte) error {
	return json.Unmarshal(buffer, i)
}

func NewIdentity(dto *usersv1.IdentityDto) *Identity {
	identity := Identity{}
	identity.FromDto(dto)
	return &identity
}
