package cfg

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"strings"
	"time"

	"codeberg.org/reiver/go-env"
	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
	"github.com/reiver/go-pckstr"

	"gamefed/cfg/boot"
	"gamefed/lib/fetch"
)

var bootStrap boot.BootStrap

func init() {
	// Skip this if we are running inside of a Go test.
	if nil != flag.Lookup("test.v") || strings.HasSuffix(os.Args[0], ".test") {
		return
	}

	gameBootURL := GameBootURL()
	if "" == gameBootURL {
		panic("empty GAME_BOOT URL environment-variable")
	}

	var bytes []byte
	func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var err error

//@TODO: make this try more than once, if it fails
		bytes, err = libfetch.Fetch(ctx, gameBootURL)
		if nil != err {
			err = erorr.Wrap(err, "failed to fetch contents of GAME_BOOT URL",
				field.String("game-boot-url", gameBootURL),
			)
			panic(err)
		}
	}()

	{
		err := json.Unmarshal(bytes, &bootStrap)
		if nil != err {
			err = erorr.Wrap(err, "failed to json-unmarshal contents of GAME_BOOT URL",
				field.String("game-boot-url", gameBootURL),
			)
			panic(err)
		}
	}
}

func GameBootURL() string {
	return env.GetElse[string]("GAME_BOOT", "")
}

func GameAcctUserName() string {
	return bootStrap.Game.Acct.UserName.GetElse("fedigame")
}

func GamePlace(placeID string) (boot.Place, bool) {
	place, found := bootStrap.Game.Places[placeID]
	return place, found
}

func GameTags() pckstr.PackedStrings {
	tags := bootStrap.Game.Tags
	if tags.IsNothing() {
		return pckstr.SomeString("PlayFedi")
	}

	return tags
}
