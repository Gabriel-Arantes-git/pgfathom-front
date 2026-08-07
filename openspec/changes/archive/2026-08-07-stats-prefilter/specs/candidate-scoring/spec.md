## ADDED Requirements

### Requirement: O score é recomponível por um único dono da regra de saturação

O sistema SHALL expor a recomposição de score a partir dos sinais de um candidato, e toda camada que acrescenta sinais após a geração MUST recompor o score por esse caminho único, nunca por implementação própria.

Duas implementações da saturação divergiriam mais cedo ou mais tarde, e o limiar passaria a significar coisas diferentes dependendo de qual camada tocou o candidato por último.

#### Scenario: Sinal acrescentado recompõe pelo mesmo mecanismo

- **WHEN** uma camada posterior à geração acrescenta um sinal a um candidato e recompõe o score
- **THEN** o resultado é idêntico ao que a geração produziria com o mesmo conjunto de sinais

#### Scenario: Score recomposto continua saturado

- **WHEN** sinais negativos acumulados levariam o score abaixo de zero
- **THEN** o score recomposto é zero, nunca negativo
