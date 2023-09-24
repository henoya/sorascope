package user

import (
	"time"

	"github.com/henoya/sorascope/enum"
	"github.com/henoya/sorascope/typedef"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id        typedef.UserId `json:"id" gorm:"type:text;primary_key"`
	Did       typedef.Did    `json:"did" gorm:"type:text;unique_index:idx_user_did"`
	Host      string         `json:"host" gorm:"type:text;unique_index:idx_user_did"`
	Role      enum.RoleType  `json:"role" gorm:"type:text;index:idx_user_role"`
	CreatedAt *time.Time     `json:"created_at" gorm:"type:datetime;nullable"`
	UpdatedAt *time.Time     `json:"updated_at" gorm:"type:datetime;nullable"`
	DeletedAt *time.Time     `json:"deleted_at" gorm:"type:datetime;nullable;index:idx_user_deleted_at"`
}

type UserHandle struct {
	Id        typedef.UserId  `json:"id" gorm:"type:text;primary_key"`
	Handle    typedef.Handle  `json:"handle" gorm:"type:text;index:idx_user_handle_handle"`
	AppPass   typedef.AppPass `json:"app_pass" gorm:"type:text"`
	CreatedAt *time.Time      `json:"created_at" gorm:"type:datetime;nullable"`
	UpdatedAt *time.Time      `json:"updated_at" gorm:"type:datetime;nullable"`
	DeletedAt *time.Time      `json:"deleted_at" gorm:"type:datetime;nullable;index:idx_user_handle_deleted_at"`
}

type UserProfile struct {
	Id        typedef.UserId `json:"id" gorm:"type:text;primary_key"`
	Name      string         `json:"name" gorm:"type:text;index:idx_user_profile_name"`
	Email     string         `json:"email" gorm:"type:text;index:idx_user_profile_email"`
	AvatarUrl string         `json:"avatar_url" gorm:"type:text;index:idx_user_profile_avatar_url"`
	CreatedAt *time.Time     `json:"created_at" gorm:"type:datetime;nullable"`
	UpdatedAt *time.Time     `json:"updated_at" gorm:"type:datetime;nullable"`
	DeletedAt *time.Time     `json:"deleted_at" gorm:"type:datetime;nullable;index:idx_user_profile_deleted_at"`
}

type UserSession struct {
	Id         typedef.UserId `json:"id" gorm:"type:text;primary_key"`
	AccessJwt  string         `json:"access_jwt" gorm:"type:text"`
	RefreshJwt string         `json:"refresh_jwt" gorm:"type:text"`
	CreatedAt  *time.Time     `json:"created_at" gorm:"type:datetime;nullable;index:idx_user_session_created_at"`
	UpdatedAt  *time.Time     `json:"updated_at" gorm:"type:datetime;nullable;index:idx_user_session_updated_at"`
	DeletedAt  *time.Time     `json:"deleted_at" gorm:"type:datetime;nullable;index:idx_user_session_deleted_at"`
}

type UserAccessSession struct {
	Id         typedef.UserId  `json:"id"`
	Did        typedef.Did     `json:"did"`
	Host       string          `json:"host"`
	Role       enum.RoleType   `json:"role"`
	Handle     typedef.Handle  `json:"handle"`
	AppPass    typedef.AppPass `json:"app_pass"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	AvatarUrl  string          `json:"avatar_url"`
	AccessJwt  string          `json:"access_jwt"`
	RefreshJwt string          `json:"refresh_jwt"`
}
