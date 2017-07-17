package addresses

import "database/sql"

type Address struct {
	ID      int     `json:"id"`
	Address string  `json:"address"`
	Lat     float32 `json:"lat"`
	Lng     float32 `json:"lng"`
}

func (addrs *Address) getAddress(db *sql.DB) error {
	return db.QueryRow("SELECT address, lat, lng FROM addresses WHERE id=$1", addrs.ID).Scan(&addrs.Address, &addrs.Lat, &addrs.Lng)
}

func (addrs *Address) updateAddress(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE addresses SET address=$1, lat=$2, lng=$3 WHERE id=$4",
			addrs.Address, addrs.Lat, addrs.Lng, addrs.ID)

	return err
}

func (addrs *Address) createAddress(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO addresses(address, lat, lng) VALUES($1, $2, $3) RETURNING id",
		addrs.Address, addrs.Lat, addrs.Lng).Scan(&addrs.ID)

	if err != nil {
		return err
	}

	return nil
}

type Addresses []Address

func GetAddresses(db *sql.DB, start, count int) (Addresses, error) {
	rows, err := db.Query(
		"SELECT id, address, lat, lng FROM addresses LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	addresses := Addresses{}

	for rows.Next() {
		var addrs Address
		if err := rows.Scan(&addrs.ID, &addrs.Address, &addrs.Lat, &addrs.Lng); err != nil {
			return nil, err
		}
		addresses = append(addresses, addrs)
	}

	return addresses, nil
}
