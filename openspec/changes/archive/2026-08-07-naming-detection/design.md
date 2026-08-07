## Context

A fase 3 entregou inferência que funciona e recupera 0,5% num banco real. O problema não é o algoritmo, é que a convenção de nomenclatura daquele schema não estava em nenhuma lista.

Três medições, três sistemas diferentes, mesmo padrão: a convenção pertence ao schema, não ao idioma. E o modo de falha é silencioso, que é o pior tipo para este projeto — a ferramenta não erra, ela devolve pouco, e o usuário conclui que o banco não tinha o que descobrir.

## Goals / Non-Goals

**Goals**

Derivar a convenção do próprio schema. Compor com o perfil base sem substituí-lo. Reportar o que foi detectado e a evidência de cada item. Nenhuma consulta nova.

**Non-Goals**

Detecção de idioma. Detecção de regra de plural — plural é morfologia, não convenção local, e continua vindo do perfil. Substituir o perfil: a detecção acrescenta, o perfil continua sendo a base.

## Decisions

### O afixo sai das chaves declaradas, não de uma lista maior

A alternativa seria continuar ampliando as listas dos perfis oficiais a cada schema novo que aparecesse. Isso não escala e nunca termina: `_idkey` foi descoberto porque alguém rodou e mediu, e existem tantas variações quanto fornecedores de sistema.

Um banco com chaves declaradas está afirmando qual é a sua convenção. Comparar o nome da coluna filha com as formas do nome da tabela alvo e olhar o resíduo é leitura direta dessa afirmação, sem palpite.

O limite é conhecido: schema sem nenhuma chave declarada não tem o que ler. Aí só o perfil base vale, que é exatamente o comportamento de hoje — a detecção nunca piora o caso.

### Prefixo de tabela exige separador

Contar prefixos por frequência bruta detectaria `ato` num schema com `atotipo`, `atoanexo`, `atotramite`. Isso não é convenção, é assunto comum, e removê-lo produziria formas sem sentido.

Exigir que o prefixo termine em separador separa os dois casos limpo: `auth_`, `django_`, `tpl_`, `tb_` são convenção; `ato` não é. O custo é não detectar convenção sem separador, que é rara e fica com o perfil base.

### Aditivo, nunca substitutivo

A decisão que torna a detecção segura. `TableForms` já devolve conjunto, então um afixo detectado errado acrescenta uma forma candidata — ruído que a pontuação penaliza e a validação derruba.

Remover ou sobrepor uma regra do perfil seria o oposto: custaria um candidato que nunca nasce, e candidato que nunca nasce é invisível. A mesma assimetria que decidiu gerar com generosidade na fase 3 decide aqui.

### Ligada por padrão

Manter desligada por padrão preservaria o comportamento atual e não resolveria nada: quem precisa da detecção é justamente quem não sabe que precisa.

O risco de ligar é a convenção mudar sem aviso entre versões. Por isso o relato do detectado é requisito, com a evidência de cada item — o usuário precisa poder reproduzir e contestar.

### Limiares por proporção, não por contagem

Contagem absoluta se comportaria de forma oposta nos dois extremos: três ocorrências num schema de dezoito tabelas é convenção, e as mesmas três em oitocentas é ruído.

Os limiares são proporção da população relevante, com um piso pequeno de ocorrências para evitar que schema minúsculo transforme um caso isolado em regra.

## Risks / Trade-offs

**Detecção pode aprender uma convenção errada de um schema inconsistente** → o resultado é aditivo, então o pior caso é forma espúria. E o relato expõe o que foi aprendido, com quantas ocorrências.

**Afixo detectado pode sombrear um do perfil base e quebrar a validação de ordem** → a mesclagem insere respeitando a ordem por comprimento, e o perfil resultante passa pela mesma `Validate` dos perfis embarcados.

**Ligar por padrão muda resultado entre versões** → assumido. A alternativa é manter todo mundo em 0,5% para preservar reprodutibilidade de um número que não serve para nada.

**Schema sem chaves declaradas não ganha nada** → é o caso onde a inferência mais importa e onde a detecção menos ajuda. Mitigado em parte pela detecção de prefixo de tabela, que não depende de chaves.

## Migration Plan

Não aplicável.

## Open Questions

Os limiares entram como estimativa, calibrados contra os schemas medidos, e são revistos na fase 8 com o corpus.
