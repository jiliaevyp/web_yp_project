
create table if not exists personals
(
    id        serial       not null
        constraint personals_pk
            primary key,
    title     varchar  not null,
    forename      varchar  not null,
    kadr      varchar  not null,
    numotdel integer      not null,
    tarif     integer      not null,
    area     varchar  not null,
    email     varchar not null,
    phone     varchar  not null,
    address   varchar not null
);

comment on table personals is 'таблица персонала';

alter table personals
    owner to yp;

create unique index personals_email_uindex
    on personals (email);

create unique index personals_id_uindex
    on personals (id);
