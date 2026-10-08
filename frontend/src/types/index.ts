export type UserRole = "student" | "author";

export type ActivityRunnerState =
  | "AVAILABLE"
  | "STARTING"
  | "ACTIVE"
  | "COMPLETING"
  | "ABANDONING";

export interface Activity {
  id: string;
  title: string;
  description: string;
  content: string;
  topic?: string;
  difficulty?: "beginner" | "intermediate" | "advanced" | string;
  estimated_duration_minutes?: number;
  correct_answer?: string;
  created_at?: string;
}

export interface Attempt {
  id: string;
  student_id: string;
  activity_id: string;
  status: "in_progress" | "completed" | "abandoned";
  started_at: string;
  completed_at?: string;
  duration?: number;
  response?: string;
  parent_attempt_id?: string;
}

export interface PersonalizedActivityData {
  id: string;
  source_activity_id: string;
  title: string;
  content: string;
  difficulty: string;
  intervention_type: string;
  personalization_reason: string;
  is_generated: boolean;
  parent_attempt_id?: string;
}

export interface AdaptationLoopData {
  id: string;
  student_id: string;
  original_attempt_id: string;
  adapted_attempt_id: string;
  assessment_id?: string;
  intervention_id?: string;
  /**
   * Type of intervention applied (e.g., INCREASE_CHALLENGE, SCAFFOLD).
   * Optional field populated when joined with intervention data or passed to UI.
   */
  intervention_type?: string;
  personalized_activity_id?: string;
  pre_state: string;
  post_state: string;
  pre_engagement_score: number;
  post_engagement_score: number;
  engagement_delta: number;
  effectiveness: "EFFECTIVE" | "NEUTRAL" | "INEFFECTIVE" | string;
  effectiveness_reasons: string[];
  /**
   * Alias / fallback for effectiveness reasoning displayed in historical cards.
   */
  reasoning?: string;
  created_at?: string;
}

export interface StudentHistoryItem {
  attempt_id: string;
  activity_id: string;
  activity_title: string;
  status: string;
  started_at: string;
  completed_at?: string;
  duration?: number;
  attempt_count: number;
  signals_count: number;
}

export interface DemoResult {
  activity_title: string;
  pre_state: string;
  pre_reasons: string[];
  intervention_type: string;
  intervention_reason: string;
  adapted_title: string;
  scenario: string;
  need: string;
  adapted_by: string;
  adapted_content: string;
  post_state: string;
  /**
   * Top-level effectiveness metric for direct display in summary cards.
   */
  effectiveness?: "EFFECTIVE" | "NEUTRAL" | "INEFFECTIVE" | string;
  loop: AdaptationLoopData;
}
