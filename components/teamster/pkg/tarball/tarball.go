/*
Copyright 2018 Pax Automa Systems, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tarball

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

type ManifestFile struct {
	Fstat   tar.Header
	Content ManifestFileContent
}

type ManifestFileContent func(interface{}, *bytes.Buffer) error

type Manifest []ManifestFile

func CreateTarPkg(manifest Manifest, ctx interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	gzwriter := gzip.NewWriter(buf)
	writer := tar.NewWriter(gzwriter)
	for _, file := range manifest {
		var data bytes.Buffer
		if err := file.Content(ctx, &data); err != nil {
			return nil, errors.Wrap(err, "cannot write manifest file")
		}
		cheader := file.Fstat
		cheader.Size = int64(data.Len())
		if err := writer.WriteHeader(&cheader); err != nil {
			return nil, errors.Wrap(err, "cannot write header")
		}

		if _, err := writer.Write(data.Bytes()); err != nil {
			return nil, errors.Wrap(err, "cannot write data")
		}
	}

	if err := writer.Close(); err != nil {
		return nil, errors.Wrap(err, "cannot close writer")
	}

	if err := gzwriter.Close(); err != nil {
		return nil, errors.Wrap(err, "cannot close gzip")
	}

	return buf.Bytes(), nil
}

func SendTarball(manifest Manifest, ctx interface{}, w http.ResponseWriter, filename string) {
	pkgBytes, err := CreateTarPkg(manifest, ctx)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(len(pkgBytes)))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.WriteHeader(http.StatusOK)
	w.Write(pkgBytes)
}
