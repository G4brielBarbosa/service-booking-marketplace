# Feature Specification: Backlog Inteligente e Priorização baseada em lacunas observadas

**Created**: 2026-02-19  
**PRD Base**: §5.2, §5.5, §9.1, §10 (R4), §6.2, §11 (RNF1–RNF3), §13

## Caso de uso *(mandatory)*

O usuário quer que o sistema mantenha uma lista de “próximos passos” que não dependa apenas de motivação momentânea, mas sim de **lacunas observáveis** (ex.: erros recorrentes, queda de rubrica, baixa consistência, gargalos de tempo/energia), para aumentar progresso real com baixa fricção.

Esta feature define como o sistema:
- cria itens de backlog a partir de sinais observáveis,
- prioriza e limita a quantidade para evitar overload,
- explica o “por quê” de cada sugestão,
- e disponibiliza itens para alimentar planejamento diário/semanal (sem definir implementação).

Referências relacionadas:
- `SPEC-016` (métricas/registros que geram sinais)
- `SPEC-003` (quality gates/evidências e estados de conclusão)
- `SPEC-009` (detecção de falso progresso → gera lacunas)
- `SPEC-010` (governança/limites de metas intensivas)
- `SPEC-002` (rotina diária consome recomendações sem “spam”)

## Scope & Non-goals *(recommended)*

**In scope (MVP slice)**:
- Backlog com itens derivados de lacunas observadas + prioridade + limite de itens ativos.
- Explicação curta e acionável para cada item (“por que isso agora”).
- Ações do usuário: aceitar, adiar, rejeitar (com motivo opcional), e ver backlog atual.

**Non-goals (agora)**:
- Não otimizar com modelos/IA avançada; pode começar com regras simples e observáveis.
- Não virar um “gerenciador de tarefas genérico” (itens não são para qualquer coisa aleatória).
- Não desenhar UI/kanban/dashboards.

## Definições *(recommended)*

- **Lacuna observada**: sinal a partir de registros/evidências (ex.: erro recorrente ≥ N, rubrica em queda, consistência baixa, gates falhando com frequência, sono irregular).
- **Item de backlog**: uma recomendação acionável e pequena, com objetivo, evidência que a motivou e critério observável de “feito”.
- **Itens ativos**: subconjunto pequeno do backlog elegível para aparecer no planejamento diário/semanal.

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Gerar item de backlog a partir de lacuna observada (Priority: P1)

O usuário quer que o sistema transforme sinais reais (ex.: repetição de erro) em um próximo passo claro.

**Why this priority**: Implementa R4 (backlog inteligente) e reduz o “o que faço agora?” com base em evidência (PRD §5.2; R4).

**Independent Test**:
- Simular 7 dias de registros com uma lacuna clara (ex.: mesmo erro 3x).
- Verificar que um item é criado com origem, motivo e critério observável.

**Acceptance Scenarios**:

1. **Scenario**: Erro recorrente em Java gera item de reforço
   - **Given** o usuário registrou o mesmo erro de Java em múltiplas sessões (ver `SPEC-016`/erros recorrentes)
   - **When** o sistema atualiza o backlog
   - **Then** um item de backlog é criado com: “reforço do erro X”, por que foi criado (contagem/recência), e critério observável de feito (evidência mínima)

2. **Scenario**: Consistência baixa em uma meta gera item de ajuste mínimo
   - **Given** a consistência semanal de uma meta caiu (ex.: 1/7)
   - **When** o sistema atualiza o backlog
   - **Then** o sistema cria um item “reduzir escopo / mínimo viável” com foco em remover fricção, não em aumentar exigência

---

### User Story 2 — Priorizar e limitar itens para evitar overload (Priority: P1)

O usuário não quer uma lista infinita; quer foco.

**Why this priority**: O problema central do PRD é overload e perda de consistência; backlog sem limite piora (PRD §6.2; RNF1/RNF2).

