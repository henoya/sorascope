package user

import (
	"context"
	"fmt"
	"github.com/henoya/sorascope/enum"
	userrecord "github.com/henoya/sorascope/repository/user"
	"github.com/henoya/sorascope/typedef"
	"github.com/henoya/sorascope/util"
	"time"

	comatapi "github.com/bluesky-social/indigo/api"
	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/util/cliutil"

	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func CreateSession(handle typedef.Handle, host string, appPass string) (*comatproto.ServerCreateSession_Output, error) {
	xrpcc, err := util.MakeBareXRPCC(host)
	auth, err := comatproto.ServerCreateSession(context.TODO(), xrpcc, &comatproto.ServerCreateSession_Input{
		Identifier: string(handle),
		Password:   appPass,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create session: %w", err)
	}
	return auth, nil
}

func ExistsUser(db *gorm.DB, did typedef.Did, host string) (exists bool, err error) {
	// User が存在するかチェック
	var count int64
	if err := db.Model(&userrecord.User{}).Where("did = ? AND host = ?", did, host).Count(&count).Error; err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func GetUserById(db *gorm.DB, id typedef.UserId) (userAccessSession *userrecord.UserAccessSession, err error) {
	// User が存在するかチェック
	var count int64
	if err := db.Model(&userrecord.User{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("user not found id: %s", id)
	}
	var user userrecord.User
	if err := db.Model(&userrecord.User{}).Where("id =?", id).First(&user).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&userrecord.UserHandle{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("user not found id: %s", id)
	}
	var userHandle userrecord.UserHandle
	if err := db.Model(&userrecord.UserHandle{}).Where("id =?", id).First(&userHandle).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&userrecord.UserProfile{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("user not found id: %s", id)
	}
	var userProfile userrecord.UserProfile
	if err := db.Model(&userrecord.UserProfile{}).Where("id =?", id).First(&userProfile).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&userrecord.UserSession{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("user not found id: %s", id)
	}
	var userSession userrecord.UserSession
	if err := db.Model(&userrecord.UserSession{}).Where("id =?", id).First(&userSession).Error; err != nil {
		return nil, err
	}

	userAccessSession = &userrecord.UserAccessSession{
		Id:         id,
		Did:        user.Did,
		Host:       user.Host,
		Role:       user.Role,
		Handle:     userHandle.Handle,
		AppPass:    userHandle.AppPass,
		Name:       userProfile.Name,
		Email:      userProfile.Email,
		AvatarUrl:  userProfile.AvatarUrl,
		AccessJwt:  userSession.AccessJwt,
		RefreshJwt: userSession.RefreshJwt,
	}
	return userAccessSession, nil
}

func AddUser(cCtx *cli.Context, paramHost string, paramDid string, paramHandle string, paramAppPass string) (id typedef.UserId, err error) {
	db := cCtx.App.Metadata["db"].(*gorm.DB)
	if db == nil || db == (*gorm.DB)(nil) {
		return "", fmt.Errorf("db is nil")
	}
	xrpcc, err := util.MakeBareXRPCC(paramHost)
	if err != nil {
		return "", fmt.Errorf("cannot create client: %w", err)
	}
	ctx := context.Background()

	handle := typedef.Handle("")
	did := typedef.Did("")
	host := paramHost
	// did がなく handle だけ与えられた場合
	if paramDid == "" && paramHandle != "" {
		// handle から did を取得
		handle = typedef.Handle(paramHandle)
		resolvHandle, err := comatproto.IdentityResolveHandle(ctx, xrpcc, string(handle))
		if err != nil {
			return "", fmt.Errorf("failed to resolve handle: %w", err)
		}
		did = typedef.Did(resolvHandle.Did)
	} else if paramDid != "" {
		// did が与えられた場合
		did = typedef.Did(paramDid)
		s := cliutil.GetDidResolver(cCtx)
		phr := &comatapi.ProdHandleResolver{}
		handleStr, serviceEndpoint, err := comatapi.ResolveDidToHandle(ctx, xrpcc, s, phr, paramDid)
		if err != nil {
			return "", fmt.Errorf("failed to resolve handle: %w", err)
		}
		handle = typedef.Handle(handleStr)
		host = serviceEndpoint
	}

	// User が存在するかチェック
	exists, err := ExistsUser(db, did, host)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("user already exists")
	}

	// createSession でログインできるかチェック
	authSession, err := CreateSession(handle, host, paramAppPass)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	if authSession == nil {
		return "", fmt.Errorf("failed to create session")
	}

	id = typedef.UserId("")
	did = typedef.Did(authSession.Did)
	handle = typedef.Handle(authSession.Handle)
	email := authSession.Email
	accessJwt := authSession.AccessJwt
	refreshJwt := authSession.RefreshJwt

	// did から User レコードを検索
	err = db.Transaction(func(tx *gorm.DB) error {
		// User が存在するかチェック
		var count int64
		if err := tx.Model(&userrecord.User{}).Where("did = ? AND host = ?", did, host).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("user already exists")
		}
		id = typedef.UserId(did)
		now := time.Now().UTC()
		user := userrecord.User{
			Id:        id,
			Did:       did,
			Host:      host,
			Role:      enum.RoleUser,
			CreatedAt: &now,
			UpdatedAt: &now,
			DeletedAt: nil,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		// UserHandle を登録
		if err := tx.Model(&userrecord.UserHandle{}).Where("id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			rows, err := tx.Model(&userrecord.UserHandle{}).Where("id = ?", id).Rows()
			var userHandle userrecord.UserHandle
			for rows.Next() {
				err := db.ScanRows(rows, &userHandle)
				if err != nil {
					return fmt.Errorf("failed to scan row: %w", err)
				}
				userHandle.Handle = handle
				userHandle.AppPass = typedef.AppPass(paramAppPass)
				userHandle.UpdatedAt = &now
				userHandle.DeletedAt = nil
			}
			if err = rows.Close(); err != nil {
				return fmt.Errorf("failed to close rows: %w", err)
			}
			if err := tx.Save(&userHandle).Error; err != nil {
				return err
			}
		} else {
			userHandle := userrecord.UserHandle{
				Id:        id,
				Handle:    handle,
				AppPass:   typedef.AppPass(paramAppPass),
				CreatedAt: &now,
				UpdatedAt: &now,
				DeletedAt: nil,
			}
			if err := tx.Create(&userHandle).Error; err != nil {
				return err
			}
		}

		// UserProfile を登録
		if err := tx.Model(&userrecord.UserProfile{}).Where("id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			rows, err := tx.Model(&userrecord.UserProfile{}).Where("id = ?", id).Rows()
			var userProfile userrecord.UserProfile
			for rows.Next() {
				err := db.ScanRows(rows, &userProfile)
				if err != nil {
					return fmt.Errorf("failed to scan row: %w", err)
				}
				userProfile.Name = string(handle)
				userProfile.Email = util.Stringp(email)
				userProfile.AvatarUrl = ""
				userProfile.UpdatedAt = &now
				userProfile.DeletedAt = nil
			}
			if err = rows.Close(); err != nil {
				return fmt.Errorf("failed to close rows: %w", err)
			}
			if err := tx.Save(&userProfile).Error; err != nil {
				return err
			}
		} else {
			userProfile := userrecord.UserProfile{
				Id:        id,
				Name:      string(handle),
				Email:     util.Stringp(email),
				CreatedAt: &now,
				UpdatedAt: &now,
				DeletedAt: nil,
			}
			if err := tx.Create(&userProfile).Error; err != nil {
				return err
			}
		}

		// UserSession を登録
		if err := tx.Model(&userrecord.UserSession{}).Where("id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			rows, err := tx.Model(&userrecord.UserSession{}).Where("id = ?", id).Rows()
			var userSession userrecord.UserSession
			for rows.Next() {
				err := db.ScanRows(rows, &userSession)
				if err != nil {
					return fmt.Errorf("failed to scan row: %w", err)
				}
				userSession.AccessJwt = accessJwt
				userSession.RefreshJwt = refreshJwt
				userSession.UpdatedAt = &now
				userSession.DeletedAt = nil
			}
			if err = rows.Close(); err != nil {
				return fmt.Errorf("failed to close rows: %w", err)
			}
			if err := tx.Save(&userSession).Error; err != nil {
				return err
			}
		} else {
			userSession := userrecord.UserSession{
				Id:         id,
				AccessJwt:  accessJwt,
				RefreshJwt: refreshJwt,
				CreatedAt:  &now,
				UpdatedAt:  &now,
				DeletedAt:  nil,
			}
			if err := tx.Create(&userSession).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func DoAddUser(cCtx *cli.Context) (err error) {
	paramHost := cCtx.String("host")
	paramHandle := cCtx.String("handle")
	paramDid := cCtx.String("did")
	paramAppPass := cCtx.String("app-pass")
	if (paramHandle == "" && paramDid == "") || paramAppPass == "" {
		cli.ShowSubcommandHelpAndExit(cCtx, 1)
	}

	db := cCtx.App.Metadata["db"].(*gorm.DB)
	if db == nil || db == (*gorm.DB)(nil) {
		return fmt.Errorf("db is nil")
	}

	id, err := AddUser(cCtx, paramHost, paramDid, paramHandle, paramAppPass)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	userAccessSession, err := GetUserById(db, id)
	if err != nil {
		return err
	}
	fmt.Printf("user added: id:%s\ndid: %s\nhandle: %s\nAccessJwt: %s\nRefreshJwt: %s\n", userAccessSession.Id, userAccessSession.Did, userAccessSession.Handle, userAccessSession.AccessJwt, userAccessSession.RefreshJwt)
	return nil
}
