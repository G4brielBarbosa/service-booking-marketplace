# Feature Specification: Onboarding e Diagnóstico Leve (1–2 semanas)

**Created**: 2026-02-17  
**PRD Base**: §5.1, §§5.2–5.3, §6.2, §9.1, §10 (R1, R6, R7), §11 (RNF1–RNF4), §14

## Caso de uso *(mandatory)*

O usuário quer um assistente (Telegram-first) que reduza fricção e aumente consistência real em metas anuais (PRD §1–§3). Para isso, o produto precisa de um onboarding que produza insumos mínimos para operar e melhorar ao longo do tempo, sem virar burocracia.

Este onboarding deve:

- Coletar **metas e restrições** suficientes para calibrar a rotina diária (PRD §5.1; §9.1).
- Estabelecer uma **baseline leve** para medir tendências (PRD §5.1).
- Definir o **Mínimo Viável Diário (MVD)** para dias ruins, mantendo a identidade/hábito vivo (PRD §5.1; RNF2).
- Aplicar limites protetivos de metas em paralelo para evitar overload (PRD §6.2; §13).
- Respeitar **privacidade por padrão**: coletar o mínimo e explicar o que é guardado e por quê (PRD RNF4).

O onboarding é observado em dois marcos:

1) **Onboarding mínimo** (rápido): baseline mínima + MVD + metas ativas do ciclo + resumo revisável + próximo passo claro para iniciar rotina diária (PRD §5.1; §14).  
2) **Diagnóstico leve (1–2 semanas)**: completar lacunas de baseline por domínio em blocos curtos, sem “protocolo clínico” (PRD §5.1).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Concluir onboarding mínimo e ficar pronto para usar o assistente (Priority: P1)

O usuário inicia o assistente pela primeira vez e quer completar um onboarding **rápido** que já o deixe pronto para usar a rotina diária. Ao final do onboarding mínimo, ele quer:

- metas ativas do ciclo atual definidas (com limites de overload),
- baseline mínima registrada (por domínio),
- um MVD para dia ruim,
- um resumo revisável do onboarding,
- e um próximo passo único (ex.: “começar a rotina diária amanhã” ou “fazer o check-in”) sem depender de detalhes técnicos.

**Why this priority**: É o slice mínimo do MVP (PRD §14) para destravar execução com baixa fricção e permitir calibração/medição desde o início (PRD §5.1; RNF1; R6).

**Independent Test**: Com um usuário “novo” (sem dados prévios), executar o onboarding mínimo e validar que todos os artefatos acima existem e são consultáveis, sem depender de revisão semanal ou de outras features.

**Acceptance Scenarios**:

1. **Scenario**: Onboarding mínimo completo em uma sessão
   - **Given** usuário novo sem dados prévios (PRD §5.1)
   - **When** informa metas anuais, restrições mínimas e baseline leve solicitada
   - **Then** o sistema registra baseline mínima, define metas ativas do ciclo, define MVD e mostra um resumo revisável com o próximo passo (PRD §5.1; §6.2; RNF1)

2. **Scenario**: Onboarding parcial (tempo curto) não bloqueia o começo
   - **Given** usuário com pouco tempo no momento (PRD §4)
   - **When** completa apenas o conjunto mínimo obrigatório
   - **Then** o sistema conclui o onboarding mínimo, registra pendências e permite completar em etapas (PRD RNF1)

3. **Scenario**: Usuário tenta ativar metas intensivas demais
   - **Given** usuário seleciona mais de 2 metas intensivas para o ciclo atual (PRD §6.2)
   - **When** tenta concluir a seleção de metas ativas
   - **Then** o sistema sinaliza o limite, explica o motivo (proteger consistência) e orienta uma escolha (manter/adiar/pausar) com linguagem não punitiva (PRD §6.2; §13; RNF3)

4. **Scenario**: Definição do MVD para dia ruim
   - **Given** usuário concluiu o onboarding mínimo (PRD §5.1)
   - **When** o sistema apresenta o MVD
   - **Then** o MVD é curto, executável em baixa energia/tempo e o usuário entende quando usar (PRD RNF2; §5.1)

