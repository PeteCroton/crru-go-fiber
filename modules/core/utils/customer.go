package module

var Name = "Customer"

type Customer struct {
	Name     string
	Age      int
	Email    string
	Password string
}

func (c Customer) Hello() string {
	return "Hello " + c.Name
}

func SumCustomer(a, b int) string {
	x := Customer{
		Name:     "test",
		Age:      10,
		Email:    "test",
		Password: "test"}

	return x.Hello()
}
