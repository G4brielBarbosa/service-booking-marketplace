# Feature Specification: Vida saudável — planejamento semanal + métricas + limites por dor/energia

**Created**: 2026-02-19  
**PRD Base**: §8.4, §6.1, §5.2, §5.3, §9.1, §9.2, §14, §11 (RNF1–RNF4)

## Caso de uso *(mandatory)*

O usuário quer melhorar saúde de forma **sustentável** sem cair no “tudo ou nada”. Esta feature define como o sistema:
- ajuda a montar um **plano semanal simples** (força + atividade moderada) adequado à rotina,
- ajusta recomendações com base em **energia/tempo** do dia e em **dor/fadiga** (segurança),
- registra métricas mínimas para ver tendência (sem burocracia),
- e mantém versão mínima viável para semanas/ dias ruins (consistência mínima).

A feature deve evitar prescrição médica e manter segurança psicológica: foco em tendência e ajustes, não em culpa.

## Scope & Non-goals *(recommended)*

**In scope (MVP slice)**:
- Planejamento semanal com 2 componentes: (a) força, (b) cardio/atividade moderada (ou equivalente).
- Ajuste por energia/tempo (Plano A/B/C) e por sinais de excesso (dor/fadiga).
- Registro mínimo de execução (minutos/sessões) e sinais (dor/fadiga percebida).
- Feedback semanal simples (“o que funcionou / próximo ajuste”).

**Non-goals (agora)**:
- Não prescrever treino nem reabilitação; sem recomendações clínicas.
- Não exigir contagem de macros, calorias, wearables ou biometria avançada.
- Não otimizar periodização/rotinas detalhadas; o foco é adesão sustentável.

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Criar plano semanal simples e realista (Priority: P1)

O usuário quer sair do “vou treinar mais” para um plano semanal concreto, curto e adaptável, que caiba na vida real.

**Why this priority**: Sem plano semanal simples a meta vira vaga; e o PRD pede sustentabilidade e adaptação (PRD §8.4; §5.2).

**Independent Test**:
- Simular um usuário com agenda cheia e um com agenda livre.
- Rodar o fluxo de planejamento semanal e validar que o resultado é um plano curto com sessões mínimas e alternativas (mínimo viável).

**Acceptance Scenarios**:

1. **Scenario**: Plano semanal criado em poucos passos
   - **Given** o usuário quer melhorar saúde e tem restrições de tempo/energia variáveis
   - **When** inicia “planejar saúde da semana”
   - **Then** o sistema propõe um plano semanal simples com: 1–2 sessões de força + 2–4 sessões de atividade moderada (ou equivalente), ajustadas ao contexto do usuário, com versão mínima viável

2. **Scenario**: Usuário rejeita ambição alta → sistema reduz escopo
   - **Given** o usuário diz que o plano proposto é demais
   - **When** pede “mais simples”
   - **Then** o sistema reduz para um mínimo viável (ex.: 1 sessão curta + 2 caminhadas curtas) e reforça que consistência mínima é sucesso

---

### User Story 2 — Dia a dia: adaptar execução por energia/tempo sem quebrar o plano (Priority: P1)

O usuário acorda e não sabe se vai conseguir cumprir o plano; ele quer um ajuste rápido, sem perder consistência.

**Why this priority**: Robustez a dias ruins é requisito explícito (RNF2) e reduz abandono.

**Independent Test**:
- Simular 3 dias: energia alta, média e baixa.
- Verificar que o sistema sempre oferece um passo executável (Plano A/B/C), registrável e alinhado com a semana.

**Acceptance Scenarios**:

1. **Scenario**: Dia bom segue plano cheio
   - **Given** hoje o usuário tem tempo e energia adequados
   - **When** pede “o que faço de saúde hoje?”
   - **Then** o sistema sugere a sessão planejada do dia (ou a melhor encaixável), com duração aproximada e critério observável de feito

2. **Scenario**: Dia ruim oferece mínimo viável
   - **Given** hoje o usuário tem pouco tempo/energia
   - **When** pede “mínimo viável”
   - **Then** o sistema oferece uma alternativa curta (ex.: caminhada curta, mobilidade leve, ou outra opção segura e simples), e registra como consistência mínima sem culpa

---

### User Story 3 — Ajustar por dor/fadiga para evitar excesso (Priority: P1)

O usuário pode sentir dor/fadiga e quer que o sistema reduza carga e evite escalada arriscada.

**Why this priority**: Sustentabilidade e segurança são centrais; excesso leva a abandono e risco físico (PRD §8.4).

**Independent Test**:
- Simular usuário reportando dor/fadiga antes de uma sessão planejada.
- Verificar que o sistema adapta ou substitui por alternativa leve e registra o sinal.

**Acceptance Scenarios**:

1. **Scenario**: Dor reportada → reduzir intensidade/alternativa segura
   - **Given** o usuário reporta dor ou fadiga acima do normal
   - **When** solicita a atividade do dia
   - **Then** o sistema reduz o esforço (ou sugere alternativa leve) e registra o sinal, deixando claro que não é prescrição médica

