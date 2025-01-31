// Package std is a package that offers standard functions
//
//	Author: Elizalde G. Baguinon
//	Created: October 17, 2019
package std

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	ssd "github.com/shopspring/decimal"
	"golang.org/x/exp/constraints"
)

type (
	FieldTypeConstraint interface {
		constraints.Ordered | time.Time | ssd.Decimal | bool | byte
	}
	SeriesOptions struct {
		Prefix string // Prefix of series
		Suffix string // Suffix of series
		Length int    // Fixed length of the series
	}
)

// AnyToString converts any variable to string
func AnyToString(value interface{}) string {
	var b string
	if value == nil {
		return ""
	}
	switch t := value.(type) {
	case string:
		b = t
	case int:
		b = strconv.FormatInt(int64(t), 10)
	case int8:
		b = strconv.FormatInt(int64(t), 10)
	case int16:
		b = strconv.FormatInt(int64(t), 10)
	case int32:
		b = strconv.FormatInt(int64(t), 10)
	case int64:
		b = strconv.FormatInt(t, 10)
	case uint:
		b = strconv.FormatUint(uint64(t), 10)
	case uint8:
		b = strconv.FormatUint(uint64(t), 10)
	case uint16:
		b = strconv.FormatUint(uint64(t), 10)
	case uint32:
		b = strconv.FormatUint(uint64(t), 10)
	case uint64:
		b = strconv.FormatUint(uint64(t), 10)
	case float32:
		b = fmt.Sprintf("%f", t)
	case float64:
		b = fmt.Sprintf("%f", t)
	case bool:
		if t {
			return "true"
		} else {
			return "false"
		}
	case time.Time:
		b = "'" + t.Format(time.RFC3339) + "'"
	case *string:
		if t == nil {
			return ""
		}
		b = *t
	case *int:
		if t == nil {
			return "0"
		}
		b = strconv.FormatInt(int64(*t), 10)
	case *int8:
		if t == nil {
			return "0"
		}
		b = strconv.FormatInt(int64(*t), 10)
	case *int16:
		if t == nil {
			return "0"
		}
		b = strconv.FormatInt(int64(*t), 10)
	case *int32:
		if t == nil {
			return "0"
		}
		b = strconv.FormatInt(int64(*t), 10)
	case *int64:
		if t == nil {
			return "0"
		}
		b = strconv.FormatInt(*t, 10)
	case *uint:
		if t == nil {
			return "0"
		}
		b = strconv.FormatUint(uint64(*t), 10)
	case *uint8:
		if t == nil {
			return "0"
		}
		b = strconv.FormatUint(uint64(*t), 10)
	case *uint16:
		if t == nil {
			return "0"
		}
		b = strconv.FormatUint(uint64(*t), 10)
	case *uint32:
		if t == nil {
			return "0"
		}
		b = strconv.FormatUint(uint64(*t), 10)
	case *uint64:
		if t == nil {
			return "0"
		}
		b = strconv.FormatUint(uint64(*t), 10)
	case *float32:
		if t == nil {
			return "0"
		}
		b = fmt.Sprintf("%f", *t)
	case *float64:
		if t == nil {
			return "0"
		}
		b = fmt.Sprintf("%f", *t)
	case *bool:
		if t == nil || !*t {
			return "false"
		}
		return "true"
	case *time.Time:
		if t == nil {
			return "'" + time.Time{}.Format(time.RFC3339) + "'"
		}
		tm := *t
		b = "'" + tm.Format(time.RFC3339) + "'"
	}

	return b
}

// BuildSeries builds series based on options
func BuildSeries(series int, opt SeriesOptions) string {
	// If length is specified, we get the difference between suffix and prefix
	if opt.Length > 0 {
		diff := opt.Length - (len(opt.Prefix) + len(opt.Suffix))
		ds := `%0` + strconv.Itoa(diff) + `d`
		return fmt.Sprintf(`%s`+ds+`%s`, opt.Prefix, series, opt.Suffix)
	}
	return fmt.Sprintf(`%s%d%s`, opt.Prefix, series, opt.Suffix)
}

// Elem returns the element of an array as specified by the index
//
// If the index exceeds the length of an array, it will return a non-nil value of the type.
// To monitor if the element exists, define a boolean value in the exists parameter
//
// Currently supported data types are
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func Elem[T any](array *[]T, index int, exists *bool) T {
	var result T
	if exists != nil {
		*exists = false
	}
	if array == nil {
		return result
	}
	arrl := len(*array)
	if arrl == 0 {
		return result
	}
	arrl--
	if arrl >= index {
		if exists != nil {
			*exists = true
		}
		return (*array)[index]
	}
	return result
}

// ElemPtr returns a pointer to the element of an array as specified by the index
//
// If the index exceeds the length of an array, it will return a non-nil value of the type.
// To monitor if the element exists, define a boolean value in the exists parameter
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func ElemPtr[T any](array *[]T, index int, exists *bool) *T {
	r := Elem(array, index, exists)
	return &r
}

