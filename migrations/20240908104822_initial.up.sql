CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.employee (
	id uuid DEFAULT uuid_generate_v4() NOT NULL,
	username varchar(50) NOT NULL,
	first_name varchar(50) NULL,
	last_name varchar(50) NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT employee_pkey PRIMARY KEY (id),
	CONSTRAINT employee_username_key UNIQUE (username)
);

CREATE TYPE public."organization_type" AS ENUM (
	'IE',
	'LLC',
	'JSC');

CREATE TABLE public.organization (
	id uuid DEFAULT uuid_generate_v4() NOT NULL,
	"name" varchar(100) NOT NULL,
	description text NULL,
	"type" public."organization_type" NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT organization_pkey PRIMARY KEY (id)
);

CREATE TABLE public.organization_responsible (
	id uuid DEFAULT uuid_generate_v4() NOT NULL,
	organization_id uuid NULL,
	user_id uuid NULL,
	CONSTRAINT organization_responsible_pkey PRIMARY KEY (id),
	CONSTRAINT organization_responsible_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE,
	CONSTRAINT organization_responsible_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.employee(id) ON DELETE CASCADE
);

CREATE TYPE service_type AS ENUM (
    'Construction', 
    'Delivery',
    'Manufacture'
);

CREATE TYPE tender_status AS ENUM (
    'Created',
    'Published',
    'Closed'
);

CREATE TABLE tender (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    type service_type NOT NULL,
    status tender_status NOT NULL DEFAULT 'Created',
    organization_id UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
    version INT NOT NULL DEFAULT 1,
    creator_username VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, version)
);

CREATE INDEX idx_tender_id_hash ON tender USING HASH (id);
CREATE INDEX idx_tender_creator_username_hash ON tender USING HASH (creator_username);

CREATE TYPE author_type AS ENUM (
    'User',
    'Organization'
);

CREATE TYPE bid_status AS ENUM (
    'Created',
    'Published',
    'Canceled'
);

CREATE TABLE bid (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    author author_type NOT NULL,
    author_id UUID NOT NULL REFERENCES employee(id) ON DELETE CASCADE,
    status bid_status NOT NULL DEFAULT 'Created',
    version INT NOT NULL DEFAULT 1,
    tender_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, version)
);

CREATE INDEX idx_bid_id_hash ON bid USING HASH (id);
CREATE INDEX idx_bid_author_id_hash ON bid USING HASH (author_id);
CREATE INDEX idx_bid_tender_id_hash ON bid USING HASH (tender_id);

INSERT INTO public.employee (id,username,first_name,last_name,created_at,updated_at) VALUES
	 ('550e8400-e29b-41d4-a716-446655440001'::uuid,'user1','First1','Last1','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-446655440002'::uuid,'user2','First2','Last2','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-446655440003'::uuid,'user3','First3','Last3','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-446655440004'::uuid,'user4','First4','Last4','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-446655440005'::uuid,'user5','First5','Last5','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-446655440006'::uuid,'user6','First6','Last6','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-446655440007'::uuid,'user7','First7','Last7','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-446655440008'::uuid,'user8','First8','Last8','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-446655440009'::uuid,'user9','First9','Last9','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-44665544000a'::uuid,'user10','First10','Last10','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011');
INSERT INTO public.employee (id,username,first_name,last_name,created_at,updated_at) VALUES
	 ('550e8400-e29b-41d4-a716-44665544000b'::uuid,'user11','First11','Last11','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-44665544000c'::uuid,'user12','First12','Last12','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-44665544000d'::uuid,'user13','First13','Last13','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-44665544000e'::uuid,'user14','First14','Last14','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-44665544000f'::uuid,'user15','First15','Last15','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-446655440010'::uuid,'user16','First16','Last16','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-446655440011'::uuid,'user17','First17','Last17','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-446655440012'::uuid,'user18','First18','Last18','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-446655440013'::uuid,'user19','First19','Last19','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-446655440014'::uuid,'user20','First20','Last20','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011');
