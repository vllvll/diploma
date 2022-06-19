package postgres

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func ConnectDatabase(uri string) (*sql.DB, error) {
	var db *sql.DB

	db, err := sql.Open("pgx", uri)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		create table users
		(
			id            serial
				constraint users_pk
					primary key,
			login         text      not null,
			password_hash integer   not null,
			created_at    timestamp not null
		);
		
		alter table users
			owner to postgres;
		
		create unique index users_id_uindex
			on users (id);
		
		create unique index users_login_uindex
			on users (login);
		
		create table orders
		(
			id          serial
				constraint orders_pk
					primary key,
			number      integer   not null,
			user_id     integer   not null
				constraint orders_users_id_fk
					references users,
			status      text      not null,
			accrual     integer   not null,
			uploaded_at timestamp not null
		);
		
		alter table orders
			owner to postgres;
		
		create unique index orders_id_uindex
			on orders (id);
		
		create table balances
		(
			id         serial
				constraint balances_pk
					primary key,
			user_id    integer   not null
				constraint balances_users_id_fk
					references users,
			sum        integer   not null,
			created_at timestamp not null,
			updated_at timestamp not null
		);
		
		alter table balances
			owner to postgres;
		
		create unique index balances_id_uindex
			on balances (id);
	`)

	if err != nil {
		return db, nil
	}

	return db, nil
}