5. **Scenario**: Resumo do onboarding é consultável depois
   - **Given** onboarding mínimo concluído
   - **When** o usuário pede para ver seu resumo de onboarding
   - **Then** o sistema apresenta metas ativas, restrições principais, baseline mínima (por domínio) e MVD, de forma curta (PRD RNF1)

6. **Scenario**: Privacidade por padrão é explicada no onboarding
   - **Given** o sistema coletou dados de baseline (PRD §5.1)
   - **When** o usuário pergunta “o que você guarda e por quê?”
   - **Then** o sistema explica o mínimo guardado para medir progresso e calibrar planos, sem jargão técnico (PRD RNF4)

7. **Scenario**: Recusa de dado sensível não impede o onboarding mínimo
   - **Given** o usuário recusa compartilhar um tipo de dado (ex.: áudio) por privacidade/ambiente (PRD RNF4)
   - **When** o sistema solicita esse item de baseline
   - **Then** o sistema conclui o onboarding mínimo sem punição, registra a lacuna como pendência/bloqueio e explica o impacto de forma curta (PRD RNF1; RNF3) **[NEEDS CLARIFICATION]**: alternativa aceitável para baseline quando áudio não é possível

8. **Scenario**: Usuário não consegue definir metas com clareza
   - **Given** o usuário não sabe descrever metas com precisão (PRD §5.1)
   - **When** tenta avançar no onboarding
   - **Then** o sistema permite escolher metas a partir da lista suportada no PRD e seguir com o mínimo necessário, sem bloquear (PRD RNF1)

---

### User Story 2 - Completar diagnóstico leve progressivo (1–2 semanas) sem virar burocracia (Priority: P2)

O usuário quer completar a baseline por domínio ao longo de 1–2 semanas, em blocos curtos e objetivos. Ele quer que o sistema:

- colete lacunas aos poucos (um passo por vez),
- evite instrumentos clínicos não especificados,
- e mantenha segurança psicológica mesmo quando a baseline indicar dificuldades.

**Why this priority**: O PRD define onboarding/diagnóstico como 1–2 semanas e lista baselines por domínio (PRD §5.1). Fazer isso progressivamente aumenta chance de conclusão e reduz abandono (RNF1; RNF3).

**Independent Test**: Após concluir onboarding mínimo com pendências, simular 3–5 dias de coleta incremental e validar que o sistema captura ao menos 1 lacuna/dia (quando o usuário aceitar) e mantém o resumo atualizado.

**Acceptance Scenarios**:

1. **Scenario**: Coleta incremental por domínio em passo único
   - **Given** onboarding mínimo concluído com pendências de baseline (PRD §5.1)
   - **When** o usuário aceita completar uma pendência
   - **Then** o sistema coleta apenas 1 item de baseline e registra como completado, mantendo as demais pendências (PRD RNF1)

2. **Scenario**: Segurança psicológica diante de baseline “baixa”
   - **Given** o usuário reporta baixa energia/alto estresse/dificuldades (PRD §4)
   - **When** responde ao diagnóstico leve
   - **Then** o sistema reconhece o contexto, evita culpa e propõe o próximo passo mínimo viável (PRD RNF3; RNF2)

---

### User Story 3 - Retomar onboarding interrompido e revisar respostas (Priority: P2)

O usuário pode interromper o onboarding e voltar depois. Ele também pode perceber que respondeu algo errado e quer corrigir, sem “resetar” tudo.

**Why this priority**: Interrupções e correções são normais; suportar retomada e revisão reduz fricção e aumenta confiabilidade do sistema (PRD RNF1).

**Independent Test**: Iniciar onboarding, interromper, retomar e revisar uma resposta; verificar que o sistema mantém histórico do que já foi coletado e atualiza o resumo.

**Acceptance Scenarios**:

1. **Scenario**: Retomada após interrupção
   - **Given** existe uma sessão de onboarding em progresso
   - **When** o usuário retorna e pede para continuar
   - **Then** o sistema retoma do ponto certo e apresenta um único próximo passo curto (PRD RNF1)

