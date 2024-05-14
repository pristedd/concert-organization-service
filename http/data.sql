CREATE EXTENSION pgcrypto;

CREATE TABLE public.type (
    id UUID PRIMARY KEY default (gen_random_uuid()),
    typeName VARCHAR(300) NOT NULL
);

CREATE TABLE public.event (
    id UUID PRIMARY KEY default (gen_random_uuid()),
    name VARCHAR(300) NOT NULL
);

CREATE TABLE public.ticket (
    id UUID PRIMARY KEY default (gen_random_uuid()),
    typeId UUID NOT NULL,
    eventId UUID NOT NULL,
    Price FLOAT NOT NULL,
    Purchased BOOLEAN NOT NULL,

    CONSTRAINT event_fk FOREIGN KEY (eventId) REFERENCES public.event(id),
    CONSTRAINT type_fk FOREIGN KEY (typeId) REFERENCES  public.ticket(id)
);

insert into public.type (typename) values ('танцпол');
insert into public.type (typename) values ('танцпол+');
insert into public.type (typename) values ('fan');
