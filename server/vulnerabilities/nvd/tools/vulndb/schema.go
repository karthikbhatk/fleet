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

package vulndb

import (
	"context"
	"database/sql"
	"encoding/base64"
)

// InitSchemaSQL is auto-generated. Executes each SQL statement from schema.sql.
func InitSchemaSQL(ctx context.Context, db *sql.DB) error {
	for _, stmt := range SchemaSQL() {
		_, err := db.ExecContext(ctx, stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// SchemaSQL is auto-generated. Returns each SQL statement from schema.sql.
func SchemaSQL() []string {
	s := make([]string, len(b64SchemaSQL))
	for i := 0; i < len(s); i++ {
		v, _ := base64.StdEncoding.DecodeString(b64SchemaSQL[i])
		s[i] = string(v)
	}
	return s
}

// b64SchemaSQL is auto-generated from schema.sql.
var b64SchemaSQL = []string{"LS0gQ29weXJpZ2h0IChjKSBGYWNlYm9vaywgSW5jLiBhbmQgaXRzIGFmZmlsaWF0ZXMuCi0tCi0tIExpY2Vuc2VkIHVuZGVyIHRoZSBBcGFjaGUgTGljZW5zZSwgVmVyc2lvbiAyLjAgKHRoZSAiTGljZW5zZSIpOwo=", "LS0geW91IG1heSBub3QgdXNlIHRoaXMgZmlsZSBleGNlcHQgaW4gY29tcGxpYW5jZSB3aXRoIHRoZSBMaWNlbnNlLgotLSBZb3UgbWF5IG9idGFpbiBhIGNvcHkgb2YgdGhlIExpY2Vuc2UgYXQKLS0KLS0gICAgIGh0dHA6Ly93d3cuYXBhY2hlLm9yZy9saWNlbnNlcy9MSUNFTlNFLTIuMAotLQotLSBVbmxlc3MgcmVxdWlyZWQgYnkgYXBwbGljYWJsZSBsYXcgb3IgYWdyZWVkIHRvIGluIHdyaXRpbmcsIHNvZnR3YXJlCi0tIGRpc3RyaWJ1dGVkIHVuZGVyIHRoZSBMaWNlbnNlIGlzIGRpc3RyaWJ1dGVkIG9uIGFuICJBUyBJUyIgQkFTSVMsCi0tIFdJVEhPVVQgV0FSUkFOVElFUyBPUiBDT05ESVRJT05TIE9GIEFOWSBLSU5ELCBlaXRoZXIgZXhwcmVzcyBvciBpbXBsaWVkLgotLSBTZWUgdGhlIExpY2Vuc2UgZm9yIHRoZSBzcGVjaWZpYyBsYW5ndWFnZSBnb3Zlcm5pbmcgcGVybWlzc2lvbnMgYW5kCi0tIGxpbWl0YXRpb25zIHVuZGVyIHRoZSBMaWNlbnNlLgoKRFJPUCBUQUJMRSBJRiBFWElTVFMKCWBzbm9vemVgLAoJYGN1c3RvbV9kYXRhYCwKCWB2ZW5kb3JfZGF0YWAsCglgdmVuZG9yYAo7Cg==", "Q1JFQVRFIFRBQkxFIGB2ZW5kb3JgICgKCWB2ZXJzaW9uYCAgSU5UICAgICAgICAgTk9UIE5VTEwgQVVUT19JTkNSRU1FTlQgQ09NTUVOVCAnSUQgb2YgdGhlIGRhdGFzZXQnLAoJYHRzYCAgICAgICBUSU1FU1RBTVAgICBOT1QgTlVMTCAgQ09NTUVOVCAnVGltZSBvZiB0aGUgZGF0YXNldCBpbXBvcnQnLAoJYHJlYWR5YCAgICBCT09MICAgICAgICBOT1QgTlVMTCAgQ09NTUVOVCAnSW5kaWNhdGVzIHRoZSBkYXRhc2V0IGlzIHJlYWR5IHRvIHVzZScsCglgb3duZXJgICAgIFZBUkNIQVIoNjQpIE5PVCBOVUxMICBDT01NRU5UICdQb2ludCBvZiBjb250YWN0IGZvciBkYXRhc2V0JywKCWBwcm92aWRlcmAgVkFSQ0hBUig2NCkgTk9UIE5VTEwgIENPTU1FTlQgJ1Nob3J0IG5hbWUgb2YgZGF0YXNldCBwcm92aWRlcicsCglQUklNQVJZIEtFWSAoYHZlcnNpb25gKSwKCUtFWSAoYHByb3ZpZGVyYCkKKQpFTkdJTkUgSW5ub0RCCkRFRkFVTFQgQ0hBUkFDVEVSIFNFVCB1dGY4bWI0CkNPTU1FTlQgJ1ZlbmRvcnMgcHJvdmlkaW5nIHZ1bG5lcmFiaWxpdHkgZGF0YXNldHMnCjsK", "Q1JFQVRFIFRBQkxFIGB2ZW5kb3JfZGF0YWAgKAoJYHZlcnNpb25gICAgIElOVCAgICAgICAgICBOT1QgTlVMTCBDT01NRU5UICdJRCBvZiB0aGUgdmVuZG9yIGRhdGFzZXQnLAoJYGN2ZV9pZGAgICAgIFZBUkNIQVIoMTI4KSBOT1QgTlVMTCBDT01NRU5UICdDb21tb24gVnVsbmVyYWJpbGl0eSBhbmQgRXhwb3N1cmUgKENWRSkgSUQnLAoJYHB1Ymxpc2hlZGAgIFRJTUVTVEFNUCAgICBOT1QgTlVMTCBDT01NRU5UICdUaW1lc3RhbXAgb2YgdnVsbmVyYWJpbGl0eSBwdWJsaWNhdGlvbicgREVGQVVMVCBDVVJSRU5UX1RJTUVTVEFNUCwKCWBtb2RpZmllZGAgICBUSU1FU1RBTVAgICAgTk9UIE5VTEwgQ09NTUVOVCAnVGltZXN0YW1wIG9mIHZ1bG5lcmFiaWxpdHkgbGFzdCBtb2RpZmljYXRpb24nIERFRkFVTFQgQ1VSUkVOVF9USU1FU1RBTVAsCglgYmFzZV9zY29yZWAgRkxPQVQoMywxKSAgIE5PVCBOVUxMIENPTU1FTlQgJ0Jhc2Ugc2NvcmUgZnJvbSBDVlNTIDMuMCBvciAyLjAgZmFsbGJhY2snLAoJYHN1bW1hcnlgICAgIFRFWFQgICAgICAgICBOT1QgTlVMTCBDT01NRU5UICdEZXNjcmlwdGlvbiBvZiB0aGUgdnVsbmVyYWJpbGl0eScsCglgY3ZlX2pzb25gICAgTUVESVVNQkxPQiAgIE5PVCBOVUxMIENPTU1FTlQgJ0pTT04gcmVjb3JkIGNvbnRhaW5pbmcgcmF3IENWRSBkYXRhJywKCVBSSU1BUlkgS0VZIChgdmVyc2lvbmAsIGBjdmVfaWRgKQopCkVOR0lORSBJbm5vREIKREVGQVVMVCBDSEFSQUNURVIgU0VUIHV0ZjhtYjQKQ09NTUVOVCAnVnVsbmVyYWJpbGl0eSBkYXRhIGZyb20gdmVuZG9ycycKOwo=", "Q1JFQVRFIFRBQkxFIGBjdXN0b21fZGF0YWAgKAoJYG93bmVyYCAgICAgICBWQVJDSEFSKDY0KSAgTk9UIE5VTEwgQ09NTUVOVCAnUG9pbnQgb2YgY29udGFjdCBmb3IgZGF0YXNldCcsCglgcHJvdmlkZXJgICAgIFZBUkNIQVIoNjQpICBOT1QgTlVMTCBDT01NRU5UICdTaG9ydCBuYW1lIG9mIGRhdGEgcHJvdmlkZXInLAoJYGN2ZV9pZGAgICAgICBWQVJDSEFSKDEyOCkgTk9UIE5VTEwgQ09NTUVOVCAnQ29tbW9uIFZ1bG5lcmFiaWxpdHkgYW5kIEV4cG9zdXJlIElEJywKCWBwdWJsaXNoZWRgICAgVElNRVNUQU1QICAgIE5PVCBOVUxMIENPTU1FTlQgJ1RpbWVzdGFtcCBvZiB2dWxuZXJhYmlsaXR5IHB1YmxpY2F0aW9uJyBERUZBVUxUIENVUlJFTlRfVElNRVNUQU1QLAoJYG1vZGlmaWVkYCAgICBUSU1FU1RBTVAgICAgTk9UIE5VTEwgQ09NTUVOVCAnVGltZXN0YW1wIG9mIGN1c3RvbWl6ZWQgbGFzdCBtb2RpZmljYXRpb24nIERFRkFVTFQgQ1VSUkVOVF9USU1FU1RBTVAsCglgYmFzZV9zY29yZWAgIEZMT0FUKDMsMSkgICBOT1QgTlVMTCBDT01NRU5UICdCYXNlIHNjb3JlIGZyb20gQ1ZTUyAzLjAgb3IgMi4wIGZhbGxiYWNrJywKCWBzdW1tYXJ5YCAgICAgVEVYVCAgICAgICAgIE5PVCBOVUxMIENPTU1FTlQgJ0Rlc2NyaXB0aW9uIG9mIHRoZSB2dWxuZXJhYmlsaXR5JywKCWBjdmVfanNvbmAgICAgTUVESVVNQkxPQiAgIE5PVCBOVUxMIENPTU1FTlQgJ0pTT04gcmVjb3JkIGNvbnRhaW5pbmcgcmF3IENWRSBkYXRhJywKCVBSSU1BUlkgS0VZIChgY3ZlX2lkYCkKKQpFTkdJTkUgSW5ub0RCCkRFRkFVTFQgQ0hBUkFDVEVSIFNFVCB1dGY4bWI0CkNPTU1FTlQgJ0N1c3RvbSB2dWxuZXJhYmlsaXR5IGRhdGEgaW5jbHVkaW5nIG92ZXJyaWRlcycKOwo=", "Q1JFQVRFIFRBQkxFIGBzbm9vemVgICgKCWBvd25lcmAgICAgIFZBUkNIQVIoNjQpICBOT1QgTlVMTCBDT01NRU5UICdQb2ludCBvZiBjb250YWN0IGZvciBzbm9vemUnLAoJYGNvbGxlY3RvcmAgdmFyY2hhcig2NCkgIE5PVCBOVUxMIENPTU1FTlQgJ1VuaXF1ZSBuYW1lIG9mIHRoZSBkYXRhIGNvbGxlY3RvcicsCglgcHJvdmlkZXJgICBWQVJDSEFSKDMyKSAgTk9UIE5VTEwgQ09NTUVOVCAnU2hvcnQgbmFtZSBvZiBkYXRhIHByb3ZpZGVyJywKCWBjdmVfaWRgICAgIFZBUkNIQVIoMTI4KSBOT1QgTlVMTCBDT01NRU5UICdDb21tb24gVnVsbmVyYWJpbGl0eSBhbmQgRXhwb3N1cmUgSUQnLAoJYGRlYWRsaW5lYCAgVElNRVNUQU1QICAgICAgICBOVUxMIENPTU1FTlQgJ1RpbWVzdGFtcCBvZiBzbm9vemUgZXhwaXJhdGlvbicgREVGQVVMVCBDVVJSRU5UX1RJTUVTVEFNUCwKCWBtZXRhZGF0YWAgIEJMT0IgICAgICAgICAgICAgTlVMTCBDT01NRU5UICdPcGFxdWUgbWV0YWRhdGEgZm9yIHNub296ZSBtYW5hZ2VtZW50JywKCVBSSU1BUlkgS0VZIChgcHJvdmlkZXJgLCBgY3ZlX2lkYCkKKQpFTkdJTkUgSW5ub0RCCkRFRkFVTFQgQ0hBUkFDVEVSIFNFVCB1dGY4bWI0CkNPTU1FTlQgJ1Z1bG5lcmFiaWxpdHkgcmVjb3JkcyB0byBpZ25vcmUgZm9yIGEgcGVyaW9kIG9mIHRpbWUnCjsK"}
