create table if not exists service
(
    service_id serial
        constraint service_pk
            primary key,
    name       text not null,
    deleted    timestamp
);

create unique index if not exists service_name_uindex
    on service (name);

create table if not exists version
(
    version_id serial
        constraint version_pk
            primary key,
    service_id integer
        constraint config_service_service_id_fk
            references service
            on update cascade on delete cascade,
    number     integer not null,
    deleted    timestamp
);

create unique index if not exists version_number_service_id_uindex
    on version (number, service_id);

create table if not exists config
(
    config_id  serial
        constraint config_pk
            primary key,
    version_id integer
        constraint config_version_version_id_fk
            references version
            on update cascade on delete cascade,
    param      text not null,
    value      text not null
);

create unique index if not exists config_version_id_param_uindex
    on config (version_id, param);