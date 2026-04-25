package verboten

import (
	"errors"
	"flag"
	"net/http"
	"os"
	"strings"

	"codeberg.org/reiver/go-accturi"
	"codeberg.org/reiver/go-field"
	"codeberg.org/reiver/go-webfinger"
	"github.com/reiver/go-errhttp"
	"github.com/reiver/go-opt"

	"gamefed/lib/refs"
	"gamefed/srv/actors"
	"gamefed/srv/http"
	"gamefed/srv/log"
)

const path string = webfinger.DefaultPath // "/.well-known/webfinger"

func init() {
	// Skip this if we are running inside of a Go test.
	if nil != flag.Lookup("test.v") || strings.HasSuffix(os.Args[0], ".test") {
		return
	}

	var webFingerHandler webfinger.Handler = webfinger.HandlerFunc(serveWebFinger)
	var handler http.Handler = webfinger.HTTPHandler(webFingerHandler)

	err := httpsrv.Mux.HandlePath(handler, path)
	if nil != err {
		panic(err)
	}
}

func serveWebFinger(resource string, rels ...string) ([]byte, error) {
	log := logsrv.Begin(field.String("www.path", path))
	defer log.End()

	{
		actor, host, err := accturi.Split(resource)
		if nil == err {
			return serveActorHost(resource, actor, host)
		}
		if !errors.Is(err, accturi.ErrAcctURISchemeNotFound) {
			return nil, errhttp.Return(http.StatusBadRequest)
		}

		log.Tracef("actor, host = %q, %q", actor, host)
	}

	{
		//@TODO: handle other types of IRI/URI/URL schemes.
	}

	return nil, errhttp.Return(http.StatusNotFound)
}

func serveActorHost(resource string, actor string, host string) ([]byte, error) {
	log := logsrv.Begin(field.String("www.path", path))
	defer log.End()

	log.Trace(
		field.String("actor", actor),
		field.String("host", host),
	)

	if !actorssrv.ExistsByUserName(actor) {
		return nil, errhttp.Return(http.StatusNotFound)
	}

	//@TODO: put the `actor` into a canonical form.

	var (
		self        string = librefs.Actor(host, actor)
		inbox       string = librefs.ActorInBox(host, actor)
		inboxShared string = librefs.SharedInBox(host)
		outbox      string = librefs.ActorOutBox(host, actor)
	)
	log.Trace(
		field.String("JRD-self", self),
		field.String("JRD-outbox", outbox),
	)

	// Return JRD (JSON Resource Descriptor) document,
	// that is expected to be returned in a WebFinger response.
	//
	// For example, for the resource "acct:joeblow:something@host.example",
	// this could return
	//
	//	{
	//		"subject" : "joeblow:something@host.example",
	//		"aliases" :
	//		[
	//			"https://host.example/gozaar/something"
	//		],
	//		"links"   :
	//		[
	//			{
	//				"rel"  : "self",
	//				"type" : "application/activity+json",
	//				"href" : "https://host.example/gozaar/something",
	//			},
	//			{
	//				"rel"  : "https://www.w3.org/TR/activitypub/#inbox",
	//				"type" : "application/activity+json",
	//				"href" : "https://host.example/gozaar/something/inbox",
	//			},
	//			{
	//				"rel"  : "https://www.w3.org/TR/activitypub/#sharedInbox",
	//				"type" : "application/activity+json",
	//				"href" : "https://host.example/inbox",
	//			},
	//			{
	//				"rel"  : "https://www.w3.org/TR/activitypub/#outbox",
	//				"type" : "application/activity+json",
	//				"href" : "https://host.example/gozaar/something/outbox",
	//			}
	//		],
	//	}
	{
		var response webfinger.Response = webfinger.Response{
			Subject: opt.Something(resource),
			Aliases: []string{
				self,
			},
			Links: []webfinger.Link{
				webfinger.Link{
					Rel:  opt.Something("self"),
					Type: opt.Something("application/activity+json"),
					HRef: opt.Something(self),
				},
				webfinger.Link{
					Rel:  opt.Something("https://www.w3.org/TR/activitypub/#inbox"),
					Type: opt.Something("application/activity+json"),
					HRef: opt.Something(inbox),
				},
				webfinger.Link{
					Rel:  opt.Something("https://www.w3.org/TR/activitypub/#sharedInbox"),
					Type: opt.Something("application/activity+json"),
					HRef: opt.Something(inboxShared),
				},
				webfinger.Link{
					Rel:  opt.Something("https://www.w3.org/TR/activitypub/#outbox"),
					Type: opt.Something("application/activity+json"),
					HRef: opt.Something(outbox),
				},
			},
		}

		return response.MarshalJSON()
	}
}
