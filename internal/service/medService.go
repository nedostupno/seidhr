package service

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nedostupno/seidhr/internal/repository"
)

type MedService struct {
	repo repository.Medicaments
}

func NewMedService(repo repository.Medicaments) *MedService {
	return &MedService{repo: repo}
}

func (s *MedService) InitMedList(medLines []string) error {
	return s.repo.InitMedList(medLines)
}

func (s *MedService) IsMedListExist() (bool, error) {
	return s.repo.IsMedListExist()
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
