package sfv

type Sfver struct {
	Id    int64
	Name  string `xorm:'vchar(100)'`
	Email string `xorm:'vchar(64)'`
}
