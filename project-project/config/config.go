package config

import (
	"bytes"
	"github.com/go-redis/redis/v8"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"log"
	"os"
	"test.com/project-common/logs"
)

var C = InitConfig()

type Config struct {
	viper       *viper.Viper
	SC          *ServerConfig
	GC          *GrpcConfig
	EtcdConfig  *EtcdConfig
	MysqlConfig *MysqlConfig
	JwtConfig   *JwtConfig
	DbConfig    DbConfig
}

type ServerConfig struct {
	Name string
	Addr string
}

type GrpcConfig struct {
	Name    string
	Addr    string
	Version string
	Weight  int64
}

type EtcdConfig struct {
	Addrs []string
}

type MysqlConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Db       string
	Name     string
}

type DbConfig struct {
	Master     MysqlConfig
	Slave      []MysqlConfig
	Separation bool
}

type JwtConfig struct {
	AccessExp     int
	RefreshExp    int
	AccessSecret  string
	RefreshSecret string
}

func InitConfig() *Config {
	conf := &Config{viper: viper.New()}
	//加入nacos
	nacos := InitNacosClient()
	configYaml, err := nacos.confClient.GetConfig(vo.ConfigParam{
		DataId: "config.yaml",
		Group:  BC.NacosConfig.Group,
	})
	if err != nil {
		log.Fatalln(err)
	}
	conf.viper.SetConfigType("yaml")
	if configYaml != "" {
		err := conf.viper.ReadConfig(bytes.NewBuffer([]byte(configYaml)))
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("load nacos config")
		err = nacos.confClient.ListenConfig(vo.ConfigParam{
			DataId: "config.yaml",
			Group:  BC.NacosConfig.Group,
			OnChange: func(namespace, group, dataId, data string) {
				log.Println("listen nacos config change", data)
				//监听变化
				err = conf.viper.ReadConfig(bytes.NewBuffer([]byte(data)))
				if err != nil {
					log.Printf("listen nacos config parse err %s \n", err.Error())
				}
				//重新载入配置
				conf.ReLoadAllConfig()
			},
		})
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		workDir, _ := os.Getwd()
		conf.viper.SetConfigName("config")
		conf.viper.AddConfigPath(workDir + "/config")
		conf.viper.AddConfigPath("F:/project/ms_project/project-project/config")
		err := conf.viper.ReadInConfig()
		if err != nil {
			log.Fatalln(err)
		}
		//workDir, _ := os.Getwd()
		//conf.viper.SetConfigName("config")
		//conf.viper.SetConfigType("yaml")
		//conf.viper.AddConfigPath("etc/ms_project/user")
		//conf.viper.AddConfigPath(workDir + "/config")
		//err := conf.viper.ReadInConfig()
		//if err != nil {
		//	log.Fatalln(err)
		//}
	}
	conf.ReLoadAllConfig()
	return conf
}

func (c *Config) ReLoadAllConfig() {
	c.ReadServerConfig()
	c.InitZapLog()
	c.ReadGrpcConfig()
	c.ReadEtcdConfig()
	c.InitMysqlConfig()
	c.InitJwtConfig()
	c.InitDbConfig()
	//重新创建相关的客户端
	c.ReConnRedis()
	c.ReConnMysql()
}

//func InitConfig() *Config {
//	conf := &Config{viper: viper.New()}
//
//	conf.ReadServerConfig()
//	conf.InitZapLog()
//	conf.ReadGrpcConfig()
//	conf.ReadEtcdConfig()
//	conf.InitMysqlConfig()
//	conf.InitJwtConfig()
//	conf.InitDbConfig()
//	return conf
//}

func (c *Config) InitZapLog() {
	// 从配置中读取日志配置，初始化日志
	lc := &logs.LogConfig{
		DebugFileName: c.viper.GetString("zap.debugFileName"),
		InfoFileName:  c.viper.GetString("zap.infoFileName"),
		WarnFileName:  c.viper.GetString("zap.warnFileName"),
		MaxSize:       c.viper.GetInt("zap.maxSize"),
		MaxAge:        c.viper.GetInt("zap.maxAge"),
		MaxBackups:    c.viper.GetInt("zap.maxBackups"),
	}
	err := logs.InitLogger(lc)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *Config) ReadServerConfig() {
	sc := &ServerConfig{}
	sc.Name = c.viper.GetString("server.name")
	sc.Addr = c.viper.GetString("server.addr")
	c.SC = sc
}

func (c *Config) ReadGrpcConfig() {
	gc := &GrpcConfig{}
	gc.Name = c.viper.GetString("grpc.name")
	gc.Addr = c.viper.GetString("grpc.addr")
	gc.Version = c.viper.GetString("grpc.version")
	gc.Weight = c.viper.GetInt64("grpc.weight")
	c.GC = gc
}

func (c *Config) ReadRedisConfig() *redis.Options {
	return &redis.Options{
		Addr:     c.viper.GetString("redis.host") + ":" + c.viper.GetString("redis.port"),
		Password: c.viper.GetString("redis.password"),
		DB:       c.viper.GetInt("redis.db"),
	}
}

func (c *Config) ReadEtcdConfig() {
	ec := &EtcdConfig{}
	var addrs []string
	err := c.viper.UnmarshalKey("etcd.addrs", &addrs)
	if err != nil {
		log.Fatalln(err)
	}
	ec.Addrs = addrs
	c.EtcdConfig = ec
}
func (c *Config) InitMysqlConfig() {
	mc := &MysqlConfig{
		Username: c.viper.GetString("mysql.username"),
		Password: c.viper.GetString("mysql.password"),
		Host:     c.viper.GetString("mysql.host"),
		Port:     c.viper.GetInt("mysql.port"),
		Db:       c.viper.GetString("mysql.db"),
	}
	c.MysqlConfig = mc
}

func (c *Config) InitJwtConfig() {
	jc := &JwtConfig{
		AccessExp:     c.viper.GetInt("jwt.accessExp"),
		RefreshExp:    c.viper.GetInt("jwt.refreshExp"),
		AccessSecret:  c.viper.GetString("jwt.accessSecret"),
		RefreshSecret: c.viper.GetString("jwt.refreshSecret"),
	}
	c.JwtConfig = jc
}

func (c *Config) InitDbConfig() {
	mc := DbConfig{}
	mc.Separation = c.viper.GetBool("db.separation")
	var slaves []MysqlConfig
	err := c.viper.UnmarshalKey("db.slave", &slaves)
	if err != nil {
		panic(err)
	}
	master := MysqlConfig{
		Username: c.viper.GetString("db.master.username"),
		Password: c.viper.GetString("db.master.password"),
		Host:     c.viper.GetString("db.master.host"),
		Port:     c.viper.GetInt("db.master.port"),
		Db:       c.viper.GetString("db.master.db"),
	}
	mc.Master = master
	mc.Slave = slaves
	c.DbConfig = mc
}
