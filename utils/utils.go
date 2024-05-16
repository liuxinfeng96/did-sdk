/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) DCPS. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"bytes"
	"encoding/json"
	"time"
)

func ISO8601Time(t int64) string {
	unix := time.Unix(t, 0)
	return unix.Format(time.RFC3339)
}

func ISO8601TimeToUnix(t string) (int64, error) {
	ti, err := time.ParseInLocation(time.RFC3339, t, time.Local)
	if err != nil {
		return 0, err
	}

	return ti.Unix(), nil
}

// CompactJson 压缩json字符串，去掉空格换行等
func CompactJson(raw []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := json.Compact(&buf, raw)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
