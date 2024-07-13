-- +goose Up
-- +goose StatementBegin
CREATE TABLE files (
    id INT NOT NULL AUTO_INCREMENT,
    dirPath VARCHAR(255) NOT NULL,
    filename VARCHAR(255) NOT NULL,
    operation VARCHAR(255) NOT NULL,
    date TIMESTAMP,
    PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE files;
-- +goose StatementEnd