2. **Scenario**: Dor persistente → recomendar cautela e simplificar por uma semana
   - **Given** dor/fadiga foi reportada repetidamente na semana
   - **When** chega a revisão semanal de saúde
   - **Then** o sistema sugere uma semana de redução de carga e incentiva buscar avaliação profissional se apropriado, sem alarmismo

---

### User Story 4 — Ver progresso semanal simples (minutos/sessões) sem burocracia (Priority: P2)

O usuário quer saber se está melhorando, sem planilha.

**Why this priority**: Métricas simples sustentam evidência > intuição e revisão semanal (PRD §5.5; §9.2).

**Independent Test**:
- Simular uma semana com 3 registros de atividade + 1 registro de dor.
- Validar que o sistema consegue consolidar e apresentar tendência simples e um ajuste.

**Acceptance Scenarios**:

1. **Scenario**: Resumo semanal de saúde
   - **Given** houve registros de atividade e/ou sinais de excesso
   - **When** o usuário pede “como foi minha saúde essa semana?”
   - **Then** o sistema retorna: sessões/minutos totais aproximados, consistência (quantos dias), e 1 ajuste proposto para a próxima semana

## Edge Cases *(mandatory)*

- What happens when o usuário não consegue fazer nada na semana?
  - Sistema evita culpa, propõe recomeço mínimo (2–3 passos curtos) e remove ambição excessiva.
- What happens when o usuário quer “compensar” tudo no fim de semana?
  - Sistema desencoraja abordagem punitiva/compensatória e sugere consistência distribuída com mínimo viável.
- What happens when dados são muito incompletos?
  - Sistema marca como “dados insuficientes” e pede 1 registro mínimo por dia/semana para melhorar recomendações.
- What happens when o usuário pede orientação médica específica?
  - Sistema reforça limitação (não é médico) e sugere procurar profissional; mantém recomendações gerais e seguras em nível de produto.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST permitir criar um plano semanal simples de vida saudável com componentes de força e atividade moderada, ajustáveis à rotina do usuário.
- **FR-002**: System MUST oferecer versão mínima viável do plano (para semanas/dias ruins) para preservar consistência.
- **FR-003**: System MUST adaptar recomendações diárias por tempo/energia, oferecendo alternativas curtas quando necessário.
- **FR-004**: System MUST registrar execução com métricas mínimas (ex.: sessões ou minutos aproximados) e consistência por semana.
- **FR-005**: System MUST registrar sinais de excesso reportados (dor/fadiga) e adaptar recomendações quando presentes.
- **FR-006**: System MUST fornecer um resumo semanal simples (o que foi feito + sinais + 1 ajuste para próxima semana), compatível com a revisão semanal (`SPEC-007`).

### Non-Functional Requirements

- **NFR-001**: System MUST evitar abordagem “tudo ou nada” e reforçar consistência mínima como sucesso (PRD §8.4; RNF3).
- **NFR-002**: System MUST manter baixa fricção (PRD RNF1): registro simples, sem burocracia.
- **NFR-003**: System MUST ser robusto a dias ruins (PRD RNF2): sempre existir uma alternativa mínima viável.
- **NFR-004**: System MUST aplicar privacidade por padrão (PRD RNF4): coletar apenas o necessário (minutos/sessões e sinais simples).

### Key Entities *(include if feature involves data)*

- **HealthWeeklyPlan**: semana; sessões planejadas (força/moderada); versão mínima; status.
- **HealthActivityLog**: data; tipo (força/moderada/alternativa_leve); duração_aprox; status (feito/parcial).
- **HealthSignal**: data; tipo (dor/fadiga); intensidade percebida (escala simples); nota curta opcional.
- **HealthWeeklySummary**: semana; total_aprox; consistência; sinais; ajuste_recomendado.

## Acceptance Criteria *(mandatory)*

- O usuário consegue criar um plano semanal simples e realista e obter uma versão mínima viável.
- Em dia ruim, o sistema oferece alternativa curta e segura sem culpa.
- Dor/fadiga reportada altera recomendações (reduz/ajusta) e é registrada como sinal.
- O usuário consegue ver um resumo semanal simples com 1 ajuste acionável para a próxima semana.

## Business Objectives *(mandatory)*

- Sustentar hábitos de saúde com consistência e segurança, sem burnout (PRD §8.4).
- Reduzir fricção e manter adaptação por evidência (PRD §5.5; RNF1).
- Proteger o usuário em semanas ruins com mínimo viável (RNF2) e tom não punitivo (RNF3).

## Error Handling *(mandatory)*

- **Dor/limitação**: reduzir carga e sugerir alternativa leve; reforçar que não é prescrição médica.
- **Ausência de registros**: retomar com defaults e mínimo viável; registrar lacuna sem penalizar.
- **Pedido de orientação clínica**: informar limitação e sugerir buscar profissional; manter recomendações gerais.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Aumento sustentado de consistência semanal (dias com alguma atividade, mesmo mínima).
- **SC-002**: Aumento de minutos/sessões semanais ao longo de 4–8 semanas sem aumento de sinais de excesso.
- **SC-003**: Redução de semanas “zero atividade” após adoção de mínimo viável.