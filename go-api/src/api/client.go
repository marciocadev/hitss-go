package api

type Client struct {
	ID           string `json:"id"`
	Nome         string `json:"nome"`
	Sobrenome    string `json:"sobrenome"`
	Contato      string `json:"contato"`
	Endereco     string `json:"endereco"`
	DtNascimento string `json:"dtNascimento"`
	CPF          string `json:"cpf"`
}

// Same as Client but CPF is removed to better security data
type ClientDTO struct {
	ID           string `json:"id"`
	Nome         string `json:"nome"`
	Sobrenome    string `json:"sobrenome"`
	Contato      string `json:"contato"`
	Endereco     string `json:"endereco"`
	DtNascimento string `json:"dtNascimento"`
}
