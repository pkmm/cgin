package core

import (
	"cgin/global"
	"cgin/schedule"
	"cgin/util"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

func Viper(path ...string) *viper.Viper {
	var config string
	if len(path) == 0 {
		flag.StringVar(&config, "c", "", "choose config file.")
		flag.Parse()
		if config == "" { // 优先级: 命令行 > 环境变量 > 默认值
			if configEnv := os.Getenv(util.ConfigEnv); configEnv == "" {
				config = util.ConfigFile
				fmt.Printf("正在使用config的默认值,config的路径为%v\n", util.ConfigFile)
			} else {
				config = configEnv
				fmt.Printf("正在使用GVA_CONFIG环境变量,config的路径为%v\n", config)
			}
		} else {
			fmt.Printf("正在使用命令行的-c参数传递的值,config的路径为%v\n", config)
		}
	} else {
		config = path[0]
		fmt.Printf("Viper func中使用的config的路径是 %v\n", config)
	}
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s\n", err))
	}
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed: ", in.Name)
		if err := v.Unmarshal(&global.Config); err != nil {
			fmt.Println(err)
		}
		// 重载任务调度
		schedule.SC.Reload()
	})

	if err := v.Unmarshal(&global.Config); err != nil {
		fmt.Println(err)
	}
	return v
}
