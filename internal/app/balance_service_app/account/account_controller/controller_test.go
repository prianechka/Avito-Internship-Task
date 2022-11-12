package account_controller

import (
	"Avito-Internship-Task/internal/app/balance_service_app/account/account_repo"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"testing"
)

// TestCheckAccountIsExist проверяет, что если аккаунт существует, то manager вернёт True
func TestCheckAccountIsExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var checkUserID int64 = 1

	rows := sqlmock.NewRows([]string{"amount"})
	expect := 293.46
	expectResult := true
	rows = rows.AddRow(expect)

	mock.
		ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(checkUserID).
		WillReturnRows(rows)

	repo := account_repo.NewAccountRepo(db)
	manager := CreateNewAccountManager(repo)

	result, execErr := manager.CheckAccountIsExist(checkUserID)

	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(result, expectResult) {
		t.Errorf("results not match, want %v, have %v", result, expectResult)
		return
	}
}

// TestCheckAccountIsNotExist проверяет, что если аккаунт не существует, то manager вернёт False
func TestCheckAccountIsNotExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var checkUserID int64 = 1
	expectResult := false

	rows := sqlmock.NewRows([]string{"amount"})

	mock.
		ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(checkUserID).
		WillReturnRows(rows)

	repo := account_repo.NewAccountRepo(db)
	manager := CreateNewAccountManager(repo)

	result, execErr := manager.CheckAccountIsExist(checkUserID)

	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(result, expectResult) {
		t.Errorf("results not match, want %v, have %v", result, expectResult)
		return
	}
}

// TestCreateNewAccount проверяет, что если аккаунт не существует, то manager создаст новый аккаунт
func TestCreateNewAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var checkUserID int64 = 1

	rows := sqlmock.NewRows([]string{"amount"})

	mock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(checkUserID).
		WillReturnRows(rows)

	mock.ExpectExec("INSERT INTO balanceApp.accounts").
		WithArgs(checkUserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := account_repo.NewAccountRepo(db)
	manager := CreateNewAccountManager(repo)

	execErr := manager.CreateNewAccount(checkUserID)

	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}

// TestCreateNewAccountWithErrorItExists проверяет, что если аккаунт существует, то manager вернет ошибку
func TestCreateNewAccountWithErrorItExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var checkUserID int64 = 1

	rows := sqlmock.NewRows([]string{"amount"})
	expect := 293.46
	rows = rows.AddRow(expect)

	mock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(checkUserID).
		WillReturnRows(rows)

	repo := account_repo.NewAccountRepo(db)
	manager := CreateNewAccountManager(repo)

	execErr := manager.CreateNewAccount(checkUserID)

	if execErr != AccountIsExistErr {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}
}

// TestCheckBalance проверяет, что менеджер корректно проверяет баланс
func TestCheckBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var checkUserID int64 = 1

	rows := sqlmock.NewRows([]string{"amount"})
	expect := 293.46
	rows = rows.AddRow(expect)

	mock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(checkUserID).
		WillReturnRows(rows)

	repo := account_repo.NewAccountRepo(db)
	manager := CreateNewAccountManager(repo)

	amount, execErr := manager.CheckBalance(checkUserID)

	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(amount, expect) {
		t.Errorf("results not match, want %v, have %v", amount, expect)
		return
	}
}

// TestCheckCanAbleToBuy проверяет, что если хватает денег для оплаты, менеджер вернет true
func TestCheckCanAbleToBuy(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var checkUserID int64 = 1
	var sum float64 = 250
	var expectCanBuy = true

	rows := sqlmock.NewRows([]string{"amount"})
	expect := 293.46
	rows = rows.AddRow(expect)

	mock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(checkUserID).
		WillReturnRows(rows)

	repo := account_repo.NewAccountRepo(db)
	manager := CreateNewAccountManager(repo)

	canBuy, execErr := manager.CheckAbleToBuyService(checkUserID, sum)

	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(canBuy, expectCanBuy) {
		t.Errorf("results not match, want %v, have %v", canBuy, expectCanBuy)
		return
	}
}

