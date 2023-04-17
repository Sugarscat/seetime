package account

type User struct {
	Name      string
	Password  string
	Identity  bool
	ClientIp  string
	LoginTime uint64
	HeartTime uint64
}
