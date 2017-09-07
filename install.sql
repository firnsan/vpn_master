create table if not exists tb_user(
		id int auto_increment not null,
		name varchar(32) default "" not null,
		password varchar(16) default "" not null,
		password_encoded varchar(32) default "" not null,
		zone varchar(16) default "" not null,
		create_time datetime default "00-00-00 00:00:00" not null,
		expire_time datetime default "00-00-00 00:00:00" not null,
		status tinyint default 0 not null,
		primary key(id),
		unique key(name)
) engine=innodb auto_increment=1 default charset=utf8mb4;
