package config

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port     string         `yaml:"port"`
	Database DatabaseConfig `yaml:"database"`
	Auth     Auth           `yaml:"auth"`
}

type Auth struct {
	JwtSecret string `yaml:"jwt_secret"`
	Expire    int32  `yaml:"expire"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

var authCfg *Auth

const _PATH = "./config.yml"

func Load() (*Config, error) {
	data, err := os.ReadFile(_PATH)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var appConfig *Config
	if err = yaml.Unmarshal(data, &appConfig); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	if appConfig == nil {
		return nil, errors.New("invalid config")
	}

	authCfg = &appConfig.Auth

	// fallback for deploy Railway
	envURL := os.Getenv("DATABASE_URL")
	if envURL != "" {
		appConfig.Database = fallbackDBConfigForRailway(envURL)
	}

	return appConfig, nil
}

func fallbackDBConfigForRailway(envURL string) DatabaseConfig {
	u, err := url.Parse(envURL)
	if err != nil {
		return DatabaseConfig{}
	}

	password, _ := u.User.Password()
	return DatabaseConfig{
		Host:     u.Hostname(),
		Port:     parsePort(u.Port()),
		User:     u.User.Username(),
		Password: password,
		DBName:   strings.TrimPrefix(u.Path, "/"),
	}
}

func parsePort(s string) int {
	p, _ := strconv.Atoi(s)
	return p
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.Host, c.User, c.Password, c.DBName, c.Port, c.SSLMode,
	)
}

func GenerateToken(userID uint) (string, time.Time, error) {
	if authCfg == nil {
		return "", time.Now(), errors.New("no auth config")
	}
	expireAt := time.Now().Add(time.Duration(authCfg.Expire) * time.Hour)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     jwt.NewNumericDate(expireAt),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(authCfg.JwtSecret))
	if err != nil {
		return "", time.Now(), errors.New(err.Error())
	}
	return tokenStr, expireAt, nil
}
