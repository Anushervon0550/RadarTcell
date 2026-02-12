DROP TRIGGER IF EXISTS trg_user_prefs_updated_at ON user_preferences;
DROP TRIGGER IF EXISTS trg_metrics_updated_at ON metrics_definitions;

DROP TABLE IF EXISTS user_preferences;
DROP TABLE IF EXISTS metrics_definitions;
