package account

import (
	"github.com/henoya/sorascope/enum"
	"time"

	"github.com/henoya/sorascope/config"
	"github.com/henoya/sorascope/post"
	"github.com/henoya/sorascope/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

type UserId string
type AppPass string
type Session string

type User struct {
	Id        UserId        `json:"id" gorm:"type:text;primary_key"`
	Did       post.Did      `json:"did" gorm:"type:text;index:idx_user_did"`
	Host      string        `json:"host" gorm:"type:text;index:idx_user_host"`
	Handle    post.Handle   `json:"handle" gorm:"type:text;index:idx_user_handle"`
	AppPass   AppPass       `json:"app_pass" gorm:"type:text"`
	Role      enum.RoleType `json:"role" gorm:"type:text;index:idx_user_role"`
	CreatedAt *time.Time    `json:"created_at" gorm:"type:datetime;nullable"`
	UpdatedAt *time.Time    `json:"updated_at" gorm:"type:datetime;nullable"`
	DeletedAt *time.Time    `json:"deleted_at" gorm:"type:datetime;nullable;index:idx_user_deleted_at"`
}

type UserProfile struct {
	Id        UserId     `json:"id" gorm:"type:text;primary_key"`
	Name      string     `json:"name" gorm:"type:text;index:idx_user_profile_name"`
	Email     string     `json:"email" gorm:"type:text;index:idx_user_profile_email"`
	AvatarUrl string     `json:"avatar_url" gorm:"type:text;index:idx_user_profile_avatar_url"`
	CreatedAt *time.Time `json:"created_at" gorm:"type:datetime;nullable"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"type:datetime;nullable"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"type:datetime;nullable;index:idx_user_profile_deleted_at"`
}

type UserSession struct {
	Id        UserId     `json:"id" gorm:"type:text;primary_key"`
	Session   Session    `json:"session" gorm:"type:text"`
	CreatedAt *time.Time `json:"created_at" gorm:"type:datetime;nullable;index:idx_user_session_created_at"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"type:datetime;nullable;index:idx_user_session_updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"type:datetime;nullable;index:idx_user_session_deleted_at"`
}

func DoAddUser(cCtx *cli.Context) (err error) {
	cfg, err := config.GetConfigFromCtx(cCtx)
	if err != nil {
		return err
	}
	fp, err := config.GetConfigFpFromCtx(cCtx)
	if err != nil {
		return err
	}
	host := cCtx.String("host")
	handle := cCtx.String("handle")
	did := cCtx.String("did")
	appPass := cCtx.String("app-pass")
	if handle == "" || appPass == "" {
		cli.ShowSubcommandHelpAndExit(cCtx, 1)
	}

	db, err := sql.InitDBConnection()
	if err != nil {
		return err
	}
	_ = cfg
	_ = fp
	_ = host
	_ = did
	_ = db
	//err = db.Transaction(func(tx *gorm.DB) error {
	//	// User が存在するかチェック
	//	var count int64
	//
	//	if err := tx.Model(&User{}).Where("Handle = ? OR (cid = ? AND did = ?)", handle, postRecord.Cid, postRecord.Did).Count(&count).Error; err != nil {
	//		return err
	//	}
	//}
	//var user User
	//var count int64
	//err = db.Model(&User{}).Where("handle =?", handle).Count(&count).Error
	//user.Id = UserId()
	//b, err := json.MarshalIndent(&cfg, "", "  ")
	//if err != nil {
	//	return fmt.Errorf("cannot make config file: %w", err)
	//}
	//err = ioutil.WriteFile(fp, b, 0644)
	//if err != nil {
	//	return fmt.Errorf("cannot write config file: %w", err)
	//}
	return nil
}
