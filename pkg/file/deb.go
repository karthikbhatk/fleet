package file

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"

	"github.com/blakesmith/ar"
	"github.com/xi2/xz"
)

// ExtractDebMetadata extracts the name and version metadata from a .deb file ,
// a debian installer package which is in archive format.
func ExtractDebMetadata(b []byte) (name, version string, err error) {
	r := ar.NewReader(bytes.NewReader(b))

	for {
		hdr, err := r.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", "", fmt.Errorf("failed to advance to next file in archive: %w", err)
		}

		name := path.Clean(hdr.Name)
		if strings.HasPrefix(name, "control.tar") {
			ext := filepath.Ext(name)
			if ext == ".tar" {
				ext = ""
			}
			return parseControl(r, ext)
		}
	}

	// no control.tar file found, return empty information
	return "", "", nil
}

// parseControl adapted from
// https://github.com/sassoftware/relic/blob/6c510a666832163a5d02587bda8be970d5e29b8c/lib/signdeb/control.go#L38-L39
//
// Copyright (c) SAS Institute Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Parse basic package info from a control.tar.* stream.
func parseControl(r io.Reader, ext string) (name, version string, err error) {
	switch ext {
	case ".gz":
		r, err = gzip.NewReader(r)
	case ".bz2":
		r = bzip2.NewReader(r)
	case ".xz":
		r, err = xz.NewReader(r, 0)
	case "":
		// uncompressed
	default:
		err = errors.New("unrecognized compression on control.tar: " + ext)
	}
	if err != nil {
		return "", "", err
	}

	tr := tar.NewReader(r)
	found := false
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", "", err
		}
		if path.Clean(hdr.Name) == "control" {
			found = true
			break
		}
	}

	if !found {
		return "", "", errors.New("control.tar has no control file")
	}

	blob, err := io.ReadAll(tr)
	if err != nil {
		return "", "", fmt.Errorf("failed to read tar file: %w", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(blob))
	for scanner.Scan() {
		line := scanner.Text()
		i := strings.IndexAny(line, " \t\r\n")
		j := strings.Index(line, ":")
		if j < 0 || i < j {
			continue
		}

		key := line[:j]
		value := strings.Trim(line[j+1:], " \t\r\n")
		switch strings.ToLower(key) {
		case "package":
			name = value
		case "version":
			version = value
		}
	}
	if err := scanner.Err(); err != nil {
		return name, version, fmt.Errorf("failed to scan control file: %w", err)
	}
	return name, version, nil
}
