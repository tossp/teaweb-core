package teageo

import (
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"github.com/oschwald/geoip2-golang"
)

var DB *geoip2.Reader = nil

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		logs.Println("[proxy]start geo db")
		db, err := geoip2.Open(Tea.Root + "/web/resources/GeoLite2-City/GeoLite2-City.mmdb")
		if err != nil {
			logs.Error(err)
		}
		DB = db
	})
}
