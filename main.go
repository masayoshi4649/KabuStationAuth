package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	kabusapi "github.com/masayoshi4649/KabuStationAPI"
	pandawin "github.com/masayoshi4649/pandalib-go/windows"
	"github.com/pelletier/go-toml/v2"
)

func main() {
	// -c もしくは --config で指定可能に
	var confPath string
	flag.StringVar(&confPath, "c", "auth.toml", "path to config file")
	flag.StringVar(&confPath, "config", "auth.toml", "path to config file (alias)")
	flag.Parse()

	cfg, err := loadConfig(confPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config (%s): %v\n", confPath, err)

		os.Exit(1)
	}

	kabusapi.SetBaseURL("http://localhost:18080/kabusapi")

	code, tok, err := kabusapi.PostAuthToken(
		kabusapi.ReqPostAuthToken{APIPassword: cfg.System.Apipw},
	)
	if err != nil {
		log.Fatalf("token error: %v (http=%d)", err, code)
	}
	fmt.Println("token:", tok.Token)
	kabusapi.SetAPIKey(tok.Token)

	fmt.Println(tok.Token)

	/*
		check
			```ps1
			Get-ItemProperty 'HKCU:\Volatile Environment\' |Select-Object APIKEY_PRD | Format-List
			```
	*/

	pandawin.SetEnv(cfg.System.EnvName, tok.Token)
}

type Config struct {
	System struct {
		Apipw   string `toml:"APIPW"`
		EnvName string `toml:"ENV_NAME"`
	} `toml:"SYSTEM"`
}

func loadConfig(path string) (Config, error) {
	var cfg Config
	b, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := toml.Unmarshal(b, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
