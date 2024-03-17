CREATE TYPE sex AS ENUM('мужчина', 'женщина');
CREATE TABLE Actors (
    a_name TEXT NOT NULL,
    a_sex sex NOT NULL,
    a_birth_date DATE NOT NULL,
    PRIMARY KEY (a_name)
);
CREATE TABLE Films (
    f_title VARCHAR(150) NOT NULL,
    f_desc VARCHAR(1000),
    f_rating INT,
    CHECK (
        f_rating >= 0
        AND f_rating <= 10
    ),
    PRIMARY KEY (f_title)
);
CREATE TABLE ActorToFilm(
    actor_name TEXT NOT NULL,
    film_title VARCHAR(150) NOT NULL,
    PRIMARY KEY (actor_name, film_title),
    FOREIGN KEY (actor_name) REFERENCES Actors(a_name) ON DELETE CASCADE,
    FOREIGN KEY (film_title) REFERENCES Films(f_title) ON DELETE CASCADE,
);