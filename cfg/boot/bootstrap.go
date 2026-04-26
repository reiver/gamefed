package boot

import (
	"github.com/reiver/go-opt"
	"github.com/reiver/go-pckstr"
)

type BootStrap struct {
	Game Game `json:"game"`
}

type Game struct {
	Acct   Acct                 `json:"acct"`
	Places map[string]Place     `json:"places"`
	Tags   pckstr.PackedStrings `json:"tags"`
}

type Acct struct {
	UserName opt.Optional[string] `json:"username"`
}

type Place struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Options     []Option `json:"options"`
}

type Option struct {
	Description string `json:"description"`
	PlaceID     string `json:"placeid"`
}
