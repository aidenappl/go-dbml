CREATE TABLE IF NOT EXISTS teams (id int PRIMARY KEY, name text NOT NULL);
CREATE TABLE IF NOT EXISTS players (
    id int PRIMARY KEY,
    team_id int,
    jersey int,
    name text NOT NULL,
    UNIQUE (team_id, jersey)
);
CREATE TABLE IF NOT EXISTS matches (
    id int PRIMARY KEY,
    home_team int NOT NULL,
    away_team int NOT NULL,
    mvp_team int,
    mvp_jersey int,
    played_on date
);
ALTER TABLE players ADD CONSTRAINT players_team_fk FOREIGN KEY (team_id) REFERENCES teams(id) ON UPDATE CASCADE ON DELETE SET NULL;
ALTER TABLE matches ADD FOREIGN KEY (home_team) REFERENCES teams(id);
ALTER TABLE matches ADD FOREIGN KEY (away_team) REFERENCES teams(id);
ALTER TABLE matches ADD CONSTRAINT matches_mvp_fk FOREIGN KEY (mvp_team, mvp_jersey) REFERENCES players(team_id, jersey);
