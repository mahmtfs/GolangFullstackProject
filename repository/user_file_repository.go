package repository

import (
	"MyRepository/helpers"
	"MyRepository/model"
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

const NAMERIGHTBORDER = 8
const EMAILRIGHTBORDER = 9
const PASSWORDRIGHTBORDER = 12
var name string
var email string
var password string

type UserRepositoryI interface{
	Create(u *model.User) (*model.User, error)
	GetByID(user *model.User, id *int32) *model.User
	GetAll()
	Delete(id int)
	Edit(id int)
}

type UserFileRepository struct{
	idMutex *sync.Mutex
}

func NewUserFileRepository() *UserFileRepository{
	return &UserFileRepository{
		idMutex: &sync.Mutex{},
	}
}

func (ufr *UserFileRepository) Create(user *model.User) (*model.User, error) {
	user.ID = ufr.GetAndWriteID()
	ufr.idMutex.Lock()
	fmt.Println("Enter the name:")
	fmt.Scanf("%s", &name)
	fmt.Println("Enter the email:")
	fmt.Scanf("%s", &email)
	fmt.Println("Enter the password:")
	fmt.Scanf("%s", &password)
	ufr.idMutex.Unlock()
	user.Name = name
	user.Email = email
	user.Password = password
	err := helpers.CreateModel("users", user)
	if err!=nil{
		return nil, err
	}
	return user, nil
}

func (ufr *UserFileRepository) GetByID(user *model.User, id *int32) *model.User {
	userByID := ufr.FindByID(*user, *id)
	return userByID
}

func (ufr *UserFileRepository) GetAll() {
	file, err := os.Open("./datastore/users.txt")
	if err != nil{
		panic(err)
	}
	defer file.Close()
	scanner:=bufio.NewScanner(file)
	for scanner.Scan(){
		line := scanner.Text()
		fmt.Println(line)
	}
}

func (ufr *UserFileRepository) Delete(id int){
	var firstHalf string
	var strDelete string
	var secondHalf string
	idFromStore := ufr.GetIDForDelete()
	if int32(id) > idFromStore{
		panic("No such id")
	}
	file, err := os.OpenFile("./datastore/users.txt", os.O_RDONLY, 0600)
	if err != nil{
		panic(err)
	}
	scanner:=bufio.NewScanner(file)
	for i:=0; i < id; i++{
		scanner.Scan()
		line := scanner.Text()
		firstHalf += line + "\n"
	}
	for i := id; i < id + 1; i++ {
		scanner.Scan()
		line := scanner.Text()
		if line[:7] == "deleted"{
			panic("user already deleted")
		}
		line = "deleted" + line
		strDelete += line + "\n"
	}
	for i:= int32(id + 1); i < idFromStore + 1; i++{
		scanner.Scan()
		line := scanner.Text()
		secondHalf += line + "\n"
	}
	terr := os.Truncate("./datastore/users.txt", 0)
	if terr != nil {
		panic(terr)
	}
	file.Close()
	file, err = os.OpenFile("./datastore/users.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, ferr := file.WriteString(firstHalf + strDelete + secondHalf)
	if ferr != nil {
		panic(ferr)
	}
}

func (ufr *UserFileRepository) Edit(id int){
	var firstHalf string
	var strEdit string
	var secondHalf string
	idFromStore := ufr.GetIDForDelete()
	if int32(id) > idFromStore{
		panic("No such id")
	}
	file, err := os.OpenFile("./datastore/users.txt", os.O_RDONLY, 0600)
	if err != nil{
		panic(err)
	}
	scanner:=bufio.NewScanner(file)
	for i:=0; i < id; i++{
		scanner.Scan()
		line := scanner.Text()
		firstHalf += line + "\n"
	}
	for i := id; i < id + 1; i++ {
		scanner.Scan()
		line := scanner.Text()
		if line[:7] == "deleted"{
			panic("user is deleted, edit is impossible")
		}
		line = ""
		line += "{\"id\":" + strconv.Itoa(id)
		fmt.Println("Enter the name:")
		fmt.Scanf("%s", &name)
		line += ",\"name\":" + "\"" + name + "\""
		fmt.Println("Enter the email:")
		fmt.Scanf("%s", &email)
		line += ",\"email\":" + "\"" + email + "\""
		fmt.Println("Enter the password:")
		fmt.Scanf("%s", &password)
		line += ",\"password\":" + "\"" + password + "\"}"
		fmt.Println(line)
		strEdit += line + "\n"
	}
	for i:= int32(id + 1); i < idFromStore + 1; i++{
		scanner.Scan()
		line := scanner.Text()
		secondHalf += line + "\n"
	}
	terr := os.Truncate("./datastore/users.txt", 0)
	if terr != nil {
		panic(terr)
	}
	file.Close()
	file, err = os.OpenFile("./datastore/users.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, ferr := file.WriteString(firstHalf + strEdit + secondHalf)
	if ferr != nil {
		panic(ferr)
	}
}

func (ufr *UserFileRepository) GetIDForDelete() int32{
	var userID int32
	ufr.idMutex.Lock()
	file, ferr := os.OpenFile("./datastore/user_id_store.txt", os.O_WRONLY, 0600)
	if ferr != nil{
		panic(ferr)
	}
	line, ferr:=ioutil.ReadFile("./datastore/user_id_store.txt")
	if ferr != nil{
		panic(ferr)
	}
	userIDInt, ierr := strconv.Atoi(string(line))
	if ierr != nil {
		userIDInt = 0
	}
	userID = int32(userIDInt)
	ufr.idMutex.Unlock()
	defer file.Close()
	return userID
}

func (ufr *UserFileRepository) GetAndWriteID() int32 {
	var userID int32
	ufr.idMutex.Lock()
	file, ferr := os.OpenFile("./datastore/user_id_store.txt", os.O_WRONLY, 0600)
	if ferr != nil{
		panic(ferr)
	}
	line, ferr:=ioutil.ReadFile("./datastore/user_id_store.txt")

	if ferr != nil{
		panic(ferr)
	}
	userIDInt, ierr := strconv.Atoi(string(line))
	if ierr != nil {
		userIDInt = 0
	}
	userID = int32(userIDInt)
	str:= strconv.Itoa(userIDInt + 1)
	_, ferr = file.WriteString(str)
	if ferr!=nil{
		panic(ferr)
	}
	ufr.idMutex.Unlock()
	defer file.Close()
	return userID
}

func (ufr *UserFileRepository) FindByID(user model.User, id int32) (*model.User){
	numStr := strconv.Itoa(int(id))
	strToCompare := "{\"id\":" + numStr
	file, ferr := os.OpenFile("./datastore/users.txt", os.O_RDONLY, 0600)
	if ferr != nil{
		panic(ferr)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		if err := scanner.Err(); err != nil{
			panic(err)
		}
		line := scanner.Text()
		items := strings.Split(line, ",")
		if items[0] == strToCompare{
			user = model.User{
				id,
				items[1][NAMERIGHTBORDER:len(items[1]) - 1],
				items[2][EMAILRIGHTBORDER:len(items[2]) - 1],
				string(items[3][PASSWORDRIGHTBORDER:(len(items[3]) - 2)]),
			}
		}
	}
	return &user
}