package shared

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/utils/string"
	"time"
)

// 缓存策略配置
type CachePolicy struct {
	Filename string `yaml:"filename" json:"filename"` // 文件名
	Id       string `yaml:"id" json:"id"`
	On       bool   `yaml:"on" json:"on"`     // 是否开启 TODO
	Name     string `yaml:"name" json:"name"` // 名称

	Key      string `yaml:"key" json:"key"`           // 每个缓存的Key规则，里面可以有变量
	Capacity string `yaml:"capacity" json:"capacity"` // 最大内容容量
	Life     string `yaml:"life" json:"life"`         // 时间
	Status   []int  `yaml:"status" json:"status"`     // 缓存的状态码列表
	MaxSize  string `yaml:"maxSize" json:"maxSize"`   // 能够请求的最大尺寸

	life     time.Duration
	maxSize  float64
	capacity float64

	Type    string                 `yaml:"type" json:"type"`       // 类型
	Options map[string]interface{} `yaml:"options" json:"options"` // 选项
}

// 获取新对象
func NewCachePolicy() *CachePolicy {
	return &CachePolicy{}
}

// 从文件中读取缓存策略
func NewCachePolicyFromFile(file string) *CachePolicy {
	if len(file) == 0 {
		return nil
	}
	reader, err := files.NewReader(Tea.ConfigFile(file))
	if err != nil {
		logs.Error(err)
		return nil
	}
	defer reader.Close()

	p := NewCachePolicy()
	err = reader.ReadYAML(p)
	if err != nil {
		logs.Error(err)
		return nil
	}

	return p
}

// 校验
func (this *CachePolicy) Validate() error {
	var err error
	this.maxSize, _ = stringutil.ParseFileSize(this.MaxSize)
	this.life, _ = time.ParseDuration(this.Life)
	this.capacity, _ = stringutil.ParseFileSize(this.Capacity)
	return err
}

// 最大数据尺寸
func (this *CachePolicy) MaxDataSize() float64 {
	return this.maxSize
}

// 容量
func (this *CachePolicy) CapacitySize() float64 {
	return this.capacity
}

// 生命周期
func (this *CachePolicy) LifeDuration() time.Duration {
	return this.life
}

// 保存
func (this *CachePolicy) Save() error {
	if len(this.Filename) == 0 {
		this.Id = stringutil.Rand(16)
		this.Filename = "cache.policy." + this.Id + ".conf"
	}
	writer, err := files.NewWriter(Tea.ConfigFile(this.Filename))
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = writer.WriteYAML(this)
	return err
}

// 删除
func (this *CachePolicy) Delete() error {
	if len(this.Filename) > 0 {
		return files.NewFile(Tea.ConfigFile(this.Filename)).Delete()
	}
	return nil
}
