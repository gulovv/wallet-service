package db

import (
    "fmt"
    "database/sql"
    _ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(connString string) error {
    var err error
    db, err = sql.Open("postgres", connString)
    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    if err := db.Ping(); err != nil {
        return fmt.Errorf("could not ping database: %w", err)
    }

    createUserQuery := `
    DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'your_db_user') THEN
            CREATE USER your_db_user WITH PASSWORD 'your_db_password';
            GRANT ALL PRIVILEGES ON DATABASE your_db_name TO your_db_user;
        END IF;
    END;
    $$;
    `
    _, err = db.Exec(createUserQuery)
    if err != nil {
        return fmt.Errorf("could not create user: %w", err)
    }

    createTableQuery := `
    CREATE TABLE IF NOT EXISTS wallets (
        id SERIAL PRIMARY KEY,
        balance DECIMAL(10, 2) NOT NULL DEFAULT 0.00
    );
    `

    _, err = db.Exec(createTableQuery)
    if err != nil {
        return fmt.Errorf("could not create table: %w", err)
    }

    var count int
    err = db.QueryRow("SELECT COUNT(*) FROM wallets").Scan(&count)
    if err != nil {
        return fmt.Errorf("could not check table: %w", err)
    }

    if count == 0 {
        initialData := []struct {
            Balance float64
        }{
            {1000},
            {500},
            {1500},
            {2500},
            {3000},
        }

        for _, wallet := range initialData {
            _, err := db.Exec("INSERT INTO wallets (balance) VALUES ($1)", wallet.Balance)
            if err != nil {
                return fmt.Errorf("could not insert wallet with balance %.2f: %w", wallet.Balance, err)
            }
        }

        fmt.Println("Initial data inserted into wallets table.")
    } else {
        fmt.Println("Table 'wallets' already contains data.")
    }

    return nil
}

func GetBalance(walletID string) (float64, error) {
    var balance float64
    err := db.QueryRow("SELECT balance FROM wallets WHERE id = $1", walletID).Scan(&balance)
    if err != nil {
        return 0, fmt.Errorf("could not get balance: %w", err)
    }
    return balance, nil
}

func UpdateBalance(walletID string, amount float64) (float64, error) {
    tx, err := db.Begin()
    if err != nil {
        return 0, fmt.Errorf("could not begin transaction: %w", err)
    }

    _, err = tx.Exec("UPDATE wallets SET balance = balance + $1 WHERE id = $2", amount, walletID)
    if err != nil {
        tx.Rollback()
        return 0, fmt.Errorf("could not update balance: %w", err)
    }

    err = tx.Commit()
    if err != nil {
        return 0, fmt.Errorf("could not commit transaction: %w", err)
    }

    balance, err := GetBalance(walletID)
    if err != nil {
        return 0, err
    }
    return balance, nil
}