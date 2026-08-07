## 1. Detecção

- [x] 1.1 Criar `internal/profile/detect.go` com os tipos de resultado e evidência
- [x] 1.2 Implementar a derivação de sufixo e prefixo de referência a partir das FKs declaradas de coluna única
- [x] 1.3 Implementar a detecção de prefixo de tabela por frequência, exigindo separador
- [x] 1.4 Aplicar os limiares por proporção, com piso de ocorrências
- [x] 1.5 Implementar a mesclagem aditiva sobre o perfil base, preservando a ordem por comprimento
- [x] 1.6 Garantir que o perfil resultante passa na mesma `Validate` dos embarcados

## 2. Integração

- [x] 2.1 Ligar a detecção ao `discover`, ativa por padrão
- [x] 2.2 Adicionar a flag que desliga a detecção
- [x] 2.3 Reportar o detectado no terminal, com a evidência de cada item
- [x] 2.4 Incluir o detectado no JSON
- [x] 2.5 Informar quando nada foi detectado, ou quando a detecção está desligada

## 3. Verificação

- [x] 3.1 Testar a derivação de sufixo com fixture no padrão dos bancos medidos
- [x] 3.2 Testar a derivação de prefixo de referência
- [x] 3.3 Testar a detecção de prefixo de tabela no padrão Django
- [x] 3.4 Testar que assunto comum sem separador não vira prefixo
- [x] 3.5 Testar que ocorrência isolada não vira convenção
- [x] 3.6 Testar que schema sem FKs declaradas não detecta afixo de referência
- [x] 3.7 Testar que nenhuma regra do perfil base se perde na mesclagem
- [x] 3.8 Testar recall ponta a ponta com fixture que reproduz o caso `_idkey`
- [x] 3.9 Rodar `golangci-lint run` e zerar os apontamentos
- [x] 3.10 Rodar `openspec validate naming-detection` e corrigir o que apontar
