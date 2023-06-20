package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jinzhu/configor"
)

// VERSION holds the current version, should be set on ldflags
var VERSION string

// BUILD holds the current build, should be set on ldflags
var BUILD string

type contextKeyType string

var contextKey = contextKeyType("config")

// C is the configuration holder
type C struct {
	// API version
	Version string
	// API Build
	Build string
	// Is server in production ?
	Production bool `default:"false"`
	// HTTP listen url
	HTTPListenAddress string `default:"127.0.0.1:14001"`
	// HTTP public url
	HTTPPublicURL string `default:"http://127.0.0.1:14001"`

	// Database
	DatabaseHost     string `default:"localhost"`
	DatabasePort     uint16 `default:"5432"`
	DatabaseUser     string `default:"Tom"`
	DatabasePassword string `default:"TomIsGreat"`
	DatabaseName     string `default:"Tchoupi"`
}

// Init initialize config and returns a new context containing config object.
func Init(ctx context.Context) (context.Context, error) {
	c := &C{
		Version: getVersion(),
		Build:   getBuild(),
	}
	err := os.Setenv("CONFIGOR_ENV_PREFIX", "ACCOUNTABILITY")
	if err != nil {
		return nil, err
	}
	if err := configor.New(&configor.Config{Debug: false}).Load(c); err != nil {
		return nil, fmt.Errorf("could not parse config file")
	}
	if !c.Production {
		fmt.Printf("\nUsing config: %+v\n", c)
	}
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, contextKey, c), nil
}

// FromContext returns the config object from a given context.
func FromContext(ctx context.Context) *C {
	cfg, ok := ctx.Value(contextKey).(*C)
	if !ok {
		panic(fmt.Errorf("calling config.FromContext from a non-config context"))
	}
	return cfg
}

// GetFullVersion returns the full version of the program
func GetFullVersion() string {
	return getVersion() + "-" + getBuild()
}

// get current API version
func getVersion() string {
	if VERSION != "" {
		return VERSION
	}
	return "0.0.0"
}

// get current API build
func getBuild() string {
	if BUILD != "" {
		return BUILD
	}
	return "test"
}
