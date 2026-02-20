package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/abriesouza/super-assistente/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

// --- UserRepository ---

func (s *Store) FindByTelegramID(ctx context.Context, telegramUserID int64) (*domain.UserProfile, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT user_id, telegram_user_id, primary_chat_id, timezone, locale, created_at, updated_at
		 FROM user_profiles WHERE telegram_user_id = $1`, telegramUserID)

	var u domain.UserProfile
	err := row.Scan(&u.UserID, &u.TelegramUserID, &u.PrimaryChatID, &u.Timezone, &u.Locale, &u.CreatedAt, &u.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found: telegram_user_id=%d", telegramUserID)
	}
	return &u, err
}

func (s *Store) SaveUser(ctx context.Context, user *domain.UserProfile) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO user_profiles (user_id, telegram_user_id, primary_chat_id, timezone, locale, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		user.UserID, user.TelegramUserID, user.PrimaryChatID, user.Timezone, user.Locale, user.CreatedAt, user.UpdatedAt)
	return err
}

func (s *Store) UpdateUser(ctx context.Context, user *domain.UserProfile) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE user_profiles SET timezone=$2, locale=$3, updated_at=now() WHERE user_id=$1`,
		user.UserID, user.Timezone, user.Locale)
	return err
}

// --- OnboardingRepository ---

func (s *Store) FindOnboardingByUserID(ctx context.Context, userID uuid.UUID) (*domain.OnboardingSession, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT session_id, user_id, status, current_step_id, answers, pending_items,
		        started_at, last_interaction_at, completed_at
		 FROM onboarding_sessions WHERE user_id = $1`, userID)

	var sess domain.OnboardingSession
	var answersJSON, pendingJSON []byte
	var completedAt *time.Time

	err := row.Scan(&sess.SessionID, &sess.UserID, &sess.Status, &sess.CurrentStepID,
		&answersJSON, &pendingJSON, &sess.StartedAt, &sess.LastInteraction, &completedAt)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("onboarding session not found")
	}
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(answersJSON, &sess.Answers)
	_ = json.Unmarshal(pendingJSON, &sess.PendingItems)
	sess.CompletedAt = completedAt
	return &sess, nil
}

func (s *Store) SaveOnboarding(ctx context.Context, session *domain.OnboardingSession) error {
	answersJSON, _ := json.Marshal(session.Answers)
	pendingJSON, _ := json.Marshal(session.PendingItems)

	_, err := s.pool.Exec(ctx,
		`INSERT INTO onboarding_sessions (session_id, user_id, status, current_step_id, answers, pending_items, started_at, last_interaction_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		session.SessionID, session.UserID, session.Status, session.CurrentStepID,
		answersJSON, pendingJSON, session.StartedAt, session.LastInteraction)
	return err
}

func (s *Store) UpdateOnboarding(ctx context.Context, session *domain.OnboardingSession) error {
	answersJSON, _ := json.Marshal(session.Answers)
	pendingJSON, _ := json.Marshal(session.PendingItems)

	_, err := s.pool.Exec(ctx,
		`UPDATE onboarding_sessions
		 SET status=$2, current_step_id=$3, answers=$4, pending_items=$5,
		     last_interaction_at=$6, completed_at=$7
		 WHERE session_id=$1`,
		session.SessionID, session.Status, session.CurrentStepID,
		answersJSON, pendingJSON, session.LastInteraction, session.CompletedAt)
	return err
}

// --- GoalCycleRepository ---

func (s *Store) FindGoalCycleByUserID(ctx context.Context, userID uuid.UUID) (*domain.ActiveGoalCycle, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT cycle_id, user_id, active_goals, paused_goals, started_at, updated_at
		 FROM active_goal_cycles WHERE user_id = $1`, userID)

	var c domain.ActiveGoalCycle
	var activeJSON, pausedJSON []byte

	err := row.Scan(&c.CycleID, &c.UserID, &activeJSON, &pausedJSON, &c.StartedAt, &c.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("goal cycle not found")
	}
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(activeJSON, &c.ActiveGoals)
	_ = json.Unmarshal(pausedJSON, &c.PausedGoals)
	return &c, nil
}

func (s *Store) SaveGoalCycle(ctx context.Context, cycle *domain.ActiveGoalCycle) error {
	activeJSON, _ := json.Marshal(cycle.ActiveGoals)
	pausedJSON, _ := json.Marshal(cycle.PausedGoals)

	_, err := s.pool.Exec(ctx,
		`INSERT INTO active_goal_cycles (cycle_id, user_id, active_goals, paused_goals, started_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		cycle.CycleID, cycle.UserID, activeJSON, pausedJSON, cycle.StartedAt, cycle.UpdatedAt)
	return err
}

