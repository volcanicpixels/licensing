package main

import (
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/appengine"
	"google.golang.org/appengine/file"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
)

type StorageContext interface {
	ReadFile(fileName string) ([]byte, error)
	//WriteFile(fileName string, data []byte) error
}

type storageContext struct {
	c      context.Context
	ctx    context.Context
	bucket string
}

func NewStorageContext(c context.Context) StorageContext {
	hc := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(c, storage.ScopeFullControl),
			Base:   &urlfetch.Transport{Context: c},
		},
	}

	ctx := cloud.NewContext(appengine.AppID(c), hc)

	return &storageContext{
		c:   c,
		ctx: ctx,
	}
}

func (sc *storageContext) Bucket() (string, error) {

	// The dev app server does not support getting the bucket name
	// IsDevAppServer has rather unhelpfully not been implemented
	if strings.HasPrefix(appengine.ServerSoftware(), "Development") {
		return "vp-licensing.appspot.com", nil
	}

	if sc.bucket == "" {
		var err error

		if sc.bucket, err = file.DefaultBucketName(sc.c); err != nil {
			return "", err
		}
	}

	return sc.bucket, nil
}

func (sc *storageContext) ReadFile(fileName string) (slurp []byte, err error) {

	bucket, err := sc.Bucket()

	if err != nil {
		return
	}

	log.Debugf(sc.c, "Reading file %v from bucket %v", fileName, bucket)

	rc, err := storage.NewReader(sc.ctx, bucket, fileName)

	if err != nil {
		return
	}

	defer rc.Close()

	slurp, err = ioutil.ReadAll(rc)
	return
}
