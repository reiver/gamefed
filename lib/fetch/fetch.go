package libfetch

import (
	"context"
	"errors"
	"io"
	"net/http"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
)

const (
	ErrFetchFailed      = erorr.Error("fetch failed")
	ErrResponseTooLarge = erorr.Error("response too large")
	ErrTimeOut          = erorr.Error("time-out")
)

const maxBodySize = 536870912

func Fetch(ctx context.Context, url string) ([]byte, error) {
	var nada []byte

	if nil == ctx {
		ctx = context.Background()
	}

//@TODO: support more than just HTTP, HTTPS.

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if nil != err {
		err = erorr.Errors{ErrFetchFailed, err}
		return nada, err
	}
	resp, err := http.DefaultClient.Do(req)
	if nil != err {
		err = erorr.Errors{ErrFetchFailed, err}
		err = erorr.Wrap(err, "failed to fetch content of URL",
			field.String("url", url),
		)
		return nada, err
	}
//@TODO: is this the correct place for it?
	defer resp.Body.Close()

	if http.StatusOK != resp.StatusCode {
		var err error = ErrFetchFailed
		err = erorr.Wrap(err, "failed to get HTTP OK response",
			field.Int("http-status-code", resp.StatusCode),
		)
		return nada, err
//@TODO
	}

	var bytes []byte
	{
		// Read up to maxBodySize + 1 to detect oversized responses.
		var limitReader io.Reader = io.LimitReader(resp.Body, int64(maxBodySize+1))

		var err error
		bytes, err = io.ReadAll(limitReader)
		if nil != err {
			if errors.Is(ctx.Err(), context.Canceled) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
/////// RETURN
				//@TODO: Should we include other errors.

				var nada []byte
				return nada, ErrTimeOut
			}
			{
				err = erorr.Errors{ErrFetchFailed, err}
/////// RETURN
				var nada []byte
				return nada, err
			}
		}

		if maxBodySize < len(bytes) {
/////// RETURN
			var nada []byte
			return nada, ErrResponseTooLarge
		}
	}

	return bytes, nil
}