func (s *Store) UpdateGoalCycle(ctx context.Context, cycle *domain.ActiveGoalCycle) error {
	activeJSON, _ := json.Marshal(cycle.ActiveGoals)
	pausedJSON, _ := json.Marshal(cycle.PausedGoals)

	_, err := s.pool.Exec(ctx,
		`UPDATE active_goal_cycles SET active_goals=$2, paused_goals=$3, updated_at=now() WHERE cycle_id=$1`,
		cycle.CycleID, activeJSON, pausedJSON)
	return err
}

// --- PrivacyPolicyRepository ---

func (s *Store) FindPrivacyPolicy(ctx context.Context, userID uuid.UUID) (*domain.PrivacyPolicy, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT user_id, opt_out_categories, retention_days, minimal_mode, updated_at
		 FROM privacy_policies WHERE user_id = $1`, userID)

	var pp domain.PrivacyPolicy
	var optOutArr []string
	var retentionJSON []byte

	err := row.Scan(&pp.UserID, &optOutArr, &retentionJSON, &pp.MinimalMode, &pp.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("privacy policy not found")
	}
	if err != nil {
		return nil, err
	}

	pp.OptOutCategories = make([]domain.SensitivityLevel, len(optOutArr))
	for i, v := range optOutArr {
		pp.OptOutCategories[i] = domain.SensitivityLevel(v)
	}

	pp.RetentionDays = map[domain.SensitivityLevel]int{}
	_ = json.Unmarshal(retentionJSON, &pp.RetentionDays)
	return &pp, nil
}

func (s *Store) SavePrivacyPolicy(ctx context.Context, policy *domain.PrivacyPolicy) error {
	optOut := make([]string, len(policy.OptOutCategories))
	for i, v := range policy.OptOutCategories {
		optOut[i] = string(v)
	}
	retJSON, _ := json.Marshal(policy.RetentionDays)

	_, err := s.pool.Exec(ctx,
		`INSERT INTO privacy_policies (user_id, opt_out_categories, retention_days, minimal_mode, updated_at)
		 VALUES ($1, $2, $3, $4, now())
		 ON CONFLICT (user_id) DO UPDATE SET
		     opt_out_categories = EXCLUDED.opt_out_categories,
		     retention_days = EXCLUDED.retention_days,
		     minimal_mode = EXCLUDED.minimal_mode,
		     updated_at = now()`,
		policy.UserID, optOut, retJSON, policy.MinimalMode)
	return err
}

func (s *Store) UpdatePrivacyPolicy(ctx context.Context, policy *domain.PrivacyPolicy) error {
	return s.SavePrivacyPolicy(ctx, policy)
}

// --- BaselineRepository ---

func (s *Store) FindBaselineByUserAndDomain(ctx context.Context, userID uuid.UUID, d domain.GoalID) (*domain.BaselineSnapshot, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT baseline_id, user_id, domain, data, completeness, captured_at, updated_at
		 FROM baseline_snapshots WHERE user_id = $1 AND domain = $2`, userID, string(d))

	var bs domain.BaselineSnapshot
	var dataJSON []byte

	err := row.Scan(&bs.BaselineID, &bs.UserID, &bs.Domain, &dataJSON, &bs.Completeness, &bs.CapturedAt, &bs.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("baseline not found")
	}
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(dataJSON, &bs.Data)
	return &bs, nil
}

