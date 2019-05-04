package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"strings"
)

// 添加新的服务
type AddAction actions.Action

func (this *AddAction) Run(params struct {
}) {
	this.Show()
}

// 提交保存
func (this *AddAction) RunPost(params struct {
	Description string
	ServerType  string
	Names       []string
	Listens     []string
	Backends    []string
	Root        string
	Must        *actions.Must
}) {
	if len(params.Description) == 0 {
		params.Description = "新代理服务"
	}

	server := teaconfigs.NewServerConfig()
	server.Http = true
	server.Description = params.Description
	server.Charset = "utf-8"
	server.Index = []string{"index.html", "index.htm", "index.php"}
	server.CacheOn = true

	if len(params.Names) > 0 {
		for _, name := range params.Names {
			name = strings.TrimSpace(name)
			if len(name) > 0 {
				server.AddName(name)
			}
		}
	}

	if len(params.Listens) > 0 {
		for _, listen := range params.Listens {
			listen = strings.TrimSpace(listen)
			if len(listen) > 0 {
				server.AddListen(listen)
			}
		}
	}

	if params.ServerType == "proxy" { // 代理服务
		for _, backend := range params.Backends {
			backend = strings.TrimSpace(backend)
			if len(backend) > 0 {
				backendObject := teaconfigs.NewBackendConfig()
				if strings.HasPrefix(backend, "http://") {
					backend = strings.TrimPrefix(backend, "http://")
					backendObject.Scheme = "http"
				} else if strings.HasPrefix(backend, "https://") {
					backend = strings.TrimPrefix(backend, "https://")
					backendObject.Scheme = "https"
				} else {
					backendObject.Scheme = "http"
				}

				backendObject.Address = backend
				backendObject.Weight = 10
				server.AddBackend(backendObject)
			}
		}
	} else if params.ServerType == "static" { // 普通服务
		server.Root = params.Root
	}

	err := server.Validate()
	if err != nil {
		this.Fail("添加时有问题发生：" + err.Error())
	}

	filename := "server." + server.Id + ".proxy.conf"
	server.Filename = filename
	err = server.Save()
	if err != nil {
		this.Fail(err.Error())
	}

	// 保存到server list
	serverList, err := teaconfigs.SharedServerList()
	if err != nil {
		logs.Error(err)
	} else {
		serverList.AddServer(server.Filename)
		err = serverList.Save()
		if err != nil {
			logs.Error(err)
		}
	}

	proxyutils.NotifyChange()

	this.Next("/proxy/detail", map[string]interface{}{
		"serverId": server.Id,
	}, "").Success("添加成功，现在去查看详细信息")
}
