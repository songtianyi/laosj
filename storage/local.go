// Copyright 2016 laosj Author @songtianyi. All Rights Reserved.
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

package storage

import (
	"github.com/songtianyi/rrframework/utils"
	"os"
)

type LocalDiskStorage struct {
	Dir string // the directory where to save binary
}

func NewLocalDiskStorage(dir string) StorageWrapper {
	s := &LocalDiskStorage{
		Dir: dir,
	}
	return s
}

func (s *LocalDiskStorage) Save(data []byte) (error, string) {
	// random name
	filename := rrutils.NewV4().String()

	//open a file for writing
	file, err := os.Create(s.Dir + filename)
	if err != nil {
		return err, filename
	}
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return err, filename
	}
	return nil, filename

}
