package constant

const (
	// Connection String
	PostgresDSN = "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta"

	// Database Pooling
	DBMaxIdleConns = "10"
	DBMaxOpenConns = "100"
	DBConnMaxLifetime = "60m"
	DBConnMaxIdleTime = "15m"
)