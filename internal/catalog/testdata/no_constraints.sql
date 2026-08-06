-- Cenário: as mesmas entidades do schema limpo, sem nenhuma FK declarada.
-- A relação existe nos dados e não existe no catálogo — é o caso central que a
-- inferência vai atacar na fase 3, e aqui serve para provar que o `audit` não
-- inventa achado onde não há constraint alguma.

CREATE TABLE municipio (
    id   bigint PRIMARY KEY,
    nome text NOT NULL
);

CREATE TABLE endereco (
    id           bigint PRIMARY KEY,
    municipio_id bigint NOT NULL,
    logradouro   text NOT NULL
);

INSERT INTO municipio (id, nome) VALUES (1, 'Sao Bernardo do Campo');
INSERT INTO endereco (id, municipio_id, logradouro) VALUES (1, 1, 'Rua das Acacias 42');
