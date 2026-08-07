## Why

As duas primeiras fases construíram tudo de que a inferência precisa e nunca a executaram. Os perfis de nomenclatura traduzem convenção mas ninguém os consome; o catálogo é lido por inteiro mas só o que está declarado vira achado.

Esta fase liga as duas pontas: produzir hipóteses sobre relacionamentos que existem nos dados e não existem no catálogo, e pontuá-las bem o bastante para que a validação da fase 5 receba um conjunto pequeno e defensável.

Ela é determinística e não acessa banco, o que a devolve ao regime barato da fase 1 — tabela de casos, sem Docker, sem rede. Isso importa porque a pontuação é onde mora a maior parte do julgamento do projeto, e julgamento precisa ser fácil de revisar e de alterar.

O que a fase entrega ainda não é resposta. É a lista de perguntas que vale a pena fazer aos dados, com a justificativa de cada uma anexada.

## What Changes

- `internal/infer`: extração de nome de entidade a partir de nome de coluna, casamento contra as formas de nome de tabela, regras de compatibilidade de tipo, e pontuação por sinais com pesos positivos e negativos.
- Corte por limiar configurável, que é o mecanismo que impede a fase 5 de disparar milhares de anti-joins.
- Registro do motivo de descarte para todo candidato eliminado, para que o usuário não fique se perguntando por que uma coluna óbvia foi ignorada.
- Comando `pgfathom discover` em versão preliminar: gera e pontua candidatos, e os reporta **sem veredito**, marcados como não validados.
- Detecção do par polimórfico pelo nome das colunas vizinhas, para reportar o padrão em vez de rejeitá-lo em silêncio.

Nenhuma leitura de dado de tabela do usuário. Nenhum veredito. `discover` nesta fase responde "o que vale investigar", não "o que é verdade".

## Capabilities

### New Capabilities

- `candidate-generation`: como um par de colunas vira hipótese — extração de entidade, casamento de nome, compatibilidade de tipo, e as formas que precisam ser reconhecidas e puladas.
- `candidate-scoring`: os sinais, seus pesos, o corte por limiar, e a exigência de que todo score seja explicável a partir dos fatos que o produziram.

### Modified Capabilities

Nenhuma. `domain-model`, `naming-profiles` e `cli-foundation` são consumidas sem alteração de requisito.

## Impact

Nenhuma dependência nova. `internal/infer` importa apenas `internal/model` e `internal/profile`, ambos puros, o que mantém a suíte desta fase rodando em milissegundos.

Fixa duas coisas que a fase 5 herda e que ficam caras de mudar depois: o conjunto de sinais com seus pesos, que determina o que chega à validação, e o valor padrão do limiar de corte, que determina o custo de uma execução em banco de produção.

O comando `discover` nasce aqui incompleto de propósito. A alternativa seria escondê-lo até a fase 5, mas expô-lo agora permite rodar a inferência contra um banco real e medir a taxa de candidatos por tabela antes de existir qualquer anti-join — que é exatamente o número necessário para calibrar o limiar com honestidade.
