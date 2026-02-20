# Feature Specification: SaaS — “aposta semanal” (bloco profundo) + microblocos

**Created**: 2026-02-19  
**PRD Base**: §8.6, §6.2, §5.2, §5.3, §9.1, §9.2, §14, §11 (RNF1–RNF3)

## Caso de uso *(mandatory)*

O usuário quer avançar no seu SaaS em paralelo sem canibalizar fundamentos (sono/saúde) e sem explodir o escopo das metas intensivas. Esta feature define como o sistema:
- organiza o SaaS como uma **aposta semanal** (1 bloco profundo),
- permite **1–2 microblocos** de manutenção (baixo custo cognitivo),
- exige **definição de pronto** e “resultado observado” para contar como feito,
- e ajusta automaticamente o compromisso quando houver overload/dias ruins.

A intenção é manter progresso consistente e mensurável (mesmo pequeno), sem transformar a semana em “mais uma meta intensiva diária”.

## Scope & Non-goals *(recommended)*

**In scope (MVP slice)**:
- Planejar 1 bloco profundo semanal (30–120 min, ajustável) com objetivo e definição de pronto.
- Planejar 1–2 microblocos (5–20 min) com definição de pronto simples.
- Registrar conclusão com evidência mínima textual (resultado/entrega/decisão tomada).
- Ajustar o compromisso em semanas ruins (reduzir escopo em vez de abandonar).

**Non-goals (agora)**:
- Não definir ferramentas, stack, repositórios, board, issues, CI/CD.
- Não gerenciar backlog completo do SaaS; apenas garantir 1 foco semanal (se precisar, isso se conecta com `SPEC-008`).
- Não detalhar execução diária do SaaS no mesmo nível de inglês/java.

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Planejar “aposta semanal” com definição de pronto e resultado esperado (Priority: P1)

O usuário quer escolher um foco do SaaS na semana e saber o que significa “feito”, evitando esforço sem entrega.

**Why this priority**: O PRD descreve SaaS como aposta semanal para não canibalizar o resto (PRD §8.6; §6.2). Sem definição de pronto, vira falso progresso.

**Independent Test**:
- Simular início de semana.
- Validar que o sistema coleta/define: foco único, duração do bloco, definição de pronto, e resultado esperado em linguagem de produto.

**Acceptance Scenarios**:

1. **Scenario**: Aposta semanal definida com clareza
   - **Given** o usuário quer avançar no SaaS esta semana
   - **When** inicia “minha aposta semanal do SaaS”
   - **Then** o sistema ajuda a definir 1 foco (uma entrega/validação/decisão), a duração do bloco profundo e uma definição de pronto observável

2. **Scenario**: Foco grande demais → sistema corta escopo
   - **Given** o usuário escolhe algo muito amplo (“construir tudo”)
   - **When** o sistema avalia a definição de pronto
   - **Then** o sistema sugere recorte (menor entrega possível) e mantém apenas 1 resultado esperado

---

### User Story 2 — Executar bloco profundo e só contar como concluído com “pronto” (Priority: P1)

O usuário quer que o sistema diferencie “trabalhei um pouco” de “entreguei algo”.

**Why this priority**: Evita ilusão de progresso e mantém consistência mensurável (PRD §8.6; princípio Qualidade > Quantidade).

**Independent Test**:
- Simular execução do bloco profundo com (a) entrega feita e (b) entrega incompleta.
- Verificar que o sistema registra corretamente concluído vs tentativa e solicita o menor passo para completar.

**Acceptance Scenarios**:

1. **Scenario**: Bloco profundo concluído com definição de pronto satisfeita
   - **Given** existe uma aposta semanal definida
   - **When** o usuário executa o bloco profundo e relata o resultado
   - **Then** o sistema marca como concluído apenas se a definição de pronto foi atendida e registra o resultado observado

2. **Scenario**: Bloco profundo incompleto vira tentativa (não concluído)
   - **Given** o usuário iniciou o bloco profundo
   - **When** não conseguiu terminar o pronto
   - **Then** o sistema registra como tentativa, salva o “próximo passo mínimo” e não conta como concluído

---

### User Story 3 — Microblocos: manter progresso sem sobrecarga (Priority: P2)

O usuário quer microblocos curtos (ex.: 5–20 min) para manutenção ou continuidade leve, sem tomar o lugar de fundamentos.

**Why this priority**: Sustenta consistência sem canibalizar o restante (PRD §8.6; §6.2; RNF1).

**Independent Test**:
- Simular uma semana com 2 microblocos executados e um não executado.
- Verificar registro simples e sem culpa.

**Acceptance Scenarios**:

1. **Scenario**: Microbloco executado com pronto simples
   - **Given** existe um microbloco planejado
   - **When** o usuário executa
   - **Then** o sistema registra como concluído com uma evidência mínima (1 frase do que foi feito)

