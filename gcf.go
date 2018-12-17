package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
	logger "github.com/go-kit/kit/log"
	"github.com/pkg/errors"

	"github.com/ahmetb/goodbye/pkg/goodbyeutil"
)

// GoodbyeHandler responds to GCF requests.
func GoodbyeHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.WithPrefix(
		logger.NewSyncLogger(logger.NewLogfmtLogger(os.Stdout)), "time", logger.DefaultTimestampUTC)
	bucket, object, err := readGCSConfig()
	if err != nil {
		log.Log("error", err)
		fmt.Fprintf(w, "ERROR: %+v", err)
		return
	}

	api, me, err := goodbyeutil.GetConfig()
	if err != nil {
		log.Log("error", err)
		fmt.Fprintf(w, "ERROR: %+v", err)
		return
	}

	prev, err := loadIDs(r.Context(), bucket, object)
	if err != nil {
		log.Log("error", err)
		fmt.Fprintf(w, "ERROR: %+v", err)
		return
	}

	cur, err := goodbyeutil.RunOnce(log, prev, api, me)
	if err != nil {
		log.Log("error", err)
		fmt.Fprintf(w, "ERROR: %+v", err)
		return
	}
	if err := saveIDs(r.Context(), bucket, object, cur); err != nil {
		log.Log("error", err)
		fmt.Fprintf(w, "ERROR: %+v", err)
		return
	}

	fmt.Fprintf(w, "ok")
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
