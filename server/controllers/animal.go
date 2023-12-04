package controllers

import (
	"encoding/csv"
	"github.com/gofiber/fiber/v2"
	"github.com/raissarib/crud-houses-of-got/models"
	"github.com/raissarib/crud-houses-of-got/server/repositories"
	"net/http"
	"strconv"
	"time"
)

func NewAnimal() *AnimalController {
	return &AnimalController{Animal: repositories.NewAnimal()}
}

type AnimalController struct {
	Animal *repositories.Animal
}

func (s *AnimalController) GetOne(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON("Invalid ID")
	}

	Animal, err := s.Animal.Get(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(Animal)
}

func (s *AnimalController) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON("Invalid ID")
	}

	err = s.Animal.Delete(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	return c.SendStatus(http.StatusOK)
}

func (s *AnimalController) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON("Invalid ID")
	}

	var newAnimal models.Animal
	if err = c.BodyParser(&newAnimal); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	err = s.Animal.Update(id, &newAnimal)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(newAnimal)
}

func (s *AnimalController) Create(c *fiber.Ctx) error {
	var Animal models.Animal
	if err := c.BodyParser(&Animal); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	err := s.Animal.Create(&Animal)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	return c.Status(http.StatusCreated).JSON(Animal)
}

func (s *AnimalController) BulkCreate(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	fileHandle, err := file.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}
	defer fileHandle.Close()

	reader := csv.NewReader(fileHandle)
	records, err := reader.ReadAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	if len(records) < 2 {
		return c.Status(http.StatusInternalServerError).JSON("csv invalido")
	}

	var animais []models.Animal
	for _, record := range records[1:] {
		id, _ := strconv.Atoi(record[0])
		donoID, _ := strconv.Atoi(record[3])
		animal := models.Animal{
			ID:      id,
			Especie: record[1],
			Nome:    record[2],
			DonoID:  donoID,
		}
		animais = append(animais, animal)
	}

	if err := s.Animal.BulkCreate(animais); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	return c.SendStatus(http.StatusCreated)
}

func (s *AnimalController) GetCSV(c *fiber.Ctx) error {
	animais, err := s.Animal.GetMany()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	c.Set(fiber.HeaderContentType, "text/csv")
	c.Set("Content-Disposition", "attachment; filename=\"animais.csv\"")

	writer := csv.NewWriter(c.Response().BodyWriter())
	defer writer.Flush()

	// Escrever cabeçalho do CSV
	if err := writer.Write([]string{"ID", "Especie", "Nome", "DonoID", "DonoNome", "DonoDataNascimento"}); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	// Escrever dados
	for _, animal := range animais {
		record := []string{
			strconv.Itoa(animal.ID),
			animal.Especie,
			animal.Nome,
			strconv.Itoa(animal.DonoID),
			animal.Dono.Nome,
			animal.Dono.DataNascimento.Format(time.RFC3339),
		}
		if err := writer.Write(record); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(err)
		}
	}

	return nil
}
