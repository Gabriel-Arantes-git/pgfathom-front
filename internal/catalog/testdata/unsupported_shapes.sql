-- Cenário: formas que a inferência de coluna única não alcança. Todas precisam
-- aparecer na cobertura com o motivo, nunca serem puladas em silêncio.

-- PK composta.
CREATE TABLE matricula (
    aluno_id  bigint NOT NULL,
    turma_id  bigint NOT NULL,
    situacao  text,
    PRIMARY KEY (aluno_id, turma_id)
);

-- Sem PK nenhuma.
CREATE TABLE log_importacao (
    ocorrido_em timestamptz NOT NULL,
    mensagem    text
);

-- Particionada: a leitura vem do pai, as partições não entram no modelo.
CREATE TABLE lancamento (
    id       bigint NOT NULL,
    competencia date NOT NULL,
    valor    numeric(12,2)
) PARTITION BY RANGE (competencia);

CREATE TABLE lancamento_2019 PARTITION OF lancamento
    FOR VALUES FROM ('2019-01-01') TO ('2020-01-01');
CREATE TABLE lancamento_2020 PARTITION OF lancamento
    FOR VALUES FROM ('2020-01-01') TO ('2021-01-01');

-- Herança, rara mas presente em base antiga.
CREATE TABLE documento (
    id   bigint PRIMARY KEY,
    tipo text
);
CREATE TABLE documento_digitalizado (
    arquivo text
) INHERITS (documento);

-- Suportada, para provar que o contador de analisadas não zera.
CREATE TABLE municipio (
    id   bigint PRIMARY KEY,
    nome text NOT NULL
);