func (s *Store) FindAllBaselinesByUser(ctx context.Context, userID uuid.UUID) ([]domain.BaselineSnapshot, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT baseline_id, user_id, domain, data, completeness, captured_at, updated_at
		 FROM baseline_snapshots WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.BaselineSnapshot
	for rows.Next() {
		var bs domain.BaselineSnapshot
		var dataJSON []byte
		if err := rows.Scan(&bs.BaselineID, &bs.UserID, &bs.Domain, &dataJSON, &bs.Completeness, &bs.CapturedAt, &bs.UpdatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(dataJSON, &bs.Data)
		results = append(results, bs)
	}
	return results, nil
}

func (s *Store) SaveBaseline(ctx context.Context, baseline *domain.BaselineSnapshot) error {
	dataJSON, _ := json.Marshal(baseline.Data)

	_, err := s.pool.Exec(ctx,
		`INSERT INTO baseline_snapshots (baseline_id, user_id, domain, data, completeness, captured_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (user_id, domain) DO UPDATE SET
		     data = EXCLUDED.data, completeness = EXCLUDED.completeness, updated_at = now()`,
		baseline.BaselineID, baseline.UserID, string(baseline.Domain), dataJSON,
		string(baseline.Completeness), baseline.CapturedAt, baseline.UpdatedAt)
	return err
}

func (s *Store) UpdateBaseline(ctx context.Context, baseline *domain.BaselineSnapshot) error {
	return s.SaveBaseline(ctx, baseline)
}

// --- MVDRepository ---

func (s *Store) FindMVDByUserID(ctx context.Context, userID uuid.UUID) (*domain.MinimumViableDaily, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT mvd_id, user_id, items, when_to_use, updated_at
		 FROM minimum_viable_dailies WHERE user_id = $1`, userID)

	var mvd domain.MinimumViableDaily
	var itemsJSON []byte

	err := row.Scan(&mvd.MVDID, &mvd.UserID, &itemsJSON, &mvd.WhenToUse, &mvd.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("mvd not found")
	}
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(itemsJSON, &mvd.Items)
	return &mvd, nil
}

func (s *Store) SaveMVD(ctx context.Context, mvd *domain.MinimumViableDaily) error {
	itemsJSON, _ := json.Marshal(mvd.Items)

	_, err := s.pool.Exec(ctx,
		`INSERT INTO minimum_viable_dailies (mvd_id, user_id, items, when_to_use, updated_at)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (user_id) DO UPDATE SET
		     items = EXCLUDED.items, when_to_use = EXCLUDED.when_to_use, updated_at = now()`,
		mvd.MVDID, mvd.UserID, itemsJSON, mvd.WhenToUse, mvd.UpdatedAt)
	return err
}

func (s *Store) UpdateMVD(ctx context.Context, mvd *domain.MinimumViableDaily) error {
	return s.SaveMVD(ctx, mvd)
}

// --- EventRepository ---

func (s *Store) AppendEvent(ctx context.Context, event domain.DomainEvent) error {
	payloadJSON, _ := json.Marshal(event.PayloadMin)

	_, err := s.pool.Exec(ctx,
		`INSERT INTO domain_event_log (event_id, user_id, timestamp, local_date, week_id, event_type, payload_min, sensitivity)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		event.EventID, event.UserID, event.Timestamp, event.LocalDate,
		event.WeekID, string(event.Type), payloadJSON, string(event.Sensitivity))
	return err
}

// --- IdempotencyStore ---

func (s *Store) CheckIdempotency(ctx context.Context, key string) (*domain.IdempotencyRecord, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT key, first_seen_at, result_ref, expires_at
		 FROM idempotency_records WHERE key = $1 AND expires_at > now()`, key)

	var rec domain.IdempotencyRecord
	err := row.Scan(&rec.Key, &rec.FirstSeen, &rec.ResultRef, &rec.ExpiresAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &rec, err
}

func (s *Store) StoreIdempotency(ctx context.Context, record domain.IdempotencyRecord) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO idempotency_records (key, first_seen_at, result_ref, expires_at)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (key) DO NOTHING`,
		record.Key, record.FirstSeen, record.ResultRef, record.ExpiresAt)
	return err
}

// --- DailyStateRepository ---

func (s *Store) FindDailyStateByUserAndDate(ctx context.Context, userID uuid.UUID, localDate string) (*domain.DailyState, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT user_id, local_date, created_at, updated_at
		 FROM daily_states WHERE user_id = $1 AND local_date = $2`, userID, localDate)

	var ds domain.DailyState
	err := row.Scan(&ds.UserID, &ds.LocalDate, &ds.CreatedAt, &ds.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("daily state not found")
	}
	if err != nil {
		return nil, err
	}

	checkIn, _ := s.findCheckIn(ctx, userID, localDate)
	ds.CheckIn = checkIn

	plan, _ := s.findLatestPlan(ctx, userID, localDate)
	ds.Plan = plan

	tasks, _ := s.findTasks(ctx, userID, localDate)
	ds.Tasks = tasks

	return &ds, nil
}

func (s *Store) findCheckIn(ctx context.Context, userID uuid.UUID, localDate string) (*domain.DailyCheckIn, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT check_in_id, user_id, local_date, time_available_min, energy_0_10,
		        mood_stress_0_10, constraints_text, created_at
		 FROM daily_check_ins WHERE user_id = $1 AND local_date = $2`, userID, localDate)

	var ci domain.DailyCheckIn
	var mood *int
	var constr *string
	err := row.Scan(&ci.CheckInID, &ci.UserID, &ci.LocalDate, &ci.TimeAvailMin,
		&ci.Energy, &mood, &constr, &ci.CreatedAt)
	if err != nil {
		return nil, err
	}
	ci.MoodStress = mood
	ci.ConstraintText = constr
	return &ci, nil
}

