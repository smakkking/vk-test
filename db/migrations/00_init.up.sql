CREATE TYPE SEX AS ENUM('мужчина', 'женщина');
CREATE TABLE Actors (
    a_id INT SERIAL,
    a_name TEXT NOT NULL,
    a_sex SEX NOT NULL,
    a_birth_date DATE NOT NULL,
    PRIMARY KEY (a_id)
);
CREATE TABLE Films (
    f_id INT SERIAL,
    f_title VARCHAR(150) NOT NULL,
    f_desc VARCHAR(1000),
    f_date_creation DATE NOT NULL,
    f_rating INT NOT NULL,
    CHECK (
        film_rating >= 0
        AND film_rating <= 10
    ),
    PRIMARY KEY (f_id)
);
CREATE TABLE ActorToFilm(
    actor_id INT,
    film_id INT,
    PRIMARY KEY (actor_id, film_id),
    FOREIGN KEY (actor_id) REFERENCES Actors(a_id) ON DELETE CASCADE,
    FOREIGN KEY (film_id) REFERENCES Films(f_id) ON DELETE CASCADE,
);