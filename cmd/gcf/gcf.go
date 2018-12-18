package gcf

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	logger "github.com/go-kit/kit/log"
	"github.com/pkg/errors"

	"github.com/ahmetb/goodbye/pkg/goodbyeutil"
)

// GoodbyeHandler responds to GCF requests.
func GoodbyeHandler(w http.ResponseWriter, r *http.Request) {
	s := time.Now()
	log := logger.WithPrefix(
		logger.NewSyncLogger(
			logger.NewJSONLogger(os.Stdout)), "timestamp", logger.DefaultTimestampUTC)
	log.Log("message", "starting run")
	if err := run(r.Context(), log); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Log("severity", "error", "error", err)
		fmt.Fprintf(w, "error: %+v", err)
		return
	}
	d := time.Since(s)
	fmt.Fprintf(w, "ok (request took %v)", d)
}

func run(ctx context.Context, log logger.Logger) error {
	bucket, object, err := readGCSConfig()
	if err != nil {
		return errors.Wrap(err, "failed to read gcs settings")
	}
	api, me, err := goodbyeutil.GetConfig()
	if err != nil {
		return errors.Wrap(err, "failed to get config")
	}
	prev, err := loadIDs(ctx, bucket, object)
	if err != nil {
		return errors.Wrap(err, "failed to load IDs")
	}
	cur, err := goodbyeutil.RunOnce(log, prev, api, me)
	if err != nil {
		return errors.Wrap(err, "goodbye run failed")
	}
	if err := saveIDs(ctx, bucket, object, cur); err != nil {
		return errors.Wrap(err, "failed to save IDs")
	}
	return nil
}

func readGCSConfig() (string, string, error) {
	gcsBucket, gcsObject := os.Getenv("GCS_BUCKET"), os.Getenv("GCS_OBJECT")
	if gcsBucket == "" {
		return "", "", errors.New("GCS_BUCKET not specified")
	}
	if gcsObject == "" {
		return "", "", errors.New("GCS_OBJECT not specified")
	}
	return gcsBucket, gcsObject, nil
}

func loadIDs(ctx context.Context, bucket, object string) ([]int64, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error initializing storage client")
	}
	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	s := string(data)
	var out []int64
	for _, v := range strings.Split(s, "\n") {
		if v == "" {
			continue
		}
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not parse id: %q", v)
		}
		out = append(out, i)
	}
	return out, nil
}

func saveIDs(ctx context.Context, bucket, object string, ids []int64) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return errors.Wrap(err, "error initializing storage client")
	}

	s := make([]string, len(ids))
	for i, v := range ids {
		s[i] = fmt.Sprintf("%d", v)
	}

	r := strings.NewReader(strings.Join(s, "\n"))
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, r); err != nil {
		return errors.Wrap(err, "write error into the object")
	}
	return errors.Wrap(wc.Close(), "write close error")
}