// TestCheckCanNotAbleToBuy проверяет, что если не хватает денег для оплаты, менеджер вернет false
func TestCheckCanNotAbleToBuy(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var checkUserID int64 = 1
	var sum float64 = 350
	var expectCanBuy = false

	rows := sqlmock.NewRows([]string{"amount"})
	expect := 293.46
	rows = rows.AddRow(expect)

	mock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(checkUserID).
		WillReturnRows(rows)

	repo := account_repo.NewAccountRepo(db)
	manager := CreateNewAccountManager(repo)

	canBuy, execErr := manager.CheckAbleToBuyService(checkUserID, sum)

	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}
	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(canBuy, expectCanBuy) {
		t.Errorf("results not match, want %v, have %v", canBuy, expectCanBuy)
		return
	}
}

// TestUpdateMoney проверяет корректное пополнение баланса
func TestUpdateMoney(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var checkUserID int64 = 1
	var sum float64 = 100

	mock.ExpectExec("UPDATE balanceApp.accounts SET amount = amoumt +").
		WithArgs(sum, checkUserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := account_repo.NewAccountRepo(db)
	manager := CreateNewAccountManager(repo)

	execErr := manager.DonateMoney(checkUserID, sum)

	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	// Проверка, что баланс действительно обновился
	rows := sqlmock.NewRows([]string{"amount"})
	expect := 393.46
	rows = rows.AddRow(expect)

	mock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(checkUserID).
		WillReturnRows(rows)

	balance, secExecError := manager.CheckBalance(checkUserID)

	if secExecError != nil {
		t.Errorf("unexpected err: %v", secExecError)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(balance, expect) {
		t.Errorf("results not match, want %v, have %v", balance, expect)
		return
	}
}

// TestSpendMoneySuccess проверяет корректную оплату услуги, если хватает на неё денег
func TestSpendMoneySuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var checkUserID int64 = 1
	var sum float64 = 100

	rows := sqlmock.NewRows([]string{"amount"})
	expect := 293.46
	rows = rows.AddRow(expect)

	mock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(checkUserID).
		WillReturnRows(rows)

	mock.ExpectExec("UPDATE balanceApp.accounts SET amount = amoumt +").
		WithArgs(-sum, checkUserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := account_repo.NewAccountRepo(db)
	manager := CreateNewAccountManager(repo)

	execErr := manager.SpendMoney(checkUserID, sum)

	if execErr != nil {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	// Проверка на то, что баланс действительно изменился

	rows = sqlmock.NewRows([]string{"amount"})
	expect = 193.46
	rows = rows.AddRow(expect)

	mock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(checkUserID).
		WillReturnRows(rows)

	balance, secExecError := manager.CheckBalance(checkUserID)

	if secExecError != nil {
		t.Errorf("unexpected err: %v", secExecError)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(balance, expect) {
		t.Errorf("results not match, want %v, have %v", balance, expect)
		return
	}
}

// TestSpendMoneyWithNotEnoughMoney проверяет, что если денег не хватает, то оплата не произойдёт, и вернется ошибка
func TestSpendMoneyWithNotEnoughMoney(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var checkUserID int64 = 1
	var sum float64 = 300

	rows := sqlmock.NewRows([]string{"amount"})
	expect := 293.46
	rows = rows.AddRow(expect)

	mock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(checkUserID).
		WillReturnRows(rows)

	repo := account_repo.NewAccountRepo(db)
	manager := CreateNewAccountManager(repo)

	execErr := manager.SpendMoney(checkUserID, sum)

	if execErr != NotEnoughMoneyErr {
		t.Errorf("unexpected err: %v", execErr)
		return
	}

	// Проверка, что баланс не изменился
	rows = sqlmock.NewRows([]string{"amount"})
	expect = 293.46
	rows = rows.AddRow(expect)

	mock.ExpectQuery("SELECT amount FROM balanceApp.accounts WHERE userID").
		WithArgs(checkUserID).
		WillReturnRows(rows)

	balance, secExecError := manager.CheckBalance(checkUserID)

	if secExecError != nil {
		t.Errorf("unexpected err: %v", secExecError)
		return
	}

	if expectationErr := mock.ExpectationsWereMet(); expectationErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", expectationErr)
		return
	}

	if !reflect.DeepEqual(balance, expect) {
		t.Errorf("results not match, want %v, have %v", balance, expect)
		return
	}
}
