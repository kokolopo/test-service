package repository

import (
	"database/sql"
	"fmt"
	"test-service/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `INSERT INTO users (name, email, created_at) VALUES (?, ?, NOW())`
	_, err := r.db.Exec(query, user.Name, user.Email)
	if err != nil {
		// Jika terjadi error (misalnya, email duplikat), kembalikan error
		return err
	}
	return err
}

func (r *userRepository) FindAll() ([]domain.User, error) {
	// Query tanpa placeholder
	query := "SELECT id, name, email FROM users"

	// Gunakan Query() karena mengharapkan banyak baris
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("gagal menjalankan query: %w", err)
	}
	// PENTING: Pastikan rows ditutup setelah selesai, atau koneksi akan terblokir!
	defer rows.Close()

	// Inisialisasi slice untuk menampung semua hasil
	users := []domain.User{}

	// Iterasi (loop) melalui setiap baris hasil
	for rows.Next() {
		var user domain.User // Struct sementara untuk menampung data per baris

		// Scan data dari baris saat ini ke struct User.
		// Urutan field harus SAMA dengan urutan kolom di query SELECT.
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, fmt.Errorf("gagal scan baris saat iterasi: %w", err)
		}

		// Tambahkan user yang sudah di-scan ke slice users
		users = append(users, user)
	}

	// Setelah loop selesai, cek apakah ada error yang terjadi selama iterasi
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error setelah iterasi rows: %w", err)
	}

	// Kembalikan slice users
	return users, nil
}

func (r *userRepository) FindByID(id int64) (*domain.User, error) {
	// Query dengan placeholder (?)
	query := "SELECT id, name, email FROM users WHERE id = ?"

	// Gunakan QueryRow() karena kita hanya mengharapkan satu baris
	row := r.db.QueryRow(query, id)

	// Inisialisasi struct User untuk menampung hasil
	var user domain.User

	// Scan hasil dari row ke field struct User.
	// Urutan field di sini harus SAMA dengan urutan kolom di query SELECT.
	err := row.Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			// Kasus khusus: Tidak ada baris yang ditemukan
			return nil, fmt.Errorf("user dengan ID %d tidak ditemukan", id)
		}
		// Error database lainnya (koneksi terputus, query salah, dll.)
		return nil, fmt.Errorf("gagal scan data user: %w", err)
	}

	// Sukses, kembalikan pointer ke user
	return &user, nil
}
