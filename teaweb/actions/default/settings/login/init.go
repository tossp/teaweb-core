package login

import (
	"github.com/TeaWeb/code/teaweb/actions/default/settings"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(helpers.UserMustAuth)).
			Helper(new(settings.Helper)).
			Prefix("/settings/login").
			Get("", new(IndexAction)).
			GetPost("/update", new(UpdateAction)).
			Post("/generateKey", new(GenerateKeyAction)).
			EndAll()
	})
}
