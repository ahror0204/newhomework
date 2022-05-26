package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
)

type Costumer struct { // table  name users
	Costumer_id  string
	FirstName    string
	LastName     string
	UserName     string
	PhoneNumbers []Phone
	Adresses     []Adress
	Products     []Product
	Email        string
	Gender       string
	Birthday     string
	Password     string //should be hashed and validate password should 8 sybols
	Status       string
}

type Adress struct {
	ID          string
	Costumer_id string
	Country     string
	City        string
	District    string
	PostalCode  string
}

type Product struct {
	ID          string
	Costumer_id string
	Name        string
	Types       []Type
	Cost        int64
	OrderNumber int64
	Amount      int64
	Currency    string
	Rating      int64
}

type Type struct {
	ID   string
	Name string
}

type Phone struct {
	ID          string
	Costumer_id string
	Numbers     []int
	Code        string
}

func main() {
	//---------------------------CONNECT DB-------------------------------------
	connSql := `user = postgres password = 1 database = project1 sslmode = disable`
	db, err := sql.Open("postgres", connSql)
	if err != nil {
		panic("error while opening db")
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("ERROR IN THE BEGINNING TRANSACTION>>>:", err)
		return
	}


	fun := Costumer{}
	
	begin:
	fmt.Println("Select the acrion:\n1 - enter\n2 - update\n3 - delete\n4 - print")
	var choice int
	fmt.Scan(&choice)
	switch choice {
		case 1:
			costumer := Costumer {
				FirstName:"John",
				LastName: "Kamolov",
				UserName: "jorj",
				PhoneNumbers: []Phone {
					{
						Costumer_id: "",
						Numbers: []int{777777},
						Code: "+998",
					},
				},
				Adresses: []Adress{
					{
						Costumer_id: "",
						Country: "Uzbekistan",
						City: "Tashkent",
						District: "Chilanzar",
						PostalCode: "111711",
					},
				},
				Products: []Product {
					{
						Costumer_id: "",
						Name: "iPhone 15 Pro MAX",
						Types: []Type {
							{
								Name: "Mobile Phones",
							},
						},
						Cost: 1500,
						OrderNumber: 777,
						Amount: 1,
						Currency: "$",
						Rating: 1,
					},
				},
				Email: "sher@gmail.com",
				Gender: "male",
				Birthday: "2000-10-11",
				Password: "1324",
				Status: "single",
			  }
			fun.insertData(*tx, costumer) // ---------------insertData FUNCTION---------------

		case 2:

			fun.updateData(*tx)// -----------------------------getData FUNCTION-----------------
			fmt.Println("Changed")
		
		case 3:

			fun.deletetData(*tx) //------------------------deletetData FUNCTION---------------

		case 4:
			alldata, err := fun.getData(*tx, *db) // -----------------------------getData FUNCTION-----------------
			if err != nil {
				fmt.Println("alldata ERROR >>>:", err)
				tx.Rollback()
				return
			}
		
			data, err := json.MarshalIndent(alldata, " -", "\t")
			fmt.Println(string(data)) 
			
			if err != nil {
				fmt.Println("MARSHAL PROCESS ERROR >>>:", err)
				tx.Rollback()
				return
			}
		default: fmt.Println("Wrong diapason enter again!"); goto begin;
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("COMMIT PROCESS ERROR >>>:", err)
		tx.Rollback()
		return
	}

}


func (data Costumer) insertData(tx sql.Tx, cust Costumer) {
	var Cust_id string
	insertQuary := `INSERT INTO Costumer (firstname, lastname, username, email, gender, birthday, password, status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING Costumer_id;`
	err := tx.QueryRow(insertQuary, cust.FirstName, cust.LastName, cust.UserName, cust.Email, cust.Gender, cust.Birthday, cust.Password, cust.Status).Scan(&Cust_id)

	if err != nil {
		fmt.Println("ERROR WHILE INSERTING costumer TABLE:", err)
		tx.Rollback()
		return
	}
	var phone_id string
	for _, number := range cust.PhoneNumbers {
		query := "INSERT INTO Phone (costumer_id, numbers, code) VALUES ($1, $2, $3) RETURNING id"
		err = tx.QueryRow(query, Cust_id, pq.Array(number.Numbers), number.Code).Scan(&phone_id)
		if err != nil {
			fmt.Println("ERROR WHILE INSERTING TO phone TABLE:", err)
			tx.Rollback()
			return
		}
	}

	var adress_id string
	for _, val := range cust.Adresses {
		query := "INSERT INTO adress (costumer_id, country, city, district, postalcode) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		err = tx.QueryRow(query, Cust_id, val.Country, val.City, val.District, val.PostalCode).Scan(&adress_id)
		if err != nil {
			fmt.Println("ERROR WHILE INSERTING TO adress TABLE:", err)
			tx.Rollback()
			return
		}
	}

	var types_id string

	var product_id string
	for _, val := range cust.Products {

		for _, val1 := range val.Types {
			query1 := "INSERT INTO type (name, costumer_id) VALUES ($1, $2) RETURNING id"
			err = tx.QueryRow(query1, val1.Name, Cust_id).Scan(&types_id)
			if err != nil {
				fmt.Println("ERROR WHILE INSERTING TO adress TABLE:", err)
				tx.Rollback()
				return
			}
		}

		query := "INSERT INTO product (costumer_id, name, types_id, cost, ordernumber, amount, currency, rating) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
		err = tx.QueryRow(query, Cust_id, val.Name, types_id, val.Cost, val.OrderNumber, val.Amount, val.Currency, val.Rating).Scan(&product_id)
		if err != nil {
			fmt.Println("ERROR WHILE INSERTING TO product TABLE:", err)
			tx.Rollback()
			return
		}
	}
}

