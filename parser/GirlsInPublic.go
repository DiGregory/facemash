package parser

import (
	"fmt"
	"github.com/himidori/golang-vk-api"
	"strconv"
	"time"
	"os"
)

//ид пабликов, собраны руками
var publics = map[string]int{
	"тс":      55059952,
	"ачс":     60289186,
	"ввх":     132177380,
	"профком": 24234717,
	"дф":      16710735,
	"студком": 43890091,
	"фф":      1149,
}
var peopleCount = 999
//логин/пароль вконтакте
var User = ""
var Password = ""

func getGirls(public int, client *vkapi.VKClient) {
	//для каждого паблика создаем текстовик с ид девушек
	IdFile, err := os.Create(strconv.Itoa(public) + ".txt")

	if err != nil {
		fmt.Println("can`t create file: ", err)
	}

	defer IdFile.Close()

	//файлы имя+ид

	//NamesFile, err := os.Create(strconv.Itoa(public) + "Names.txt")
	//
	//if err != nil {
	//	fmt.Println("can`t create file: ", err)
	//}
	//
	//defer NamesFile.Close()

LOOP:
	for j := 1; j <= 16; j++ { //1000*15 человек в паблике максимум

		_, users, err := client.GroupGetMembers(public, peopleCount, peopleCount*j-1)

		if err != nil {
			fmt.Println("error was happened at group members getting: ", err)
		}
		if len(users) == 0 {
			break LOOP
		}
		for i := 0; i < len(users); i++ {
			if users[i].Sex == 1 {
				fmt.Println("Девушка: " + users[i].FirstName + " c id: " + strconv.Itoa(users[i].UID))
				//пишем в файлы данные
				IdFile.WriteString(strconv.Itoa(users[i].UID) + "\r\n")
				//NamesFile.WriteString(users[i].FirstName + " " + users[i].LastName + " " + strconv.Itoa(users[i].UID) + "\r\n")
			}
		}
		if j%3 == 0 {
			time.Sleep(1100 * time.Millisecond)
		}
	}
}

func main() {

	client, err := vkapi.NewVKClient(vkapi.PlatformWindows, User, Password)

	if err != nil {
		fmt.Println("error was happened at vk connect: ", err)
	}

	for _, v := range publics {
		getGirls(v, client)
	}
}
