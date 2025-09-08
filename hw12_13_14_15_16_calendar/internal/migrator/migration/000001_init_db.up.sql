create table account (id serial primary key, 
	login varchar(100) not null,
	password varchar(255) not null,
	created_at timestamp not null default now(),
	unique(login));

create table event 
    (id serial primary key,
	title varchar(255) not null,
	created_at timestamp not null default now(),
    date_start timestamp not null default now(),
	date_end timestamp not null default now(),
	description text,
	account_id integer,
	notification timestamp not null default now(),
    foreign key (account_id) references account (id),
    unique(account_id));