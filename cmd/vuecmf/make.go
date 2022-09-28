package main

import "fmt"

func Make(atype, name string)  {
	switch atype {
	case "app":
		fmt.Println("app make ...", name)
	case "controller":
		fmt.Println(" controller make ...", name)
	case "model":
		fmt.Println("model make ...", name)
	case "service":
		fmt.Println("service make ...", name)
	}
}