2. **Scenario**: Revisão de resposta altera o resumo
   - **Given** onboarding mínimo concluído
   - **When** o usuário revisa uma informação (ex.: restrição de horário)
   - **Then** o resumo revisável reflete a alteração e o sistema explica, de forma concisa, impacto esperado (PRD §5.2; RNF1) **[NEEDS CLARIFICATION]**: quais campos podem ser revisados sem afetar métricas históricas.

### Edge Cases *(mandatory)*

- What happens when o usuário **não sabe definir metas anuais** com clareza (PRD §5.1)? O sistema deve permitir começar com metas provisórias guiadas pela lista do PRD. **[NEEDS CLARIFICATION]**: existe uma forma preferida de enquadrar metas (ex.: “resultado trimestral”) além do texto do PRD?
- How does system handle respostas conflitantes (ex.: “pouco tempo” e “muitas metas intensivas”) (PRD §6.2; RNF2)? O sistema deve priorizar proteção contra overload e pedir uma única escolha. **[NEEDS CLARIFICATION]**: bloquear ou apenas recomendar quando o limite é violado.
- What happens when o usuário **recusa compartilhar certos dados** (ex.: áudio de speaking) (PRD RNF4)? **[NEEDS CLARIFICATION]**: alternativas aceitáveis para baseline de inglês sem áudio e se isso conta como “baseline suficiente”.
- What happens when o usuário quer pular tudo (“só me manda o plano”) (PRD RNF1)? **[NEEDS CLARIFICATION]**: qual é o mínimo absoluto de onboarding para destravar rotina diária sem comprometer qualidade.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST coletar as metas anuais suportadas e permitir selecionar quais metas estarão **ativas no ciclo atual** (PRD §5.1).
- **FR-002**: System MUST aplicar e comunicar o limite de **no máximo 2 metas intensivas por ciclo**, orientando o usuário a manter/adiar/pausar metas para evitar overload (PRD §6.2; §13). **[NEEDS CLARIFICATION]**: política de classificação de “meta intensiva”.
- **FR-003**: System MUST coletar restrições operacionais mínimas (tempo típico, horários, limitações relevantes) para calibrar planos (PRD §5.1).
- **FR-004**: System MUST capturar uma baseline mínima por domínio conforme PRD §5.1 (sono/inglês/java/autoestima/contexto), sem inventar instrumentos clínicos. **[NEEDS CLARIFICATION]**: quais campos são obrigatórios no onboarding mínimo vs no diagnóstico progressivo.
- **FR-005**: System MUST suportar baseline de sono em formato de diário mínimo (PRD §5.1; §8.3).
- **FR-006**: System MUST suportar baseline de inglês por amostras curtas de fala e compreensão (PRD §5.1; §8.1). **[NEEDS CLARIFICATION]**: evidência alternativa quando áudio não é possível.
- **FR-007**: System MUST suportar baseline de Java por exercícios curtos de conceitos + código (PRD §5.1; §8.2).
- **FR-008**: System MUST suportar baseline de autoestima por registros de padrões de autocrítica e gatilhos (PRD §5.1; §8.5).
- **FR-009**: System MUST definir um **Mínimo Viável Diário (MVD)** executável em dia ruim e explicar quando usar (PRD §5.1; RNF2).
- **FR-010**: System MUST produzir um **resumo revisável** do onboarding (metas ativas, restrições, baseline por domínio e MVD) (PRD §5.1).
- **FR-011**: System MUST permitir concluir onboarding mínimo em uma sessão curta e completar lacunas em etapas ao longo de 1–2 semanas (PRD §5.1).
- **FR-012**: System MUST permitir retomar onboarding interrompido sem perda de dados e permitir revisão de respostas, refletindo no resumo (PRD RNF1; §5.2).
- **FR-013**: System MUST operar com privacidade por padrão durante onboarding: coletar o mínimo, explicar o que é guardado e por quê (PRD RNF4). **[NEEDS CLARIFICATION]**: política de retenção/remoção dos dados de baseline.

### Non-Functional Requirements

