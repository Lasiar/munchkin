package model

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
)

type UserLogon struct {
	Id    int
	Login string
	Pass  string
}

func (db *Database) createTableAppLogon() error {
	_, err := db.DB.Exec(`
CREATE TABLE users.logon (
	id serial NOT NULL,
	login varchar NOT NULL,
	pass varchar NOT NULL,
	CONSTRAINT logon_pk PRIMARY KEY (id),
	CONSTRAINT logon_un_login UNIQUE (login)
);

-- Permissions

ALTER TABLE users.logon OWNER TO app;
GRANT ALL ON TABLE users.logon TO app;
`)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) GetAllLogon() []UserLogon {
	rows, err := d.DB.Query("select * from users.logon")
	if err != nil {
		log.Printf("[DB] GetAllLogon query %v", err)
	}

	ul := new([]UserLogon)

	for rows.Next() {
		var (
			id    int
			login string
			pass  string
		)

		if err := rows.Scan(&id, &login, &pass); err != nil {
			log.Printf("[DB] GetlAllLogon rows scan %v", err)
		}
		*ul = append(*ul, UserLogon{id, login, pass})
	}
	return *ul
}

func (d *Database) GetUserByCookie(cookie ...string) ([]UserLogon, error) {
	rows, err := d.DB.Query("select distinct ulg.id, ulg.login, ulg.pass from  users.session uss  inner join users.logon as ulg on ulg.id = uss.id_user and uss.session_cookie = $1", cookie[0]) // TODO: fix in array search
	if err != nil {
		return nil, err
	}

	ul := new([]UserLogon)

	for rows.Next() {

		var (
			id    int
			login string
			pass  string
		)

		if err := rows.Scan(&id, &login, &pass); err != nil {
			return nil, err
		}

		*ul = append(*ul, UserLogon{id, login, pass})

	}
	return *ul, nil

}

func (d *Database) GetUserByLogin(login ...string) ([]UserLogon, error) {
	rows, err := d.DB.Query("select * from users.logon  where login = $1", login[0]) // TODO: fix in array search
	if err != nil {
		return nil, err
	}

	fmt.Println(login[0])

	ul := new([]UserLogon)

	for rows.Next() {

		var (
			id    int
			login string
			pass  string
		)

		if err := rows.Scan(&id, &login, &pass); err != nil {
			return nil, err
		}

		*ul = append(*ul, UserLogon{id, login, pass})

	}
	return *ul, nil
}

func (d *Database) SetUser(login, password string) error {

	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		//log.Printf("[Encode json] registration %v", err)
		return err
	}

	if _, err := d.DB.Exec("INSERT INTO users.logon (login, pass) 	VALUES($1, $2)", login, string(hashPass)); err != nil {
		return err
	}
	return nil
}

func (d *Database) Authentications(login, password string) (string, bool, error) {

	ul, err := d.GetUserByLogin(login)
	if err != nil {
		return "", false, err
	}

	fmt.Println(ul)

	err = bcrypt.CompareHashAndPassword([]byte(ul[0].Pass), []byte(password))
	if err != nil {
		switch err.(type) {
		case error:
			return "", false, nil
		default:
			return "", false, err
		}
	}

	h := hmac.New(sha256.New, []byte(ul[0].Login+strconv.Itoa(ul[0].Id)))

	shaCookie := hex.EncodeToString(h.Sum(nil))

	if err := d.SetCookie(ul[0].Id, shaCookie); err != nil {
		return "", false, fmt.Errorf("[db] set cookie %v", err)
	}

	return shaCookie, true, nil
}

func (d *Database) SetCookie(idUser int, cookie string) error {

	if _, err := d.DB.Exec("insert into users.session (id_user, session_cookie) values($1,$2)", idUser, cookie); err != nil {
		return err
	}
	return nil
}