// If is a basic ternary operator to return whatever is set in
// truthy and falsey parameter.
// If the subject is nil, empty string, 0, -0 or false, it will return the falsey parameter
//
// This function requires version 1.18+
func If[T constraints.Ordered](subject any, truthy T, falsey T) T {
	if subject == nil {
		return falsey
	}
	switch t := subject.(type) {
	case string:
		if t == "" {
			return falsey
		}
	case *string:
		if t == nil || *t == "" {
			return falsey
		}
	case
		int8, int16, int32, int64, int,
		uint8, uint16, uint32, uint64, uint,
		float32, float64, complex64, complex128:
		if t == 0 || t == -0 {
			return falsey
		}
	case
		*int8, *int16, *int32, *int64, *int,
		*uint8, *uint16, *uint32, *uint64, *uint,
		*float32, *float64, *complex64, *complex128:
		vo := reflect.ValueOf(t)
		tx := vo.Elem()
		if !tx.IsValid() || tx.IsZero() {
			return falsey
		}
	case bool:
		if !t {
			return falsey
		}
	case *bool:
		if !*t {
			return falsey
		}
	}
	return truthy
}

// In checks if the seek parameter is in the list parameter
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func In[T comparable](seek T, list ...T) bool {
	for _, li := range list {
		if li == seek {
			return true
		}
	}
	return false
}

// IsNullOrEmpty checks for emptiness of a pointer variable ignoring nullity
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func IsEmpty[T FieldTypeConstraint](value *T) bool {
	// if value == nil {
	// 	return false
	// }
	// if *value == GetZero[T]() {
	// 	return true
	// }
	// return false
	return value != nil && *value == getZero[T]()
}

// IsNullOrEmpty checks for nullity and emptiness of a pointer variable
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func IsNullOrEmpty[T FieldTypeConstraint](value *T) bool {
	return value == nil || *value == getZero[T]()
}

// IsNumeric checks if a string is numeric
func IsNumeric(value string) error {
	if value == "" {
		return fmt.Errorf("is empty")
	}
	if _, err := strconv.ParseFloat(value, 64); err != nil {
		return fmt.Errorf(`is not a number (%s)`, err)
	}
	return nil
}

// IsInterfaceNil checks if an interface is nil
func IsInterfaceNil(i interface{}) bool {
	if i == nil {
		return true
	}
	iv := reflect.ValueOf(i)
	if !iv.IsValid() {
		return true
	}
	switch iv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
		return iv.IsNil()
	default:
		return false
	}
}

// MapVal retrieves a value from a map by a key and converts it to the type indicated by T.
// Returns a pointer to the value if found, or nil if not found.
//
// The third parameter, dateLayout can be set with many time layouts. The specified layouts
// will be the only one to try parsing.
//
// If not set, the built-in date layouts of ParseDate are used. See function for supported layouts
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func MapVal[T FieldTypeConstraint](kvmap *map[string]any, key string, dateLayout ...string) *T {
	var (
		ok bool
	)
	if kvmap == nil {
		return nil
	}
	miv, ok := (*kvmap)[key]
	if !ok {
		return nil
	}
	mv, ok := miv.(T)
	if ok {
		return &mv
	}
	t := new(T)  // initialize a variable to the return generic type
	a := any(*t) // initialize a variable to the value type of the generic type
	switch a.(type) {
	case time.Time:
		if s, ok := miv.(string); ok {
			var dlo *string
			if len(dateLayout) > 0 {
				dlo = &dateLayout[0]
			}
			v, _, err := ParseDate(s, dlo)
			if err != nil {
				return nil
			}
			*t = any(v).(T) // Convert the parsed value to any before asserting the type for the return
			return t
		}
	case ssd.Decimal:
		if s, ok := miv.(string); ok {
			v, err := ssd.NewFromString(s)
			if err != nil {
				return nil
			}
			*t = any(v).(T)
			return t
		}
	default:

	}
	return nil
}

// New initializes a variable and returns a pointer of its type
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func New[T FieldTypeConstraint](value T) *T {
	n := new(T)
	*n = value
	return n
}

// NonNullComp compares two parameters when both are not nil.
//
//   - When one or both of the parameters is nil, the function returns -1
//   - When the parameters are equal, the function returns 0.
//   - else it returns 1
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func NonNullComp[T FieldTypeConstraint](param1 *T, param2 *T) int {
	if param1 == nil || param2 == nil {
		return -1
	}
	if *param1 == *param2 {
		return 0
	}
	return 1
}

