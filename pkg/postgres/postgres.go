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
			password_hash text      not null,
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
			number      bigint            not null,
			user_id     integer           not null
				constraint orders_users_id_fk
					references users,
			status      text              not null,
			uploaded_at timestamp         not null,
			accrual     integer default 0 not null
		);
		
		alter table orders
			owner to postgres;
		
		create unique index orders_id_uindex
			on orders (id);
		
		create unique index orders_number_uindex
			on orders (number);
		
		create table balances
		(
			id         serial
				constraint balances_pk
					primary key,
			user_id    integer                    not null
				constraint balances_users_id_fk
					references users,
			sum        double precision default 0 not null,
			created_at timestamp                  not null,
			updated_at timestamp                  not null
		);
		
		alter table balances
			owner to postgres;
		
		create unique index balances_id_uindex
			on balances (id);
		
		create unique index balances_user_id_uindex
			on balances (user_id);
		
		create table tokens
		(
			id         serial
				constraint tokens_pk
					primary key,
			token      text      not null,
			user_id    integer   not null
				constraint tokens_users_id_fk
					references users,
			last_login timestamp not null
		);
		
		alter table tokens
			owner to postgres;
		
		create unique index tokens_id_uindex
			on tokens (id);
		
		create table withdraw
		(
			id         serial
				constraint withdraw_pk
					primary key,
			user_id    integer          not null,
			number     bigint           not null,
			sum        double precision not null,
			created_at timestamp        not null
		);
		
		alter table withdraw
			owner to postgres;
		
		create unique index withdraw_id_uindex
			on withdraw (id);
	`)

	if err != nil {
		return db, nil
	}

	return db, nil
}
