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

func (repo *AccountRepo) AddNewAccount(accountID int64) error {
	repo.mutex.Lock()
	query := MySQLAddNewAccount{}.GetString()
	_, err := repo.conn.Exec(query, accountID)
	repo.mutex.Unlock()
	return err
}

func (repo *AccountRepo) GetCurrentAmount(accountID int64) (amount float64, err error) {
	repo.mutex.Lock()
	query := MySQLGetCurrentAmount{}.GetString()
	row := repo.conn.QueryRow(query, accountID)
	err = row.Scan(&amount)
	repo.mutex.Unlock()

	if err == sql.ErrNoRows {
		err = AccountNotExist
	}
	
	return amount, err
}

func (repo *AccountRepo) ChangeAmount(accountID int64, delta float64) error {
	repo.mutex.Lock()
	query := MySQLChangeAmount{}.GetString()
	_, err := repo.conn.Exec(query, delta, accountID)
	repo.mutex.Unlock()

	if err == sql.ErrTxDone {
		err = AccountNotExist
	}

	return err
}

func (repo *AccountRepo) DeleteAccount(accountID int64) error {
	repo.mutex.Lock()
	query := MySQLDeleteAccount{}.GetString()
	_, err := repo.conn.Exec(query, accountID)
	repo.mutex.Unlock()

	if err == sql.ErrTxDone {
		err = AccountNotExist
	}

	return err
}
