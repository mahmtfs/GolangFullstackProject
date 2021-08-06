package main

import (
	"MyRepository/model"
	"MyRepository/repository"
	"fmt"
)

var command string
var id int32
var run = true

func main(){
	u := model.User{}

	userRepository:=repository.NewUserFileRepository()
	for run {
		fmt.Println("Enter the command:")
		fmt.Scanf("%s", &command)
		switch command {
		case "create":
			storedUser, err := userRepository.Create(&u)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Printf("The user created: %v\n", storedUser)
			break
		case "get":
			fmt.Println("Enter the id:")
			fmt.Scanf("%d", &id)
			storedUser:= userRepository.GetByID(&u, &id)
			fmt.Printf("The user found: %v\n", storedUser)
			break
		case "list":
			userRepository.GetAll()
			break
		case "delete":
			fmt.Println("Enter the id:")
			fmt.Scanf("%d", &id)
			userRepository.Delete(int(id))
			break
		case "edit":
			fmt.Println("Enter the id:")
			fmt.Scanf("%d", &id)
			userRepository.Edit(int(id))
			break
		case "exit":
			run = false
			break
		}
	}
}