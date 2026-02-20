# Feature Specification: Privacidade por padrão — dados mínimos, transparência e controle do usuário

**Created**: 2026-02-19  
**PRD Base**: §11 (RNF4), §5.1, §5.3, §5.4, §9.1, §9.2, §10 (R6), §2 (consultar progresso), §14

## Caso de uso *(mandatory)*

O usuário quer um assistente pessoal Telegram-first que coleta dados para funcionar (planos, evidências, métricas), mas que por padrão:
- colete **o mínimo necessário**,
- seja **transparente** sobre o que é guardado e por quê,
- dê **controle real** ao usuário (revisar, exportar, apagar, limitar retenção, opt-out por tipo),
- e permita operar em “modo mínimo” quando privacidade/ambiente forem limitantes, sem punição.

Esta SPEC define comportamentos observáveis e políticas de produto (não define banco, criptografia, provedores, nem detalhes técnicos).

> Cross-cutting: esta SPEC governa coleta/retensão/controle para fluxos descritos em `SPEC-001/002/003/005/006/007/008/009/012/013/014/016`.

## Scope & Non-goals *(recommended)*

**In scope (MVP slice)**:
- Transparência: o usuário consegue perguntar “o que você guarda?” e receber resposta clara por fluxo.
- Minimização: cada fluxo registra apenas o necessário para cumprir seus objetivos.
- Controles: revisar dados, apagar dados, e configurar retenção/opt-out por categoria.
- “Modo mínimo” para situações sensíveis (ambiente sem privacidade, áudio impossível, etc.).

**Non-goals (agora)**:
- Não detalhar implementação de segurança (criptografia, infra, logs).
- Não prometer anonimização perfeita; focar em controle e minimização.

## Princípios (produto) *(recommended)*

- **Minimização**: se um dado não for necessário para uma funcionalidade do PRD, não coletar/guardar.
- **Finalidade explícita**: todo dado guardado deve ter “por que” claro (planejamento, métricas, gates, revisão semanal).
- **Retenção limitada**: guardar pelo tempo mínimo que preserve utilidade.
- **Controle do usuário**: ver, apagar, e limitar.
- **Sem punição**: recusar dados (ex.: áudio) não vira “bronca”; vira alternativa mínima ou bloqueio explícito do que depende do dado.

## Classificação de dados (para políticas) *(recommended)*

Categorias (linguagem de produto):
- **C1 — Check-ins e planos**: tempo/energia do dia, plano A/B/C, status de tarefas (`SPEC-002`).
- **C2 — Evidências de aprendizagem**: respostas, rubricas, tentativas e resultados de gates (`SPEC-003`).
- **C3 — Conteúdo sensível**: áudio de speaking, textos com conteúdo emocional/privado (pode aparecer em `SPEC-003/013`).
- **C4 — Métricas agregadas**: consistência semanal, médias de rubrica, tendências de sono/energia (`SPEC-016`).
- **C5 — Governança e preferências**: metas ativas/pausadas, limites, configurações e opt-outs (`SPEC-001/010`).

## Política padrão (defaults) *(recommended)*

> Defaults devem ser claros e fáceis de explicar; o usuário pode alterar.

- **Padrão de coleta**: coletar apenas o que o fluxo precisa para funcionar hoje + comparações simples (sem “capturar tudo”).
- **Padrão de retenção**:
  - **C4 (métricas agregadas)**: retenção mais longa (necessária para tendência/semana vs semana).
  - **C1/C2**: retenção moderada (necessária para consultar progresso recente e revisão semanal).
  - **C3 (conteúdo sensível bruto)**: retenção curta por padrão, com opção de “não guardar” (apenas usar para validação do gate e descartar).
- **Padrão de transparência**: sempre que um fluxo pedir dado sensível, explicar antes “o mínimo guardado e por quê” em 1–2 frases.
- **Padrão de opt-out**: permitir limitar armazenamento de categorias (especialmente C3), aceitando operar em modo mínimo quando possível.

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Perguntar “o que você guarda?” e receber resposta clara por fluxo (Priority: P1)

**Why this priority**: Confiança é pré-requisito; RNF4 exige transparência.

**Independent Test**:
- Em onboarding, em envio de evidência e em revisão semanal, perguntar “o que você guarda?”.
- Validar resposta consistente, curta e alinhada às categorias.

