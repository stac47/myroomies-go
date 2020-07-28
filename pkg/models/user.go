package models

type PasswordType string

type User struct {
	Firstname string       `bson:"firstname"`
	Lastname  string       `bson:"lastname"`
	IsAdmin   bool         `bson:"is_admin"`
	Login     string       `bson:"login"`
	Password  PasswordType `bson:"password"`
}

// Make sure the password in never marshalled
func (PasswordType) MarshalJSON() ([]byte, error) {
	return []byte(`""`), nil
}

func (p PasswordType) String() string {
	return string(p)
}
