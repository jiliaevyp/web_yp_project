
create table if not exists personals
(
    Id        serial       not null
        constraint personals_pk
            primary key,
    Title     varchar  not null,
    Forename      varchar  not null,
    Kadr      varchar  not null,
    Numotdel integer      not null,
    Tarif     integer      not null,
    Email     varchar not null,
    Phone     varchar  not null,
    Address   varchar not null
);

comment on table personals is 'таблица персонала';

alter table personals
    owner to yp;

create unique index personals_email_uindex
    on personals (Email);

create unique index personals_id_uindex
    on personals (Id);
