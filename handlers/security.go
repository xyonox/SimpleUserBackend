package handlers

func HashPassword(password string) string {
	return password
}

func CheckPasswordHash(password string) bool {
	return password == HashPassword(password)
}
