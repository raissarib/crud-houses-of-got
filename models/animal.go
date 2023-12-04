package models

type Animal struct {
	ID      int    `json:"id" db:"codigo_animal"`
	Especie string `json:"especie" db:"especie"`
	Nome    string `json:"nome" db:"nome_animal"`
	DonoID  int    `json:"dono_id" db:"codigo_pessoa"`
	Dono    Pessoa `json:"dono" db:"pessoa"`
}
