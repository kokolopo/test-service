package database

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"log"
	"os"
	"test-service/internal/config"

	"github.com/go-sql-driver/mysql" // Driver MySQL untuk Go
	// Pengganti dotenv
)

// getEnv adalah fungsi helper untuk mengambil environment variable
// dengan nilai default.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// --- Fungsi Baru yang Dipisahkan ---

// NewTiDBConnection membuat dan mengkonfigurasi koneksi ke TiDB.
// Fungsi ini mengembalikan handle database (*sql.DB) atau error jika gagal.
func NewTiDBConnection(cfg *config.Config) (*sql.DB, error) {
	// 1. Ambil konfigurasi dari environment variables
	host := cfg.DBHost
	port := cfg.DBPort
	user := cfg.DBUser
	password := cfg.DBPassword
	dbName := cfg.DBName

	tlsConfigName := "tidb-tls"
	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: false, // Sama dengan rejectUnauthorized: true
	}

	customCaPath := cfg.SSLKey
	log.Printf("TIDB_CA_PATH terdeteksi, membaca CA kustom dari: %s", customCaPath)

	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile(customCaPath)
	if err != nil {
		// Jangan gunakan log.Fatal di sini! Kembalikan error ke pemanggil.
		return nil, fmt.Errorf("gagal membaca file CA (%s): %w", customCaPath, err)
	}

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return nil, fmt.Errorf("gagal menambahkan sertifikat CA ke pool")
	}
	tlsConfig.RootCAs = rootCertPool

	// Daftarkan konfigurasi TLS ini ke driver mysql
	if err := mysql.RegisterTLSConfig(tlsConfigName, tlsConfig); err != nil {
		return nil, fmt.Errorf("gagal mendaftarkan TLS config: %w", err)
	}

	// 4. Buat DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=%s",
		user,
		password,
		host,
		port,
		dbName,
		tlsConfigName,
	)

	// 5. Buka koneksi (sql.Open tidak langsung terhubung)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal menyiapkan koneksi database: %w", err)
	}

	// Kembalikan handle database yang siap pakai
	return db, nil
}
