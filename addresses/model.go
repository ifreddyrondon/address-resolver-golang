package addresses

import (
	"database/sql"
	"log"
)

const (
	GetAddressQuery    = "SELECT address, lat, lng FROM addresses WHERE id=$1"
	UpdateAddressQuery = "UPDATE addresses SET address=$1, lat=$2, lng=$3 WHERE id=$4"
	DeleteAddressQuery = "DELETE FROM addresses WHERE id=$1"
	CreateAddressQuery = "INSERT INTO addresses(address, lat, lng) VALUES($1, $2, $3) RETURNING id"
	GetAddressesQuery  = "SELECT id, address, lat, lng FROM addresses LIMIT $1 OFFSET $2"
)

type Address struct {
	ID      int     `json:"id"`
	Address string  `json:"address"`
	Lat     float32 `json:"lat"`
	Lng     float32 `json:"lng"`
}

func (addrs *Address) getAddress(db *sql.DB) error {
	return db.QueryRow(GetAddressQuery, addrs.ID).Scan(&addrs.Address, &addrs.Lat, &addrs.Lng)
}

func (addrs *Address) updateAddress(db *sql.DB) error {
	_, err := db.Exec(UpdateAddressQuery, addrs.Address, addrs.Lat, addrs.Lng, addrs.ID)
	return err
}

func (addrs *Address) deleteAddress(db *sql.DB) error {
	_, err := db.Exec(DeleteAddressQuery, addrs.ID)
	return err
}

func (addrs *Address) createAddress(db *sql.DB) error {
	err := db.QueryRow(CreateAddressQuery, addrs.Address, addrs.Lat, addrs.Lng).Scan(&addrs.ID)
	if err != nil {
		return err
	}
	return nil
}

type Addresses []Address

func GetAddresses(db *sql.DB, start, count int) (Addresses, error) {
	rows, err := db.Query(GetAddressesQuery, count, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := Addresses{}
	for rows.Next() {
		var addrs Address
		if err := rows.Scan(&addrs.ID, &addrs.Address, &addrs.Lat, &addrs.Lng); err != nil {
			return nil, err
		}
		list = append(list, addrs)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return list, nil
}
