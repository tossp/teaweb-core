package cache

import (
	"fmt"
	"github.com/TeaWeb/code/teacache"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

type UpdatePolicyAction actions.Action

// 修改缓存策略
func (this *UpdatePolicyAction) Run(params struct {
	Filename string
}) {
	policy := shared.NewCachePolicyFromFile(params.Filename)
	if policy == nil {
		this.Fail("找不到要修改的缓存策略")
	}

	this.Data["types"] = teacache.AllCacheTypes()

	policy.Validate()

	this.Data["policy"] = maps.Map{
		"filename": policy.Filename,
		"name":     policy.Name,
		"key":      policy.Key,
		"type":     policy.Type,
		"options":  policy.Options,
		"life":     policy.Life,
		"status":   policy.Status,
		"maxSize":  policy.MaxSize,
		"capacity": policy.Capacity,
	}

	this.Show()
}

func (this *UpdatePolicyAction) RunPost(params struct {
	Filename string
	Name     string
	Key      string
	Type     string

	Capacity     float64
	CapacityUnit string
	Life         int
	LifeUnit     string
	StatusList   []int
	MaxSize      float64
	MaxSizeUnit  string

	FileDir string

	RedisNetwork  string
	RedisHost     string
	RedisPort     int
	RedisSock     string
	RedisPassword string

	LeveldbDir string

	Must *actions.Must
}) {
	policy := shared.NewCachePolicyFromFile(params.Filename)
	if policy == nil {
		this.Fail("找不到要修改的缓存策略")
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入策略名称").
		Field("key", params.Key).
		Require("请输入缓存Key")

	policy.Name = params.Name
	policy.Key = params.Key
	policy.Type = params.Type

	policy.Capacity = fmt.Sprintf("%.2f%s", params.Capacity, params.CapacityUnit)
	policy.Life = fmt.Sprintf("%d%s", params.Life, params.LifeUnit)
	for _, status := range params.StatusList {
		i := types.Int(status)
		if i >= 0 {
			policy.Status = append(policy.Status, i)
		}
	}
	policy.MaxSize = fmt.Sprintf("%.2f%s", params.MaxSize, params.MaxSizeUnit)
	policy.Status = params.StatusList

	// 选项
	// 选项
	switch policy.Type {
	case "file":
		params.Must.
			Field("fileDir", params.FileDir).
			Require("请输入缓存存放目录")
		policy.Options = map[string]interface{}{
			"dir": params.FileDir,
		}
	case "redis":
		params.Must.
			Field("redisNetwork", params.RedisNetwork).
			Require("请选择Redis连接协议").
			Field("redisHost", params.RedisHost).
			Require("请输入Redis服务器地址")
		policy.Options = map[string]interface{}{
			"network":  params.RedisNetwork,
			"host":     params.RedisHost,
			"port":     params.RedisPort,
			"password": params.RedisPassword,
			"sock":     params.RedisSock,
		}
	case "leveldb":
		params.Must.
			Field("leveldbDir", params.LeveldbDir).
			Require("请输入数据库存放目录")
		policy.Options = map[string]interface{}{
			"dir": params.LeveldbDir,
		}
	}

	err := policy.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 重置缓存策略实例
	teacache.ResetCachePolicyManager(policy.Filename)

	this.Success("保存成功")
}
