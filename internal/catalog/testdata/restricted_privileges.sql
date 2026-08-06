-- Cenário: papel sem SELECT em parte das tabelas.
--
-- É o caso comum no ambiente alvo, não a exceção: conseguir privilégio em órgão
-- público é processo político, e uma concessão parcial é o que se encontra na
-- prática. A ferramenta precisa listar o que não conseguiu ler, senão um
-- relatório incompleto passa por limpo.

CREATE TABLE municipio (
    id   bigint PRIMARY KEY,
    nome text NOT NULL
);

CREATE TABLE folha_pagamento (
    id         bigint PRIMARY KEY,
    servidor   text NOT NULL,
    liquido    numeric(12,2) NOT NULL
);

CREATE TABLE prontuario (
    id       bigint PRIMARY KEY,
    paciente text NOT NULL,
    cid      text
);

INSERT INTO municipio (id, nome) VALUES (1, 'Sao Bernardo do Campo');
INSERT INTO folha_pagamento (id, servidor, liquido) VALUES (1, 'Maria Aparecida Silva', 4210.55);
INSERT INTO prontuario (id, paciente, cid) VALUES (1, 'Joao Carlos Pereira', 'E11');

-- O papel de análise enxerga o cadastro e não enxerga o que é sensível.
CREATE ROLE pgfathom_reader LOGIN PASSWORD 'pgfathom_reader';
GRANT USAGE ON SCHEMA public TO pgfathom_reader;
GRANT SELECT ON municipio TO pgfathom_reader;
REVOKE ALL ON folha_pagamento FROM pgfathom_reader;
REVOKE ALL ON prontuario FROM pgfathom_reader;
