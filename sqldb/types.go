// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
)

var ErrDataInvalid = errors.New("invalid data")
var errNotByteArray = errors.Wrap(ErrDataInvalid, "expecting []byte")

type MapStrStr map[string]string

// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a MapStrStr) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *MapStrStr) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errNotByteArray
	}

	return json.Unmarshal(b, &a)
}

type Bytes []byte

func (b Bytes) Value() (driver.Value, error) {
	data := make([]byte, len(b)*2+3)
	data[0] = 'x'
	data[1] = '\''
	hex.Encode(data[2:], b)
	data[len(b)*2+2] = '\''
	return data, nil
}

func (b *Bytes) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return errNotByteArray
	}

	if data[0] != 'x' || data[1] != '\'' || data[len(data)-1] != '\'' {
		return errors.New("incorrect format")
	}

	*b = make([]byte, (len(data)-3)/2)
	_, err := hex.Decode(*b, data[2:len(data)-1])
	return err
}
