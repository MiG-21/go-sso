package internal

import (
	"crypto/rsa"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/dig"
)

var (
	gitHash   string
	gitBranch string
	gitUrl    string
	version   = "development"
)

type (
	Config struct {
		Env            string
		AppName        string
		GitHash        string
		GitBranch      string
		GitUrl         string
		Version        string
		Debug          bool   `yaml:"debug"`
		Port           int    `yaml:"port" env-default:"8080"`
		PrivateKeyPath string `yaml:"private_key_path"`
		PublicKeyPath  string `yaml:"public_key_path"`

		Logger   ConfigLogger   `yaml:"logger"`
		Http     ConfigHttp     `yaml:"http"`
		Frontend ConfigFrontend `yaml:"frontend"`
		Mysql    ConfigMysql    `yaml:"mysql"`
		Cookie   ConfigCookie   `yaml:"cookie"`
		Smtp     ConfigSmtp     `yaml:"smtp"`
		Crypto   ConfigCrypto
	}

	ConfigLogger struct {
		Path  string `yaml:"path" env:"APP_LOG_PATH"`
		Level string `yaml:"level" env:"APP_LOG_LEVEL"`
	}

	ConfigHttp struct {
		ReadBufferSize  int `yaml:"read_buffer_size" env:"APP_HTTP_READ_BUFFER_SIZE" env-default:"16384"`
		WriteBufferSize int `yaml:"write_buffer_size" env:"APP_HTTP_WRITE_BUFFER_SIZE" env-default:"16384"`
		ReadTimeout     int `yaml:"read_timeout" env:"APP_HTTP_READ_TIMEOUT" env-default:"18"`
		WriteTimeout    int `yaml:"write_timeout" env:"APP_HTTP_WRITE_TIMEOUT" env-default:"18"`
		IdleTimeout     int `yaml:"idle_timeout" env:"APP_HTTP_IDLE_TIMEOUT" env-default:"60"`
		ReqTimeout      int `yaml:"request_timeout" env:"APP_HTTP_REQUEST_TIMEOUT" env-default:"5"`
	}

	ConfigFrontend struct {
		Path  string `yaml:"path" env:"APP_FRONTEND_PATH" env-default:"/app/web/"`
		Index string `yaml:"index" env:"APP_FRONTEND_INDEX" env-default:"index.html"`
	}

	ConfigMysql struct {
		Dsn          string `yaml:"dsn" env:"APP_MYSQL_DSN"`
		MaxLifetime  int    `yaml:"max_lifetime" env:"APP_MYSQL_MAX_LIFETIME" env-default:"20"`
		MaxOpenConns int    `yaml:"max_open_conns" env:"APP_MYSQL_MAX_OPEN_CONNS" env-default:"100"`
		MaxIdleConns int    `yaml:"max_idle_conns" env:"APP_MYSQL_MAX_IDLE_CONNS" env-default:"100"`
	}

	ConfigCookie struct {
		Name       string `yaml:"cookie_name" env:"APP_COOKIE_NAME" env-default:"SSO_C"`
		Domain     string `yaml:"cookie_domain" env:"APP_COOKIE_DOMAIN" env-default:"127.0.0.1"`
		ValidHours int64  `yaml:"cookie_valid_hours" env:"APP_VALID_HOURS" env-default:"20"`
	}

	ConfigSmtp struct {
		SmtpUser     string `yaml:"smtp_user" env:"APP_SMTP_USER"`
		SmtpPassword string `yaml:"smtp_password" env:"APP_SMTP_PASSWORD"`
		SmtpHost     string `yaml:"smtp_host" env:"APP_SMTP_HOST"`
		SmtpPort     int    `yaml:"smtp_port" env:"APP_SMTP_PORT"`
		SmtpSsl      bool   `yaml:"smtp_ssl" env:"APP_SMTP_SSL"`
	}

	ConfigCrypto struct {
		PrivateKey *rsa.PrivateKey
		PublicKey  *rsa.PublicKey
	}

	SetupConfigResult struct {
		dig.Out

		Config *Config
		Error  error `group:"errors"`
	}
)

func SetupConfig() SetupConfigResult {
	sr := SetupConfigResult{}
	env := os.Getenv("APP_ENV")
	cfgPath := strings.TrimRight(os.Getenv("APP_CONFIG_PATH"), "/")
	files := []string{cfgPath + "/app.yaml"}
	if env != "" {
		files = append(files, cfgPath+"/app."+env+".yaml")
	}

	config := &Config{
		Env:       env,
		AppName:   "go-sso",
		GitHash:   gitHash,
		GitBranch: gitBranch,
		GitUrl:    gitUrl,
		Version:   version,
	}

	for i := 0; i < len(files); i++ {
		err := cleanenv.ReadConfig(files[i], config)
		if err != nil {
			sr.Error = err
			return sr
		}
	}

	privateKeyData, err := ioutil.ReadFile(config.PrivateKeyPath)
	if err != nil {
		sr.Error = err
		return sr
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		sr.Error = err
		return sr
	}

	publicKeyData, err := ioutil.ReadFile(config.PublicKeyPath)
	if err != nil {
		sr.Error = err
		return sr
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		sr.Error = err
		return sr
	}

	config.Crypto = ConfigCrypto{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	sr.Config = config
	return sr
}