**Independent Test**:
- Simular múltiplas lacunas em metas diferentes.
- Verificar que o sistema limita itens ativos e explica a priorização.

**Acceptance Scenarios**:

1. **Scenario**: Muitas lacunas ao mesmo tempo → sistema limita itens ativos
   - **Given** existem lacunas em múltiplas metas intensivas e fundação
   - **When** o backlog é atualizado
   - **Then** o sistema mantém apenas um número pequeno de itens ativos (ex.: 3–7) e mantém o restante como “não ativo”, sem empurrar tudo para o usuário

2. **Scenario**: Prioridade baseada em impacto e contexto
   - **Given** existem itens de backlog com diferentes origens (erro recorrente, baixa consistência, rubrica em queda)
   - **When** o sistema escolhe quais itens ficam ativos
   - **Then** o sistema prioriza os que têm maior impacto esperado no progresso real e que cabem no contexto típico do usuário, explicando em 1 frase por item (“por que agora”)

---

### User Story 3 — Usuário aceita, adia ou rejeita sugestão (Priority: P2)

O usuário quer autonomia sem quebrar o sistema.

**Why this priority**: Mantém segurança psicológica e reduz fricção/abandono (RNF1; RNF3).

**Independent Test**:
- Criar 3 itens ativos e simular ações: aceitar/adiar/rejeitar.
- Verificar que o backlog e o “próximo passo” mudam sem insistência excessiva.

**Acceptance Scenarios**:

1. **Scenario**: Aceitar item de backlog
   - **Given** existe um item ativo recomendado
   - **When** o usuário aceita
   - **Then** o item muda para estado “aceito/planejável” e fica elegível para aparecer no plano diário/semanal

2. **Scenario**: Adiar item (sem culpa)
   - **Given** existe um item ativo
   - **When** o usuário adia por falta de tempo/energia
   - **Then** o sistema registra adiamento, reduz insistência, e oferece alternativa mínima se existir

3. **Scenario**: Rejeitar item com motivo opcional
   - **Given** existe um item recomendado
   - **When** o usuário rejeita
   - **Then** o sistema registra a rejeição (motivo opcional) e evita sugerir o mesmo item repetidamente no curto prazo, mantendo tom não punitivo

---

### User Story 4 — Consultar backlog atual e entender “por quê” (Priority: P2)

O usuário quer transparência: “por que você está sugerindo isso?”

**Why this priority**: Aumenta confiança e reduz sensação de aleatoriedade (PRD “evidência > intuição”).

**Independent Test**:
- Com itens ativos existentes, pedir “meu backlog”.
- Validar que cada item mostra origem e motivo de forma curta.

**Acceptance Scenarios**:

1. **Scenario**: Consultar backlog
   - **Given** existe backlog com itens ativos e não ativos
   - **When** o usuário solicita “meu backlog”
   - **Then** o sistema lista itens ativos primeiro, cada um com: título curto + motivo (1 frase) + critério observável de feito

## Edge Cases *(mandatory)*

- What happens when existem itens demais (explosão de backlog)?
  - Sistema limita itens ativos, agrupa por meta/domínio e guarda o restante como não ativo; não despeja tudo no usuário.
- What happens when faltam dados suficientes para inferir lacunas?
  - Sistema não inventa; marca como insuficiente e sugere coletar o mínimo necessário no próximo check-in/revisão.
- How does system handle lacunas em múltiplas metas intensivas ao mesmo tempo?
  - Sistema respeita limites de governança (`SPEC-010`) e prioriza 1–2 frentes; o resto vira “fila” não ativa.
- What happens when o usuário tenta usar backlog como lista genérica (“coloca ‘pagar contas’ aí”)?
  - Sistema pode registrar como nota separada (fora do backlog inteligente) ou recusar, explicando escopo (o backlog é para progresso observável das metas).
