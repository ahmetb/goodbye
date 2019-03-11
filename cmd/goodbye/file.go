package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type file interface {
	ReadCloser() (io.ReadCloser, error)
	WriteCloser() (io.WriteCloser, error)
}

func loadIDs(f file) ([]int64, error) {
	rc, err := f.ReadCloser()
	if err != nil {
		return nil, errors.Wrap(err, "failed to open for reading")
	}
	defer rc.Close()
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read from reader")
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

func saveIDs(f file, ids []int64) error {
	wc, err := f.WriteCloser()
	if err != nil {
		return errors.Wrap(err, "could not get a writer")
	}
	defer wc.Close()

	s := make([]string, len(ids))
	for i, v := range ids {
		s[i] = fmt.Sprintf("%d", v)
	}
	r := strings.NewReader(strings.Join(s, "\n"))
	if _, err = io.Copy(wc, r); err != nil {
		return errors.Wrap(err, "write error into the object")
	}
	return errors.Wrap(wc.Close(), "write close error") // important to handle close for GCS as it surfaces significant errors
}
