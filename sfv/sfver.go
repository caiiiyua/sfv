package sfv

import (
	"time"

	"github.com/go-xorm/xorm"
)

type Sfver struct {
	Id        int64
	Name      string    `xorm:'vchar(100)' json:'name'`
	Email     string    `xorm:'vchar(64)' json:'email'`
	Days      string    `xorm:'vchar(32)' json:'days'`
	CreatedAt time.Time `xorm:'created'`
}

func DefaultSfvService(engine *xorm.Engine) *SfvService {
	return &SfvService{engine}
}

type SfvService struct {
	engine *xorm.Engine
}

func (s *SfvService) Insert(sfver Sfver) bool {
	items, err := s.engine.Insert(&sfver)
	if err != nil {
		return false
	}
	return items > 0
}

func (s *SfvService) Delete(email string) bool {
	var sfver Sfver
	s.engine.Where("email=?", email).Get(&sfver)
	s.engine.Delete(&sfver)
	return true
}

func (s *SfvService) Count() int64 {
	counts, err := s.engine.Count(&Sfver{})
	if err != nil {
		return 0
	}
	return counts
}

func (s *SfvService) List() (sfvers []Sfver) {
	s.engine.Find(&sfvers)
	return
}
