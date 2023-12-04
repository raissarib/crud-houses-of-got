package models

import (
	"time"
)

type Pessoa struct {
	ID             int       `json:"id" db:"codigo_pessoa"`
	Nome           string    `json:"nome" db:"nome_pessoa"`
	DataNascimento time.Time `json:"data_nascimento"`
}
