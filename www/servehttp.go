package verboten

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"io"
	"net/http"
	"os"
	"strings"

	"codeberg.org/reiver/go-field"
	"github.com/reiver/go-etag"
	"github.com/reiver/go-http500"

	"gamefed/srv/http"
	"gamefed/srv/log"
)

const path string = "/"

var digest string

func init() {
	// Skip this if we are running inside of a Go test.
	if nil != flag.Lookup("test.v") || strings.HasSuffix(os.Args[0], ".test") {
		return
	}

	{
		digestBytes := sha256.Sum256([]byte(webpage))
		digest = hex.EncodeToString(digestBytes[:])
        }

	{
		var handler http.Handler = http.HandlerFunc(serveHTTP)

		err := httpsrv.Mux.HandlePath(handler, path)
		if nil != err {
			panic(err)
		}
	}
}

func serveHTTP(responsewriter http.ResponseWriter, request *http.Request) {
	log := logsrv.Begin()
	defer log.End()

	if nil == responsewriter {
		log.Error(field.S("nil response-writer"))
		return
	}
	if nil == request {
		http500.InternalServerError(responsewriter, request)
		log.Error(field.S("nil request"))
		return
	}

	log.Debug(field.String("digest", digest))

	var eTag string = "sha256-" + digest
	log.Debug(field.String("etag", eTag))

	if etag.Handle(responsewriter, request, eTag) {
		log.Debug(
			field.S("etag caching HIT"),
			field.String("path", path),
			field.String("etag", eTag),
			field.String("digest", digest),
		)
		return
	} else {
		log.Debug(
			field.S("etag caching MISS"),
			field.String("path", path),
			field.String("etag", eTag),
			field.String("digest", digest),
		)
	}

	_, err := io.WriteString(responsewriter, webpage)
	if nil != err {
		log.Error(
			field.S("problem writing HTML content to client"),
			field.E(err),
		)
	}
}
