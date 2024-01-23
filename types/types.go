// package main

// import (
// 	"fmt"
// 	"reflect"
// )

// type Email = string
// type Email1 string

// //	func (e Email1) Validate() bool {
// //		return true
// //	}
// func Do() {
// 	var e Email = "email"
// 	var s string = "email"
// 	var ee Email1 = "email1"
// 	// ee.Validate()
// 	fmt.Println(reflect.TypeOf(e))
// 	fmt.Println(reflect.TypeOf(s))
// 	fmt.Println(reflect.TypeOf(ee))
// 	fmt.Println(ee)
// }

//	func main() {
//		Do()
//	}

// package main

// import "fmt"

// type Email = string
// type Email1 string

// func main() {
// 	// Using type alias
// 	var emailAlias Email = "john@example.com"
// 	fmt.Println(emailAlias)

// 	// Using new type
// 	var emailNewType Email1 = "jane@example.com"
// 	fmt.Println(emailNewType)

// 	// Cannot assign string to Email1 directly, need explicit conversion
// 	var invalidAssignment Email1 = "bob@example.com" // This would cause a compilation error
// 	// var validAssignment Email1 = Email1("bob@example.com") // This is valid

// 	fmt.Println(invalidAssignment)
// }

package main

import "fmt"

type Duration int64

func main() {

	var dur Duration

	// dur = int64(1000)
	// dur = Duration(1000)
	dur = 1001
	// dur = Duration("1000")
	fmt.Println(dur)

}
