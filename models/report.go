package models

type Report struct {
	TotalRevenue   int               `json:"total_revenue"`
	TotalTransaksi int               `json:"total_transaksi"`
	ProdukTerlaris BestSellerProduct `json:"produk_terlaris"`
}

type BestSellerProduct struct {
	Name    string `json:"name"`
	QtySold int    `json:"qty_sold"`
}
