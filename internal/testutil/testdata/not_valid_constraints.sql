-- Cenário: constraints declaradas NOT VALID e nunca validadas, com órfãos
-- preexistentes. É o caso que engana toda ferramenta de diagrama: o \d mostra
-- a FK, a seta é desenhada, e as linhas antigas nunca foram verificadas.
--
-- Os valores plantados aqui são reconhecíveis de propósito. A varredura de
-- vazamento procura por eles em toda saída da ferramenta.

CREATE TABLE cliente (
    id       bigint PRIMARY KEY,
    cpf      text NOT NULL,
    nome     text NOT NULL
);

COMMENT ON TABLE cliente IS 'cadastro de clientes';
COMMENT ON COLUMN cliente.cpf IS 'documento do cliente';

CREATE TABLE pedido (
    id         bigint PRIMARY KEY,
    cliente_id bigint NOT NULL,
    valor      numeric(12,2) NOT NULL
);

CREATE TABLE item_pedido (
    id        bigint PRIMARY KEY,
    pedido_id bigint NOT NULL,
    descricao text
);

INSERT INTO cliente (id, cpf, nome) VALUES
    (1, '529.318.470-11', 'Maria Aparecida Silva'),
    (2, '145.892.663-04', 'Joao Carlos Pereira');

-- Órfãos plantados: cliente_id 98 e 99 não existem em cliente.
INSERT INTO pedido (id, cliente_id, valor) VALUES
    (1, 1, 150.00),
    (2, 2, 89.90),
    (3, 98, 12.00),
    (4, 99, 45.50);

INSERT INTO item_pedido (id, pedido_id, descricao) VALUES
    (1, 1, 'servico de manutencao'),
    (2, 1, 'peca de reposicao');

-- NOT VALID: aceita as linhas órfãs que já estavam lá.
ALTER TABLE pedido
    ADD CONSTRAINT pedido_cliente_fk
    FOREIGN KEY (cliente_id) REFERENCES cliente (id) NOT VALID;

-- Íntegra e validada, para provar que a ferramenta distingue as duas.
ALTER TABLE item_pedido
    ADD CONSTRAINT item_pedido_pedido_fk
    FOREIGN KEY (pedido_id) REFERENCES pedido (id);

CREATE INDEX ix_pedido_cliente_id ON pedido (cliente_id);
CREATE INDEX ix_item_pedido_pedido_id ON item_pedido (pedido_id);
