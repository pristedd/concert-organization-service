CREATE EXTENSION pgcrypto;

CREATE TABLE public.location (
                                 id UUID PRIMARY KEY default (gen_random_uuid()),
                                 name VARCHAR(300) NOT NULL,
                                 address VARCHAR(500) NOT NULL,
                                 seats_num INT NOT NULL,
                                 comment VARCHAR(500)
);

CREATE TABLE public.artist (
                               id UUID PRIMARY KEY default (gen_random_uuid()),
                               name VARCHAR(300) NOT NULL,
                               email VARCHAR(300) NOT NULL,
                               phone VARCHAR(11) NOT NULL
);

CREATE TABLE public.booking (
    id UUID PRIMARY KEY default (gen_random_uuid()),
    location_id UUID NOT NULL,
    artist_id UUID NOT NULL,
    date DATE NOT NULL,
    comment VARCHAR(500),

    CONSTRAINT loc_fk FOREIGN KEY (location_id) REFERENCES public.location(id),
    CONSTRAINT art_fk FOREIGN KEY (artist_id) REFERENCES  public.artist(id)
);

INSERT INTO location (name, address, seats_num)
VALUES
    ('AdrenalineRushStadium', 'Kirova Pr., bld. 407, appt. 52', 20000),
    ('Лужники', 'ул. Лужники, 24, стр. 1, Москва', 81000),
    ('Аврора', 'Sudostroitelnaya Ul., bld. 52, appt. 37', 10000);