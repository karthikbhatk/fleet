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

// OSS MySQL connector

package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // register mysql driver
)

const mysqlDriver = "mysql"

func openRead(dsn string) (*sql.DB, error)  { return openDB(dsn) }
func openWrite(dsn string) (*sql.DB, error) { return openDB(dsn) }
