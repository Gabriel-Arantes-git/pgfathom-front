-- Cenário de inferência: relações que existem nos dados e não no catálogo,
-- misturadas com uma armadilha e com formas que a inferência deve pular.
--
-- Nomenclatura em português com prefixo de convenção antiga e plural, que é
-- exatamente o que o perfil pt-br existe para atravessar.

CREATE TABLE tb_clientes (
    id   bigint PRIMARY KEY,
    cpf  text NOT NULL,
    nome text NOT NULL
);

CREATE TABLE municipios (
    id   bigint PRIMARY KEY,
    nome text NOT NULL
);

-- Tabela de domínio pequena: a armadilha. `status_id` casa por nome, e a
-- penalidade precisa incidir.
CREATE TABLE status (
    id        bigint PRIMARY KEY,
    descricao text NOT NULL
);

CREATE TABLE pedido (
    id           bigint PRIMARY KEY,
    cliente_id   bigint NOT NULL,   -- inferível: tb_clientes via prefixo e plural
    municipio_id bigint,            -- inferível: municipios via plural
    status_id    bigint,            -- armadilha: nome genérico, alvo pequeno
    valor        numeric(12,2)
);

CREATE INDEX ix_pedido_cliente_id ON pedido (cliente_id);

-- Par polimórfico: precisa ser reconhecido, não apenas rejeitado.
CREATE TABLE anexo (
    id            bigint PRIMARY KEY,
    entidade_id   bigint NOT NULL,
    entidade_tipo text NOT NULL
);

-- Alvo com chave composta: precisa ser pulado com registro.
CREATE TABLE matricula (
    aluno_id bigint NOT NULL,
    turma_id bigint NOT NULL,
    PRIMARY KEY (aluno_id, turma_id)
);

CREATE TABLE historico (
    id            bigint PRIMARY KEY,
    matricula_id  bigint NOT NULL
);

INSERT INTO tb_clientes (id, cpf, nome) VALUES
    (1, '529.318.470-11', 'Maria Aparecida Silva');
INSERT INTO municipios (id, nome) VALUES (1, 'Sao Bernardo do Campo');
INSERT INTO status (id, descricao) VALUES (1, 'aberto'), (2, 'fechado');
INSERT INTO pedido (id, cliente_id, municipio_id, status_id, valor)
    VALUES (1, 1, 1, 1, 150.00);
INSERT INTO anexo (id, entidade_id, entidade_tipo) VALUES (1, 1, 'pedido');
