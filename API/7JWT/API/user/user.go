package user

type User struct {
	Id   int    `json:"userid"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Err struct {
	Message string `json:"message"`
}
