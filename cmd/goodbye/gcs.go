package main

import (
	"context"
	"io"
	"net/url"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
)

type gcsObject struct{ o *storage.ObjectHandle }

func (g *gcsObject) WriteCloser() (io.WriteCloser, error) { return g.o.NewWriter(context.TODO()), nil }

func (g *gcsObject) ReadCloser() (io.ReadCloser, error) { return g.o.NewReader(context.TODO()) }

func openGCSObject(path string) (file, error) {
	client, err := storage.NewClient(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize gcs client")
	}
	b, o, err := splitGCSPath(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse path")
	}
	return &gcsObject{o: client.Bucket(b).Object(o)}, nil
}

func splitGCSPath(s string) (string, string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to parse file url")
	}
	if u.Scheme != "gs" {
		return "", "", errors.Errorf("path doesn't have gs:// prefix (%s)", u.Scheme)
	}
	return u.Host, strings.TrimPrefix(u.Path, "/"), nil
}