// Null accepts a value to test and the default value
// if it fails. It returns a non-pointer value of T.
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func Null[T any](testValue any, defaultValue any) T {
	var (
		def T
	)
	if defaultValue == nil {
		defaultValue = def
	}
	if testValue == nil {
		return defaultValue.(T)
	}
	vo := reflect.ValueOf(testValue)
	if k := vo.Kind(); k == reflect.Map ||
		k == reflect.Func ||
		k == reflect.Ptr ||
		k == reflect.Slice ||
		k == reflect.Interface {
		if vo.IsZero() && vo.IsNil() {
			return defaultValue.(T)
		}
		ifv := vo.Elem().Interface()
		return ifv.(T)
	}
	return def
}

// NullPtr accepts a value to test and the default value
// if it fails. It returns a pointer value of T.
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func NullPtr[T any](testValue any, defaultValue any) *T {
	val := Null[T](testValue, defaultValue)
	return &val
}

// ParseDate parses a string as date.
//
// If dateLayout is not provided, this function try all layout combinations.
// The following date layouts has been provided:
//   - "2006-01-02"
//   - "2006-02-01"
//   - "01-02-2006"
//   - "02-01-2006"
//   - "2006/01/02"
//   - "2006/02/01"
//   - "01/02/2006"
//   - "02/01/2006"
//   - "1/2/2006"
//   - "2/1/2006"
//   - "2006/1/2"
//   - "2006/2/1"
//   - "06-01-02"
//   - "06-02-01"
//   - "01-02-06"
//   - "02-01-06"
//   - "06/01/02"
//   - "06/02/01"
//   - "01/02/06"
//   - "02/01/06"
//
// The date layout partitions means:
//   - Anything with 1, with or without zero is the month
//   - Anything with 2, with or without zero is the day
//   - Anything with 06, with or without the prefix 20 is the year
func ParseDate(dtText string, dateLayout *string) (time.Time, string, error) {
	var (
		rtm time.Time
		rlo string
		err error
	)
	if dtText == "" {
		return rtm, rlo, fmt.Errorf("invalid date or time input")
	}
	dlo :=
		[]string{
			// Dashed full year
			"2006-01-02",
			"2006-02-01",
			"01-02-2006",
			"02-01-2006",
			// Forward-slashed full year
			"2006/01/02",
			"2006/02/01",
			"01/02/2006",
			"02/01/2006",
			// Forward-slashed single digit day and month
			"1/2/2006",
			"2/1/2006",
			"2006/1/2",
			"2006/2/1",
			// 2-digit dashed year with full digit day and month
			"06-01-02",
			"06-02-01",
			"01-02-06",
			"02-01-06",
			// 2-digit forward-slashed year with full digit day and month
			"06/01/02",
			"06/02/01",
			"01/02/06",
			"02/01/06",
		}
	// Try to parse using layout provided
	// The function will return upon failure or success
	if dateLayout != nil {
		// Check if layout is in the array
		if !In(*dateLayout, dlo...) {
			return rtm, *dateLayout, fmt.Errorf("layout provided not supported")
		}
		rtm, err = time.Parse(*dateLayout, dtText)
		return rtm, *dateLayout, err
	}
	// Try each layout until it succeeds
	for _, lo := range dlo {
		rtm, err = time.Parse(lo, dtText)
		if err != nil {
			continue
		}
		return rtm, lo, err
	}
	return rtm, rlo, fmt.Errorf("date parsing failed")
}

// SafeMapWrite allows writing to maps by locking, preventing the library from crashing
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func SafeMapWrite[T any](ptrMap *map[string]T, key string, value T, rw *sync.RWMutex) bool {
	defer func() {
		recover()
	}()
	// Prepare mutex
	// attempt writing to map
	if rw.TryLock() {
		defer rw.Unlock()
		(*ptrMap)[key] = value
	}
	return true
}

// SafeMapRead allows reading maps by locking it, preventing the library from crashing
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func SafeMapRead[T any](ptrMap *map[string]T, key string, rw *sync.RWMutex) T {
	var result T
	defer func() {
		recover()
	}()
	if rw.TryLock() {
		defer rw.Unlock()
		result = (*ptrMap)[key]
	}
	return result
}

// Seek checks if the seek parameter is in the list parameter and returns it.
// If the value is not found in the list, the function returns nil
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func Seek[T comparable](seek T, list ...T) *T {
	for _, li := range list {
		if li == seek {
			return &li
		}
	}
	return nil
}

// ToInterfaceArray converts a value to interface array
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func ToInterfaceArray[T FieldTypeConstraint](value T) []any {
	var val [1]any
	val[0] = value
	return val[:]
}

// Val gets the value of a pointer in order
//
// Currently supported data types are:
//   - constraints.Ordered (Integer | Float | ~string)
//   - time.Time
//   - bool
//   - shopspring/decimal
//
// This function requires version 1.18+
func Val[T FieldTypeConstraint](value *T) T {
	if value == nil {
		return getZero[T]()
	}
	return *value
}

func getZero[T FieldTypeConstraint | bool]() T {
	var r T
	return r
}
