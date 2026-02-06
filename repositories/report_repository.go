package repositories

import (
	"context"
	"database/sql"
	"kasir-api/models"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetSalesSummary(
	ctx context.Context,
	startDate, endDate time.Time,
) (*models.Report, error) {

	var summary models.Report

	// total revenue & transaksi
	err := r.db.QueryRowContext(ctx, `
		SELECT 
			COALESCE(SUM(total_amount), 0) AS total_revenue,
			COUNT(*) AS total_transaksi
		FROM transactions
		WHERE created_at >= $1 AND created_at < $2
	`, startDate, endDate).Scan(
		&summary.TotalRevenue,
		&summary.TotalTransaksi,
	)
	if err != nil {
		return nil, err
	}

	// produk terlaris
	err = r.db.QueryRowContext(ctx, `
		SELECT p.name, SUM(td.quantity) AS qty
		FROM transaction_details td
		JOIN products p ON p.id = td.product_id
		JOIN transactions t ON t.id = td.transaction_id
		WHERE t.created_at >= $1 AND t.created_at < $2
		GROUP BY p.name
		ORDER BY qty DESC
		LIMIT 1
	`, startDate, endDate).Scan(
		&summary.ProdukTerlaris.Name,
		&summary.ProdukTerlaris.QtySold,
	)

	if err == sql.ErrNoRows {
		summary.ProdukTerlaris = models.BestSellerProduct{}
		return &summary, nil
	}

	if err != nil {
		return nil, err
	}

	return &summary, nil
}
