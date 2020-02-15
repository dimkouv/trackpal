package conf

type Argon2Conf struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var Argon2Params = Argon2Conf{
	Memory:      32 * 1024,
	Iterations:  4,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}
