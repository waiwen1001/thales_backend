package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"waiwen.com/thales-backend/utils"
)

type Product struct {
	Id          int       `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Type        string    `json:"type" binding:"required"`
	PictureUrl  string    `json:"picture_url"`
	Price       float64   `json:"price" binding:"required"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductRequest struct {
	Search    string
	SortKey   string // name
	SortOrder string // asc / desc
	Page      int
	PageSize  int
}

func (s *Product) GetAllProducts(req ProductRequest, ctx context.Context, db *sql.DB) (data []Product, paginate *utils.Paginate, err error) {
	// get total products count
	var totalCount int
	countQuery := `SELECT COUNT(*) from products`
	err = db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, nil, err
	}

	// get products
	if req.SortKey == "" {
		req.SortKey = "id"
	}

	sortOrder := "ASC"
	if req.SortOrder == "desc" {
		sortOrder = `DESC`
	}

	query := "SELECT * FROM products"
	if req.Search != "" {
		query += " WHERE name ILIKE $1"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", req.SortKey, sortOrder)
	offset := (req.Page - 1) * req.PageSize

	var rows *sql.Rows
	if req.Search != "" {
		query += " LIMIT $2 OFFSET $3"
		rows, err = db.QueryContext(ctx, query, "%"+req.Search+"%", req.PageSize, offset)
	} else {
		query += " LIMIT $1 OFFSET $2"
		rows, err = db.QueryContext(ctx, query, req.PageSize, offset)
	}

	if err != nil {
		log.Println("Error querying products:", err)
		return nil, nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(
			&p.Id, &p.Name, &p.Type, &p.PictureUrl, &p.Price,
			&p.Description, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("rows iteration error: %w", err)
	}

	paginate = &utils.Paginate{
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalCount: totalCount,
	}

	return products, paginate, nil
}

func (s *Product) CreateProduct(ctx context.Context, db *sql.DB) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error starting transaction:", err)
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO products (name, type, picture_url, price, description) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	err = tx.QueryRowContext(ctx, query, s.Name, s.Type, s.PictureUrl, s.Price, s.Description).Scan(&s.Id, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		log.Println("Error inserting product:", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		return err
	}
	return nil
}

func (s *Product) CheckProductExist(ctx context.Context, db *sql.DB) (err error) {
	query := `SELECT EXISTS(SELECT 1 from products WHERE id = $1)`
	var exists bool
	err = db.QueryRowContext(ctx, query, s.Id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking product existence: %v", err)
	}
	if !exists {
		return fmt.Errorf("product %d does not exist", s.Id)
	}

	return nil
}

func (s *Product) UpdateProduct(ctx context.Context, db *sql.DB) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error starting transaction:", err)
		return err
	}
	defer tx.Rollback()

	s.UpdatedAt = time.Now()

	query := `UPDATE products SET name = $1, type = $2, picture_url = $3, price = $4, description = $5, updated_at = $6 WHERE id = $7 RETURNING created_at`
	err = tx.QueryRowContext(ctx, query, s.Name, s.Type, s.PictureUrl, s.Price, s.Description, s.UpdatedAt, s.Id).Scan(&s.CreatedAt)
	if err != nil {
		log.Println("Error updating product:", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		return err
	}

	return nil
}

func (s *Product) GetProductById(ctx context.Context, db *sql.DB) (err error) {
	query := "SELECT id, name, type, picture_url, price, description, created_at, updated_at FROM products WHERE id = $1"
	row := db.QueryRowContext(ctx, query, s.Id)

	if err := row.Scan(&s.Id, &s.Name, &s.Type, &s.PictureUrl, &s.Price, &s.Description, &s.CreatedAt, &s.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no product found with id %d", s.Id)
		}
		return fmt.Errorf("failed to fetch product: %w", err)
	}

	return nil
}

func (s *Product) DeleteProduct(ctx context.Context, db *sql.DB) (err error) {
	query := "DELETE FROM products WHERE id = $1"
	result, err := db.ExecContext(ctx, query, s.Id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no product found with id %d", s.Id)
	}

	return nil
}
