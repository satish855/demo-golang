package domain

import (
	"database/sql"
	"encoding/json"
	"github.com/byteintellect/go_commons/entity"
	"github.com/byteintellect/protos_go/commons/v1"
	"github.com/byteintellect/protos_go/users/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type User struct {
	entity.BaseDomain
	FirstName    string
	LastName     string
	Gender       commonsv1.Gender
	Dob          *time.Time
	Addresses    []Address  `gorm:"foreignKey:UserID;references:ExternalId"`
	Identities   []Identity `gorm:"foreignKey:UserID;references:ExternalId"`
	Relations    []User     `gorm:"foreignKey:ParentID;references:ExternalId"`
	ParentID     *string
	RelationType commonsv1.Relation
}

func (u *User) GetTable() entity.DomainName {
	return "users"
}

func (u *User) ToDto() interface{} {
	return &usersv1.UserDto{
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Gender:     u.Gender,
		Dob:        timestamppb.New(*u.Dob),
		ExternalId: u.ExternalId,
		Relation:   u.RelationType,
		Status:     commonsv1.Status(int32(u.Status)),
	}
}

func (u *User) FromDto(dto interface{}) (entity.Base, error) {
	uDto := dto.(*usersv1.UserDto)
	u.FirstName = uDto.FirstName
	u.LastName = uDto.LastName
	u.Gender = uDto.Gender
	if uDto.Dob != nil {
		dobTime := uDto.Dob.AsTime()
		u.Dob = &dobTime
	}
	u.RelationType = uDto.Relation
	return u, nil
}

func (u *User) Merge(other interface{}) {
	uDto := other.(*User)
	if uDto.FirstName != "" {
		u.FirstName = uDto.FirstName
	}
	if uDto.LastName != "" {
		u.LastName = uDto.LastName
	}
	if uDto.Dob != nil {
		u.Dob = uDto.Dob
	}
	if uDto.Gender != commonsv1.Gender_GENDER_INVALID {
		u.Gender = uDto.Gender
	}
	if uDto.RelationType != commonsv1.Relation_RELATION_INVALID {
		u.RelationType = uDto.RelationType
	}
	u.Status = uDto.Status
}

func (u *User) FromSqlRow(rows *sql.Rows) (entity.Base, error) {
	for rows.Next() {
		err := rows.Scan(&u.Id, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt, &u.ExternalId, &u.Status, &u.FirstName, &u.LastName, &u.Gender, &u.Dob, &u.RelationType)
		if err != nil {
			return nil, err
		}
	}
	return u, nil
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(buffer []byte) error {
	return json.Unmarshal(buffer, u)
}

func (u *User) Block() {
	u.Status = int(commonsv1.Status_STATUS_BLOCKED)
}

func NewUser(uDto *usersv1.UserDto) *User {
	u := User{}
	u.FirstName = uDto.FirstName
	u.LastName = uDto.LastName
	u.Gender = uDto.Gender
	u.RelationType = uDto.Relation
	if uDto.Dob != nil {
		dobTime := uDto.Dob.AsTime()
		u.Dob = &dobTime
	}
	u.ParentID = nil
	return &u
}
