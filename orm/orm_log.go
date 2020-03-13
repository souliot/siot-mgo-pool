// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package orm

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

// Log implement the log.Logger
type Log struct {
	*log.Logger
}

//costomer log func
var LogFunc func(query map[string]interface{})

// NewLog set io.Writer to create a Logger.
func NewLog(out io.Writer) *Log {
	d := new(Log)
	d.Logger = log.New(out, "[ORM]", log.LstdFlags)
	return d
}

func debugLogQueies(alias *alias, operaton, query string, t time.Time, err error, args ...interface{}) {
	var logMap = make(map[string]interface{})
	sub := time.Now().Sub(t) / 1e5
	elsp := float64(int(sub)) / 10.0
	logMap["cost_time"] = elsp
	flag := "  OK"
	if err != nil {
		flag = "FAIL"
	}
	logMap["flag"] = flag
	con := fmt.Sprintf(" -[Queries/%s] - [%s / %11s / %7.1fms] - [%s]", alias.Name, flag, operaton, elsp, query)
	cons := make([]string, 0, len(args))
	for _, arg := range args {
		cons = append(cons, fmt.Sprintf("%v", arg))
	}
	if len(cons) > 0 {
		con += fmt.Sprintf(" - `%s`", strings.Join(cons, "`, `"))
	}
	if err != nil {
		con += " - " + err.Error()
	}
	logMap["sql"] = fmt.Sprintf("%s-`%s`", query, strings.Join(cons, "`, `"))
	if LogFunc != nil {
		LogFunc(logMap)
	}
	DebugLog.Println(con)
}