func (data Costumer) getData(tx sql.Tx, db sql.DB ) (Costumer, error) {

	var name string
	var costumerID []uint8
	fmt.Println("Enter name of the person who you want to select from tables:")
	fmt.Scan(&name)
	selesctQuery := `SELECT costumer_id FROM costumer WHERE firstname = $1 limit 1`
	err1 := tx.QueryRow(selesctQuery, name).Scan(&costumerID)
	if err1 != nil {
		fmt.Println("ENTER PERSON`S NAME PROCESS ERROR >>>",err1 )
	}
	var (
		cust     Costumer
		err      error
		adresses []Adress
		phones   []Phone
		products []Product
	)
	queryCostumer := `SELECT costumer_id, firstname, lastname, username, gender, email, password, status FROM costumer WHERE costumer_id = $1;`
	row := tx.QueryRow(queryCostumer, costumerID)
	if err != nil {fmt.Println("COSTUMER QUERY ERROR: ", err)}
	err = row.Scan(
		&cust.Costumer_id,
		&cust.FirstName,
		&cust.LastName,
		&cust.UserName,
		&cust.Gender,
		&cust.Email,
		&cust.Password,
		&cust.Status,
	)
	if err != nil {
		fmt.Println("ERROR WHILE SELECTING CUSTUMER DATA", err)
		tx.Rollback()
		return cust, err
	}

	adressQuery := `SELECT id, country, city, district, postalcode FROM adress WHERE costumer_id = $1;`
	rows, err := tx.Query(adressQuery, costumerID)
	if err != nil {fmt.Println("ADRESS QUERY ERROR: ", err)}
	defer rows.Close()
	for rows.Next() {
		var adress Adress
		err = rows.Scan(
			&adress.ID,
			&adress.Country,
			&adress.City,
			&adress.District,
			&adress.PostalCode,
		)
		if err != nil {
			fmt.Println("ERROR WHILE SELECTING ADRESS DATA", err)
			tx.Rollback()
			return cust, err
		}
		adresses = append(adresses, adress)
	}

	phoneQuery := `SELECT id, numbers, code FROM phone WHERE costumer_id = $1;`
	rows, err = db.Query(phoneQuery, costumerID)
	if err != nil {fmt.Println("PHONE QUERY ERROR: ", err)}
	defer rows.Close()
	for rows.Next() {
		var phone Phone
		err = rows.Scan(
			&phone.ID,
			pq.Array(&phone.Numbers),
			&phone.Code,
		)
		if err != nil {
			fmt.Println("ERROR WHILE SELECTING PHONE DATA", err)
			tx.Rollback()
			return cust, err
		}
		phones = append(phones, phone)
	}

	productQuery := `SELECT id, name, cost, ordernumber, amount, currency, rating FROM product WHERE costumer_id = $1;`
	rows, err = db.Query(productQuery, costumerID)
	
	if err != nil {
		fmt.Println("PRODUCT QUERY ERROR: ", err)
		return cust, err
	}
	
	defer rows.Close()
	
	for rows.Next() {
		var product Product
		var types    []Type

		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Cost,
			&product.OrderNumber,
			&product.Amount,
			&product.Currency,
			&product.Rating,
		)
		if err != nil {
			fmt.Println("ERROR WHILE SELECTING PRODUCT DATA", err)
			tx.Rollback()
			return cust, err
		}

		typeQuery := `SELECT id, name FROM type WHERE costumer_id = $1`
		
		rows, err = db.Query(typeQuery, costumerID)
		
		if err != nil {
			fmt.Println("TYPE QUERY ERROR: ", err)
			return cust, err
		}

		defer rows.Close()

		for rows.Next() {
			var type1 Type
			err = rows.Scan(
				&type1.ID,
				&type1.Name,
			)
		
			if err != nil {
				fmt.Println("ERROR WHILE SELECTING TYPE DATA", err)
				return cust, err
			}

			types = append(types, type1)
			fmt.Println("***********************",types)
		}
		product.Types = types
		products = append(products, product)
	}
	cust.Adresses = adresses
	cust.PhoneNumbers = phones
	cust.Products = products

	return cust, err
}

func (data Costumer) deletetData(tx sql.Tx) {
	var name string
	fmt.Println("Enter name of the person who you want to delete from tables:")
	fmt.Scan(&name)
	deleteQuery := `DELETE FROM costumer WHERE firstname = $1`
	_, err := tx.Exec(deleteQuery, name)
	if err != nil {
		fmt.Println("DELETE FUNCTION ERROR >>>", err)
		tx.Rollback()
		return
	}
}

func (data Costumer) updateData (tx sql.Tx) {
	var oldName, newName string

	fmt.Print("Enter the name: ")
	fmt.Scan((&oldName))
	fmt.Print("Enter new name: ")
	fmt.Scan((&newName))

	updateQuery := `UPDATE costumer SET firstname  = $1 WHERE firstname = $2`
	_, err := tx.Exec(updateQuery, newName, oldName)
	if err != nil {
		fmt.Println("UPDATE DATA ERROR >>>", err)
		tx.Rollback()
		return
	}
}