func (s *Store) findLatestPlan(ctx context.Context, userID uuid.UUID, localDate string) (*domain.DailyPlan, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT plan_id, user_id, local_date, plan_type, rationale,
		        priority_task_id, complementary_ids, foundation_task_id, version, created_at
		 FROM daily_plans WHERE user_id = $1 AND local_date = $2
		 ORDER BY version DESC LIMIT 1`, userID, localDate)

	var dp domain.DailyPlan
	var compIDs []uuid.UUID
	var foundationID *uuid.UUID
	err := row.Scan(&dp.PlanID, &dp.UserID, &dp.LocalDate, &dp.PlanType, &dp.Rationale,
		&dp.PriorityTaskID, &compIDs, &foundationID, &dp.Version, &dp.CreatedAt)
	if err != nil {
		return nil, err
	}
	dp.ComplementaryIDs = compIDs
	dp.FoundationTaskID = foundationID
	return &dp, nil
}

func (s *Store) findTasks(ctx context.Context, userID uuid.UUID, localDate string) ([]domain.PlannedTask, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT task_id, user_id, local_date, title, goal_domain, estimated_min,
		        instructions, done_criteria, status, block_reason, note, gate_ref,
		        created_at, updated_at
		 FROM planned_tasks WHERE user_id = $1 AND local_date = $2`, userID, localDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.PlannedTask
	for rows.Next() {
		var t domain.PlannedTask
		if err := rows.Scan(&t.TaskID, &t.UserID, &t.LocalDate, &t.Title, &t.GoalDomain,
			&t.EstimatedMin, &t.Instructions, &t.DoneCriteria, &t.Status,
			&t.BlockReason, &t.Note, &t.GateRef, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *Store) SaveDailyState(ctx context.Context, state *domain.DailyState) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO daily_states (user_id, local_date, created_at, updated_at)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (user_id, local_date) DO UPDATE SET updated_at = EXCLUDED.updated_at`,
		state.UserID, state.LocalDate, state.CreatedAt, state.UpdatedAt)
	if err != nil {
		return err
	}

	if state.CheckIn != nil {
		_, err = tx.Exec(ctx,
			`INSERT INTO daily_check_ins (check_in_id, user_id, local_date, time_available_min, energy_0_10, mood_stress_0_10, constraints_text, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			 ON CONFLICT (user_id, local_date) DO UPDATE SET
			     time_available_min = EXCLUDED.time_available_min,
			     energy_0_10 = EXCLUDED.energy_0_10,
			     mood_stress_0_10 = EXCLUDED.mood_stress_0_10,
			     constraints_text = EXCLUDED.constraints_text`,
			state.CheckIn.CheckInID, state.CheckIn.UserID, state.CheckIn.LocalDate,
			state.CheckIn.TimeAvailMin, state.CheckIn.Energy,
			state.CheckIn.MoodStress, state.CheckIn.ConstraintText, state.CheckIn.CreatedAt)
		if err != nil {
			return err
		}
	}

	if state.Plan != nil {
		compIDs := state.Plan.ComplementaryIDs
		if compIDs == nil {
			compIDs = []uuid.UUID{}
		}
		_, err = tx.Exec(ctx,
			`INSERT INTO daily_plans (plan_id, user_id, local_date, plan_type, rationale, priority_task_id, complementary_ids, foundation_task_id, version, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			state.Plan.PlanID, state.Plan.UserID, state.Plan.LocalDate,
			string(state.Plan.PlanType), state.Plan.Rationale,
			state.Plan.PriorityTaskID, compIDs,
			state.Plan.FoundationTaskID, state.Plan.Version, state.Plan.CreatedAt)
		if err != nil {
			return err
		}
	}

	for _, t := range state.Tasks {
		_, err = tx.Exec(ctx,
			`INSERT INTO planned_tasks (task_id, user_id, local_date, title, goal_domain, estimated_min, instructions, done_criteria, status, block_reason, note, gate_ref, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			 ON CONFLICT (task_id) DO UPDATE SET
			     status = EXCLUDED.status,
			     block_reason = EXCLUDED.block_reason,
			     note = EXCLUDED.note,
			     updated_at = EXCLUDED.updated_at`,
			t.TaskID, t.UserID, t.LocalDate, t.Title, string(t.GoalDomain),
			t.EstimatedMin, t.Instructions, t.DoneCriteria, string(t.Status),
			t.BlockReason, t.Note, t.GateRef, t.CreatedAt, t.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *Store) UpdateDailyState(ctx context.Context, state *domain.DailyState) error {
	return s.SaveDailyState(ctx, state)
}

// --- EvidenceRepository ---

func (s *Store) SaveEvidence(ctx context.Context, ev *domain.Evidence) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO evidences (evidence_id, user_id, task_id, kind, sensitivity, storage_policy, content_ref, summary, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		ev.EvidenceID, ev.UserID, ev.TaskID,
		string(ev.Kind), string(ev.Sensitivity), string(ev.StoragePolicy),
		ev.ContentRef, ev.Summary, ev.Timestamp)
	return err
}

func (s *Store) FindEvidenceByTaskID(ctx context.Context, taskID uuid.UUID) ([]domain.Evidence, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT evidence_id, user_id, task_id, kind, sensitivity, storage_policy, content_ref, summary, created_at
		 FROM evidences WHERE task_id = $1 ORDER BY created_at`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Evidence
	for rows.Next() {
		var ev domain.Evidence
		if err := rows.Scan(&ev.EvidenceID, &ev.UserID, &ev.TaskID,
			&ev.Kind, &ev.Sensitivity, &ev.StoragePolicy,
			&ev.ContentRef, &ev.Summary, &ev.Timestamp); err != nil {
			return nil, err
		}
		results = append(results, ev)
	}
	return results, nil
}

// --- GateResultRepository ---

func (s *Store) SaveGateResult(ctx context.Context, gr *domain.GateResult) error {
	metricsJSON, _ := json.Marshal(gr.DerivedMetrics)
	eids := gr.EvidenceIDs
	if eids == nil {
		eids = []uuid.UUID{}
	}

	_, err := s.pool.Exec(ctx,
		`INSERT INTO gate_results (gate_result_id, user_id, task_id, gate_status, failure_reason_code, reason_short, next_min_step, evidence_ids, derived_metrics, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 ON CONFLICT (gate_result_id) DO UPDATE SET
		     gate_status = EXCLUDED.gate_status,
		     failure_reason_code = EXCLUDED.failure_reason_code,
		     reason_short = EXCLUDED.reason_short,
		     next_min_step = EXCLUDED.next_min_step,
		     evidence_ids = EXCLUDED.evidence_ids,
		     derived_metrics = EXCLUDED.derived_metrics`,
		gr.GateResultID, gr.UserID, gr.TaskID,
		string(gr.GateStatus), string(gr.FailureReason),
		gr.ReasonShort, gr.NextMinStep, eids, metricsJSON, gr.Timestamp)
	return err
}

func (s *Store) FindGateResultByTaskID(ctx context.Context, taskID uuid.UUID) (*domain.GateResult, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT gate_result_id, user_id, task_id, gate_status, failure_reason_code, reason_short, next_min_step, evidence_ids, derived_metrics, created_at
		 FROM gate_results WHERE task_id = $1 ORDER BY created_at DESC LIMIT 1`, taskID)

	var gr domain.GateResult
	var metricsJSON []byte

	err := row.Scan(&gr.GateResultID, &gr.UserID, &gr.TaskID,
		&gr.GateStatus, &gr.FailureReason, &gr.ReasonShort, &gr.NextMinStep,
		&gr.EvidenceIDs, &metricsJSON, &gr.Timestamp)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	gr.DerivedMetrics = map[string]any{}
	_ = json.Unmarshal(metricsJSON, &gr.DerivedMetrics)
	return &gr, nil
}

func (s *Store) FindGateResultsByUserAndDate(ctx context.Context, userID uuid.UUID, localDate string) ([]domain.GateResult, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT gr.gate_result_id, gr.user_id, gr.task_id, gr.gate_status, gr.failure_reason_code, gr.reason_short, gr.next_min_step, gr.evidence_ids, gr.derived_metrics, gr.created_at
		 FROM gate_results gr
		 JOIN planned_tasks pt ON gr.task_id = pt.task_id
		 WHERE gr.user_id = $1 AND pt.local_date = $2
		 ORDER BY gr.created_at`, userID, localDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.GateResult
	for rows.Next() {
		var gr domain.GateResult
		var metricsJSON []byte
		if err := rows.Scan(&gr.GateResultID, &gr.UserID, &gr.TaskID,
			&gr.GateStatus, &gr.FailureReason, &gr.ReasonShort, &gr.NextMinStep,
			&gr.EvidenceIDs, &metricsJSON, &gr.Timestamp); err != nil {
			return nil, err
		}
		gr.DerivedMetrics = map[string]any{}
		_ = json.Unmarshal(metricsJSON, &gr.DerivedMetrics)
		results = append(results, gr)
	}
	return results, nil
}

// --- RubricRepository ---

func (s *Store) SaveRubric(ctx context.Context, r *domain.RubricScore) error {
	dimsJSON, _ := json.Marshal(r.Dimensions)
	_, err := s.pool.Exec(ctx,
		`INSERT INTO rubric_scores (rubric_id, user_id, task_id, domain, dimensions, total, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, now())
		 ON CONFLICT (rubric_id) DO UPDATE SET
		     dimensions = EXCLUDED.dimensions,
		     total = EXCLUDED.total,
		     status = EXCLUDED.status`,
		r.RubricID, r.UserID, r.TaskID, string(r.Domain), dimsJSON, r.Total, r.Status)
	return err
}

func (s *Store) FindRubricByTaskID(ctx context.Context, taskID uuid.UUID) (*domain.RubricScore, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT rubric_id, user_id, task_id, domain, dimensions, total, status
		 FROM rubric_scores WHERE task_id = $1 LIMIT 1`, taskID)

	var r domain.RubricScore
	var dimsJSON []byte
	err := row.Scan(&r.RubricID, &r.UserID, &r.TaskID, &r.Domain, &dimsJSON, &r.Total, &r.Status)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	r.Dimensions = map[string]int{}
	_ = json.Unmarshal(dimsJSON, &r.Dimensions)
	return &r, nil
}

// --- EnglishInputRepository ---

func (s *Store) SaveEnglishInput(ctx context.Context, sess *domain.EnglishInputSession) error {
	answersJSON, _ := json.Marshal(sess.ComprehensionAnswers)
	_, err := s.pool.Exec(ctx,
		`INSERT INTO english_input_sessions (session_id, user_id, task_id, local_date, duration_est_min, content_descriptor, comprehension_answers, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 ON CONFLICT (session_id) DO UPDATE SET
		     comprehension_answers = EXCLUDED.comprehension_answers,
		     status = EXCLUDED.status`,
		sess.SessionID, sess.UserID, sess.TaskID, sess.LocalDate,
		sess.DurationEstMin, sess.ContentDescriptor, answersJSON, sess.Status, sess.CreatedAt)
	return err
}

func (s *Store) FindEnglishInputByTaskID(ctx context.Context, taskID uuid.UUID) (*domain.EnglishInputSession, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT session_id, user_id, task_id, local_date, duration_est_min, content_descriptor, comprehension_answers, status, created_at
		 FROM english_input_sessions WHERE task_id = $1 LIMIT 1`, taskID)

	var sess domain.EnglishInputSession
	var answersJSON []byte
	err := row.Scan(&sess.SessionID, &sess.UserID, &sess.TaskID, &sess.LocalDate,
		&sess.DurationEstMin, &sess.ContentDescriptor, &answersJSON, &sess.Status, &sess.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(answersJSON, &sess.ComprehensionAnswers)
	return &sess, nil
}

// --- EnglishRetrievalRepository ---

func (s *Store) SaveEnglishRetrieval(ctx context.Context, r *domain.EnglishRetrieval) error {
	targetsJSON, _ := json.Marshal(r.Targets)
	_, err := s.pool.Exec(ctx,
		`INSERT INTO english_retrievals (retrieval_id, user_id, task_id, local_date, items_answered, items_total, status, targets, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 ON CONFLICT (retrieval_id) DO UPDATE SET
		     items_answered = EXCLUDED.items_answered,
		     status = EXCLUDED.status,
		     targets = EXCLUDED.targets`,
		r.RetrievalID, r.UserID, r.TaskID, r.LocalDate,
		r.ItemsAnswered, r.ItemsTotal, r.Status, targetsJSON, r.CreatedAt)
	return err
}

func (s *Store) FindEnglishRetrievalByTaskID(ctx context.Context, taskID uuid.UUID) (*domain.EnglishRetrieval, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT retrieval_id, user_id, task_id, local_date, items_answered, items_total, status, targets, created_at
		 FROM english_retrievals WHERE task_id = $1 LIMIT 1`, taskID)

	var r domain.EnglishRetrieval
	var targetsJSON []byte
	err := row.Scan(&r.RetrievalID, &r.UserID, &r.TaskID, &r.LocalDate,
		&r.ItemsAnswered, &r.ItemsTotal, &r.Status, &targetsJSON, &r.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(targetsJSON, &r.Targets)
	return &r, nil
}

func (s *Store) FindEnglishRetrievalsByUserAndDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]domain.EnglishRetrieval, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT retrieval_id, user_id, task_id, local_date, items_answered, items_total, status, targets, created_at
		 FROM english_retrievals WHERE user_id = $1 AND local_date >= $2 AND local_date <= $3
		 ORDER BY created_at`, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.EnglishRetrieval
	for rows.Next() {
		var r domain.EnglishRetrieval
		var targetsJSON []byte
		if err := rows.Scan(&r.RetrievalID, &r.UserID, &r.TaskID, &r.LocalDate,
			&r.ItemsAnswered, &r.ItemsTotal, &r.Status, &targetsJSON, &r.CreatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(targetsJSON, &r.Targets)
		results = append(results, r)
	}
	return results, nil
}

// --- EnglishErrorLogRepository ---

func (s *Store) SaveEnglishErrorLog(ctx context.Context, entry *domain.EnglishErrorLogEntry) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO english_error_log (error_id, user_id, local_date, label, note_short, recurring_count_14d, is_recurring, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 ON CONFLICT (error_id) DO UPDATE SET
		     recurring_count_14d = EXCLUDED.recurring_count_14d,
		     is_recurring = EXCLUDED.is_recurring`,
		entry.ErrorID, entry.UserID, entry.LocalDate, entry.Label,
		entry.NoteShort, entry.RecurringCount14d, entry.IsRecurring, entry.CreatedAt)
	return err
}

