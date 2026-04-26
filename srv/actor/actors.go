package actorsrv

import (
	"gamefed/cfg"
)

func ExistsByUserName(username string) bool {
	if cfg.GameAcctUserName() == username {
		return true
	}

	//@TODO: check database, too

	return false
}
