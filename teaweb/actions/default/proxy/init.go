package proxy

import (
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.Module("").
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantProxy,
			}).
			Helper(new(Helper)).
			Prefix("/proxy").
			Get("", new(IndexAction)).
			Get("/status", new(StatusAction)).
			GetPost("/add", new(AddAction)).
			GetPost("/delete", new(DeleteAction)).
			GetPost("/update", new(UpdateAction)).
			Post("/move", new(MoveAction)).
			Get("/detail", new(DetailAction)).
			Get("/localPath", new(LocalPathAction)).
			Get("/frontend", new(FrontendAction)).
			Get("/restart", new(RestartAction)).
			Get("/cache", new(CacheAction)).
			Post("/updateCache", new(UpdateCacheAction)).
			Post("/startHttp", new(StartHttpAction)).
			Post("/shutdownHttp", new(ShutdownHttpAction)).
			EndAll()
	})
}