**Acceptance Scenarios**:

1. **Scenario**: Transparência durante onboarding
   - **Given** o usuário está no onboarding (`SPEC-001`)
   - **When** pergunta “o que você guarda e por quê?”
   - **Then** o sistema lista as categorias relevantes (metas/restrições/baseline) e a finalidade em linguagem simples

2. **Scenario**: Transparência ao solicitar evidência sensível
   - **Given** uma tarefa exige evidência (ex.: áudio) (`SPEC-003`)
   - **When** o sistema pede a evidência
   - **Then** antes (ou junto) explica o mínimo guardado, por quanto tempo (em termos de política), e a alternativa/opt-out disponível

---

### User Story 2 — Revisar meus dados e apagar por categoria/período (Priority: P1)

O usuário quer controle real: poder ver e apagar.

**Why this priority**: RNF4; reduz ansiedade e aumenta adesão.

**Independent Test**:
- Criar dados de 1 semana (check-ins, rubricas, sono).
- Executar: “mostrar meus dados” e “apagar dados de X”.

**Acceptance Scenarios**:

1. **Scenario**: Revisar resumo dos dados registrados
   - **Given** existem dados registrados nas categorias C1–C5
   - **When** o usuário pede “quais dados você tem sobre mim?”
   - **Then** o sistema mostra um resumo por categoria (sem despejar conteúdo sensível por padrão) e oferece opções de detalhar

2. **Scenario**: Apagar conteúdo sensível bruto
   - **Given** existem evidências sensíveis (C3)
   - **When** o usuário pede “apague meus áudios/textos sensíveis”
   - **Then** o sistema confirma a ação (irreversível), apaga e informa o que pode deixar de funcionar/ficar menos preciso

3. **Scenario**: Apagar por período (ex.: última semana)
   - **Given** existem registros na última semana
   - **When** o usuário pede para apagar dados de um período
   - **Then** o sistema executa e atualiza consultas/revisões futuras para refletir ausência de dados, sem “penalizar” o usuário

---

### User Story 3 — Opt-out e modo mínimo (sem guardar C3) sem perder o produto (Priority: P1)

O usuário quer usar o assistente mesmo quando não pode compartilhar/guardar dados sensíveis.

**Why this priority**: Muitos fluxos pedem evidência e registros; sem modo mínimo o produto vira inviável em dias reais.

**Independent Test**:
- Ativar opt-out de C3.
- Tentar executar tarefas que normalmente usariam C3 e verificar comportamento (alternativa, bloqueio explícito, ou versão mínima).

**Acceptance Scenarios**:

1. **Scenario**: Opt-out de armazenamento de conteúdo sensível
   - **Given** o usuário quer privacidade máxima
   - **When** configura “não guardar conteúdo sensível”
   - **Then** o sistema confirma e explica efeitos (ex.: menos histórico detalhado), mas mantém métricas agregadas quando possível

2. **Scenario**: Evidência exigida mas usuário não pode enviar/guardar
   - **Given** uma tarefa de aprendizagem exige evidência sensível (`SPEC-003`)
   - **When** o usuário diz “não posso enviar áudio agora” ou “não quero que guarde”
   - **Then** o sistema oferece alternativa equivalente **se** definida; caso contrário registra bloqueio/pendência de forma não punitiva e oferece mínimo viável do dia (Plano C)

---

### User Story 4 — Retenção configurável e expiração automática (Priority: P2)

O usuário quer limites claros: o que expira, o que fica, e poder ajustar.

**Independent Test**:
- Configurar retenção curta para C1/C2 e “não guardar” para C3.
- Verificar que o sistema comunica o efeito nas consultas e revisões.

**Acceptance Scenarios**:

1. **Scenario**: Ajustar retenção por categoria
   - **Given** o usuário quer reduzir retenção
   - **When** ajusta retenção de uma categoria
   - **Then** o sistema confirma, explica impacto e passa a operar conforme a nova política

2. **Scenario**: Dados expirados não quebram o sistema
   - **Given** parte dos dados expirou
   - **When** o usuário pede progresso/revisão semanal
   - **Then** o sistema explicita limitações (“não tenho X porque expirou”), usa métricas agregadas quando disponíveis e sugere coleta mínima para recuperar

