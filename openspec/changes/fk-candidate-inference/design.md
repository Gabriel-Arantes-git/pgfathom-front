## Context

As fases 1 e 2 entregaram tudo de que a inferência precisa. Os perfis traduzem convenção de nomenclatura e ninguém os consome; o catálogo é lido por inteiro e só o que está declarado vira achado.

Esta camada é determinística e não acessa banco, o que a devolve ao regime de teste da fase 1: tabela de casos, milissegundos, sem Docker. Isso não é conveniência — a pontuação concentra quase todo o julgamento do projeto, e julgamento precisa ser barato de revisar e de mudar de ideia.

Restrição herdada que domina o desenho: **falso negativo é aceitável, falso positivo confirmado não é.** Aqui isso se traduz em gerar candidato com generosidade e cortar com rigor, porque o corte tem recurso — o usuário baixa o limiar — enquanto o candidato que nunca nasceu é invisível.

## Goals / Non-Goals

**Goals**

Produzir hipóteses a partir de nome e tipo. Pontuá-las de modo que o score seja sempre explicável pelos sinais. Cortar por limiar configurável, que é o mecanismo que protege a fase 5 de custar caro. Registrar todo descarte com motivo.

**Non-Goals**

Nenhuma leitura de dado. Nenhum veredito — `discover` nesta fase responde "o que investigar", não "o que é verdade". Nenhum sinal vindo de view ou de função: isso é a fase 6, e a ordem é deliberada. Chave composta segue fora de escopo, agora com registro. Pré-filtro estatístico é a fase 4.

## Decisions

### Gerar com generosidade, cortar com rigor

Cada etapa poderia filtrar mais cedo. Compatibilidade de tipo poderia rejeitar antes de tentar o casamento de nome; ambiguidade poderia escolher o alvo mais provável em vez de gerar vários.

A escolha é o contrário: gerar tudo que passa nos filtros duros — tipo compatível, alvo com chave única — e deixar a pontuação decidir. A razão é assimetria de recurso. Um candidato cortado pelo limiar aparece na lista de descartados com o motivo, e o usuário que discorda baixa o limiar. Um candidato que nunca foi gerado não existe em lugar nenhum, e ninguém pode discordar do que não vê.

O custo é uma lista maior em memória, o que para centenas de tabelas é irrelevante.

### Ambiguidade gera todos os candidatos

Quando um nome casa com várias tabelas, a alternativa seria escolher por heurística — a maior, a do mesmo schema, a com mais colunas. Toda heurística aqui seria um palpite disfarçado de decisão, e esconderia do usuário que havia incerteza.

Gerar todos e penalizar com o sinal de alvo ambíguo mantém a incerteza visível. A validação da fase 5 resolve com dados, que é a única coisa que pode resolver isso de verdade.

### Score normalizado por saturação, não por soma crua

Somar pesos livremente faria o intervalo depender do número de sinais, e o limiar deixaria de significar a mesma coisa entre um schema com comentários ricos e outro sem nenhum.

A pontuação combina os pesos e satura em zero e um. O limiar passa a ser um número que o usuário consegue raciocinar — "0.5 significa metade da confiança possível" — em vez de um valor sem escala.

### Nome genérico penalizado, não excluído

`status_id` apontando para `status` quase sempre é uma relação real. Também quase sempre é a relação menos interessante do banco, e num schema grande esse padrão aparece dezenas de vezes.

Excluir seria errado: é relação de verdade, e algumas dessas tabelas de domínio têm órfãos. Não penalizar também seria errado: por volume, elas empurrariam para o fim do relatório os achados que motivam a existência da ferramenta.

A penalidade incide quando o nome é genérico **e** a tabela alvo é pequena. Nome genérico apontando para tabela grande não é tabela de domínio, e não deve ser penalizado.

### Par polimórfico reconhecido antes de virar rejeição

`documento_id` só faz sentido junto de `documento_tipo`. Sem tratamento, a validação encontra contenção baixa e rejeita — resultado correto pelo caminho errado, porque o usuário lê "coincidência de nome" onde existe uma relação real que esta versão não modela.

A detecção é sintática e barata: coluna irmã na mesma tabela com o mesmo prefixo e sufixo de tipo. Ela não resolve o relacionamento, apenas o nomeia. Transformar um falso descarte numa observação é o que separa uma ferramenta que parece burra de uma que admite o próprio limite.

### `discover` nasce sem veredito, e diz isso

A alternativa seria esconder o comando até a fase 5. Expô-lo agora permite rodar contra um banco real e medir candidatos por tabela **antes** de existir qualquer anti-join — que é exatamente o número necessário para calibrar o limiar padrão com honestidade em vez de por chute.

O risco é o usuário tomar candidato bem pontuado por relação confirmada e criar uma constraint a partir de casamento de nome, que é precisamente o erro que a ferramenta existe para evitar. Por isso todo candidato sai com veredito de não avaliado e a saída afirma que os dados não foram consultados.

### Sem sinal de uso nesta fase

Minerar junção de view é o que quebra o teto do casamento de nome, e é tentador antecipar.

Fica para a fase 6 porque essa é a única ordem em que o ganho é mensurável. Com as fases 1 a 5 fechadas existe uma linha de base de recall usando só nome; ligar o probe e medir de novo produz a decomposição que o README promete. Implementado agora, esse número não existiria e a afirmação principal do projeto ficaria sem prova.

## Risks / Trade-offs

**Os pesos são julgamento, e o primeiro conjunto vai estar errado** → por isso ficam em constantes nomeadas num único lugar, com teste de ordenação relativa em vez de teste de valor absoluto. O que precisa se manter verdadeiro é "exato pontua mais que normalizado", não "exato vale 0.4".

**O limiar padrão só pode ser calibrado com dados reais** → a contagem de candidatos entra na cobertura desde já, e o valor padrão é revisto na fase 8 com o corpus de benchmark. Até lá é uma estimativa declarada como tal.

**Gerar com generosidade pode explodir em schema com muitas colunas genéricas** → o corte por limiar é o mecanismo de contenção, e a contagem na cobertura é o alarme. Se a razão candidatos por tabela ficar alta no corpus, a resposta é ajustar pesos, não filtrar mais cedo.

**Detecção polimórfica por sufixo pode ter falso positivo** — uma tabela com `documento_id` e `documento_tipo` onde as duas colunas são independentes → o padrão só é registrado como observação, nunca suprime a geração do candidato. O pior caso é uma nota a mais no relatório.

## Migration Plan

Não aplicável.

## Open Questions

Nenhuma bloqueante. O valor padrão do limiar entra como estimativa e é revisto na fase 8.
