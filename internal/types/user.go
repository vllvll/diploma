package types

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	Id    int
	Login string
	Hash  string
}
