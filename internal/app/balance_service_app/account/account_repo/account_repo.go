package account_repo

import (
	"database/sql"
	"sync"
)

type AccountRepo struct {
	mutex sync.RWMutex
	conn  *sql.DB
}

func NewAccountRepo(conn *sql.DB) *AccountRepo {
	return &AccountRepo{conn: conn}
}

func (repo *AccountRepo) AddNewAccount(userID int) error {
	repo.mutex.Lock()
	query := MySQLAddNewAccount{}.GetString()
	_, err := repo.conn.Exec(query, userID)
	repo.mutex.Unlock()
	return err
}

func (repo *AccountRepo) GetCurrentAmount(userID int) (amount float64, err error) {
	repo.mutex.Lock()
	query := MySQLGetCurrentAmount{}.GetString()
	row := repo.conn.QueryRow(query, userID)
	err = row.Scan(&amount)
	repo.mutex.Unlock()

	if err == sql.ErrNoRows {
		err = AccountNotExist
	}

	return amount, err
}

func (repo *AccountRepo) ChangeAmount(accountID int, delta float64) error {
	repo.mutex.Lock()
	query := MySQLChangeAmount{}.GetString()
	_, err := repo.conn.Exec(query, delta, accountID)
	repo.mutex.Unlock()

	if err == sql.ErrTxDone {
		err = AccountNotExist
	}

	return err
}
