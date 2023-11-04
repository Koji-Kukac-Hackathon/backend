package product

import "fmt"

type controller struct{}

var Controller controller = controller{}

func (controller) Test() {
	fmt.Println("controller test")
}
