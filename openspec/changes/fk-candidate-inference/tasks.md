## 1. Compatibilidade de tipo

- [x] 1.1 Criar `internal/infer` com a função de compatibilidade entre tipo base filho e tipo base da chave alvo
- [x] 1.2 Implementar a família de inteiros com direção: menor para maior aceita, maior para menor rejeita
- [x] 1.3 Implementar a intercambialidade entre `text` e `varchar` de qualquer tamanho
- [x] 1.4 Implementar a estrita igualdade para `uuid`
- [x] 1.5 Garantir que numérico nunca casa com textual
- [x] 1.6 Distinguir tipo idêntico de apenas compatível no valor de retorno, porque os dois produzem sinais diferentes
- [x] 1.7 Escrever tabela de casos cobrindo cada regra e seus limites

## 2. Geração de candidatos

- [x] 2.1 Implementar a seleção de colunas elegíveis, excluindo chave primária da própria tabela e colunas com FK declarada
- [x] 2.2 Construir o índice de tabelas alvo por forma de nome, usando o conjunto de formas do perfil ativo
- [x] 2.3 Gerar candidato quando entidade, forma de tabela e compatibilidade de tipo coincidem
- [x] 2.4 Registrar a forma que casou e emitir sinal de nome exato ou normalizado conforme a origem
- [x] 2.5 Gerar um candidato por alvo quando o nome é ambíguo, marcando todos
- [x] 2.6 Pular alvo com chave composta ou sem chave, registrando o motivo
- [x] 2.7 Permitir auto-referência dentro da mesma tabela
- [x] 2.8 Garantir ordenação determinística da saída
- [x] 2.9 Escrever tabela de casos com fixtures de `model.Schema` montadas à mão

## 3. Detecção de par polimórfico

- [x] 3.1 Implementar a detecção de coluna irmã com mesmo prefixo e sufixo de tipo
- [x] 3.2 Tornar os sufixos de tipo parte do perfil de nomenclatura, não constante em código
- [x] 3.3 Registrar o padrão como observação sem suprimir a geração do candidato
- [x] 3.4 Escrever casos cobrindo detecção e não-detecção

## 4. Pontuação

- [x] 4.1 Definir os pesos como constantes nomeadas num único lugar
- [x] 4.2 Emitir os sinais positivos: nome exato, tipo idêntico, alvo único, coluna indexada, menção em comentário, `NOT NULL`
- [x] 4.3 Emitir os sinais negativos: nome genérico com alvo pequeno, alvo ambíguo, tipo apenas compatível
- [x] 4.4 Implementar a combinação com saturação em zero e um
- [x] 4.5 Garantir que o score é reconstruível a partir dos sinais registrados
- [x] 4.6 Implementar a detecção de nome genérico, penalizando apenas quando a tabela alvo é pequena
- [x] 4.7 Escrever testes de ordenação relativa entre pares de candidatos, não de valor absoluto
- [x] 4.8 Escrever teste provando saturação nos dois extremos

## 5. Corte e relato de descarte

- [x] 5.1 Implementar o corte por limiar configurável, com padrão declarado como estimativa
- [x] 5.2 Registrar motivo de descarte em todo candidato eliminado
- [x] 5.3 Preencher na cobertura a contagem de candidatos gerados e de sobreviventes
- [x] 5.4 Escrever teste provando que o limiar altera o conjunto sobrevivente

## 6. Comando discover

- [x] 6.1 Criar o comando `discover` com as flags de conexão, escopo, perfil, limiar e formato
- [x] 6.2 Encadear leitura de catálogo, geração e pontuação
- [x] 6.3 Marcar todo candidato com o veredito de não avaliado
- [x] 6.4 Afirmar na saída que os dados não foram consultados nesta versão
- [x] 6.5 Renderizar candidatos em terminal, ordenados por score decrescente, com os sinais visíveis
- [x] 6.6 Implementar a flag que mostra também os descartados, com score e motivo
- [x] 6.7 Estender a saída JSON com os candidatos
- [x] 6.8 Reportar o perfil de nomenclatura ativo no cabeçalho

## 7. Verificação

- [x] 7.1 Escrever o cenário-armadilha: coluna `status_id` numa base com tabela `status` onde a coluna guarda outra coisa, conferindo que a penalidade de nome genérico incide
- [x] 7.2 Escrever teste de determinismo rodando a geração repetidas vezes
- [x] 7.3 Escrever teste provando que coluna com FK declarada não gera candidato
- [ ] 7.4 Escrever fixture de integração com relação inferível e conferir a geração ponta a ponta
- [ ] 7.5 Varrer a saída de `discover` procurando valor de dado do usuário
- [x] 7.6 Confirmar que `internal/infer` importa apenas `internal/model` e `internal/profile`
- [x] 7.7 Rodar `golangci-lint run` e zerar os apontamentos
- [x] 7.8 Revisar a densidade de comentário, mantendo o padrão das fases anteriores
- [x] 7.9 Rodar `openspec validate fk-candidate-inference` e corrigir o que apontar
