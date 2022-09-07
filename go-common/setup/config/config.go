package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/Calmantara/go-common/logger"
	"github.com/spf13/viper"

	serviceutil "github.com/Calmantara/go-common/service/util"
	_ "github.com/spf13/viper/remote"
)

type ConfigSetup interface {
	GetConfig(key string, configModel any) (err error)
	GetRawConfig() map[string]any
}

type ConfigSetupImpl struct {
	FileName string
	FilePath string
	FileType string

	config map[string]any
	sugar  logger.CustomLogger
	util   serviceutil.UtilService
}

type Option func(*ConfigSetupImpl)

func NewConfigSetup(sugar logger.CustomLogger, util serviceutil.UtilService, ops ...Option) ConfigSetup {
	// default value
	c := &ConfigSetupImpl{
		sugar:    sugar,
		util:     util,
		FilePath: "../manifest",
		FileType: "yaml",
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	c.FileName = fmt.Sprintf("config.%v.%v", strings.ToLower(env), c.FileType)

	// loop over option function
	for _, val := range ops {
		val(c)
	}

	// getting file from viper
	viper.SetConfigName(c.FileName)
	viper.SetConfigType(c.FileType)
	viper.AddConfigPath(c.FilePath)

	// check error
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	c.config = viper.AllSettings()
	return c
}

func (c *ConfigSetupImpl) GetConfig(key string, configModel any) (err error) {
	c.sugar.Logger().Infof("getting config for:%v", key)
	conf := c.config[key]
	if err = c.util.ObjectMapper(&conf, &configModel); err != nil {
		c.sugar.Logger().Errorf("error mapping config %v with error:%v", key, err)
	}
	return err
}

func (c *ConfigSetupImpl) GetRawConfig() map[string]any {
	return c.config
}

func WithFileName(fileName string) Option {
	return func(csi *ConfigSetupImpl) { csi.FileName = fileName }
}

func WithFilePath(filePath string) Option {
	return func(csi *ConfigSetupImpl) { csi.FilePath = filePath }
}

func WithFileType(fileType string) Option {
	return func(csi *ConfigSetupImpl) { csi.FileType = fileType }
}
