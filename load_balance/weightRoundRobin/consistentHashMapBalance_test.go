package weightRoundRobin

import (
	"fmt"
	"testing"
)

func TestConsistentHashMapBalance(t *testing.T) {
	rb := NewConsistentHashBanlance(10, nil)
	rb.Add("127.0.0.1:2001")
	rb.Add("127.0.0.1:2002")
	rb.Add("127.0.0.1:2003")
	rb.Add("127.0.0.1:2004")
	rb.Add("127.0.0.1:2005")

	fmt.Println(rb.Get("http://127.0.0.1:2000/base/aaaaaaa"))
	fmt.Println(rb.Get("http://127.0.0.1:2000/base/bbbbbbb"))
	fmt.Println(rb.Get("http://127.0.0.1:2000/base/ccccccc"))
	fmt.Println(rb.Get("http://127.0.0.1:2000/base/ddddddd"))
	fmt.Println(rb.Get("http://127.0.0.1:2000/base/eeeeeee"))
	fmt.Println(rb.Get("http://127.0.0.1:2000/base/aaaaaaa"))

	fmt.Println(rb.Get("127.0.0.1:3000"))
	fmt.Println(rb.Get("127.0.0.1:3001"))
	fmt.Println(rb.Get("127.0.0.1:3002"))
	fmt.Println(rb.Get("127.0.0.1:3003"))
	fmt.Println(rb.Get("127.0.0.1:3002"))
}
