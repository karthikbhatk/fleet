// Copyright (c) Facebook, Inc. and its affiliates.
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

package nvd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestCPE(t *testing.T) {
	td, err := ioutil.TempDir("", "nvdsync-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(td)

	handler := &cpeTestServer{}
	ts, src := httptestNewServer(handler)
	defer ts.Close()

	cases := make([]CPE, 0, len(SupportedCPE))
	for _, cve := range SupportedCPE {
		cases = append(cases, cve)
	}

	for _, cpe := range cases {
		label := []string{"CreateSync", "UseExistingSync"}
		for i := 0; i < 2; i++ {
			info := fmt.Sprintf("%s/%s", label[i], cpe)
			t.Run(info, func(t *testing.T) {
				err = cpe.Sync(context.Background(), src, td)
				if err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}

type cpeTestServer struct{}

func (ts cpeTestServer) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Etag", "foobar")
	fmt.Fprintf(w, "hello, world")
}
