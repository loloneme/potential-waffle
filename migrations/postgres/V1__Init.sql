CREATE TABLE teams(
    team_name VARCHAR(255) NOT NULL,

    PRIMARY KEY (team_name)
);


CREATE TABLE users(
   user_id VARCHAR(255) NOT NULL,
   username VARCHAR(255) NOT NULL,
   is_active BOOLEAN NOT NULL,
   team_name VARCHAR(255) NOT NULL,

   PRIMARY KEY (user_id),
   FOREIGN KEY (team_name) REFERENCES teams(team_name) ON DELETE CASCADE
);


CREATE TABLE statuses(
     status_id SERIAL NOT NULL,
     status_name VARCHAR(255) NOT NULL,

     PRIMARY KEY (status_id)
);

CREATE TABLE pull_requests(
   pr_id VARCHAR(255) NOT NULL,
   pr_name TEXT NOT NULL,
   author_id VARCHAR(255) NOT NULL ,
   status_id INTEGER NOT NULL,
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   merged_at TIMESTAMP,

    PRIMARY KEY (pr_id),
    FOREIGN KEY (author_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (status_id) REFERENCES statuses(status_id)
);

CREATE TABLE reviewers(
    pr_id VARCHAR(255) NOT NULL,
    reviewer_id VARCHAR(255) NOT NULL,

    PRIMARY KEY (pr_id, reviewer_id),
    FOREIGN KEY (pr_id) REFERENCES pull_requests(pr_id) ON DELETE CASCADE,
    FOREIGN KEY (reviewer_id) REFERENCES users(user_id) ON DELETE CASCADE

);