func (s *Store) FindEnglishErrorsByUserAndDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]domain.EnglishErrorLogEntry, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT error_id, user_id, local_date, label, note_short, recurring_count_14d, is_recurring, created_at
		 FROM english_error_log WHERE user_id = $1 AND local_date >= $2 AND local_date <= $3
		 ORDER BY created_at`, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.EnglishErrorLogEntry
	for rows.Next() {
		var e domain.EnglishErrorLogEntry
		if err := rows.Scan(&e.ErrorID, &e.UserID, &e.LocalDate, &e.Label,
			&e.NoteShort, &e.RecurringCount14d, &e.IsRecurring, &e.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, e)
	}
	return results, nil
}

func (s *Store) CountEnglishErrorsByUserLabelSince(ctx context.Context, userID uuid.UUID, label string, since string) (int, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM english_error_log WHERE user_id = $1 AND label = $2 AND local_date >= $3`,
		userID, label, since)

	var count int
	err := row.Scan(&count)
	return count, err
}

// --- JavaPracticeRepository ---

func (s *Store) SaveJavaPractice(ctx context.Context, sess *domain.JavaPracticeSession) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO java_practice_sessions (session_id, user_id, task_id, local_date, duration_est_min, objective_constraint, evidence_short, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 ON CONFLICT (session_id) DO UPDATE SET
		     evidence_short = EXCLUDED.evidence_short,
		     status = EXCLUDED.status`,
		sess.SessionID, sess.UserID, sess.TaskID, sess.LocalDate,
		sess.DurationEstMin, sess.ObjectiveConstraint, sess.EvidenceShort, sess.Status, sess.CreatedAt)
	return err
}

func (s *Store) FindJavaPracticeByTaskID(ctx context.Context, taskID uuid.UUID) (*domain.JavaPracticeSession, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT session_id, user_id, task_id, local_date, duration_est_min, objective_constraint, evidence_short, status, created_at
		 FROM java_practice_sessions WHERE task_id = $1 LIMIT 1`, taskID)

	var sess domain.JavaPracticeSession
	err := row.Scan(&sess.SessionID, &sess.UserID, &sess.TaskID, &sess.LocalDate,
		&sess.DurationEstMin, &sess.ObjectiveConstraint, &sess.EvidenceShort, &sess.Status, &sess.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

// --- JavaRetrievalRepository ---

func (s *Store) SaveJavaRetrieval(ctx context.Context, r *domain.JavaRetrieval) error {
	targetsJSON, _ := json.Marshal(r.Targets)
	_, err := s.pool.Exec(ctx,
		`INSERT INTO java_retrievals (retrieval_id, user_id, task_id, local_date, items_answered, items_total, status, targets, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 ON CONFLICT (retrieval_id) DO UPDATE SET
		     items_answered = EXCLUDED.items_answered,
		     status = EXCLUDED.status,
		     targets = EXCLUDED.targets`,
		r.RetrievalID, r.UserID, r.TaskID, r.LocalDate,
		r.ItemsAnswered, r.ItemsTotal, r.Status, targetsJSON, r.CreatedAt)
	return err
}

func (s *Store) FindJavaRetrievalByTaskID(ctx context.Context, taskID uuid.UUID) (*domain.JavaRetrieval, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT retrieval_id, user_id, task_id, local_date, items_answered, items_total, status, targets, created_at
		 FROM java_retrievals WHERE task_id = $1 LIMIT 1`, taskID)

	var r domain.JavaRetrieval
	var targetsJSON []byte
	err := row.Scan(&r.RetrievalID, &r.UserID, &r.TaskID, &r.LocalDate,
		&r.ItemsAnswered, &r.ItemsTotal, &r.Status, &targetsJSON, &r.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(targetsJSON, &r.Targets)
	return &r, nil
}

func (s *Store) FindJavaRetrievalsByUserAndDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]domain.JavaRetrieval, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT retrieval_id, user_id, task_id, local_date, items_answered, items_total, status, targets, created_at
		 FROM java_retrievals WHERE user_id = $1 AND local_date >= $2 AND local_date <= $3
		 ORDER BY created_at`, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.JavaRetrieval
	for rows.Next() {
		var r domain.JavaRetrieval
		var targetsJSON []byte
		if err := rows.Scan(&r.RetrievalID, &r.UserID, &r.TaskID, &r.LocalDate,
			&r.ItemsAnswered, &r.ItemsTotal, &r.Status, &targetsJSON, &r.CreatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(targetsJSON, &r.Targets)
		results = append(results, r)
	}
	return results, nil
}

// --- JavaLearningLogRepository ---

func (s *Store) SaveJavaLearningLog(ctx context.Context, entry *domain.JavaLearningLogEntry) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO java_learning_log (entry_id, user_id, task_id, local_date, error_or_learning, fix_or_note, category, recurring_count_14d, is_recurring, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 ON CONFLICT (entry_id) DO UPDATE SET
		     recurring_count_14d = EXCLUDED.recurring_count_14d,
		     is_recurring = EXCLUDED.is_recurring`,
		entry.EntryID, entry.UserID, entry.TaskID, entry.LocalDate, entry.ErrorOrLearning,
		entry.FixOrNote, entry.Category, entry.RecurringCount14d, entry.IsRecurring, entry.CreatedAt)
	return err
}

func (s *Store) FindJavaLearningLogByUserAndDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]domain.JavaLearningLogEntry, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT entry_id, user_id, task_id, local_date, error_or_learning, fix_or_note, category, recurring_count_14d, is_recurring, created_at
		 FROM java_learning_log WHERE user_id = $1 AND local_date >= $2 AND local_date <= $3
		 ORDER BY created_at`, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.JavaLearningLogEntry
	for rows.Next() {
		var e domain.JavaLearningLogEntry
		if err := rows.Scan(&e.EntryID, &e.UserID, &e.TaskID, &e.LocalDate, &e.ErrorOrLearning,
			&e.FixOrNote, &e.Category, &e.RecurringCount14d, &e.IsRecurring, &e.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, e)
	}
	return results, nil
}

func (s *Store) CountJavaLearningLogByUserLabelSince(ctx context.Context, userID uuid.UUID, label string, since string) (int, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM java_learning_log WHERE user_id = $1 AND error_or_learning = $2 AND local_date >= $3`,
		userID, label, since)

	var count int
	err := row.Scan(&count)
	return count, err
}
