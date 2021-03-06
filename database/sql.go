package database

import (
	"database/sql/driver"
	"encoding/json"
)

// Map type is an alias for map[string]interface{}
//
// It implements native SQL driver interfaces and hence can
// be used for SQL json or jsonb types as a drop in replacement
// of golang native maps
type Map map[string]interface{}

// Interface implements the nulls.nullable interface.
func (m Map) Interface() interface{} {
	return map[string]interface{}(m)
}

// Scan implements the sql.Scanner interface.
// It allows to read the map from the database value.
func (m *Map) Scan(src interface{}) error {
	var b []byte

	switch t := src.(type) {
	case []byte:
		b = t
	case string:
		b = []byte(t)
	default:
		return ErrSQLMapInvalidScan
	}

	if err := json.Unmarshal(b, m); err != nil {
		return ErrSQLMapUnmarshalScanned(err)
	}

	return nil
}

// Value implements the driver.Valuer interface.
// It allows to convert the map to a driver.value.
func (m Map) Value() (driver.Value, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, ErrSQLMapMarshalValue(err)
	}

	return string(b), nil
}

// UnmarshalJSON will unmarshall JSON value into
// the map representation of this value.
func (m *Map) UnmarshalJSON(b []byte) error {
	var stuff map[string]interface{}

	if err := json.Unmarshal(b, &stuff); err != nil {
		return ErrSQLMapUnmarshalJSON(err)
	}

	if *m == nil {
		*m = Map{}
	}

	for key, value := range stuff {
		(*m)[key] = value
	}

	return nil
}

// UnmarshalText will unmarshall text value into
// the map representation of this value.
func (m Map) UnmarshalText(text []byte) error {
	if err := json.Unmarshal(text, &m); err != nil {
		return ErrSQLMapUnmarshalText(err)
	}

	return nil
}
