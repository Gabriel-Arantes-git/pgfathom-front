-- Cenário: FKs declaradas sem índice utilizável do lado filho.
--
-- O caso decisivo é `nota_fiscal`: a coluna filha aparece num índice composto,
-- mas em posição não inicial. Esse índice não serve para a busca que o DELETE
-- no pai dispara, e contá-lo como cobertura seria um falso negativo — a
-- ferramenta diria que está tudo bem numa tabela onde o problema é real.

CREATE TABLE fornecedor (
    id   bigint PRIMARY KEY,
    nome text NOT NULL
);

CREATE TABLE contrato (
    id            bigint PRIMARY KEY,
    fornecedor_id bigint NOT NULL REFERENCES fornecedor (id),
    numero        text NOT NULL
);

CREATE TABLE nota_fiscal (
    id            bigint PRIMARY KEY,
    fornecedor_id bigint NOT NULL REFERENCES fornecedor (id),
    emitida_em    date NOT NULL
);

CREATE TABLE empenho (
    id            bigint PRIMARY KEY,
    fornecedor_id bigint NOT NULL REFERENCES fornecedor (id),
    exercicio     int NOT NULL
);

-- contrato: sem índice nenhum na coluna filha.

-- nota_fiscal: coluna filha em posição NÃO inicial. Inútil para o lookup.
CREATE INDEX ix_nota_fiscal_composto ON nota_fiscal (emitida_em, fornecedor_id);

-- empenho: coberta corretamente, em posição inicial de um índice composto.
CREATE INDEX ix_empenho_composto ON empenho (fornecedor_id, exercicio);

INSERT INTO fornecedor (id, nome) VALUES (1, 'Construtora Horizonte LTDA');
INSERT INTO contrato (id, fornecedor_id, numero) VALUES (1, 1, 'CT-2019-0041');
INSERT INTO nota_fiscal (id, fornecedor_id, emitida_em) VALUES (1, 1, '2019-03-14');
INSERT INTO empenho (id, fornecedor_id, exercicio) VALUES (1, 1, 2019);