- What happens when o usuário sente sobrecarga com sugestões?
  - Sistema reduz itens ativos e oferece modo “mínimo viável” por uma semana, sem culpa.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST manter um backlog de itens derivados de lacunas observadas (ex.: erros recorrentes, baixa consistência, rubrica em queda), referenciando a evidência que motivou o item (via `SPEC-016`).
- **FR-002**: System MUST atribuir a cada item: origem (lacuna), prioridade, status, e critério observável de “feito” (alinhado com `SPEC-003` quando envolver evidência/gates).
- **FR-003**: System MUST priorizar itens com base em impacto esperado no progresso real e recência/frequência do sinal, evitando recomendar “mais coisas” quando o problema é overload.
- **FR-004**: System MUST limitar a quantidade de itens **ativos** para evitar overload, mantendo o restante como não ativo/arquivado conforme política definida.
- **FR-005**: System MUST permitir ao usuário: aceitar, adiar e rejeitar itens, registrando o estado e reduzindo repetição excessiva.
- **FR-006**: System MUST permitir consulta do backlog e MUST explicar “por que” cada item está recomendado (1 frase).
- **FR-007**: System MUST respeitar governança de metas em paralelo (PRD §6.2 / `SPEC-010`) ao recomendar itens, não aumentando escopo em múltiplas metas intensivas simultaneamente.

### Non-Functional Requirements

- **NFR-001**: System MUST manter simplicidade e baixa fricção (PRD RNF1): backlog deve ser curto e acionável, não um relatório.
- **NFR-002**: System MUST ser robusto a semanas ruins (PRD RNF2): reduzir sugestões e focar em mínimo viável quando consistência/energia estiverem baixas.
- **NFR-003**: System MUST manter segurança psicológica (PRD RNF3): recomendações firmes, não punitivas.

### Key Entities *(include if feature involves data)*

- **BacklogItem**: id; título curto; meta/domínio; origem (lacuna); evidência_resumo; prioridade; status (ativo/aceito/adiado/rejeitado/concluído/arquivado); critério_de_feito; criado_em; atualizado_em.
- **BacklogSignal**: tipo (erro_recorrente/queda_rubrica/baixa_consistencia/energia_baixa/gate_falhando); intensidade (baixo/médio/alto); janela; referência a registros relevantes.
- **BacklogDecision**: item; ação do usuário (aceitar/adiar/rejeitar); motivo opcional; timestamp.

## Acceptance Criteria *(mandatory)*

- O sistema cria itens de backlog a partir de lacunas observáveis e registra a origem/evidência de forma auditável.
- O sistema limita itens ativos e prioriza foco, evitando overload.
- O usuário consegue aceitar/adiar/rejeitar itens sem culpa e sem repetição excessiva.
- Cada item recomendado inclui uma explicação curta (“por que agora”) e critério observável de “feito”.

## Business Objectives *(mandatory)*

- Priorizar próximos passos com base em evidência (Evidência > Intuição) e reduzir decisões difíceis no dia a dia (PRD §5.2; R4).
- Evitar overload e sustentar consistência com foco e limites (PRD §6.2; §13).
- Manter baixa fricção e confiança do usuário (RNF1; RNF3).

## Error Handling *(mandatory)*

- **Dados insuficientes**: não sugerir com falsa certeza; pedir o mínimo de dados no próximo check-in e marcar a recomendação como “indeterminada”.
- **Explosão de itens**: reduzir ativos, arquivar/adiar automaticamente itens de baixo impacto e comunicar foco.
- **Usuário sobrecarregado**: oferecer modo mínimo (menos sugestões) por um período curto, mantendo apenas fundação + 1 frente.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Aumento da taxa de execução de itens aceitos (itens aceitos → concluídos com gate quando aplicável).
- **SC-002**: Redução de recorrência de lacunas alvo (ex.: queda de erros recorrentes) após 2–4 semanas.
- **SC-003**: Backlog permanece pequeno e acionável: maioria dos usuários mantém ≤ limite de itens ativos sem sensação de burocracia.