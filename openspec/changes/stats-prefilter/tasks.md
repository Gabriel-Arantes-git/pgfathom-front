## 1. Modelo e pontuação

- [ ] 1.1 Adicionar o sinal de estatística indisponível ao conjunto fechado de `SignalKind`, com peso zero
- [ ] 1.2 Adicionar à `Coverage` os contadores do funil: avaliados pelo pré-filtro, rejeitados por estatística, sem estatística, e se a camada rodou
- [ ] 1.3 Exportar em `internal/infer` a recomposição de score a partir dos sinais, e fazer a geração usá-la internamente para que exista um único caminho
- [ ] 1.4 Teste provando que recompor com os mesmos sinais reproduz o score da geração

## 2. Leitura de pg_stats

- [ ] 2.1 Criar `internal/stats` com a consulta dirigida: `null_frac` e `n_distinct` das colunas dos candidatos, numa única ida ao servidor
- [ ] 2.2 Ler os limites do histograma como `float8`, convertidos no servidor, apenas para colunas da família numérica
- [ ] 2.3 Guardar limites em struct não exportada, sem tag de serialização, morrendo com a avaliação
- [ ] 2.4 Teste que falha se alguma consulta da camada referenciar relação fora de `pg_catalog` e visões de estatística

## 3. Avaliação

- [ ] 3.1 Separar a avaliação pura da leitura: a avaliação recebe candidatos e estatísticas e devolve mantidos, rejeitados e não avaliáveis, sem tocar em banco
- [ ] 3.2 Checagem de cardinalidade via `EstimatedDistinct` contra `EstimatedRowCount` do pai: penalidade dentro da margem, rejeição com motivo além dela
- [ ] 3.3 Checagem de faixa numérica: sinal negativo quando os limites da filha saem dos do pai; nunca rejeita sozinha
- [ ] 3.4 Estatística ausente ou não interpretável: sinal de peso zero no candidato, contagem na cobertura, score intacto
- [ ] 3.5 Motivo de rejeição com contagens declaradas como estimativas e nomes de objeto, nunca valores
- [ ] 3.6 Tabela de casos da avaliação pura: violação dentro e além da margem, faixa deslocada, faixa não numérica, estatística ausente, `n_distinct` negativo

## 4. Integração no discover

- [ ] 4.1 Flag `--no-stats` removendo a camada da execução
- [ ] 4.2 Aplicar o pré-filtro após o corte por limiar, movendo rejeitados para os descartados reportáveis
- [ ] 4.3 Renderizar o funil no bloco de cobertura do terminal, incluindo a declaração de camada desligada
- [ ] 4.4 Estender o JSON com os campos aditivos da cobertura

## 5. Verificação

- [ ] 5.1 Fixture `stats_prefilter` com candidato aritmeticamente impossível, candidato de faixa deslocada, candidato legítimo e tabela sem `ANALYZE`, com valores plantados reconhecíveis e `ANALYZE` explícito nas demais
- [ ] 5.2 Teste de integração medindo a redução: o impossível é rejeitado, o legítimo sobrevive, o sem estatística atravessa com registro
- [ ] 5.3 Teste de vazamento serializando o resultado inteiro da camada e a saída de `discover` com e sem `--no-stats`, varrendo contra os valores plantados
- [ ] 5.4 Teste provando que `--no-stats` produz candidatos sem nenhum sinal estatístico e cobertura declarando a camada desligada
- [ ] 5.5 Rodar `golangci-lint run` e zerar os apontamentos
- [ ] 5.6 Confirmar que `go test ./...` segue sem Docker e sem rede
- [ ] 5.7 Revisar a densidade de comentário, mantendo o padrão das fases anteriores
- [ ] 5.8 Rodar `openspec validate stats-prefilter` e corrigir o que apontar
