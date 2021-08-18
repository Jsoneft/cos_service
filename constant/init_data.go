package constant

import (
	"time"

	"github.com/Jsoneft/cos_service/constant/config"
	"github.com/Jsoneft/cos_service/constant/setting"
)

func SetupSetting() error {
	ASetting, err := setting.NewSettings()
	if err != nil {
		return err
	}
	err = ASetting.ReadSection("PrivateCOS", &config.Cos)
	if err != nil {
		return err
	}
	config.Cos.COSTimeOut *= time.Second
	return nil
}
