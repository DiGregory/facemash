package storage

import (
	"fmt"
	"math"
)

type Girl struct {
	VkId      string
	Id        int
	Rating    float64
	Wins      int
	Losses    int
	FirstName string
	LastName  string
}

func (s *Storage) GetGirlsFromDb() (*Girl, *Girl, error) {

	var FirstGirl = new(Girl)
	var SecondGirl = new(Girl)

	FirstRecord := s.DB.QueryRow("SELECT * FROM Girls ORDER BY RANDOM() LIMIT 1")
	err := FirstRecord.Scan(&FirstGirl.Id, &FirstGirl.VkId, &FirstGirl.Rating, &FirstGirl.Wins, &FirstGirl.Losses, &FirstGirl.FirstName, &FirstGirl.LastName)

	if err != nil {
		return nil, nil, err
	}

	SecondRecord := s.DB.QueryRow("SELECT * FROM Girls ORDER BY RANDOM() LIMIT 1")
	err = SecondRecord.Scan(&SecondGirl.Id, &SecondGirl.VkId, &SecondGirl.Rating, &SecondGirl.Wins, &SecondGirl.Losses, &SecondGirl.FirstName, &SecondGirl.LastName)

	if err != nil {
		return nil, nil, err
	}

	if SecondGirl.VkId == FirstGirl.VkId {

		SecondRecord := s.DB.QueryRow("SELECT * FROM Girls ORDER BY RANDOM() LIMIT 1")

		err = SecondRecord.Scan(&SecondGirl.Id, &SecondGirl.VkId, &SecondGirl.Rating, &SecondGirl.Wins, &SecondGirl.Losses, &SecondGirl.FirstName, &SecondGirl.LastName)

		if err != nil {
			return nil, nil, err
		}
	}

	return FirstGirl, SecondGirl, nil
}

func (s *Storage) GetTop() ([]*Girl, error) {

	girls := make([]*Girl, 0)

	girlsRecords, err := s.DB.Query("SELECT * FROM Girls ORDER BY Rating DESC LIMIT 500")

	if err != nil {
		return nil, err
	}
	defer girlsRecords.Close()

	for girlsRecords.Next() {
		OneGirl := new(Girl)
		err := girlsRecords.Scan(&OneGirl.Id, &OneGirl.VkId, &OneGirl.Rating, &OneGirl.Wins, &OneGirl.Losses, &OneGirl.FirstName, &OneGirl.LastName)
		if err != nil {
			return nil, err
		}
		girls = append(girls, OneGirl)

	}

	return girls, nil
}
func (s *Storage) GetGirlById(VkId string) (*Girl, error) {
	G := new(Girl)
	Record := s.DB.QueryRow("SELECT * FROM Girls WHERE VkId=$1", VkId)
	err := Record.Scan(&G.Id, &G.VkId, &G.Rating, &G.Wins, &G.Losses, &G.FirstName, &G.LastName)
	if err != nil {
		return nil, err
	}
	return G, nil
}
func (s *Storage) UpdateGirls(WinId, LooseId string) (err error) {

	WinGirl, err := s.GetGirlById(WinId)

	if err != nil {
		return err
	}

	LooseGirl, err := s.GetGirlById(LooseId)

	if err != nil {
		return err
	}

	E := 1 / (1 + math.Pow(10, (WinGirl.Rating-LooseGirl.Rating)/400))
	NewWinRating := WinGirl.Rating + 10*(1-E)
	NewLooseRating := LooseGirl.Rating - 10*(1-E)

	_, err = s.DB.Exec("UPDATE Girls SET Rating=$1, wins=wins+1 WHERE VkId=$2", NewWinRating, WinGirl.VkId)

	if err != nil {
		return err
	}

	_, err = s.DB.Exec("UPDATE Girls SET Rating=$1, losses=losses+1 WHERE VkId=$2", NewLooseRating, LooseGirl.VkId)
	fmt.Println("girls updated")
	if err != nil {
		return err
	}

	return nil

}
