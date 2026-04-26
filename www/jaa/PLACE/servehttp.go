package verboten

import (
	"flag"
	"net/http"
	"os"
	"strings"

	"codeberg.org/reiver/go-asns"
	"codeberg.org/reiver/go-field"
	"github.com/reiver/go-http404"
	"github.com/reiver/go-http500"
	"github.com/reiver/go-nul"
	"github.com/reiver/go-opt"
	"github.com/reiver/go-pathmux"

	"gamefed/cfg"
	"gamefed/lib/refs"
	"gamefed/srv/http"
	"gamefed/srv/log"
	"gamefed/srv/place"
)

const pattern string = "/jaa/{placeid}"

func init() {
	// Skip this if we are running inside of a Go test.
	if nil != flag.Lookup("test.v") || strings.HasSuffix(os.Args[0], ".test") {
		return
	}

	var handler pathmux.PatternHandler = pathmux.PatternHandlerFunc(serveHTTP)

	err := httpsrv.Mux.HandlePattern(handler, pattern)
	if nil != err {
		panic(err)
	}
}

func serveHTTP(responseWriter http.ResponseWriter, request *pathmux.ParameterizedRequest) {
	log := logsrv.Begin(field.String("www.pattern", pattern))
	defer log.End()

	if nil == responseWriter {
		log.Error(field.S("nil HTTP response-writer"))
		return
	}
	if nil == request {
		http500.InternalServerError(responseWriter, nil)
		log.Error(field.S("nil HTTP path-mux request"))
		return
	}

	placeID, found := request.ParameterByIndex(0)
	if !found {
		http500.InternalServerError(responseWriter, request.HTTPRequest())
		log.Error(field.S("missing 'actorname' (this should never happen)"))
		return
	}
	log.Trace(field.String("place-id", placeID))

	var gamePlace placesrv.Place
	{
		var found bool
		gamePlace, found = placesrv.Get(placeID)

		if !found {
			log.Warn(
				field.S("not found because invalid place-id"),
				field.String("place-id", placeID),
			)
			http404.NotFound(responseWriter, request.HTTPRequest())
			return
		}

		log.Trace(field.Any("game-place", gamePlace))
	}

	var tags []asns.ProtoObjectOrProtoLink
	{
		for _, tag := range cfg.GameTags().Strings() {
			if "" == tag {
				continue
			}

			var hashtag = asns.HashTag{
				//HRef: ???,
				Name: opt.Something("#"+tag),
			}

			tags = append(tags, hashtag)
		}
	}

	var question asns.Question
	{
		var host string = request.HTTPRequest().Host

		for _, option := range gamePlace.Options {
			var note = asns.Note{
				Name: opt.Something(option.Description),
				URL: opt.Something(librefs.Place(host, option.PlaceID)),
			}

			question.OneOf = append(question.OneOf, note)
		}

	}

	{
		var host string = request.HTTPRequest().Host

		var (
			name    opt.Optional[string] = opt.Something(gamePlace.Name)
			summary nul.Nullable[string] = nul.Something(gamePlace.Description)
		)

		var place = asns.Place{
			ID: opt.Something(librefs.Place(host, placeID)),

			Name:    name,
			Summary: summary,
		}
		place.Attachments = append(place.Attachments, question)
		place.Tags = append(place.Tags, tags...)

		bytes, err := asns.Marshal(place)
		if nil != err {
			http500.InternalServerError(responseWriter, request.HTTPRequest())
			log.Error(
				field.S("failed to jsonld-marshal ActivityPub / ActivityStreams 'Place'"),
				field.E(err),
			)
			return
		}

		asns.ServeHTTP(responseWriter, request.HTTPRequest(), bytes)
	}
}
