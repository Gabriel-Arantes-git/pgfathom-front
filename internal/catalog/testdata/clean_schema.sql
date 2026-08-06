-- Cenário de controle: tudo declarado, validado e indexado. A ferramenta deve
-- afirmar que analisou e está limpo, em vez de apenas não listar nada.

CREATE TABLE municipio (
    id   bigint PRIMARY KEY,
    nome text NOT NULL
);

CREATE TABLE endereco (
    id           bigint PRIMARY KEY,
    municipio_id bigint NOT NULL REFERENCES municipio (id),
    logradouro   text NOT NULL
);

CREATE INDEX ix_endereco_municipio_id ON endereco (municipio_id);

INSERT INTO municipio (id, nome) VALUES (1, 'Sao Bernardo do Campo');
INSERT INTO endereco (id, municipio_id, logradouro) VALUES (1, 1, 'Rua das Acacias 42');
