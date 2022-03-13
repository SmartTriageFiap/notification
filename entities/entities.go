package entities

type Event struct {
	UnityId    string     `json:"unity_id,omitempty"`
	EmployeeId string     `json:"employee_id,omitempty"`
	Patient    Patient    `json:"Patient,omitempty"`
	Questions  []Question `json:"questions,omitempty"`
}

type Patient struct {
	Cpf      string `json:"cpf,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

type Question struct {
	Id     string `json:"id,omitempty"`
	Text   string `json:"text,omitempty"`
	Answer int    `json:"answer,omitempty"`
	Type   int    `json:"type,omitempty"`
}

type Partner struct {
	Url string `json:"url,omitempty"`
}