2. **Scenario**: Microbloco omitido em semana ruim
   - **Given** semana com pouco tempo/energia
   - **When** o usuário não executa microblocos
   - **Then** o sistema não penaliza; mantém apenas a aposta semanal (ou reduz) e reforça mínimo viável

---

### User Story 4 — Semana ruim: reduzir compromisso sem abandonar (Priority: P2)

O usuário tem uma semana caótica e quer adaptar o SaaS para não virar culpa e não sabotar fundação.

**Why this priority**: RNF2 e risco de overload (PRD §13).

**Independent Test**:
- Simular semana de baixa energia.
- Verificar que o sistema reduz escopo e mantém consistência mínima.

**Acceptance Scenarios**:

1. **Scenario**: Downgrade da aposta semanal
   - **Given** o usuário reporta semana ruim
   - **When** chega o momento do bloco profundo
   - **Then** o sistema oferece redução para um bloco menor ou um “resultado mínimo viável” (ex.: decidir X, escrever 3 bullets, validar 1 hipótese), sem culpa

## Edge Cases *(mandatory)*

- What happens when o usuário quer transformar SaaS em meta intensiva diária?
  - Sistema explica o racional (proteção de consistência) e mantém o SaaS como aposta semanal + microblocos opcionais.
- What happens when o usuário não consegue definir “pronto”?
  - Sistema oferece templates de pronto em linguagem natural (sem stack) e pede 1 confirmação curta.
- What happens when o usuário só “estudou/planejou” e quer contar como feito?
  - Sistema só conta como feito se a definição de pronto permitir; caso contrário registra como tentativa e sugere o menor passo para tornar observável.
- What happens when a aposta semanal entra em conflito com saúde/sono?
  - Sistema reduz escopo do SaaS e preserva fundação (PRD §6.1; §6.2).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST suportar definir 1 aposta semanal (foco único) para o SaaS com: resultado esperado e definição de pronto observável.
- **FR-002**: System MUST suportar planejar 1 bloco profundo semanal com duração aproximada e critério de “feito” (pronto).
- **FR-003**: System MUST suportar 1–2 microblocos por semana com pronto simples e evidência mínima textual.
- **FR-004**: System MUST registrar conclusão do bloco profundo apenas quando a definição de pronto for atendida; caso contrário registrar tentativa e próximo passo mínimo.
- **FR-005**: System MUST permitir reduzir escopo em semanas ruins (downgrade para mínimo viável) preservando segurança psicológica.
- **FR-006**: System MUST manter o SaaS fora do limite de “metas intensivas diárias” e respeitar governança de metas em paralelo (PRD §6.2).

### Non-Functional Requirements

- **NFR-001**: System MUST manter baixa fricção: planejamento e registro devem ser curtos (PRD RNF1).
- **NFR-002**: System MUST ser robusto a semanas/dias ruins (PRD RNF2).
- **NFR-003**: System MUST manter tom não punitivo (PRD RNF3).

### Key Entities *(include if feature involves data)*

- **SaasWeeklyBet**: semana; foco; resultado_esperado; definição_de_pronto; status (planejada/concluída/tentativa/adiada).
- **DeepWorkBlock**: semana; duração_aprox; pronto; resultado_observado; status.
- **MicroBlock**: data/semana; pronto_simples; evidência_mínima; status.
- **SaasProgressNote**: nota curta do que mudou/decidiu/validou (1–3 bullets).

## Acceptance Criteria *(mandatory)*

- O usuário consegue definir uma aposta semanal com pronto e resultado esperado.
- O bloco profundo só conta como concluído quando o pronto é atendido; caso contrário vira tentativa com próximo passo mínimo.
- Microblocos são opcionais e não geram overload; em semanas ruins o sistema reduz escopo sem culpa.
- O SaaS permanece como aposta semanal e não canibaliza fundamentos/metas intensivas.

## Business Objectives *(mandatory)*

- Progresso consistente no SaaS sem sabotar fundamentos (PRD §8.6; §6.2).
- Evitar falso progresso exigindo pronto e resultado observável (qualidade > quantidade).
- Sustentar hábito de entrega pequena e contínua com baixa fricção (RNF1).

## Error Handling *(mandatory)*

- **Overload/semana ruim**: reduzir escopo automaticamente e sugerir mínimo viável.
- **Pronto indefinido**: oferecer exemplos e pedir 1 confirmação curta.
- **Usuário some**: manter estado e permitir retomar com próximo passo mínimo.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: % de semanas com aposta semanal concluída (ou mínima viável concluída) ao longo de 4–8 semanas.
- **SC-002**: Aumento de entregas/validações pequenas do SaaS sem queda em fundamentos (sono/saúde).
- **SC-003**: Redução de semanas “zero” (sem qualquer avanço) por existir downgrade para mínimo viável.