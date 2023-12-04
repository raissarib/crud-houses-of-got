package repositories

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	"github.com/raissarib/crud-houses-of-got/models"
	"github.com/raissarib/crud-houses-of-got/services"
	"time"
)

func NewAnimal() *Animal {
	return &Animal{
		db:      services.MySQL(),
		timeout: time.Second * 5,
	}
}

type Animal struct {
	db      *sqlx.DB
	timeout time.Duration
}

func (s *Animal) Create(animal *models.Animal) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `INSERT INTO animal (especie, nome_animal, codigo_pessoa) VALUES (?, ?, ?)`
	_, err = tx.ExecContext(ctx, query, animal.Especie, animal.Nome, animal.DonoID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Info("erro no roolback", rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return err
}

func (s *Animal) Get(id int) (*models.Animal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	var animal models.Animal
	query := `
	SELECT a.codigo_animal, a.especie, a.nome_animal, a.codigo_pessoa,
	       p.codigo_pessoa, p.nome_pessoa, p.data_nascimento
	FROM animal a
	INNER JOIN pessoa p ON a.codigo_pessoa = p.codigo_pessoa
	WHERE a.codigo_animal = ?
	`
	row := s.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&animal.ID, &animal.Especie, &animal.Nome, &animal.DonoID,
		&animal.Dono.ID, &animal.Dono.Nome, &animal.Dono.DataNascimento)

	return &animal, err
}

func (s *Animal) Update(id int, newAnimal *models.Animal) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	query := `UPDATE animal SET especie = ?, nome_animal = ?, codigo_pessoa = ? WHERE codigo_animal = ?`
	_, err := s.db.ExecContext(ctx, query, newAnimal.Especie, newAnimal.Nome, newAnimal.DonoID, id)

	return err
}

func (s *Animal) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	query := `DELETE FROM animal WHERE codigo_animal = ?`
	_, err := s.db.ExecContext(ctx, query, id)

	return err
}

func (s *Animal) BulkCreate(animais []models.Animal) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `INSERT INTO animal (especie, nome_animal, codigo_pessoa) VALUES (?, ?, ?)`
	for _, animal := range animais {
		_, err := tx.ExecContext(ctx, query, animal.Especie, animal.Nome, animal.DonoID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *Animal) GetMany() ([]models.Animal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	var animais []models.Animal
	query := `
    SELECT a.codigo_animal, a.especie, a.nome_animal, a.codigo_pessoa,
           p.codigo_pessoa, p.nome_pessoa, p.data_nascimento
    FROM animal a
    INNER JOIN pessoa p ON a.codigo_pessoa = p.codigo_pessoa
    `
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var animal models.Animal
		var dono models.Pessoa
		err := rows.Scan(&animal.ID, &animal.Especie, &animal.Nome, &animal.DonoID,
			&dono.ID, &dono.Nome, &dono.DataNascimento)
		if err != nil {
			return nil, err
		}
		animal.Dono = dono
		animais = append(animais, animal)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return animais, nil
}
