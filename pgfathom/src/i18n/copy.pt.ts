import type { Copy } from './types'

export const pt: Copy = {
  docTitle: 'pgfathom — sonde a profundidade de um schema PostgreSQL legado',
  docDescription:
    'O pgfathom encontra os relacionamentos que seu banco PostgreSQL tem mas nunca declarou, e os prova contra os dados em vez de adivinhar pelos nomes das colunas.',

  navProblem: 'Problema',
  navVerdicts: 'O que faz',
  navSafety: 'Segurança',
  navHow: 'Como funciona',
  navBenchmark: 'Benchmark',
  navMenu: 'Abrir menu',
  navClose: 'Fechar menu',

  langSwitchLabel: 'Idioma',
  toEnglish: 'Switch to English',
  toPortuguese: 'Mudar para português',

  heroTitle: 'Sonde a profundidade de um schema PostgreSQL legado.',
  heroLead:
    'O pgfathom encontra os relacionamentos que seu banco tem mas nunca declarou — e os prova contra os dados, em vez de adivinhar pelos nomes das colunas.',
  ctaDesign: 'Ler o design',
  copyCommand: 'Copiar comando',
  copiedCommand: 'Copiado',
  noticeTitle: 'Pré-release. Em desenvolvimento ativo.',
  noticeBody:
    'O que roda hoje: o pgfathom audit ponta a ponta, e o pgfathom discover passando por geração de candidatos, pontuação, pré-filtro estatístico e validação contra os dados. A mineração de joins e os formatos finais de saída ainda estão pela frente. A saída de terminal mostrada aqui é o design alvo, não uma gravação. Os benchmarks de taxa de recuperação serão publicados quando a ferramenta rodar contra o corpus de referência — nenhum número é afirmado até lá.',
  facts: [
    { label: 'Somente leitura', detail: 'Sem modo de escrita, sob nenhuma flag, em nenhuma fase' },
    { label: 'Fica em memória', detail: 'Saem contagens e nomes — nunca um valor das suas tabelas' },
    { label: 'Quatro dependências', detail: 'Sem cgo. O go.mod é uma leitura curta' },
    { label: 'Apache-2.0', detail: 'Escolhida pela concessão explícita de patente' },
  ],
  terminalCaption: 'pgfathom — saída alvo',
  graph: {
    pedido: { name: 'pedido', columns: ['id', 'cliente_id', 'emitido_em'] },
    cliente: { name: 'cliente', columns: ['id', 'razao_social', 'ativo'] },
    item_pedido: { name: 'item_pedido', columns: ['id', 'pedido_id', 'quantidade'] },
    os_servico: { name: 'os_servico', columns: ['id', 'resp_tecnico', 'aberta_em'] },
    funcionario: { name: 'funcionario', columns: ['id', 'matricula'] },
    endereco: { name: 'endereco', columns: ['id', 'municipio_id', 'cep'] },
    municipio: { name: 'municipio', columns: ['id', 'nome'] },
    conta: { name: 'conta', columns: ['id', 'codigo'] },
    lancamento: { name: 'lancamento', columns: ['id', 'conta_id', 'valor'] },
    produto: { name: 'produto', columns: ['id', 'codigo', 'preco_unit'] },
    nota_fiscal: { name: 'nota_fiscal', columns: ['id', 'numero'] },
    usuario: { name: 'usuario', columns: ['id', 'login', 'criado_em'] },
    filial: { name: 'filial', columns: ['id', 'cnpj'] },
    movimento: { name: 'movimento', columns: ['id', 'tipo', 'movido_em'] },
    sessao: { name: 'sessao', columns: ['id', 'usuario_id', 'expira_em'] },
    log_evento: { name: 'log_evento', columns: ['id', 'ator_id'] },
  },

  problemEyebrow: 'O problema',
  problemTitle: 'Bancos PostgreSQL antigos carregam mais estrutura do que declaram.',
  problemP1:
    'Nenhuma foreign key. Mas cliente_id aponta para cliente.id em cada uma das linhas — o ORM da época nunca criou a constraint, ou alguém removeu as constraints para uma carga em massa e nunca as recolocou.',
  problemP2:
    'Duas coisas decorrem disso. Ninguém consegue ler o modelo, porque o \\d não mostra nada e sua ferramenta de ERD desenha uma página de caixas desconectadas. E como a constraint nunca esteve lá, nada jamais impediu a entrada de linhas órfãs — então elas provavelmente já entraram, anos atrás, em silêncio, e ninguém foi olhar.',
  notValidLabel: 'A variante mais cruel',
  notValidBody:
    'Uma foreign key pode estar declarada e ainda assim não garantir nada, se foi criada NOT VALID e nunca validada. Ela aparece no \\d. Ela desenha uma seta no seu diagrama. Ela nunca verificou uma única linha preexistente.',

  verdictsEyebrow: 'O que faz',
  verdictsTitle: 'Três veredictos, três respostas diferentes.',
  verdictsLead:
    'Cada relacionamento inferido carrega um veredicto e a métrica por trás dele. Candidatos fracos e rejeitados também são reportados, para você nunca ficar se perguntando por que uma coluna óbvia foi ignorada.',
  verdicts: [
    {
      tag: 'BROKEN',
      headline: 'O ponto da ferramenta.',
      body: 'Um bug de dados que está em produção há anos e ninguém sabe.',
      output:
        '→ a query que lista os órfãos, porque eles têm que ser resolvidos antes que qualquer constraint possa ser adicionada.',
    },
    {
      tag: 'CONFIRMED',
      headline: 'Uma foreign key que alguém esqueceu de declarar.',
      body: 'A contenção se sustenta por linha e por valor distinto, sem órfãos encontrados.',
      output:
        '→ o DDL, com o VALIDATE CONSTRAINT separado para que o ALTER TABLE inicial não segure um lock pesado — mais o CREATE INDEX CONCURRENTLY quando a coluna filha não tem índice.',
    },
    {
      tag: 'WEAK',
      headline: 'Evidência insuficiente para concluir.',
      body: 'Reportado em vez de descartado, com o motivo anexado.',
      output:
        '→ a métrica que ficou aquém e qualquer padrão detectado — um par polimórfico, um alvo ambíguo.',
    },
  ],
  joinMining:
    'Nenhuma heurística de correspondência de nomes no mundo encontra essa — resp_tecnico não se parece em nada com funcionario. O pgfathom a encontra lendo os predicados de join nas definições das suas próprias views e funções.',

  safetyEyebrow: 'Segurança',
  safetyTitle:
    'Feito para ser apontado para um banco de produção de alguém que está nervoso com isso.',
  safetyLead: 'Cada garantia abaixo é um requisito rígido na especificação, não uma meta.',
  safety: [
    {
      name: 'Somente leitura, estruturalmente',
      detail:
        'O pgfathom nunca emite uma instrução que modifique o banco sob análise — não há modo de escrita, sob nenhuma flag, em nenhuma fase. A sessão define default_transaction_read_only, e uma role somente-leitura é a configuração recomendada. Ele emite arquivos .sql para você revisar e executar.',
    },
    {
      name: 'Seus dados nunca saem da memória',
      detail:
        'O que sai são contagens, proporções e nomes de objetos — nunca um valor das suas tabelas, em nenhuma saída, log, campo JSON ou mensagem de erro. Garantido por um teste que serializa cada estrutura e varre o resultado, não por revisão de código.',
    },
    {
      name: 'Não vai derrubar seu servidor',
      detail:
        'Toda query de validação roda sob statement_timeout, lock_timeout e idle_in_transaction_session_timeout. A concorrência é limitada e o padrão é baixo. A conexão se anuncia como pgfathom no pg_stat_activity.',
    },
    {
      name: 'Nenhuma afirmação sem evidência',
      detail:
        'Execuções amostradas nunca podem reportar um relacionamento confirmado — só o --full pode provar a ausência de órfãos.',
    },
    {
      name: 'Silêncio nunca é atestado de saúde',
      detail:
        'Tabelas puladas por falta de privilégio, candidatos que estouraram o tempo, schemas não cobertos — tudo isso aparece no bloco de cobertura em toda execução. Um relatório limpo significa "eu olhei e está limpo", nunca "eu não consegui olhar".',
    },
  ],

  howEyebrow: 'Como funciona',
  howTitle: 'Seis passagens, do catálogo à validação.',
  stages: [
    {
      num: '01',
      name: 'Ler o catálogo',
      detail:
        'Tabelas, colunas, chaves, índices, comentários, foreign keys declaradas e seu estado de validação, estatísticas de uso com o timestamp de reset.',
    },
    {
      num: '02',
      name: 'Minerar evidência de uso',
      detail:
        'Predicados de join extraídos de definições de views, corpos de funções e do pg_stat_statements. Uma view que junta duas colunas é prova de que seu código as trata como relacionadas. Puro catálogo: sem dados de usuário, sem custo.',
    },
    {
      num: '03',
      name: 'Gerar candidatos',
      detail:
        'Afixos de nomes de coluna são removidos e comparados com nomes de tabela no singular usando um perfil de nomenclatura — um arquivo de configuração, não regras fixas no código.',
    },
    {
      num: '04',
      name: 'Pontuar só com metadados',
      detail:
        'Correspondência exata de nome, identidade de tipo, ambiguidade do alvo, índice existente, menções em comentários. Candidatos fracos caem antes de qualquer coisa tocar os dados.',
    },
    {
      num: '05',
      name: 'Pré-filtrar por estatísticas',
      detail:
        'Se a coluna filha tem mais valores distintos do que a tabela pai tem linhas, a contenção total é aritmeticamente impossível. De graça, do pg_stats, sem I/O.',
    },
    {
      num: '06',
      name: 'Validar contra os dados',
      detail:
        'Um agregado por candidato sobrevivente — nunca buscando linhas, só contagens. A contenção é reportada por linha e por valor distinto, porque um valor ruim repetido um milhão de vezes e um milhão de valores ruins raros são problemas diferentes.',
    },
  ],

  profilesEyebrow: 'Perfis de nomenclatura',
  profilesTitle:
    'A maioria das ferramentas de schema assume inglês. Os bancos que mais precisam desta ferramenta, frequentemente, não são.',
  profilesP1:
    'Regras de afixo e plural vivem em TOML, não em Go, então ensinar uma nova convenção ao pgfathom é um arquivo de configuração, não um patch.',
  profilesP2:
    'A normalização retorna um conjunto de formas candidatas em vez de uma só, então plurais ambíguos não custam nada em recall — toda forma é tentada, e a que casou é reportada. Adicionar um perfil para o seu idioma é a contribuição mais fácil possível: um arquivo TOML e uma tabela de casos de teste.',
  profilesYours: 'seu idioma',

  ciEyebrow: 'Integração contínua',
  ciTitle: 'Quebre o build quando o schema derivar.',
  ciBadge: 'planejado · depois da v0.1',
  ciP1:
    'O pgfathom check --baseline compara uma execução com um modelo versionado e sai com código diferente de zero quando o schema mudou. Ainda não disponível — esta é a forma que está especificada.',
  ciRules: [
    { code: 'exit 1', detail: 'Um novo relacionamento não declarado apareceu desde o baseline' },
    { code: 'exit 1', detail: 'A contagem de órfãos cresceu em um relacionamento já quebrado' },
    { code: 'exit 0', detail: 'Cobertura reportada, nada mudou' },
  ],

  priorEyebrow: 'Trabalhos anteriores',
  priorTitle: 'Não é ciência nova, e diz isso.',
  priorLead:
    'Contenção é conhecida na literatura de data profiling como inclusion dependency — a parte automaticamente testável de uma foreign key. O pgfathom deliberadamente não compete com ferramentas que já resolveram bem os seus problemas.',
  colTool: 'Ferramenta',
  colDoes: 'O que faz bem',
  colOverlap: 'Onde o pgfathom difere',
  priorRows: [
    {
      tool: 'Atlas',
      does: 'Drift de schema e gestão de migrações',
      overlap:
        'Compara com um estado desejado declarado; o pgfathom infere o estado que nunca foi declarado',
    },
    {
      tool: 'SchemaSpy · Azimutt',
      does: 'Diagramas e exploração de schema',
      overlap:
        'Desenham o que o catálogo contém; o Azimutt sinaliza colunas _id, mas nada é validado contra as linhas',
    },
    {
      tool: 'Squawk',
      does: 'Lint de migrações',
      overlap: 'Revisa o DDL que você escreve; o pgfathom escreve o DDL que você revisa',
    },
    {
      tool: 'Metanome (SPIDER, BINDER, MIND)',
      does: 'Pesquisa em descoberta de inclusion dependencies',
      overlap:
        'Acadêmico, genérico, offline; o pgfathom é um CLI nativo de PostgreSQL que minera o próprio catálogo',
    },
    {
      tool: 'pgfathom',
      does: 'Valida suas inferências contra os dados reais e entrega DDL revisável',
      overlap: 'Fala convenções de nomenclatura fora do inglês e reporta a própria cobertura honestamente',
      highlight: true,
    },
  ],
  priorNote:
    'O modelo JSON versionado é o ponto de integração: consuma-o e gere o que quiser a partir de um schema que finalmente conhece os próprios relacionamentos. Geração de código explicitamente não está no roadmap.',

  benchEyebrow: 'Como a corretude é medida',
  benchTitle: 'Remova todas as foreign keys e conte quantas voltam.',
  benchBadge: 'protótipo · números de exemplo',
  benchLead:
    'O corpus é público — GitLab, Odoo, Discourse, Redmine, Mastodon — para qualquer um poder reproduzir a execução. Os resultados são divididos entre o que a correspondência de nomes recupera sozinha e o que a evidência de uso soma em cima, porque essa diferença é exatamente o que a mineração de joins existe para fechar.',
  colSchema: 'Schema',
  colRecovered: 'Foreign keys recuperadas',
  benchLegend: {
    byName: 'Correspondência de nomes',
    byEvidence: 'Evidência de uso',
  },
  benchAggregate: 'média do corpus',
  metricLabel: 'Taxa de recuperação',
  metricBody:
    'Pegue um schema com foreign keys completas, remova todas elas, rode o pgfathom e conte quantas voltam. Contra um corpus público — GitLab, Odoo, Discourse, Redmine, Mastodon — para qualquer um poder reproduzir.',
  fpLabel: 'Zero falsos positivos confirmados',
  fpBody:
    'O recall vai ficar bem abaixo de 100%, e isso é esperado. A métrica sem tolerância é a outra: um relacionamento perdido custa um achado, um errado confirmado custa a ferramenta.',

  ctaTitle: 'O design é o momento mais barato para mudá-lo.',
  ctaBody:
    'Se você já topou com esse problema em um banco legado real, abra uma issue: como era o schema, que convenção de nomes ele usava, e o que uma ferramenta precisaria ter encontrado. As duas contribuições mais valiosas são perfis de nomenclatura para outros idiomas e schemas reais para o corpus de benchmark.',
  ctaIssue: 'Abrir uma issue',
  ctaRepo: 'Ver o repositório',

  tagline: 'Arqueologia de schema PostgreSQL',
  linkDesign: 'Documento de design',
  linkRoadmap: 'Roadmap',
  footLeft: 'Apache-2.0 · pré-release',
  footRight: 'Sem afiliação com o PostgreSQL Global Development Group.',
}
