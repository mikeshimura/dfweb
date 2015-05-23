package dfweb

import (
	"database/sql"
	"os"
)

func PgOpen(urlstring string, sslmode string) (*sql.DB, error) {
	if len(sslmode) > 0 {
		return sql.Open("postgres", urlstring+"?sslmode="+sslmode)
	} else {
		return sql.Open("postgres", urlstring)
	}
}
func PgOpenDatabaseUrl(sslmode string) (*sql.DB, error) {
	return PgOpen(os.Getenv("DATABASE_URL"), sslmode)
}
