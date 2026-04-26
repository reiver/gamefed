package placesrv

import (
	"gamefed/cfg"
	"gamefed/cfg/boot"
)

type Place = boot.Place

func Get(placeID string) (Place, bool)  {
	place, found := cfg.GamePlace(placeID)
	if found {
		return place, true
	}

	//@TODO: check database, too

	{
		var nada Place
		return nada, false
	}
}
