package auth

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"zgrabi-mjesto.hr/backend/src/entities/auth/roles"
	"zgrabi-mjesto.hr/backend/src/providers/database"
)

type User struct {
	gorm.Model

	Name          string     `json:"name"`
	Email         string     `json:"email" gorm:"not null;unique"`
	EmailVerified *time.Time `json:"emailVerifiedAt"`
	Image         *string    `json:"image"`
	Role          string     `json:"role" gorm:"not null"`
	Accounts      []Account  `json:"accounts"`
	Sessions      []Session  `json:"sessions"`
}

func (u *User) BeforeSave() error {
	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(u.Email)

	if !roles.IsRole(u.Role) {
		return fmt.Errorf("invalid role")
	}

	return nil
}

type Account struct {
	gorm.Model

	UserId                uint    `json:"-" gorm:"not null;constraint:OnDelete:CASCADE;"`
	User                  User    `json:"user"`
	AccountType           string  `json:"type" gorm:"not null"`     // oauth | local
	Provider              string  `json:"provider" gorm:"not null"` // credentials | google | facebook | ...
	ProviderAccountId     string  `json:"-" gorm:"not null"`        // for credentials equals UserId
	RefreshToken          *string `json:"refreshToken"`
	RefreshTokenExpiresIn *uint   `json:"-"`
	AccessToken           *string `json:"-"`
	ExpiresAt             *uint   `json:"-"`
	TokenType             *string `json:"-"`
	Scope                 *string `json:"-"`
	IdToken               *string `json:"-"`
	SessionState          *string `json:"-"`
}

type Session struct {
	gorm.Model

	SesionToken string    `json:"-" gorm:"not null"`
	UserId      uint      `json:"-" gorm:"not null;constraint:OnDelete:CASCADE;"`
	User        User      `json:"user"`
	Expires     time.Time `json:"expires" gorm:"not null"`
	Meta        *string   `json:"meta"`
}

type VerificationToken struct {
	gorm.Model

	Identifier string    `gorm:"not null"`
	Token      string    `gorm:"not null"`
	Expires    time.Time `gorm:"not null"`
}

func Init() {
	database.DatabaseProvider().Client().AutoMigrate(&User{}, &Account{}, &Session{}, &VerificationToken{})
}
