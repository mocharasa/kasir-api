package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll() ([]models.Product, error) {
	query := `SELECT p.id, p.category_id, c.name as category_name, p.name, p.price, p.stock 
			  FROM products p JOIN categories c 
			  ON p.category_id = c.id`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.CategoryID, &p.CategoryName, &p.Name, &p.Price, &p.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO products (category_id, name, price, stock) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(query, product.CategoryID, product.Name, product.Price, product.Stock).Scan(&product.ID)
	return err
}

// GetByID - ambil produk by ID
func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `SELECT p.id, p.category_id, c.name as category_name, p.name, p.price, p.stock 
			  FROM products p JOIN categories c 
			  ON p.category_id = c.id
			  WHERE p.id = $1`

	var p models.Product
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.CategoryID, &p.CategoryName, &p.Name, &p.Price, &p.Stock)
	if err == sql.ErrNoRows {
		return nil, errors.New("produk tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE products SET category_id = $1, name = $2, price = $3, stock = $4 WHERE id = $5"
	result, err := repo.db.Exec(query, product.CategoryID, product.Name, product.Price, product.Stock, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return err
}
