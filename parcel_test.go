package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	number, err := store.Add(parcel)
	if err != nil {
		t.Errorf("Error adding: %v", err)
	}
	if number <= 0 {
		t.Errorf("Expected positive parcel number, got: %d", number)
	}

	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	parc, err := store.Get(number)
	if err != nil {
		t.Errorf("Error getting: %v", err)
	}
	if parc.Number != number {
		t.Errorf("Parcel number mismatch. expected: %d, got: %d", number, parc.Number)
	}

	if parc.Client != parcel.Client {
		t.Errorf("Client mismatch. expected: %d, got: %d", parcel.Client, parc.Client)
	}

	if parc.Status != parcel.Status {
		t.Errorf("Status mismatch. expected: %s, got: %s", parcel.Status, parc.Status)
	}

	if parc.Address != parcel.Address {
		t.Errorf("Address mismatch. expected: %s, got: %s", parcel.Address, parc.Address)
	}

	if parc.CreatedAt != parcel.CreatedAt {
		t.Errorf("CreatedAt mismatch. expected: %s, got: %s", parcel.CreatedAt, parc.CreatedAt)
	}

	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что посылку больше нельзя получить из БД
	err = store.Delete(number)
	if err != nil {
		t.Errorf("Error deleting: %v", err)
	}

	_, err = store.Get(number)
	if err == nil {
		t.Errorf("Expected error when getting deleted parcel, got nil")
	}
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	number, err := store.Add(parcel)
	if err != nil {
		t.Errorf("Error adding: %v", err)
	}
	if number <= 0 {
		t.Errorf("Expected positive parcel number, got: %d", number)
	}

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(number, newAddress)
	if err != nil {
		t.Errorf("Error setting address: %v", err)
	}

	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	parc, err := store.Get(number)
	if err != nil {
		t.Errorf("Error getting: %v", err)
	}
	if parc.Address != newAddress {
		t.Errorf("Address mismatch. expected: %s, got: %s", newAddress, parc.Address)
	}
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	number, err := store.Add(parcel)
	if err != nil {
		t.Errorf("Error adding: %v", err)
	}
	if number <= 0 {
		t.Errorf("Expected positive parcel number, got: %d", number)
	}

	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	status := "registered"
	err = store.SetStatus(number, status)
	if err != nil {
		t.Errorf("Error setting status: %v", err)
	}

	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	parc, err := store.Get(number)
	if err != nil {
		t.Errorf("Error getting: %v", err)
	}
	if parc.Status != status {
		t.Errorf("Status mismatch. expected: %s, got: %s", status, parc.Status)
	}
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
		id, err := store.Add(parcels[i])
		if err != nil {
			t.Errorf("Error adding: %v", err)
		}
		if id <= 0 {
			t.Errorf("Expected positive parcel number, got: %d", id)
		}

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	// получите список посылок по идентификатору клиента, сохранённого в переменной client
	// убедитесь в отсутствии ошибки
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных
	storedParcels, err := store.GetByClient(client)
	if err != nil {
		t.Errorf("Error getting by client: %v", err)
	}
	if len(storedParcels) != len(parcels) {
		t.Errorf("Getting by client mismatch. expected: %d, got: %d", len(parcels), len(storedParcels))
	}

	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
		originalParcel, exists := parcelMap[parcel.Number]
		if !exists {
			t.Errorf("parcel with number %d not found in original map", parcel.Number)
			continue
		}

		// Проверяем значения полей
		if parcel.Client != originalParcel.Client {
			t.Errorf("parcel %d: client mismatch. expected %d, got %d",
				parcel.Number, originalParcel.Client, parcel.Client)
		}

		if parcel.Status != originalParcel.Status {
			t.Errorf("parcel %d: status mismatch. expected %s, got %s",
				parcel.Number, originalParcel.Status, parcel.Status)
		}

		if parcel.Address != originalParcel.Address {
			t.Errorf("parcel %d: address mismatch. expected %s, got %s",
				parcel.Number, originalParcel.Address, parcel.Address)
		}

		if parcel.Address != originalParcel.Address {
			t.Errorf("parcel %d: address mismatch. expected %s, got %s",
				parcel.Number, originalParcel.Address, parcel.Address)
		}

		// Проверяем, что это посылка нужного клиента
		if parcel.Client != client {
			t.Errorf("parcel %d: wrong client. expected %d, got %d",
				parcel.Number, client, parcel.Client)
		}
	}
}
