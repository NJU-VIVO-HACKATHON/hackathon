package global

import (
	"github.com/NJU-VIVO-HACKATHON/hackathon/util"
	"sync"
)

var (
	jwt     *util.Jwt
	jwtOnce sync.Once
)

func GetJwt() *util.Jwt {
	jwtOnce.Do(func() {
		conf := GetConfig()
		jwt = util.NewJwt(conf.Jwt.SecretKey, conf.Jwt.ExpiresDuration)
	})
	return jwt
}