INSERT INTO public.employee (id,username,first_name,last_name,created_at,updated_at) VALUES
	 ('550e8400-e29b-41d4-a716-446655440015'::uuid,'user21','First21','Last21','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-446655440016'::uuid,'user22','First22','Last22','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-446655440017'::uuid,'user23','First23','Last23','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-446655440018'::uuid,'user24','First24','Last24','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-446655440019'::uuid,'user25','First25','Last25','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-44665544001a'::uuid,'user26','First26','Last26','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-44665544001b'::uuid,'user27','First27','Last27','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-44665544001c'::uuid,'user28','First28','Last28','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-44665544001d'::uuid,'user29','First29','Last29','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011'),
	 ('550e8400-e29b-41d4-a716-44665544001e'::uuid,'user30','First30','Last30','2024-09-09 17:47:19.92011','2024-09-09 17:47:19.92011');

INSERT INTO public.organization (id,"name",description,"type",created_at,updated_at) VALUES
	 ('550e8400-e29b-41d4-a716-446655440020'::uuid,'Organization 1','Description 1','LLC'::public."organization_type",'2024-09-09 17:47:19.955391','2024-09-09 17:47:19.955391'),
	 ('550e8400-e29b-41d4-a716-446655440021'::uuid,'Organization 2','Description 2','IE'::public."organization_type",'2024-09-09 17:47:19.955391','2024-09-09 17:47:19.955391'),
	 ('550e8400-e29b-41d4-a716-446655440022'::uuid,'Organization 3','Description 3','JSC'::public."organization_type",'2024-09-09 17:47:19.955391','2024-09-09 17:47:19.955391'),
	 ('550e8400-e29b-41d4-a716-446655440023'::uuid,'Organization 4','Description 4','LLC'::public."organization_type",'2024-09-09 17:47:19.955391','2024-09-09 17:47:19.955391');

INSERT INTO public.organization_responsible (id,organization_id,user_id) VALUES
	 ('550e8400-e29b-41d4-a716-446655440030'::uuid,'550e8400-e29b-41d4-a716-446655440020'::uuid,'550e8400-e29b-41d4-a716-446655440001'::uuid),
	 ('550e8400-e29b-41d4-a716-446655440031'::uuid,'550e8400-e29b-41d4-a716-446655440020'::uuid,'550e8400-e29b-41d4-a716-446655440002'::uuid),
	 ('550e8400-e29b-41d4-a716-446655440032'::uuid,'550e8400-e29b-41d4-a716-446655440020'::uuid,'550e8400-e29b-41d4-a716-446655440003'::uuid),
	 ('550e8400-e29b-41d4-a716-446655440033'::uuid,'550e8400-e29b-41d4-a716-446655440021'::uuid,'550e8400-e29b-41d4-a716-446655440004'::uuid),
	 ('550e8400-e29b-41d4-a716-446655440034'::uuid,'550e8400-e29b-41d4-a716-446655440021'::uuid,'550e8400-e29b-41d4-a716-446655440005'::uuid),
	 ('550e8400-e29b-41d4-a716-446655440035'::uuid,'550e8400-e29b-41d4-a716-446655440021'::uuid,'550e8400-e29b-41d4-a716-446655440006'::uuid),
	 ('550e8400-e29b-41d4-a716-446655440036'::uuid,'550e8400-e29b-41d4-a716-446655440022'::uuid,'550e8400-e29b-41d4-a716-446655440007'::uuid),
	 ('550e8400-e29b-41d4-a716-446655440037'::uuid,'550e8400-e29b-41d4-a716-446655440022'::uuid,'550e8400-e29b-41d4-a716-446655440008'::uuid),
	 ('550e8400-e29b-41d4-a716-446655440038'::uuid,'550e8400-e29b-41d4-a716-446655440022'::uuid,'550e8400-e29b-41d4-a716-446655440009'::uuid),
	 ('550e8400-e29b-41d4-a716-446655440039'::uuid,'550e8400-e29b-41d4-a716-446655440023'::uuid,'550e8400-e29b-41d4-a716-44665544000a'::uuid);
INSERT INTO public.organization_responsible (id,organization_id,user_id) VALUES
	 ('550e8400-e29b-41d4-a716-44665544003a'::uuid,'550e8400-e29b-41d4-a716-446655440023'::uuid,'550e8400-e29b-41d4-a716-44665544000b'::uuid),
	 ('550e8400-e29b-41d4-a716-44665544003b'::uuid,'550e8400-e29b-41d4-a716-446655440023'::uuid,'550e8400-e29b-41d4-a716-44665544000c'::uuid);
