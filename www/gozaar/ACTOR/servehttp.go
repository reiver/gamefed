package verboten

import (
	"flag"
	"net/http"
	"os"
	"strings"

	"codeberg.org/reiver/go-activitypub"
	"codeberg.org/reiver/go-field"
	"github.com/reiver/go-http404"
	"github.com/reiver/go-http500"
	"github.com/reiver/go-nul"
	"github.com/reiver/go-opt"
	"github.com/reiver/go-pathmux"

	"gamefed/lib/refs"
	"gamefed/srv/actor"
	"gamefed/srv/http"
	"gamefed/srv/log"
)

const pattern string = "/gozaar/{actorname}"

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

	actorName, found := request.ParameterByIndex(0)
	if !found {
		http500.InternalServerError(responseWriter, request.HTTPRequest())
		log.Error(field.S("missing 'actorname' (this should never happen)"))
		return
	}
	log.Trace(field.String("actor-name", actorName))

	if !actorsrv.ExistsByUserName(actorName) {
		log.Warn(
			field.S("not found because invalid actor user-name"),
			field.String("actor-name", actorName),
		)
		http404.NotFound(responseWriter, request.HTTPRequest())
		return
	}

	{
		var host string = request.HTTPRequest().Host

		var (
			name    opt.Optional[string] = opt.Something(actorName)
			summary nul.Nullable[string] = nul.Something("Shall we play a game?")
		)

		var application activitypub.Application
		application.ID        = opt.Something(librefs.Actor(host, actorName))
		application.Name      = name
		application.Summary   = summary
		application.InBox     =  opt.Something(librefs.ActorInBox(host, actorName))
		application.OutBox    = opt.Something(librefs.ActorOutBox(host, actorName))
		application.EndPoints = activitypub.EndPoints{
			SharedInBox: opt.Something(librefs.SharedInBox(host)),
		}

		bytes, err := activitypub.Marshal(application)
		if nil != err {
			http500.InternalServerError(responseWriter, request.HTTPRequest())
			log.Error(
				field.S("failed to jsonld-marshal ActivityPub / ActivityStreams 'Appliation'"),
				field.E(err),
			)
			return
		}

		activitypub.ServeHTTP(responseWriter, request.HTTPRequest(), bytes)
	}
}
