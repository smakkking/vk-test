CREATE TYPE sex AS ENUM('мужчина', 'женщина');
CREATE TABLE Actors (
    a_id UUID,
    a_name TEXT NOT NULL,
    a_sex sex NOT NULL,
    a_birth_date DATE NOT NULL,
    PRIMARY KEY (a_id)
);
CREATE TABLE Films (
    film_id UUID,
    film_title VARCHAR(150) NOT NULL,
    film_desc VARCHAR(1000),
    film_rating INT,
    CHECK (
        film_rating >= 0
        AND film_rating <= 10
    ),
    PRIMARY KEY (film_id)
);
CREATE TABLE ActorToFilm(
    actor_id UUID,
    film_id UUID,
    PRIMARY KEY (actor_id, film_id),
    FOREIGN KEY (actor_id) REFERENCES Actors(a_id) ON DELETE CASCADE,
    FOREIGN KEY (film_id) REFERENCES Films(film_id) ON DELETE CASCADE,
);