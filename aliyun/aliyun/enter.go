package aliyun

import (
	"aliyun/global"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Exc() {
	ip, err := ipv6()
	if err != nil {
		global.L.Error("ip 获取失败", zap.Error(err))
		return
	}
	// ip一致不需要处理
	if ip == viper.GetString("lastip") {
		global.L.Info("ip和上次相同，不需要更新", zap.String("ip", ip))
		return
	}
	if err = initClient(); err != nil {
		global.L.Error("client 初始化失败", zap.Error(err))
		return
	}
	lastip, id, err := getAliyunIP()
	if err != nil {
		global.L.Error("获取阿里云上次ip失败", zap.Error(err))
		return
	}
	if err = updateAliyun(ip, id); err != nil {
		global.L.Error("更新阿里云上次ip失败", zap.String("ip", ip), zap.Error(err))
		return
	}
	global.L.Info("设置阿里云域名解析完成", zap.String("newip", ip), zap.String("oldip", lastip))
	viper.Set("lastip", ip)
	if err = viper.WriteConfig(); err != nil {
		global.L.Error("viper重写ip失败", zap.Error(err))
		return
	}
}
