package model

// Apothecary - в данной структуре содержится информация об аптеке
// и колличестве лекрств в ней
type Apothecary struct {
	Name         string `json:"name"`
	Addr         string `json:"address"`
	FedExemption string `json:"ost1"`
	RegExemption string `json:"ost2"`
	PsyExemption string `json:"ost3"`
	VZN          string `json:"ost4"`
	Date         string `json:"date"`
}

//District - в данной струткуре содержится информация о районе и массив с аптеками
type District struct {
	Name         string       `json:"name"`
	ID           string       `json:"id"`
	Apothecaries []Apothecary `json:"apothecaries"`
}

// Result - содержит информацию о районах, где есть нужное лекарство
type Result struct {
	Name      string     `json:"name"`
	Districts []District `json:"districts"`
}

// Model - промежуточная структура ответа
type Model struct {
	Result []Result `json:"result"`
}

// Jsn - основная структура ответа
type Jsn struct {
	Status string `json:"status"`
	Model  Model  `json:"model,omitempty"`
	Errors string `json:"errors,omitempty"`
}

// Medicoment - структура описывающая лекарство
type Medicoment struct {
	ID           int
	Title        string
	Avaliability bool
}
