-- Flyway Migration: Initial schema

-- =====================
-- Language (先に定義)
-- =====================
CREATE TABLE language (
    language_ID SERIAL PRIMARY KEY,
    language_short TEXT,
    language_common TEXT
);

INSERT INTO language VALUES
(1, 'en', 'English'),
(2, 'jp', '日本語');

-- =====================
-- Wiki pages (先に定義)
-- =====================
CREATE TABLE wiki_pages (
    page_ID SERIAL PRIMARY KEY,
    title TEXT,
    userID INT, -- usersテーブルをまだ作っていないので、REFERENCESは後で追加する
    created_date TIMESTAMP DEFAULT now(),
    updated_date TIMESTAMP,
    content_path TEXT
);

-- =====================
-- Users & Roles
-- =====================
CREATE TABLE user_roles (
    role_ID SERIAL PRIMARY KEY,
    role_name TEXT NOT NULL UNIQUE
);

INSERT INTO user_roles (role_ID, role_name) VALUES
(1, 'admin'),
(2, 'editor'),
(3, 'viewer'),
(4, 'guest');


CREATE TABLE users (
    userID SERIAL PRIMARY KEY,
    user_name VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    mail_address VARCHAR(255) UNIQUE,
    password TEXT ,
    role_ID INT REFERENCES user_roles(role_ID),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    timezone SMALLINT NOT NULL
);

-- usersテーブルができたので、wiki_pagesテーブルに外部キー制約を追加
ALTER TABLE wiki_pages ADD FOREIGN KEY (userID) REFERENCES users(userID);

INSERT INTO users (userID, user_name, display_name, mail_address, password, role_ID, created_at, timezone) VALUES
(1, 'admin', 'admin', 'admin@mail.com', 'admin-user', 1, '2001-01-01 00:00:00+00', 0);


CREATE TABLE users_defaults (
    userID INT PRIMARY KEY REFERENCES users(userID),
    theme TEXT NOT NULL DEFAULT 'white'
);

INSERT INTO users_defaults (userID, theme) VALUES
(1, 'white');

-- =====================
-- Projects
-- =====================
CREATE TABLE projects (
    project_ID SERIAL PRIMARY KEY,
    project_name TEXT NOT NULL,
    disscription TEXT,
    start_day DATE,
    finished_day DATE,
    updated_day DATE,
    note TEXT
);

CREATE TABLE project_members (
    project_member_ID SERIAL PRIMARY KEY,
    project_ID INT REFERENCES projects(project_ID),
    userID INT REFERENCES users(userID),
    join_day DATE,
    finish_day DATE
);

-- =====================
-- Places
-- =====================
CREATE TABLE place_names_json (
    place_name_ID SERIAL PRIMARY KEY,
    class_place_name jsonb
);

CREATE EXTENSION IF NOT EXISTS postgis; -- IF NOT EXISTS をつけると安全

CREATE TABLE places (
    place_ID SERIAL PRIMARY KEY,
    coordinates GEOGRAPHY(Point, 4326),
    place_name_ID INT REFERENCES place_names_json(place_name_ID),
    accuracy NUMERIC
);

-- =====================
-- Classification
-- =====================
CREATE TABLE classification_json (
    classification_ID SERIAL PRIMARY KEY,
    class_classification jsonb
);

-- =====================
-- File Attachments (先に定義)
-- =====================
CREATE TABLE file_types (
    file_type_ID SERIAL PRIMARY KEY,
    type_name TEXT
);

CREATE TABLE file_extensions(
    extension_ID SERIAL PRIMARY KEY,
    extension_text VARCHAR(255),
    file_type_ID INT REFERENCES file_types(file_type_ID)
);

CREATE TABLE attachments (
    attachment_ID SERIAL PRIMARY KEY,
    file_path TEXT NOT NULL,
    extension_ID INT REFERENCES file_extensions(extension_ID),
    user_ID INT REFERENCES users(userID)
);

-- =====================
-- Occurrence (attachmentsを参照するので先に定義が必要)
-- =====================
CREATE TABLE occurrence (
    occurrence_ID SERIAL PRIMARY KEY,
    project_ID INT REFERENCES projects(project_ID),
    user_ID INT REFERENCES users(userID),
    individual_ID INT,
    lifestage TEXT,
    sex TEXT,
    classification_ID INT REFERENCES classification_json(classification_ID),
    place_ID INT REFERENCES places(place_ID),
    attachment_group_ID INT,
    body_length NUMERIC,
    language_ID INT REFERENCES language(language_ID),
    note TEXT, -- カンマ追加
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    timezone SMALLINT NOT NULL
);

-- =====================
-- Attachment Group (occurrenceとattachmentsの両方が必要)
-- =====================
CREATE TABLE attachment_goup (
    occurrence_ID INT REFERENCES occurrence(occurrence_ID),
    attachment_ID INT REFERENCES attachments(attachment_ID),
    priority INT,

    PRIMARY KEY(occurrence_ID, attachment_ID)
);

-- =====================
-- Specimen
-- =====================
CREATE TABLE institution_ID_code (
    institution_ID SERIAL PRIMARY KEY,
    institution_code TEXT
);

CREATE TABLE collection_ID_code (
    collection_ID SERIAL PRIMARY KEY,
    collection_code TEXT
);

CREATE TABLE specimen_methods (
    specimen_methods_ID SERIAL PRIMARY KEY,
    method_common_name TEXT,
    page_ID INT REFERENCES wiki_pages(page_ID)
);

CREATE TABLE specimen (
    specimen_ID SERIAL PRIMARY KEY,
    occurrence_ID INT REFERENCES occurrence(occurrence_ID),
    specimen_method_ID INT REFERENCES specimen_methods(specimen_methods_ID),
    institution_ID INT REFERENCES institution_ID_code(institution_ID),
    collectionID INT REFERENCES collection_ID_code(collection_ID)
);

CREATE TABLE make_specimen (
    make_specimen_ID SERIAL PRIMARY KEY,
    occurrence_ID INT REFERENCES occurrence(occurrence_ID),
    userID INT REFERENCES users(userID),
    specimen_ID INT REFERENCES specimen(specimen_ID),
    date DATE,
    specimen_method_ID INT REFERENCES specimen_methods(specimen_methods_ID), -- カンマ追加
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    timezone SMALLINT NOT NULL
);

-- =====================
-- Observations
-- =====================
CREATE TABLE observation_methods (
    observation_method_ID SERIAL PRIMARY KEY,
    method_common_name TEXT,
    pageID INT REFERENCES wiki_pages(page_ID)
);

CREATE TABLE observations (
    observations_ID SERIAL PRIMARY KEY,
    userID INT REFERENCES users(userID),
    occurrence_ID INT REFERENCES occurrence(occurrence_ID),
    observation_method_ID INT REFERENCES observation_methods(observation_method_ID),
    behavior TEXT, -- カンマ追加
    observed_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    timezone SMALLINT NOT NULL
);

-- =====================
-- Identifications
-- =====================
CREATE TABLE identifications (
    identification_ID SERIAL PRIMARY KEY,
    userID INT REFERENCES users(userID),
    occurrence_ID INT REFERENCES occurrence(occurrence_ID),
    source_info TEXT, -- カンマ追加
    identificated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    timezone SMALLINT NOT NULL
);

-- =====================
-- Change logs
-- =====================
CREATE TABLE change_logs (
    log_ID SERIAL PRIMARY KEY,
    type TEXT,
    changed_ID INT,
    before_value TEXT,
    after_value TEXT,
    user_ID INT REFERENCES users(userID),
    date TIMESTAMP DEFAULT now(),
    Row TEXT
);
