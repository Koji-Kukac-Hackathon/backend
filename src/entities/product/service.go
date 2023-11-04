package product

import "fmt"

type service struct{}

var Service service = service{}

func (service) Test() {
	fmt.Println("service test")
}
