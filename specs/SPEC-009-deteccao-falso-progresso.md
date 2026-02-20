# Feature Specification: Detecção de falhas reais (“falso progresso”) e reforço automático

**Created**: 2026-02-19  
**PRD Base**: §5.4, §5.5, §9.1, §9.2, §10 (R3, R2), §13, §11 (RNF1–RNF3)

## Caso de uso *(mandatory)*

O usuário quer evitar “falso progresso”: dias em que ele até “fez”, mas **não houve aprendizagem/competência real** (Inglês/Java) ou o hábito/fundação virou marcação vazia. Esta feature define como o sistema:
- detecta sinais observáveis de falso progresso a partir de evidências e registros,
- reage com um **reforço pequeno e verificável** (não burocrático),
- ajusta recomendações do plano para reduzir repetição do mesmo erro,
- e mantém um tom firme, não punitivo, focado em processo.

Esta SPEC não define stack nem algoritmos avançados; no MVP a detecção pode ser baseada em regras simples usando dados já capturados.

Referências relacionadas:
- `SPEC-003` (quality gates/evidências e estados de conclusão)
- `SPEC-016` (métricas, rubricas, erros recorrentes, consistência)
- `SPEC-008` (backlog inteligente que pode receber itens de reforço)
- `SPEC-007` (revisão semanal onde padrões ficam explícitos)

## Scope & Non-goals *(recommended)*

**In scope (MVP slice)**:
- Definir sinais mínimos de “falso progresso” (aprendizagem) e “qualidade baixa persistente”.
- Acionar reforço curto (micro-tarefa) com evidência mínima.
- Registrar evento de detecção + ação tomada + resultado.

**Non-goals (agora)**:
- Não fazer diagnóstico educacional/psicológico; é um sistema de suporte.
- Não exigir análise perfeita/ML; começar com heurísticas observáveis.
- Não criar um “curso completo”; reforço é pequeno e focado.

## Definições *(recommended)*

- **Falso progresso (aprendizagem)**: completar tarefas mas (a) falhar repetidamente no mesmo erro, (b) queda persistente de rubrica/qualidade, (c) retrieval fraco recorrente, (d) alta taxa de gates falhando por evidência insuficiente.
- **Sinal**: um padrão observável derivado de registros (`SPEC-016`) e resultados de gates (`SPEC-003`).
- **Reforço**: micro-intervenção verificável (1–10 min) para atacar o sinal (ex.: reexplicar conceito, mini-quiz, repetir um padrão específico, micro-drill).
- **Janela**: período recente considerado para detectar padrões (definida por política simples, ex.: “últimos 7 dias” ou “últimas 5 sessões do domínio”).

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Detectar erro recorrente e acionar reforço curto (Priority: P1)

O usuário repete o mesmo erro e quer que o sistema perceba e faça ele atacar o erro, não só “seguir em frente”.

**Why this priority**: R3 é a mitigação central do risco “fiz mas não aprendi” (PRD §13) e diferencia o produto de checklist (PRD §5.4).

**Independent Test**:
- Simular 5 sessões de Java ou Inglês com o mesmo erro registrado ≥ 3 vezes.
- Verificar que o sistema detecta o sinal, solicita reforço, aplica gate mínimo ao reforço e registra resultado.

**Acceptance Scenarios**:

1. **Scenario**: Erro recorrente dispara reforço
   - **Given** o mesmo erro foi registrado em múltiplas ocasiões recentes (ver `SPEC-016`)
   - **When** o sistema processa uma nova conclusão de tarefa relacionada
   - **Then** o sistema sinaliza o padrão e solicita um reforço curto focado no erro, com critério observável de “feito”

2. **Scenario**: Reforço não vira burocracia (dia com pouco tempo)
   - **Given** o usuário tem pouco tempo/energia hoje
   - **When** o reforço é solicitado
   - **Then** o sistema oferece uma versão mínima do reforço (1 passo) que ainda seja observável e útil

3. **Scenario**: Usuário não completa reforço
   - **Given** o reforço foi solicitado
   - **When** o usuário não conclui / some
   - **Then** o sistema registra tentativa/pendência e reoferece em momento futuro apropriado sem insistência agressiva

---

### User Story 2 — Detectar queda de qualidade (rubrica) e ajustar foco (Priority: P1)

O usuário está “fazendo”, mas a qualidade está caindo e o sistema deve ajustar para proteger progresso real.

**Why this priority**: Rubricas e evidência são parte explícita do PRD (ex.: speaking) e devem orientar adaptação (PRD §5.5).

**Independent Test**:
- Simular rubricas semanais em queda (ex.: média 6→4).
- Verificar que o sistema detecta tendência e sugere ajuste/redução de carga + reforço.

**Acceptance Scenarios**:

1. **Scenario**: Queda de rubrica dispara alerta protetivo
   - **Given** existe uma sequência recente de rubricas abaixo do padrão anterior
   - **When** o sistema compila sinais da semana/dias recentes
   - **Then** o sistema comunica o sinal com tom não punitivo e sugere ajuste de estratégia (ex.: reduzir carga, focar em base, atacar 1 erro)

2. **Scenario**: Ajuste sugerido é pequeno e acionável
   - **Given** o sistema detectou queda de qualidade
   - **When** sugere mudança
   - **Then** a sugestão vem como 1 mudança por vez (ex.: “esta semana, foco em X”) com critério observável e sem expandir escopo

---

### User Story 3 — Detectar “conclusões vazias” (gates falhando) e reforçar evidência mínima (Priority: P2)

O usuário tenta marcar como feito sem evidência mínima; isso é um tipo de falso progresso operacional.

