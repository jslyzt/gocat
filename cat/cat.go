package cat

import "C"
import (
	"github.com/jslyzt/gocat/ccat"
)

type Config struct {
	EncoderType     int
	EnableHeartbeat int
	EnableSampling  int
	EnableDebugLog  int
}

func DefaultConfig() Config {
	return Config{
		ENCODER_BINARY,
		1,
		1,
		0,
	}
}

func DefaultConfigForCat2() Config {
	return Config{
		ENCODER_TEXT,
		1,
		0,
		0,
	}
}


func Init(domain string, configs ...Config) {
	var config Config;
	if len(configs) > 1 {
		panic("Only 1 config can be specified while initializing cat.")
	} else if len(configs) == 1 {
		config = configs[0]
	} else {
		config = DefaultConfig()
	}

	ccat.InitWithConfig(domain, ccat.BuildConfig(
		config.EncoderType,
		config.EnableHeartbeat,
		config.EnableSampling,
		config.EnableDebugLog,
	))
	go ccat.Background()
}

func Shutdown() {
	ccat.Shutdown()
}

func Wait() {
	ccat.Wait()
}
