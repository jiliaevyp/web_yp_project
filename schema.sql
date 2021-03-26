create table if not exists monds
(
    id              serial               not null
        constraint monds_pk
            primary key,
    yahre           integer default 2021 not null,
    monat           varchar(16)          not null,
    num_monat       integer              not null,
    tag             integer default 21   not null,
    hour            integer default 168  not null,
    kf_oberhour     integer default 2    not null,
    block_monat     integer default 0    not null,
    block_personal  integer default 0    not null,
    block_timetabel integer default 0    not null,
    block_tabel     integer default 0    not null,
    block_buchtabel integer default 0    not null,
    time_stamp      timestamp with time zone
);

alter table monds
    owner to yp;

create unique index monds_id_uindex
    on monds (id);

create table if not exists personals
(
    id        serial       not null
        constraint personals_pk
            primary key,
    forename      varchar(32)  not null,
    title     varchar(32)  not null,
    kadr      varchar(64)  not null,
    otdel     varchar(64)  not null,
    num_otdel integer      not null,
    email     varchar(64)  not null,
    phone     varchar(16)  not null,
    address   varchar(256) not null,
    timestamp timestamp with time zone,
    tarif     integer      not null
);

comment on table personals is 'таблица персонала';

alter table personals
    owner to yp;

create unique index personals_email_uindex
    on personals (email);

create unique index personals_id_uindex
    on personals (id);
