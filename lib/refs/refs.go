package librefs

import (
	gourl "net/url"
)

// gozaar (noun): one who carries out / performs / executes
const (
	ActorPathPrefix string = "/gozaar/"
	PlacePathPrefix string = "/jaa/"
)

func Actor(host string, actor string) string {
	var url = gourl.URL{
		Scheme: "https",
		Host:   host,
//@TODO: make the path join safer.
		Path:   ActorPathPrefix + actor,
	}

	return url.String()
}

func ActorInBox(host string, actor string) string {
//@TODO: make the path join safer.
	return Actor(host, actor) + "/inbox"
}

func ActorOutBox(host string, actor string) string {
//@TODO: make the path join safer.
	return Actor(host, actor) + "/outbox"
}

func Place(host string, placeID string) string {
	var url = gourl.URL{
		Scheme: "https",
		Host:   host,
//@TODO: make the path join safer.
		Path:   PlacePathPrefix + placeID,
	}

	return url.String()
}

func SharedInBox(host string) string {
	var url = gourl.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/inbox",
	}

	return url.String()
}
