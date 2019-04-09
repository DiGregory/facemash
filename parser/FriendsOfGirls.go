package parser

import (
	"fmt"
	"github.com/himidori/golang-vk-api"
	"strconv"
	"os"
	"bufio"
	"time"
	"strings"
)

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

func getFriends(GirlId int, client *vkapi.VKClient, file1 *os.File, file2 *os.File) {

	_, users, err := client.FriendsGet(GirlId, peopleCount)

	if err != nil {
		fmt.Println("error was happened at friends getting: ", err)
	}

	for i := 0; i < len(users); i++ {
		if users[i].Sex == 1 {
			fmt.Println("Девушка: " + users[i].FirstName + " c id: " + strconv.Itoa(users[i].UID))

			file1.WriteString(strconv.Itoa(users[i].UID) + "\r\n")
			file2.WriteString(users[i].FirstName + " " + users[i].LastName + " " + strconv.Itoa(users[i].UID) + "\r\n")
		}
	}

}

func main() {
	IdFile, err := os.Create("FriendIds.txt")

	if err != nil {
		fmt.Println("can`t create file: ", err)
	}

	defer IdFile.Close()

	//файл имя+ид
	NamesFile, err := os.Create("FriendIdsNames.txt")

	if err != nil {
		fmt.Println("can`t create file: ", err)
	}

	defer NamesFile.Close()

	client, err := vkapi.NewVKClient(vkapi.PlatformWindows, User, Password)

	if err != nil {
		fmt.Println("error was happened at vk connect: ", err)
	}

	n := 1 //итерация

	for _, v := range publics {

		file, _ := os.Open(strconv.Itoa(v) + ".txt")
		f := bufio.NewReader(file)

		for {
			n++
			newId, _ := f.ReadString('\n')
			read_lines := strings.Split(newId, "\r")

			Id, err := strconv.Atoi(read_lines[0])
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(Id)
			if n%3 == 0 {
				time.Sleep(1100 * time.Millisecond)
			}
			getFriends(Id, client, IdFile, NamesFile)
		}
	}
}
