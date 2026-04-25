package boot

import (
	"github.com/reiver/go-opt"
	"github.com/reiver/go-pckstr"
)

type BootStrap struct {
	Game Game `json:"game"`
}

type Game struct {
	Acct Acct                 `json:"acct"`
	Tags pckstr.PackedStrings `json:"tags"`
}

type Acct struct {
	UserName opt.Optional[string] `json:"username"`
}
