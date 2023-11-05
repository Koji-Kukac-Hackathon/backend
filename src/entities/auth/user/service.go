package user

import "golang.org/x/crypto/bcrypt"

type service struct{}

var Service service = service{}

func (service) HashPassword(password string) (hashedPassword string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return
	}
	hashedPassword = string(hash)

	return
}

func (service) CheckPasswordHash(password string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
