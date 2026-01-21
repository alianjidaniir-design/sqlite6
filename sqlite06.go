package sqlite6

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Filename = ""
)

//Userdata is for holding full user data
//Userdata table + Username

type Userdata struct {
	ID          int
	Username    string
	Name        string
	Surname     string
	Description string
}

func openConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", Filename)
	if err != nil {
		return nil, err
	}
	return db, nil
}

//the function return the User ID of the username
// -1 if the user does not  exist

func exists(username string) int {
	username = strings.ToLower(username)
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userId := -1
	statement := fmt.Sprintf(`SELECT ID FROM Users Where Username = '%s'`, username)
	rows, err := db.Query(statement)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("exists() Scan", err)
			return -1
		}
		userId = id
	}
	return userId

}

//Adduser adds a new user to the database
// Return new User ID
// -1 if there was an error

func AddUser(d Userdata) int {
	d.Username = strings.ToLower(strings.TrimSpace(d.Username))
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()
	userId := exists(d.Username)
	if userId != -1 {
		fmt.Println("User already exists: ", d.Username)
		return -1
	}
	insertStatment := `INSERT INTO Users values (Null , ?)`
	_, err = db.Exec(insertStatment, d.Username)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	userId = exists(d.Username)
	if userId == -1 {
		return userId
	}

	insertStatment = `INSERT INTO Userdata values (?,?,?,?,?)`
	_, err = db.Exec(insertStatment, userId, d.Username, d.Username, d.Name, d.Description)
	if err != nil {
		fmt.Println("db.Exec()", err)
		return -1
	}
	return userId
}

//DeleteUser deletes an existing user

func DeleteUser(id int) error {
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	statement := fmt.Sprintf(`SELECT Username FROM Users WHERE ID = %d`, id)
	rows, err := db.Query(statement)
	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			fmt.Println("rows.Scan()", err)
		}
	}
	if exists(username) != -1 {
		return fmt.Errorf("User with ID %d does not exist", id)
	}
	//Delete from Userdata
	deletestatement := `DELETE from Userdata WHERE UserID = ?`
	_, err = db.Exec(deletestatement, id)
	if err != nil {
		return err
	}
	//Delete from Users
	deletestatement = `DELETE from users where ID = ?`
	_, err = db.Exec(deletestatement, id)
	if err != nil {
		return err
	}
	return nil

}

func ListUsers() ([]Userdata, error) {
	Data := []Userdata{}
	db, err := openConnection()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(`SELECT ID, Username, Name, Surname, Description FROM Users , Userdata WHERE Users.ID = Userdata.ID`)
	defer rows.Close()
	if err != nil {
		return Data, err
	}
	for rows.Next() {
		var id int
		var username string
		var name string
		var surname string
		var desc string
		err = rows.Scan(&id, &username, &name, &surname, &desc)
		temp := Userdata{ID: id, Username: username, Name: name, Surname: surname, Description: desc}
		Data = append(Data, temp)
		if err != nil {
			fmt.Println("rows.Scan()", err)
		}
	}

	return Data, nil

}

// UpdateUser is for updating an existing user
func UpdateUser(d Userdata) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	userId := exists(d.Username)
	if userId == -1 {
		return errors.New("User does not exist")
	}
	d.ID = userId
	updateStatement := `UPDATE Userdata set Name = ?, Surname = ?, Description = ? where ID = ?`
	_, err = db.Exec(updateStatement, d.Name, d.Surname, d.Description, d.ID)
	if err != nil {
		return err
	}

	return nil

}
