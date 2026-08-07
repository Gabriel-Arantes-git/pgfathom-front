-- Cenário da validação contra dados: cada relação abaixo é inferível por nome
-- e desenhada para produzir um veredito específico em modo completo. As
-- cardinalidades ficam dentro do que o pré-filtro aceita, para que todos os
-- candidatos cheguem à validação.
--
-- Órfãos usam valores 900+ para nunca colidirem com chave válida. Os valores
-- de texto plantados alimentam a varredura de vazamento.

-- Confirmada: contenção total, mais de um valor distinto.
CREATE TABLE pedido (
    id         bigint PRIMARY KEY,
    cliente_id bigint NOT NULL,
    descricao  text
);

CREATE TABLE item_pedido (
    id        bigint PRIMARY KEY,
    pedido_id bigint NOT NULL,
    descricao text
);

-- Quebrada: 47 de 50 linhas contidas (94%), 3 órfãs em 2 valores distintos.
CREATE TABLE cliente (
    id   bigint PRIMARY KEY,
    cpf  text NOT NULL,
    nome text NOT NULL
);

-- Quebrada com volume, para o modo amostrado: 380 de 400 contidas (95%).
CREATE TABLE conta (
    id   bigint PRIMARY KEY,
    nome text NOT NULL
);

CREATE TABLE lancamento (
    id       bigint PRIMARY KEY,
    conta_id bigint NOT NULL
);

-- Rejeitada: 9 de 30 contidas (30%).
CREATE TABLE contrato (
    id     bigint PRIMARY KEY,
    objeto text NOT NULL
);

CREATE TABLE fatura (
    id          bigint PRIMARY KEY,
    contrato_id bigint NOT NULL
);

-- Zona morta: 7 de 10 contidas (70%).
CREATE TABLE equipe (
    id   bigint PRIMARY KEY,
    nome text NOT NULL
);

CREATE TABLE chamado (
    id        bigint PRIMARY KEY,
    equipe_id bigint NOT NULL
);

-- Fraca por valor único: toda linha aponta para a mesma moeda.
CREATE TABLE moeda (
    id     bigint PRIMARY KEY,
    codigo text NOT NULL
);

CREATE TABLE nota (
    id       bigint PRIMARY KEY,
    moeda_id bigint NOT NULL
);

-- Fraca por filha vazia.
CREATE TABLE fornecedor (
    id   bigint PRIMARY KEY,
    nome text NOT NULL
);

CREATE TABLE pagamento (
    id            bigint PRIMARY KEY,
    fornecedor_id bigint NOT NULL
);

-- Identificadores exóticos: a citação precisa funcionar contra objeto real.
CREATE TABLE "Unidade Gestora" (
    id   bigint PRIMARY KEY,
    nome text NOT NULL
);

CREATE TABLE "Empenho 2024" (
    id                 bigint PRIMARY KEY,
    "unidade gestora"  bigint NOT NULL
);

INSERT INTO cliente (id, cpf, nome)
SELECT g, '529.318.470-' || lpad(g::text, 2, '0'), 'Maria Aparecida Silva'
FROM generate_series(1, 30) AS g;

INSERT INTO pedido (id, cliente_id, descricao)
SELECT g, 1 + (g % 25), 'Rua das Acacias 42'
FROM generate_series(1, 47) AS g;
INSERT INTO pedido (id, cliente_id, descricao) VALUES
    (48, 900, 'Joao Carlos Pereira'),
    (49, 900, 'Joao Carlos Pereira'),
    (50, 901, 'Joao Carlos Pereira');

INSERT INTO item_pedido (id, pedido_id, descricao)
SELECT g, 1 + (g % 25), 'servico de manutencao'
FROM generate_series(1, 100) AS g;

INSERT INTO conta (id, nome)
SELECT g, 'Conta Corrente' FROM generate_series(1, 300) AS g;

INSERT INTO lancamento (id, conta_id)
SELECT g, 1 + (g % 250) FROM generate_series(1, 380) AS g;
INSERT INTO lancamento (id, conta_id)
SELECT 380 + g, 900 + (g % 20) FROM generate_series(1, 20) AS g;

INSERT INTO contrato (id, objeto)
SELECT g, 'Construtora Horizonte LTDA' FROM generate_series(1, 40) AS g;

INSERT INTO fatura (id, contrato_id)
SELECT g, g FROM generate_series(1, 9) AS g;
INSERT INTO fatura (id, contrato_id)
SELECT 9 + g, 900 + g FROM generate_series(1, 21) AS g;

INSERT INTO equipe (id, nome)
SELECT g, 'Equipe de Campo' FROM generate_series(1, 20) AS g;

INSERT INTO chamado (id, equipe_id) VALUES
    (1, 1), (2, 2), (3, 3), (4, 4), (5, 5), (6, 6), (7, 7),
    (8, 900), (9, 901), (10, 902);

INSERT INTO moeda (id, codigo) VALUES (1, 'BRL'), (2, 'USD'), (3, 'EUR');

INSERT INTO nota (id, moeda_id)
SELECT g, 1 FROM generate_series(1, 20) AS g;

INSERT INTO fornecedor (id, nome)
SELECT g, 'Fornecedor Municipal' FROM generate_series(1, 5) AS g;

INSERT INTO "Unidade Gestora" (id, nome) VALUES (1, 'Secretaria de Obras'), (2, 'Secretaria de Saude');

INSERT INTO "Empenho 2024" (id, "unidade gestora") VALUES (1, 1), (2, 2), (3, 1);

ANALYZE;
