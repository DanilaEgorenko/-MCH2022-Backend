package cookieRepository

import (
	"context"
	"fmt"
	"snakealive/m/pkg/domain"

	pgxpool "github.com/jackc/pgx/v4/pgxpool"
)

type cookieStorage struct {
	dataHolder *pgxpool.Pool
}

func NewCookieStorage(DB *pgxpool.Pool) domain.CookieStorage {
	return &cookieStorage{dataHolder: DB}
}

func (c *cookieStorage) Add(key string, userId int) error {
	conn, err := c.dataHolder.Acquire(context.Background())
	if err != nil {
		fmt.Printf("Connection error while adding cookie ", err)
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(),
		`INSERT INTO Cookies ("hash", "user_id") VALUES ($1, $2)`,
		key,
		userId,
	)
	return err
}

func (c *cookieStorage) Get(value string) (user domain.Company, err error) {
	conn, err := c.dataHolder.Acquire(context.Background())
	if err != nil {
		fmt.Printf("Error while getting cookie")
		return user, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT U.id, U.name, U.surname, U.password, U.email
		FROM Users AS U
		JOIN Cookies AS C ON U.id = C.user_id
		WHERE C.hash = $1`,
		value,
	).Scan(&user.Id, &user.Name, &user.Password, &user.Email)

	return user, err

}
func (c *cookieStorage) Delete(value string) error {
	conn, err := c.dataHolder.Acquire(context.Background())
	if err != nil {
		fmt.Printf("Connection error while deleting cookie", err)
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(),
		`DELETE FROM Cookies WHERE hash = $1`,
		value,
	)
	return err
}