**Why this priority**: Conecta com `SPEC-003` (gates) e reduz “faz de conta” (PRD §5.4).

**Independent Test**:
- Simular 3 tentativas de concluir tarefas de aprendizagem sem evidência mínima.
- Verificar que o sistema detecta o padrão e aplica reforço de processo (ensinar o menor caminho para evidenciar) sem humilhar.

**Acceptance Scenarios**:

1. **Scenario**: Padrão de gate falhando aciona intervenção de processo
   - **Given** o usuário falhou gates por ausência/incompletude repetidamente
   - **When** falha novamente
   - **Then** o sistema oferece uma intervenção curta de processo (“como evidenciar em 30s”) e registra que o usuário precisa de ajuste de fricção/clareza

## Edge Cases *(mandatory)*

- What happens when faltam dados suficientes para afirmar falso progresso?
  - Sistema não acusa; marca como “sinal fraco” e pede 1 dado mínimo no próximo ciclo.
- What happens when o usuário está em dia ruim e a qualidade cai por energia baixa?
  - Sistema prioriza Plano C/MVD e reduz escopo, evitando interpretar como incapacidade; sugere foco em consistência mínima.
- What happens when múltiplos sinais aparecem ao mesmo tempo (erro recorrente + rubrica em queda + pouca consistência)?
  - Sistema escolhe 1 sinal principal por vez (maior impacto) e sugere uma intervenção, para não criar overload.
- What happens when usuário rejeita repetidamente reforços?
  - Sistema reduz insistência, registra preferência e leva o tema para revisão semanal como escolha consciente (sem culpa).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST identificar sinais de falso progresso a partir de evidências/registro existentes (ex.: erros recorrentes, rubricas em queda, retrieval fraco recorrente, gates falhando frequentemente).
- **FR-002**: System MUST definir uma janela/limiar simples para detecção no MVP (ex.: últimos N eventos/sessões ou últimos X dias), consistente em todo o sistema. *(Pode ser política global, mas a SPEC deve ter um default observável.)*
- **FR-003**: System MUST, ao detectar um sinal, acionar **um reforço curto** com critério observável de “feito” e evidência mínima, alinhado a `SPEC-003`.
- **FR-004**: System MUST limitar a intervenção para evitar overload: no máximo 1 reforço principal por ciclo curto (ex.: por dia) e 1 foco por vez quando múltiplos sinais existirem.
- **FR-005**: System MUST registrar: sinal detectado, motivo resumido, reforço solicitado, resultado (feito/não feito/tentativa), para uso em revisão semanal (`SPEC-007`) e métricas (`SPEC-016`).
- **FR-006**: System MUST ajustar recomendações futuras com base no sinal (ex.: priorizar reforço no backlog `SPEC-008` ou sugerir reduzir carga), sem aumentar escopo.
- **FR-007**: System MUST manter comunicação clara e não punitiva: foco em tendência e próximo passo mínimo, não em culpa.

### Non-Functional Requirements

- **NFR-001**: System MUST manter baixa fricção (PRD RNF1): reforços devem ser pequenos, e a explicação deve caber em 1–2 frases.
- **NFR-002**: System MUST ser robusto a dias ruins (PRD RNF2): oferecer versão mínima do reforço e/ou adiar sem penalizar.
- **NFR-003**: System MUST manter segurança psicológica (PRD RNF3): alertas protetivos, sem humilhação.

### Key Entities *(include if feature involves data)*

- **FalseProgressSignal**: tipo (erro_recorrente/queda_rubrica/retrieval_fraco/gate_falhando/baixa_consistencia); evidência_resumo; janela; severidade (baixa/média/alta); timestamp.
- **ReinforcementAction**: tipo; descrição curta; critério_de_feito; evidência_minima; versão_minima (dia ruim); status (solicitado/feito/não_feito/tentativa).
- **SignalOutcome**: sinal; ação; resultado; nota curta do usuário (opcional); efeito observado (se houver).

## Acceptance Criteria *(mandatory)*

- O sistema detecta pelo menos os sinais MVP (erro recorrente e queda de qualidade) usando dados já coletados.
- Ao detectar sinal, o sistema solicita um reforço curto com critério observável e registra resultado.
- O sistema não cria overload: escolhe 1 foco por vez e oferece versão mínima em dias ruins.
- Linguagem é firme e não punitiva; o usuário entende “por que” e “qual o próximo passo mínimo”.

## Business Objectives *(mandatory)*

- Reduzir “fiz mas não aprendi” com detecção e reforço baseado em evidência (PRD R3; §5.4).
- Sustentar adaptação contínua do plano com sinais reais (PRD §5.5).
- Manter consistência sem burnout, com intervenções pequenas (RNF1/RNF2).

## Error Handling *(mandatory)*

- **Evidência insuficiente**: não afirmar falso progresso; pedir 1 evidência mínima no próximo passo e marcar como “sinal fraco”.
- **Sinais conflitantes**: priorizar o de maior impacto esperado e reduzir escopo.
- **Usuário some**: registrar pendência e retomar com o menor próximo passo quando voltar.
- **Rejeição do usuário**: registrar rejeição e reduzir repetição; levar para revisão semanal se persistente.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Redução na recorrência de erros-alvo ao longo de 2–4 semanas após reforços.
- **SC-002**: Aumento da taxa de tarefas com evidência válida e redução de falhas repetidas de gate.
- **SC-003**: Melhora ou estabilização de rubricas após ajustes/redução de escopo quando queda é detectada.
- **SC-004**: Manter fricção baixa: reforços geralmente concluíveis em ≤ 10 min (ou versão mínima em ≤ 2–5 min em dia ruim).