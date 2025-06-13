package repository

import (
	"book_rest_api/models"
	"database/sql"
)

type BookRepository interface {
	GetAllBooks() ([]*models.Book, error)
	GetBookByID(int) (*models.Book, error)
	Create(*models.Book) error
	Update(*models.Book) error
	Delete(int) error
}

type bookRepo struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) BookRepository {
	return &bookRepo{
		db: db,
	}
}

func (br *bookRepo) GetAllBooks() ([]*models.Book, error) {

	var books []*models.Book

	rows, err := br.db.Query(`
						SELECT id, title, author 
						FROM books`)
	if err != nil {
		return []*models.Book{}, err
	}
	defer rows.Close()

	for rows.Next() {
		book := models.Book{}
		if err := rows.Scan(&book.ID, &book.Title, &book.Author); err != nil {
			return []*models.Book{}, err
		}
		books = append(books, &book)
	}

	return books, nil
}

func (br *bookRepo) GetBookByID(id int) (*models.Book, error) {

	row := br.db.QueryRow(`
	SELECT id, title, author 
	FROM books 
	WHERE id=?`, id)

	book := models.Book{}

	err := row.Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (br *bookRepo) Create(book *models.Book) error {

	_, err := br.db.Exec(`
	INSERT INTO books (title,author) 
	VALUES (?, ?) 
	`, book.Title, book.Author)
	if err != nil {
		return err
	}

	return nil
}

func (br *bookRepo) Update(book *models.Book) error {
	_, err := br.db.Exec(`
	UPDATE books
	SET title = ?, author = ?
	WHERE id = ?
	`, book.Title, book.Author, book.ID)
	if err != nil {
		return err
	}

	return nil
}

func (br *bookRepo) Delete(id int) error {
	_, err := br.db.Exec(`
	DELETE
	FROM books
	where id = ?
	`, id)
	if err != nil {
		return err
	}
	return nil
}
