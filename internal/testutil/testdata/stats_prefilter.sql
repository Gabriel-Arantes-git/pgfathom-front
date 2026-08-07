-- Cenário do pré-filtro estatístico: candidatos que o casamento de nome gera
-- com convicção e que as estatísticas do planner desmontam — ou não podem
-- avaliar. Todas as relações abaixo são inferíveis por nome; a diferença está
-- nos dados.
--
-- Os valores plantados são reconhecíveis de propósito, para a varredura de
-- vazamento. Os ids fora de faixa (1001+ e 9001+) fazem o papel dos valores
-- reais que um histogram_bounds carregaria.

-- Alvo minúsculo do candidato impossível.
CREATE TABLE unidade (
    id   bigint PRIMARY KEY,
    nome text NOT NULL
);

-- 300 valores distintos contra 3 linhas em unidade: contenção total é
-- aritmeticamente impossível, e a violação passa longe da margem de 2x.
CREATE TABLE leitura (
    id         bigint PRIMARY KEY,
    unidade_id bigint NOT NULL,
    observacao text
);

-- Par da checagem de faixa: cardinalidade passa, histograma inteiro deslocado.
CREATE TABLE ponto (
    id   bigint PRIMARY KEY,
    nome text NOT NULL
);

CREATE TABLE medicao (
    id       bigint PRIMARY KEY,
    ponto_id bigint NOT NULL
);

-- Relação legítima: contida, e deve atravessar sem nenhum sinal estatístico.
CREATE TABLE cliente (
    id   bigint PRIMARY KEY,
    cpf  text NOT NULL,
    nome text NOT NULL
);

CREATE TABLE pedido (
    id         bigint PRIMARY KEY,
    cliente_id bigint NOT NULL
);

-- Deliberadamente sem ANALYZE: o pré-filtro não pode opinar, e o silêncio
-- precisa ficar registrado. Autovacuum desligado para a janela do teste não
-- depender de sorte.
CREATE TABLE sem_estatistica (
    id         bigint PRIMARY KEY,
    unidade_id bigint NOT NULL
);
ALTER TABLE sem_estatistica SET (autovacuum_enabled = false);

INSERT INTO unidade (id, nome) VALUES
    (1, 'Matriz Central'),
    (2, 'Filial Leste'),
    (3, 'Filial Oeste');

INSERT INTO leitura (id, unidade_id, observacao)
SELECT g, 1000 + g, 'Leitura Manual Bloco C'
FROM generate_series(1, 300) AS g;

INSERT INTO ponto (id, nome)
SELECT g, 'Ponto de Coleta'
FROM generate_series(1, 500) AS g;

INSERT INTO medicao (id, ponto_id)
SELECT g, 9000 + g
FROM generate_series(1, 300) AS g;

INSERT INTO cliente (id, cpf, nome) VALUES
    (1, '529.318.470-11', 'Maria Aparecida Silva'),
    (2, '145.892.663-04', 'Joao Carlos Pereira');

INSERT INTO pedido (id, cliente_id)
SELECT g, 1 + (g % 2)
FROM generate_series(1, 50) AS g;

INSERT INTO sem_estatistica (id, unidade_id) VALUES (1, 1), (2, 2);

-- ANALYZE explícito e por tabela, porque um ANALYZE sem argumento cobriria
-- sem_estatistica e destruiria o cenário de ausência.
ANALYZE unidade;
ANALYZE leitura;
ANALYZE ponto;
ANALYZE medicao;
ANALYZE cliente;
ANALYZE pedido;
