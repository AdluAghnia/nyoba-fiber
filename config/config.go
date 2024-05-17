package config

import "os"

var Secret string = os.Getenv("SECRETKEY")
