-- Fixture mínima: o pacote db testa a sessão, não o schema. Basta existir uma
-- tabela para que a tentativa de escrita tenha alvo.
CREATE TABLE municipio (
    id   bigint PRIMARY KEY,
    nome text NOT NULL
);

INSERT INTO municipio (id, nome) VALUES (1, 'Sao Bernardo do Campo');
