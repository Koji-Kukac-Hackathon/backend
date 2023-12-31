package auth

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"
	"zgrabi-mjesto.hr/backend/src/entities/auth/account"
	"zgrabi-mjesto.hr/backend/src/entities/auth/roles"
	"zgrabi-mjesto.hr/backend/src/entities/auth/user"
	"zgrabi-mjesto.hr/backend/src/providers/database"
)

type service struct{}

var Service service = service{}

func (service) Register(registerData *RegisterData) (err error) {
	fmt.Printf("Register: %+v", registerData)

	db := database.DatabaseProvider().Client()

	err = db.Transaction(func(tx *gorm.DB) error {
		// da, znam, postoji oracle za enumeraciju registriranih korisnika,
		// ali to ne smatramo problemom

		var existingUser User
		res := tx.Model(&User{}).Where("email = ?", registerData.Email).First(&existingUser)
		if res.Error == nil {
			return fmt.Errorf("user with email %s already exists", registerData.Email)
		}

		passwordHash, err := user.Service.HashPassword(registerData.Password)
		if err != nil {
			return err
		}

		newUser := User{
			Email: registerData.Email,
			Name:  registerData.Name,
			Role:  roles.RoleUser,
		}
		if tx.Create(&newUser).Error != nil {
			return err
		}

		fmt.Printf("\n===============================================\nUser created: %+v\n\n===============================================\n", newUser)

		newAccount := Account{
			UserId:            newUser.ID,
			AccountType:       account.AccountTypeLocal,
			Provider:          "credentials",
			ProviderAccountId: strconv.FormatUint(uint64(newUser.ID), 10),
			RefreshToken:      &passwordHash,
		}
		if tx.Create(&newAccount).Error != nil {
			return err
		}
		fmt.Printf("\n===============================================\nAccount created: %+v\n\n===============================================\n", newAccount)

		return nil
	})

	return
}

var providerHandlers = map[string]func(*Account, *LoginData) bool{
	"credentials": func(a *Account, ld *LoginData) bool {
		return user.Service.CheckPasswordHash(*ld.Password, *a.RefreshToken)
	},
}

func (service) Login(loginData *LoginData) (user *User, err error) {
	db := database.DatabaseProvider().Client()

	var existingUser User
	res := db.Model(&User{}).Preload("Accounts").Where("email = ?", loginData.Email).First(&existingUser)
	if err = res.Error; err != nil {
		return
	}

	handler, ok := providerHandlers[loginData.Provider]
	if !ok {
		err = fmt.Errorf("provider %s is not supported", loginData.Provider)

		return
	}

	var account *Account
	for _, a := range existingUser.Accounts {
		if a.Provider == loginData.Provider {
			account = &a
			break
		}
	}

	if account == nil {
		err = fmt.Errorf("account for provider %s not found", loginData.Provider)
	}

	if !handler(account, loginData) {
		err = fmt.Errorf("invalid credentials")
	}

	user = &existingUser

	return
}

func (service) EditUser(userId uint, editUserData *UpdateUserData) (err error) {
	db := database.DatabaseProvider().Client()

	err = db.Transaction(func(tx *gorm.DB) (err error) {
		var dbUser *User
		err = db.Model(&User{}).Preload("Accounts").Where("id = ?", userId).First(&dbUser).Error
		if dbUser == nil {
			return
		}

		dbUser.Name = editUserData.Name
		dbUser.Email = editUserData.Email
		dbUser.Role = editUserData.Role

		err = db.Save(dbUser).Error
		if err != nil {
			return
		}

		var credAccount *Account
		for _, a := range dbUser.Accounts {
			if a.Provider == "credentials" {
				credAccount = &a
				break
			}
		}

		if credAccount == nil {
			return fmt.Errorf("account for provider %s not found", "credentials")
		}

		hashedPassword, err := user.Service.HashPassword(editUserData.Password)
		if err != nil {
			return
		}

		credAccount.RefreshToken = &hashedPassword
		err = db.Save(credAccount).Error

		return
	})

	return
}

func (service) DeleteUser(userId uint) (err error) {
	db := database.DatabaseProvider().Client()

	delUser := User{}
	delUser.ID = userId
	err = db.Unscoped().Delete(&delUser).Error

	return
}
