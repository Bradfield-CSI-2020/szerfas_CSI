package sandbox

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

func ToJSON(x interface{}) string {
	var buf bytes.Buffer
	marshal(reflect.ValueOf(x), &buf)
	return buf.String()
}

func marshal(v reflect.Value, w io.Writer) {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Fprint(w, "null")
		} else {
			marshal(v.Elem(), w)
		}

	case reflect.Slice:
		if v.IsNil() {
			fmt.Fprint(w, "null")
		} else {
			fmt.Fprint(w, "[")
			for i := 0; i < v.Len(); i++ {
				if i > 0 {
					fmt.Fprint(w, ", ")
				}
				marshal(v.Index(i), w)
			}
			fmt.Fprint(w, "]")
		}

	case reflect.Struct:
		fmt.Fprint(w, "{")
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				fmt.Fprint(w, ", ")
			}
			f := v.Type().Field(i)
			fmt.Fprintf(w, `"%s": `, f.Name)
			marshal(v.Field(i), w)
		}
		fmt.Fprint(w, "}")

	case reflect.Int:
		fmt.Fprint(w, v.Int())

	case reflect.String:
		fmt.Fprintf(w, "%q", v.String())

	default:
		panic(fmt.Sprintf("Unexpected kind %v", v.Kind()))
	}
}
