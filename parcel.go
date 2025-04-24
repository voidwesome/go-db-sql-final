package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec(`INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)`, p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}

	// верните идентификатор последней добавленной записи
	id, err := res.LastInsertId()

	return int(id), err
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row := s.db.QueryRow(`SELECT number, client, status, address, created_at FROM parcel WHERE number = ?`, number)

	// заполните объект Parcel данными из таблицы
	var p Parcel
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	return p, err
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query(`SELECT number, client, status, address, created_at FROM parcel WHERE client = ?`, client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		var p Parcel
		if err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec(`UPDATE parcel SET status = ? WHERE number = ?`, status, number)
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	parcel, err := s.Get(number)
	if err != nil {
		return fmt.Errorf("не удалось получить посылку: %w", err)
	}

	if parcel.Status != ParcelStatusRegistered {
		return errors.New("нельзя изменить адрес: посылка уже не зарегистрирована")
	}

	_, err = s.db.Exec(`UPDATE parcel SET address = ? WHERE number = ?`, address, number)
	return err
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	parcel, err := s.Get(number)
	if err != nil {
		return err
	}

	if parcel.Status != ParcelStatusRegistered {
		return errors.New("посылка уже отправлена или доставлена: нельзя удалить")
	}

	_, err = s.db.Exec(`DELETE FROM parcel WHERE number = ?`, number)
	return err
}
