package uuid

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

var _ driver.Valuer = UUID("")
var _ sql.Scanner = (*UUID)(nil)

// Value implements the driver.Valuer interface.
func (u UUID) Value() (driver.Value, error) {
	return u.String(), nil
}

// Scan implements the sql.Scanner interface.
// A 16-byte slice will be handled by UnmarshalBinary, while
// a longer byte slice or a string will be handled by UnmarshalText.
func (u *UUID) Scan(src interface{}) error {
	switch src := src.(type) {
	case UUID:
		*u = src
		return nil
	case []byte:
		if len(src) == Size {
			return u.UnmarshalBinary(src)
		}
		return u.UnmarshalText(src)
	case string:
		uu, err := Parse(src)
		*u = uu
		return err
	}

	return fmt.Errorf("cannot convert %T to UUID", src)
}
