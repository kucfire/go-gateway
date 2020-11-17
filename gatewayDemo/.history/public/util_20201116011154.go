package public

import "crypto/sha256"

// GenSaltPassword : get saltPassword
func GenSaltPassword(salt, password string) string {
	s1 := sha256.New()
	s1.Write([]byte(password))
	return
}