- **NFR-001**: System MUST minimizar digitação e manter interação curta e previsível durante onboarding (PRD RNF1).
- **NFR-002**: System MUST ser robusto a dias ruins: permitir onboarding mínimo com baixa energia/tempo e completar pendências depois (PRD RNF2).
- **NFR-003**: System MUST manter segurança psicológica: feedback firme, sem culpa/humilhação, especialmente quando baseline indicar dificuldades (PRD RNF3).
- **NFR-004**: System MUST aplicar privacidade por padrão: coletar o mínimo necessário e ser transparente sobre uso/retensão (PRD RNF4).

### Key Entities *(include if feature involves data)*

- **OnboardingSession**: status (novo/em progresso/concluído mínimo/concluído); respostas; pendências; timestamps; próximo passo de retomada (PRD §5.1).
- **ActiveGoalCycle**: metas ativas; metas pausadas/adiadas; indicação de metas intensivas; justificativa curta (PRD §6.2).
- **BaselineSnapshot**: baseline por domínio (sono/inglês/java/autoestima/contexto) + completude por domínio (PRD §5.1).
- **MinimumViableDaily (MVD)**: lista curta de ações; condições de uso (PRD §5.1; RNF2).
- **PrivacyDisclosure**: declaração do que é guardado e por quê; preferências/recusas por tipo de dado (PRD RNF4). **[NEEDS CLARIFICATION]**: formato/política de opt-out por dado.

## Acceptance Criteria *(mandatory)*

- Ao final do onboarding mínimo, existe: metas ativas (respeitando limite de metas intensivas), baseline mínima por domínio, MVD definido, resumo revisável e próximo passo claro (PRD §5.1; §6.2; §14).
- O onboarding mínimo pode ser concluído mesmo com pouco tempo, registrando pendências para completar em etapas (PRD RNF1).
- O usuário consegue retomar onboarding interrompido sem perder respostas e consegue revisar informações, com resumo atualizado (PRD RNF1; §5.2).
- O sistema explica privacidade por padrão e o propósito dos dados coletados (PRD RNF4).

## Business Objectives *(mandatory)*

- **Baixa fricção inicial**: reduzir esforço para começar e chegar rapidamente à rotina diária (PRD §5.1; RNF1; §14).
- **Calibração e adaptação**: baseline e restrições habilitam planos A/B/C e ajustes ao longo do tempo (PRD §5.2; R1; R7).
- **Proteção contra overload**: impor/explicitar limites de metas intensivas desde o onboarding (PRD §6.2; §13).
- **Confiança e privacidade**: transparência e minimização de dados desde a primeira interação (PRD RNF4).

## Error Handling *(mandatory)*

- **Usuário some**: manter sessão e permitir retomar; ao retornar, mostrar estado e pedir apenas o próximo passo (PRD RNF1; §5.1).
- **Dados insuficientes**: concluir onboarding mínimo quando possível e registrar lacunas para completar no diagnóstico progressivo; sem culpa (PRD RNF3).
- **Respostas ambíguas**: fazer no máximo 1 pergunta de confirmação por vez e seguir (PRD RNF1).
- **Conflito/overload**: priorizar proteção e orientar escolha/adiamento com linguagem não punitiva (PRD §6.2; §13; RNF3). **[NEEDS CLARIFICATION]**: se o sistema pode bloquear ativação até ajuste.
- **Recusa por privacidade**: explicar impacto, oferecer alternativa (quando definida) ou registrar pendência/bloqueio sem punição (PRD RNF4; RNF3). **[NEEDS CLARIFICATION]**: alternativas para baseline de inglês e política de retenção.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Usuário completa onboarding mínimo em \u2264 10 minutos totais (somando interações) (PRD RNF1; §5.1).
- **SC-002**: Usuário consegue revisar o resumo do onboarding a qualquer momento (PRD §5.1).
- **SC-003**: Em até 14 dias, a baseline progressiva atinge alta completude sem virar burocracia. **[NEEDS CLARIFICATION]**: definição objetiva de “completude” por domínio.
- **SC-004**: A maioria dos ciclos iniciados respeita o limite de 2 metas intensivas (PRD §6.2).

