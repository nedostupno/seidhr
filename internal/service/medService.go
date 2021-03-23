package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/nedostupno/seidhr/internal/model"
	"github.com/nedostupno/seidhr/internal/repository"
)

type MedService struct {
	repo repository.Medicaments
}

func NewMedService(repo repository.Medicaments) *MedService {
	return &MedService{repo: repo}
}

func (s *MedService) IsExist(medName string) (bool, error) {
	return s.repo.IsExist(medName)
}

func (s *MedService) GetTrueName(medName string) (string, error) {
	return s.repo.GetTrueName(medName)
}

func (s *MedService) GetID(medTitle string) (int, error) {
	return s.repo.GetID(medTitle)
}

func (s *MedService) GetTitle(medicamentID int) (string, error) {
	return s.repo.GetTitle(medicamentID)
}

func (s *MedService) InitMedList(medLines []string) error {
	return s.repo.InitMedList(medLines)
}

func (s *MedService) IsMedListExist() (bool, error) {
	return s.repo.IsMedListExist()
}

func (s *MedService) AreTheAnySubscriptions() (bool, error) {
	return s.repo.AreTheAnySubscriptions()
}

func (s *MedService) GetAllMedicamentsWithSub() ([]int, error) {
	return s.repo.GetAllMedicamentsWithSub()
}

func (s *MedService) GetAvailability(medicamentID int) (bool, error) {
	return s.repo.GetAvailability(medicamentID)
}

func (s *MedService) ChangeAvailability(medicamentID int, value bool) error {
	return s.repo.ChangeAvailability(medicamentID, value)
}

func (s *MedService) GetSubscribers(medicamentID int) ([]int, error) {
	return s.repo.GetSubscribers(medicamentID)
}

// ReadFileWithMeds - считывает данные из файла drugs.txt и подготавливает их для
// передачи в функцию InitMedList, котрая заполнит ими базу данных
func (s *MedService) ReadFileWithMeds() ([]string, error) {
	file, err := os.Open("drugs.txt")
	if err != nil {
		err := fmt.Errorf("Ошибка открытия файла drugs.txt")
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

// ReqMedInfo - опрашивает сторонний сервис о наличии лекарства,
// анмаршалит полученный json в структуру Jsn и возвращает ее
func (s *MedService) ReqMedInfo(medTitle string) (model.Jsn, error) {
	hh := url.QueryEscape(medTitle)

	client := &http.Client{}
	req, err := http.NewRequest(
		"GET", "https://eservice.gu.spb.ru/portalFront/proxy/async?filter="+hh+"&operation=getMedicament", nil,
	)
	// добавляем заголовки
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:74.0) Gecko/20100101 Firefox/74.0")

	resp, err := client.Do(req)
	if err != nil {
		return model.Jsn{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := fmt.Errorf("код ответа: %v", resp.StatusCode)
		return model.Jsn{}, err
	}

	j := model.Jsn{}
	data, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(data, &j)

	return j, nil
}

// ParseJSON - парсит json структуру и возвращает готовый текст сообщения для пользователя
func (s *MedService) ParseJSON(j model.Jsn) string {
	var text []string

	title := fmt.Sprintln("Название: ", j.Model.Result[0].Name)

	for name := range j.Model.Result[0].Districts {
		district := fmt.Sprint("\n\n[[", j.Model.Result[0].Districts[name].Name, " ]]\n\n")
		text = append(text, district)

		for _, apothecary := range j.Model.Result[0].Districts[name].Apothecaries {
			name := apothecary.Name
			addr := apothecary.Addr
			a := strings.Trim(addr, "  * На момент обращения в аптеку не гарантируется наличие лекарственного препарата к выдаче, в связи с ограничением количества препарата в аптеке. Информацию о наличии препарата необходимо уточнить по телефону")

			s := strings.Split(a, ",")

			// index := fmt.Sprintln("Индекс: ", s[0])
			street := strings.TrimPrefix(s[2], " ")
			house := s[3]
			address := street + " " + house

			fedExemption := fmt.Sprintln("Федеральная льгота: ", apothecary.FedExemption)
			//Региональная льгота
			regExemption := fmt.Sprintln("Региональнальная льгота: ", apothecary.RegExemption)
			//Писхиатрическая льгота
			psyExemption := fmt.Sprintln("Психиатрическая льгота: ", apothecary.PsyExemption)
			//ВЗН
			vzn := fmt.Sprintln("ВЗН: ", apothecary.VZN)

			apoth := name + "\n" + address + "\n\n" + fedExemption + regExemption + psyExemption + vzn

			text = append(text, apoth)
		}
	}

	msg := title
	for _, i := range text {
		msg += i
	}

	return msg
}

// IsErrExistInJSON - Проверяет json структуру на наличие ошибок
func (s *MedService) IsErrExistInJSON(j model.Jsn) bool {
	if u := j.Errors; u != "" {
		return true
	}
	return false
}
