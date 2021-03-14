drop table if exists users;

create table users (
	user_name varchar primary key,
	first_name varchar(40) not null,
	last_name varchar(40) not null,
	email varchar not null,
	pass_hash varchar not null,  
	created_at timestamptz not null DEFAULT (now()),
	
	unique(email)
);