## Edge Cases *(mandatory)*

- What happens when o usuário pede “apaga tudo”?
  - Sistema confirma irreversibilidade, executa e volta para modo “usuário novo”, sem julgamento.
- What happens when o usuário quer privacidade máxima mas exige features que dependem de dados (ex.: tendência semanal)?
  - Sistema explica trade-off e oferece configuração intermediária (guardar apenas agregados).
- What happens when o usuário escreve mensagens com dados pessoais acidentalmente?
  - Sistema sugere minimizar conteúdo e oferece opção rápida “não guardar esta mensagem como registro”.
- What happens when há conflito entre “apagar dados” e “preciso de histórico para tendência”?
  - Sistema respeita o pedido; tendências futuras passam a refletir dados faltantes sem penalizar.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST responder “o que você guarda e por quê?” com explicação clara, por categoria e por fluxo (check-in, evidência, métricas, revisão).
- **FR-002**: System MUST minimizar coleta: cada fluxo registra apenas o necessário para suas funções declaradas no PRD/SPECS.
- **FR-003**: System MUST permitir ao usuário revisar seus dados por categoria e por período, com visualização resumida por padrão.
- **FR-004**: System MUST permitir apagar dados por categoria e por período, com confirmação de irreversibilidade e explicação de impacto.
- **FR-005**: System MUST permitir configurar opt-out por categoria (especialmente conteúdo sensível C3) e operar em modo mínimo quando opt-out estiver ativo.
- **FR-006**: System MUST suportar retenção configurável e expiração automática por categoria (com defaults claros).
- **FR-007**: System MUST evitar que falta/remoção de dados gere punição: o sistema ajusta recomendações e explicita limitações de forma não punitiva.

### Non-Functional Requirements

- **NFR-001**: System MUST manter fricção mínima: configurações e explicações cabem em mensagens curtas (PRD RNF1).
- **NFR-002**: System MUST ser robusto a dias ruins e ambientes sem privacidade (PRD RNF2).
- **NFR-003**: System MUST manter segurança psicológica ao falar de privacidade e falhas de evidência (PRD RNF3).
- **NFR-004**: System MUST aplicar privacidade por padrão (PRD RNF4): minimização + transparência + controle.

### Key Entities *(include if feature involves data)*

- **PrivacyPolicy**: configurações por categoria (retenção, opt-out, modo mínimo).
- **DataInventorySummary**: resumo do que existe por categoria e período.
- **DeletionRequest**: categoria/período; confirmação; status; impacto_resumido.
- **SensitiveEvidenceHandling**: preferência do usuário (guardar por quanto tempo / não guardar) e comportamento esperado.

## Acceptance Criteria *(mandatory)*

- O usuário consegue entender (em linguagem simples) o que é guardado e por quê em cada fluxo.
- O usuário consegue revisar e apagar dados por categoria/período com confirmação clara.
- Existe opt-out para conteúdo sensível e modo mínimo, com comportamento não punitivo.
- Retenção padrão é limitada e configurável; expiração não quebra o sistema (apenas reduz confiança/precisão e isso é comunicado).

## Business Objectives *(mandatory)*

- Aumentar confiança e adesão reduzindo medo de coleta excessiva (RNF4).
- Reduzir fricção e permitir uso em condições reais (ambiente/privacidade) (RNF1/RNF2).
- Sustentar qualidade/evidência com controle e alternativas (PRD §5.4).

## Error Handling *(mandatory)*

- **Pedido ambíguo (“apaga”)**: pedir 1 clarificação curta (categoria/período) antes de executar.
- **Apagar conteúdo sensível**: confirmar irreversibilidade e explicar impacto em features dependentes.
- **Dados expirados**: comunicar limitação e sugerir o mínimo de coleta para recuperar utilidade.
- **Opt-out ativo**: não solicitar repetidamente o mesmo dado; oferecer alternativas/versão mínima.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Aumento de uso contínuo sem queda de confiança (proxy: menor taxa de abandono após pedidos de evidência).
- **SC-002**: Alta taxa de entendimento: usuários conseguem responder corretamente “o que o sistema guarda” após explicação curta (checar via pergunta simples opcional).
- **SC-003**: Redução de fricção em ambientes sensíveis: uso de modo mínimo sem abandono do sistema.