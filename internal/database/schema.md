# Table "public.access_tokens"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| created_at | timestamp with time zone | No | now() |  |
| creator_user_id | integer | No |  |  |
| deleted_at | timestamp with time zone | Yes |  |  |
| id | bigint | No | nextval('access_tokens_id_seq'::regclass) |  |
| internal | boolean | Yes | false |  |
| last_used_at | timestamp with time zone | Yes |  |  |
| note | text | No |  |  |
| scopes | text[] | No |  |  |
| subject_user_id | integer | No |  |  |
| value_sha256 | bytea | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| access_tokens_lookup | no | no | no | no | CREATE INDEX access_tokens_lookup ON access_tokens USING hash (value_sha256) WHERE deleted_at IS NULL |
| access_tokens_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX access_tokens_pkey ON access_tokens USING btree (id) |
| access_tokens_value_sha256_key | no | Yes | no | no | CREATE UNIQUE INDEX access_tokens_value_sha256_key ON access_tokens USING btree (value_sha256) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| access_tokens_creator_user_id_fkey | users | FOREIGN KEY (creator_user_id) REFERENCES users(id) |
| access_tokens_subject_user_id_fkey | users | FOREIGN KEY (subject_user_id) REFERENCES users(id) |
### References
| Name | Definition |
| --- | --- |
| batch_spec_workspace_execution_jobs | batch_spec_workspace_execution_jobs_access_token_id_fkey |
# Table "public.batch_changes"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| batch_spec_id | bigint | No |  |  |
| closed_at | timestamp with time zone | Yes |  |  |
| created_at | timestamp with time zone | No | now() |  |
| creator_id | integer | Yes |  |  |
| description | text | Yes |  |  |
| id | bigint | No | nextval('batch_changes_id_seq'::regclass) |  |
| last_applied_at | timestamp with time zone | Yes |  |  |
| last_applier_id | bigint | Yes |  |  |
| name | text | No |  |  |
| namespace_org_id | integer | Yes |  |  |
| namespace_user_id | integer | Yes |  |  |
| updated_at | timestamp with time zone | No | now() |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| batch_changes_namespace_org_id | no | no | no | no | CREATE INDEX batch_changes_namespace_org_id ON batch_changes USING btree (namespace_org_id) |
| batch_changes_namespace_user_id | no | no | no | no | CREATE INDEX batch_changes_namespace_user_id ON batch_changes USING btree (namespace_user_id) |
| batch_changes_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX batch_changes_pkey ON batch_changes USING btree (id) |
| batch_changes_unique_org_id | no | Yes | no | no | CREATE UNIQUE INDEX batch_changes_unique_org_id ON batch_changes USING btree (name, namespace_org_id) WHERE namespace_org_id IS NOT NULL |
| batch_changes_unique_user_id | no | Yes | no | no | CREATE UNIQUE INDEX batch_changes_unique_user_id ON batch_changes USING btree (name, namespace_user_id) WHERE namespace_user_id IS NOT NULL |
### Check constraints
| Name | Definition |
| --- | --- |
| batch_changes_has_1_namespace | CHECK ((namespace_user_id IS NULL) <> (namespace_org_id IS NULL)) |
| batch_changes_name_not_blank | CHECK (name <> ''::text) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| batch_changes_batch_spec_id_fkey | batch_specs | FOREIGN KEY (batch_spec_id) REFERENCES batch_specs(id) DEFERRABLE |
| batch_changes_initial_applier_id_fkey | users | FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE SET NULL DEFERRABLE |
| batch_changes_last_applier_id_fkey | users | FOREIGN KEY (last_applier_id) REFERENCES users(id) ON DELETE SET NULL DEFERRABLE |
| batch_changes_namespace_org_id_fkey | orgs | FOREIGN KEY (namespace_org_id) REFERENCES orgs(id) ON DELETE CASCADE DEFERRABLE |
| batch_changes_namespace_user_id_fkey | users | FOREIGN KEY (namespace_user_id) REFERENCES users(id) ON DELETE CASCADE DEFERRABLE |
### Triggers
| Name | Definition |
| --- | --- |
| trig_delete_batch_change_reference_on_changesets | CREATE TRIGGER trig_delete_batch_change_reference_on_changesets AFTER DELETE ON batch_changes FOR EACH ROW EXECUTE FUNCTION delete_batch_change_reference_on_changesets() |
### References
| Name | Definition |
| --- | --- |
| changeset_jobs | changeset_jobs_batch_change_id_fkey |
| changesets | changesets_owned_by_batch_spec_id_fkey |
# Table "public.batch_changes_site_credentials"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| created_at | timestamp with time zone | No | now() |  |
| credential | bytea | No |  |  |
| encryption_key_id | text | No | ''::text |  |
| external_service_id | text | No |  |  |
| external_service_type | text | No |  |  |
| id | bigint | No | nextval('batch_changes_site_credentials_id_seq'::regclass) |  |
| updated_at | timestamp with time zone | No | now() |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| batch_changes_site_credentials_credential_idx | no | no | no | no | CREATE INDEX batch_changes_site_credentials_credential_idx ON batch_changes_site_credentials USING btree ((encryption_key_id = ANY (ARRAY[''::text, 'previously-migrated'::text]))) |
| batch_changes_site_credentials_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX batch_changes_site_credentials_pkey ON batch_changes_site_credentials USING btree (id) |
| batch_changes_site_credentials_unique | no | Yes | no | no | CREATE UNIQUE INDEX batch_changes_site_credentials_unique ON batch_changes_site_credentials USING btree (external_service_type, external_service_id) |
# Table "public.batch_spec_execution_cache_entries"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| created_at | timestamp with time zone | No | now() |  |
| id | bigint | No | nextval('batch_spec_execution_cache_entries_id_seq'::regclass) |  |
| key | text | No |  |  |
| last_used_at | timestamp with time zone | Yes |  |  |
| user_id | integer | No |  |  |
| value | text | No |  |  |
| version | integer | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| batch_spec_execution_cache_entries_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX batch_spec_execution_cache_entries_pkey ON batch_spec_execution_cache_entries USING btree (id) |
| batch_spec_execution_cache_entries_user_id_key_unique | no | Yes | no | no | CREATE UNIQUE INDEX batch_spec_execution_cache_entries_user_id_key_unique ON batch_spec_execution_cache_entries USING btree (user_id, key) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| batch_spec_execution_cache_entries_user_id_fkey | users | FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE DEFERRABLE |
# Table "public.batch_spec_resolution_jobs"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| batch_spec_id | integer | Yes |  |  |
| created_at | timestamp with time zone | No | now() |  |
| execution_logs | json[] | Yes |  |  |
| failure_message | text | Yes |  |  |
| finished_at | timestamp with time zone | Yes |  |  |
| id | bigint | No | nextval('batch_spec_resolution_jobs_id_seq'::regclass) |  |
| last_heartbeat_at | timestamp with time zone | Yes |  |  |
| num_failures | integer | No | 0 |  |
| num_resets | integer | No | 0 |  |
| process_after | timestamp with time zone | Yes |  |  |
| queued_at | timestamp with time zone | Yes | now() |  |
| started_at | timestamp with time zone | Yes |  |  |
| state | text | Yes | 'queued'::text |  |
| updated_at | timestamp with time zone | No | now() |  |
| worker_hostname | text | No | ''::text |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| batch_spec_resolution_jobs_batch_spec_id_unique | no | Yes | no | no | CREATE UNIQUE INDEX batch_spec_resolution_jobs_batch_spec_id_unique ON batch_spec_resolution_jobs USING btree (batch_spec_id) |
| batch_spec_resolution_jobs_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX batch_spec_resolution_jobs_pkey ON batch_spec_resolution_jobs USING btree (id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| batch_spec_resolution_jobs_batch_spec_id_fkey | batch_specs | FOREIGN KEY (batch_spec_id) REFERENCES batch_specs(id) ON DELETE CASCADE DEFERRABLE |
# Table "public.batch_spec_workspace_execution_jobs"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| access_token_id | bigint | Yes |  |  |
| batch_spec_workspace_id | integer | Yes |  |  |
| cancel | boolean | No | false |  |
| created_at | timestamp with time zone | No | now() |  |
| execution_logs | json[] | Yes |  |  |
| failure_message | text | Yes |  |  |
| finished_at | timestamp with time zone | Yes |  |  |
| id | bigint | No | nextval('batch_spec_workspace_execution_jobs_id_seq'::regclass) |  |
| last_heartbeat_at | timestamp with time zone | Yes |  |  |
| num_failures | integer | No | 0 |  |
| num_resets | integer | No | 0 |  |
| process_after | timestamp with time zone | Yes |  |  |
| queued_at | timestamp with time zone | Yes | now() |  |
| started_at | timestamp with time zone | Yes |  |  |
| state | text | Yes | 'queued'::text |  |
| updated_at | timestamp with time zone | No | now() |  |
| worker_hostname | text | No | ''::text |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| batch_spec_workspace_execution_jobs_cancel | no | no | no | no | CREATE INDEX batch_spec_workspace_execution_jobs_cancel ON batch_spec_workspace_execution_jobs USING btree (cancel) |
| batch_spec_workspace_execution_jobs_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX batch_spec_workspace_execution_jobs_pkey ON batch_spec_workspace_execution_jobs USING btree (id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| batch_spec_workspace_execution_job_batch_spec_workspace_id_fkey | batch_spec_workspaces | FOREIGN KEY (batch_spec_workspace_id) REFERENCES batch_spec_workspaces(id) ON DELETE CASCADE DEFERRABLE |
| batch_spec_workspace_execution_jobs_access_token_id_fkey | access_tokens | FOREIGN KEY (access_token_id) REFERENCES access_tokens(id) ON DELETE SET NULL DEFERRABLE |
# Table "public.batch_spec_workspaces"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| batch_spec_id | integer | Yes |  |  |
| branch | text | No |  |  |
| cached_result_found | boolean | No | false |  |
| changeset_spec_ids | jsonb | Yes | '{}'::jsonb |  |
| commit | text | No |  |  |
| created_at | timestamp with time zone | No | now() |  |
| file_matches | text[] | No |  |  |
| id | bigint | No | nextval('batch_spec_workspaces_id_seq'::regclass) |  |
| ignored | boolean | No | false |  |
| only_fetch_workspace | boolean | No | false |  |
| path | text | No |  |  |
| repo_id | integer | Yes |  |  |
| skipped | boolean | No | false |  |
| step_cache_results | jsonb | No | '{}'::jsonb |  |
| unsupported | boolean | No | false |  |
| updated_at | timestamp with time zone | No | now() |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| batch_spec_workspaces_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX batch_spec_workspaces_pkey ON batch_spec_workspaces USING btree (id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| batch_spec_workspaces_batch_spec_id_fkey | batch_specs | FOREIGN KEY (batch_spec_id) REFERENCES batch_specs(id) ON DELETE CASCADE DEFERRABLE |
| batch_spec_workspaces_repo_id_fkey | repo | FOREIGN KEY (repo_id) REFERENCES repo(id) DEFERRABLE |
### References
| Name | Definition |
| --- | --- |
| batch_spec_workspace_execution_jobs | batch_spec_workspace_execution_job_batch_spec_workspace_id_fkey |
# Table "public.batch_specs"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| allow_ignored | boolean | No | false |  |
| allow_unsupported | boolean | No | false |  |
| created_at | timestamp with time zone | No | now() |  |
| created_from_raw | boolean | No | false |  |
| id | bigint | No | nextval('batch_specs_id_seq'::regclass) |  |
| namespace_org_id | integer | Yes |  |  |
| namespace_user_id | integer | Yes |  |  |
| no_cache | boolean | No | false |  |
| rand_id | text | No |  |  |
| raw_spec | text | No |  |  |
| spec | jsonb | No | '{}'::jsonb |  |
| updated_at | timestamp with time zone | No | now() |  |
| user_id | integer | Yes |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| batch_specs_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX batch_specs_pkey ON batch_specs USING btree (id) |
| batch_specs_rand_id | no | no | no | no | CREATE INDEX batch_specs_rand_id ON batch_specs USING btree (rand_id) |
### Check constraints
| Name | Definition |
| --- | --- |
| batch_specs_has_1_namespace | CHECK ((namespace_user_id IS NULL) <> (namespace_org_id IS NULL)) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| batch_specs_user_id_fkey | users | FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL DEFERRABLE |
### References
| Name | Definition |
| --- | --- |
| batch_changes | batch_changes_batch_spec_id_fkey |
| batch_spec_resolution_jobs | batch_spec_resolution_jobs_batch_spec_id_fkey |
| batch_spec_workspaces | batch_spec_workspaces_batch_spec_id_fkey |
| changeset_specs | changeset_specs_batch_spec_id_fkey |
# Table "public.changeset_events"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| changeset_id | bigint | No |  |  |
| created_at | timestamp with time zone | No | now() |  |
| id | bigint | No | nextval('changeset_events_id_seq'::regclass) |  |
| key | text | No |  |  |
| kind | text | No |  |  |
| metadata | jsonb | No | '{}'::jsonb |  |
| updated_at | timestamp with time zone | No | now() |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| changeset_events_changeset_id_kind_key_unique | no | Yes | no | no | CREATE UNIQUE INDEX changeset_events_changeset_id_kind_key_unique ON changeset_events USING btree (changeset_id, kind, key) |
| changeset_events_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX changeset_events_pkey ON changeset_events USING btree (id) |
### Check constraints
| Name | Definition |
| --- | --- |
| changeset_events_key_check | CHECK (key <> ''::text) |
| changeset_events_kind_check | CHECK (kind <> ''::text) |
| changeset_events_metadata_check | CHECK (jsonb_typeof(metadata) = 'object'::text) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| changeset_events_changeset_id_fkey | changesets | FOREIGN KEY (changeset_id) REFERENCES changesets(id) ON DELETE CASCADE DEFERRABLE |
# Table "public.changeset_jobs"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| batch_change_id | integer | No |  |  |
| bulk_group | text | No |  |  |
| changeset_id | integer | No |  |  |
| created_at | timestamp with time zone | No | now() |  |
| execution_logs | json[] | Yes |  |  |
| failure_message | text | Yes |  |  |
| finished_at | timestamp with time zone | Yes |  |  |
| id | bigint | No | nextval('changeset_jobs_id_seq'::regclass) |  |
| job_type | text | No |  |  |
| last_heartbeat_at | timestamp with time zone | Yes |  |  |
| num_failures | integer | No | 0 |  |
| num_resets | integer | No | 0 |  |
| payload | jsonb | Yes | '{}'::jsonb |  |
| process_after | timestamp with time zone | Yes |  |  |
| queued_at | timestamp with time zone | Yes | now() |  |
| started_at | timestamp with time zone | Yes |  |  |
| state | text | Yes | 'queued'::text |  |
| updated_at | timestamp with time zone | No | now() |  |
| user_id | integer | No |  |  |
| worker_hostname | text | No | ''::text |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| changeset_jobs_bulk_group_idx | no | no | no | no | CREATE INDEX changeset_jobs_bulk_group_idx ON changeset_jobs USING btree (bulk_group) |
| changeset_jobs_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX changeset_jobs_pkey ON changeset_jobs USING btree (id) |
| changeset_jobs_state_idx | no | no | no | no | CREATE INDEX changeset_jobs_state_idx ON changeset_jobs USING btree (state) |
### Check constraints
| Name | Definition |
| --- | --- |
| changeset_jobs_payload_check | CHECK (jsonb_typeof(payload) = 'object'::text) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| changeset_jobs_batch_change_id_fkey | batch_changes | FOREIGN KEY (batch_change_id) REFERENCES batch_changes(id) ON DELETE CASCADE DEFERRABLE |
| changeset_jobs_changeset_id_fkey | changesets | FOREIGN KEY (changeset_id) REFERENCES changesets(id) ON DELETE CASCADE DEFERRABLE |
| changeset_jobs_user_id_fkey | users | FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE DEFERRABLE |
# Table "public.changeset_specs"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| batch_spec_id | bigint | Yes |  |  |
| created_at | timestamp with time zone | No | now() |  |
| diff_stat_added | integer | Yes |  |  |
| diff_stat_changed | integer | Yes |  |  |
| diff_stat_deleted | integer | Yes |  |  |
| external_id | text | Yes |  |  |
| fork_namespace | citext | Yes |  |  |
| head_ref | text | Yes |  |  |
| id | bigint | No | nextval('changeset_specs_id_seq'::regclass) |  |
| rand_id | text | No |  |  |
| repo_id | integer | No |  |  |
| spec | jsonb | No | '{}'::jsonb |  |
| title | text | Yes |  |  |
| updated_at | timestamp with time zone | No | now() |  |
| user_id | integer | Yes |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| changeset_specs_external_id | no | no | no | no | CREATE INDEX changeset_specs_external_id ON changeset_specs USING btree (external_id) |
| changeset_specs_head_ref | no | no | no | no | CREATE INDEX changeset_specs_head_ref ON changeset_specs USING btree (head_ref) |
| changeset_specs_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX changeset_specs_pkey ON changeset_specs USING btree (id) |
| changeset_specs_rand_id | no | no | no | no | CREATE INDEX changeset_specs_rand_id ON changeset_specs USING btree (rand_id) |
| changeset_specs_title | no | no | no | no | CREATE INDEX changeset_specs_title ON changeset_specs USING btree (title) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| changeset_specs_batch_spec_id_fkey | batch_specs | FOREIGN KEY (batch_spec_id) REFERENCES batch_specs(id) ON DELETE CASCADE DEFERRABLE |
| changeset_specs_repo_id_fkey | repo | FOREIGN KEY (repo_id) REFERENCES repo(id) DEFERRABLE |
| changeset_specs_user_id_fkey | users | FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL DEFERRABLE |
### References
| Name | Definition |
| --- | --- |
| changesets | changesets_changeset_spec_id_fkey |
| changesets | changesets_previous_spec_id_fkey |
# Table "public.changesets"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| batch_change_ids | jsonb | No | '{}'::jsonb |  |
| closing | boolean | No | false |  |
| created_at | timestamp with time zone | No | now() |  |
| current_spec_id | bigint | Yes |  |  |
| diff_stat_added | integer | Yes |  |  |
| diff_stat_changed | integer | Yes |  |  |
| diff_stat_deleted | integer | Yes |  |  |
| execution_logs | json[] | Yes |  |  |
| external_branch | text | Yes |  |  |
| external_check_state | text | Yes |  |  |
| external_deleted_at | timestamp with time zone | Yes |  |  |
| external_fork_namespace | citext | Yes |  |  |
| external_id | text | Yes |  |  |
| external_review_state | text | Yes |  |  |
| external_service_type | text | No |  |  |
| external_state | text | Yes |  |  |
| external_title | text | Yes |  | Normalized property generated on save using Changeset.Title() |
| external_updated_at | timestamp with time zone | Yes |  |  |
| failure_message | text | Yes |  |  |
| finished_at | timestamp with time zone | Yes |  |  |
| id | bigint | No | nextval('changesets_id_seq'::regclass) |  |
| last_heartbeat_at | timestamp with time zone | Yes |  |  |
| log_contents | text | Yes |  |  |
| metadata | jsonb | Yes | '{}'::jsonb |  |
| num_failures | integer | No | 0 |  |
| num_resets | integer | No | 0 |  |
| owned_by_batch_change_id | bigint | Yes |  |  |
| previous_spec_id | bigint | Yes |  |  |
| process_after | timestamp with time zone | Yes |  |  |
| publication_state | text | Yes | 'UNPUBLISHED'::text |  |
| queued_at | timestamp with time zone | Yes | now() |  |
| reconciler_state | text | Yes | 'queued'::text |  |
| repo_id | integer | No |  |  |
| started_at | timestamp with time zone | Yes |  |  |
| sync_state | jsonb | No | '{}'::jsonb |  |
| syncer_error | text | Yes |  |  |
| ui_publication_state | batch_changes_changeset_ui_publication_state | Yes |  |  |
| updated_at | timestamp with time zone | No | now() |  |
| worker_hostname | text | No | ''::text |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| changesets_batch_change_ids | no | no | no | no | CREATE INDEX changesets_batch_change_ids ON changesets USING gin (batch_change_ids) |
| changesets_external_state_idx | no | no | no | no | CREATE INDEX changesets_external_state_idx ON changesets USING btree (external_state) |
| changesets_external_title_idx | no | no | no | no | CREATE INDEX changesets_external_title_idx ON changesets USING btree (external_title) |
| changesets_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX changesets_pkey ON changesets USING btree (id) |
| changesets_publication_state_idx | no | no | no | no | CREATE INDEX changesets_publication_state_idx ON changesets USING btree (publication_state) |
| changesets_reconciler_state_idx | no | no | no | no | CREATE INDEX changesets_reconciler_state_idx ON changesets USING btree (reconciler_state) |
| changesets_repo_external_id_unique | no | Yes | no | no | CREATE UNIQUE INDEX changesets_repo_external_id_unique ON changesets USING btree (repo_id, external_id) |
### Check constraints
| Name | Definition |
| --- | --- |
| changesets_batch_change_ids_check | CHECK (jsonb_typeof(batch_change_ids) = 'object'::text) |
| changesets_external_id_check | CHECK (external_id <> ''::text) |
| changesets_external_service_type_not_blank | CHECK (external_service_type <> ''::text) |
| changesets_metadata_check | CHECK (jsonb_typeof(metadata) = 'object'::text) |
| external_branch_ref_prefix | CHECK (external_branch ~~ 'refs/heads/%'::text) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| changesets_changeset_spec_id_fkey | changeset_specs | FOREIGN KEY (current_spec_id) REFERENCES changeset_specs(id) DEFERRABLE |
| changesets_owned_by_batch_spec_id_fkey | batch_changes | FOREIGN KEY (owned_by_batch_change_id) REFERENCES batch_changes(id) ON DELETE SET NULL DEFERRABLE |
| changesets_previous_spec_id_fkey | changeset_specs | FOREIGN KEY (previous_spec_id) REFERENCES changeset_specs(id) DEFERRABLE |
| changesets_repo_id_fkey | repo | FOREIGN KEY (repo_id) REFERENCES repo(id) ON DELETE CASCADE DEFERRABLE |
### References
| Name | Definition |
| --- | --- |
| changeset_events | changeset_events_changeset_id_fkey |
| changeset_jobs | changeset_jobs_changeset_id_fkey |
# Table "public.cm_action_jobs"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| email | bigint | Yes |  | The ID of the cm_emails action to execute if this is an email job. Mutually exclusive with webhook and slack_webhook |
| execution_logs | json[] | Yes |  |  |
| failure_message | text | Yes |  |  |
| finished_at | timestamp with time zone | Yes |  |  |
| id | integer | No | nextval('cm_action_jobs_id_seq'::regclass) |  |
| last_heartbeat_at | timestamp with time zone | Yes |  |  |
| log_contents | text | Yes |  |  |
| num_failures | integer | No | 0 |  |
| num_resets | integer | No | 0 |  |
| process_after | timestamp with time zone | Yes |  |  |
| queued_at | timestamp with time zone | Yes | now() |  |
| slack_webhook | bigint | Yes |  | The ID of the cm_slack_webhook action to execute if this is a slack webhook job. Mutually exclusive with email and webhook |
| started_at | timestamp with time zone | Yes |  |  |
| state | text | Yes | 'queued'::text |  |
| trigger_event | integer | Yes |  |  |
| webhook | bigint | Yes |  | The ID of the cm_webhooks action to execute if this is a webhook job. Mutually exclusive with email and slack_webhook |
| worker_hostname | text | No | ''::text |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| cm_action_jobs_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX cm_action_jobs_pkey ON cm_action_jobs USING btree (id) |
### Check constraints
| Name | Definition |
| --- | --- |
| cm_action_jobs_only_one_action_type | CHECK ((
CASE
    WHEN email IS NULL THEN 0
    ELSE 1
END +
CASE
    WHEN webhook IS NULL THEN 0
    ELSE 1
END +
CASE
    WHEN slack_webhook IS NULL THEN 0
    ELSE 1
END) = 1) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| cm_action_jobs_email_fk | cm_emails | FOREIGN KEY (email) REFERENCES cm_emails(id) ON DELETE CASCADE |
| cm_action_jobs_slack_webhook_fkey | cm_slack_webhooks | FOREIGN KEY (slack_webhook) REFERENCES cm_slack_webhooks(id) ON DELETE CASCADE |
| cm_action_jobs_trigger_event_fk | cm_trigger_jobs | FOREIGN KEY (trigger_event) REFERENCES cm_trigger_jobs(id) ON DELETE CASCADE |
| cm_action_jobs_webhook_fkey | cm_webhooks | FOREIGN KEY (webhook) REFERENCES cm_webhooks(id) ON DELETE CASCADE |
# Table "public.cm_emails"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| changed_at | timestamp with time zone | No | now() |  |
| changed_by | integer | No |  |  |
| created_at | timestamp with time zone | No | now() |  |
| created_by | integer | No |  |  |
| enabled | boolean | No |  |  |
| header | text | No |  |  |
| id | bigint | No | nextval('cm_emails_id_seq'::regclass) |  |
| include_results | boolean | No | false |  |
| monitor | bigint | No |  |  |
| priority | cm_email_priority | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| cm_emails_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX cm_emails_pkey ON cm_emails USING btree (id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| cm_emails_changed_by_fk | users | FOREIGN KEY (changed_by) REFERENCES users(id) ON DELETE CASCADE |
| cm_emails_created_by_fk | users | FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE |
| cm_emails_monitor | cm_monitors | FOREIGN KEY (monitor) REFERENCES cm_monitors(id) ON DELETE CASCADE |
### References
| Name | Definition |
| --- | --- |
| cm_action_jobs | cm_action_jobs_email_fk |
| cm_recipients | cm_recipients_emails |
# Table "public.cm_last_searched"


The last searched commit hashes for the given code monitor and unique set of search arguments

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| args_hash | bigint | No |  | A unique hash of the gitserver search arguments to identify this search job |
| commit_oids | text[] | No |  | The set of commit OIDs that was previously successfully searched and should be excluded on the next run |
| monitor_id | bigint | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| cm_last_searched_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX cm_last_searched_pkey ON cm_last_searched USING btree (monitor_id, args_hash) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| cm_last_searched_monitor_id_fkey | cm_monitors | FOREIGN KEY (monitor_id) REFERENCES cm_monitors(id) ON DELETE CASCADE |
# Table "public.cm_monitors"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| changed_at | timestamp with time zone | No | now() |  |
| changed_by | integer | No |  |  |
| created_at | timestamp with time zone | No | now() |  |
| created_by | integer | No |  |  |
| description | text | No |  |  |
| enabled | boolean | No | true |  |
| id | bigint | No | nextval('cm_monitors_id_seq'::regclass) |  |
| namespace_org_id | integer | Yes |  | DEPRECATED: code monitors cannot be owned by an org |
| namespace_user_id | integer | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| cm_monitors_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX cm_monitors_pkey ON cm_monitors USING btree (id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| cm_monitors_changed_by_fk | users | FOREIGN KEY (changed_by) REFERENCES users(id) ON DELETE CASCADE |
| cm_monitors_created_by_fk | users | FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE |
| cm_monitors_org_id_fk | orgs | FOREIGN KEY (namespace_org_id) REFERENCES orgs(id) ON DELETE CASCADE |
| cm_monitors_user_id_fk | users | FOREIGN KEY (namespace_user_id) REFERENCES users(id) ON DELETE CASCADE |
### References
| Name | Definition |
| --- | --- |
| cm_emails | cm_emails_monitor |
| cm_last_searched | cm_last_searched_monitor_id_fkey |
| cm_slack_webhooks | cm_slack_webhooks_monitor_fkey |
| cm_queries | cm_triggers_monitor |
| cm_webhooks | cm_webhooks_monitor_fkey |
# Table "public.cm_queries"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| changed_at | timestamp with time zone | No | now() |  |
| changed_by | integer | No |  |  |
| created_at | timestamp with time zone | No | now() |  |
| created_by | integer | No |  |  |
| id | bigint | No | nextval('cm_queries_id_seq'::regclass) |  |
| latest_result | timestamp with time zone | Yes |  |  |
| monitor | bigint | No |  |  |
| next_run | timestamp with time zone | Yes | now() |  |
| query | text | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| cm_queries_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX cm_queries_pkey ON cm_queries USING btree (id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| cm_triggers_changed_by_fk | users | FOREIGN KEY (changed_by) REFERENCES users(id) ON DELETE CASCADE |
| cm_triggers_created_by_fk | users | FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE |
| cm_triggers_monitor | cm_monitors | FOREIGN KEY (monitor) REFERENCES cm_monitors(id) ON DELETE CASCADE |
### References
| Name | Definition |
| --- | --- |
| cm_trigger_jobs | cm_trigger_jobs_query_fk |
# Table "public.cm_recipients"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| email | bigint | No |  |  |
| id | bigint | No | nextval('cm_recipients_id_seq'::regclass) |  |
| namespace_org_id | integer | Yes |  |  |
| namespace_user_id | integer | Yes |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| cm_recipients_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX cm_recipients_pkey ON cm_recipients USING btree (id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| cm_recipients_emails | cm_emails | FOREIGN KEY (email) REFERENCES cm_emails(id) ON DELETE CASCADE |
| cm_recipients_org_id_fk | orgs | FOREIGN KEY (namespace_org_id) REFERENCES orgs(id) ON DELETE CASCADE |
| cm_recipients_user_id_fk | users | FOREIGN KEY (namespace_user_id) REFERENCES users(id) ON DELETE CASCADE |
# Table "public.cm_slack_webhooks"


Slack webhook actions configured on code monitors

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| changed_at | timestamp with time zone | No | now() |  |
| changed_by | integer | No |  |  |
| created_at | timestamp with time zone | No | now() |  |
| created_by | integer | No |  |  |
| enabled | boolean | No |  |  |
| id | bigint | No | nextval('cm_slack_webhooks_id_seq'::regclass) |  |
| include_results | boolean | No | false |  |
| monitor | bigint | No |  | The code monitor that the action is defined on |
| url | text | No |  | The Slack webhook URL we send the code monitor event to |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| cm_slack_webhooks_monitor | no | no | no | no | CREATE INDEX cm_slack_webhooks_monitor ON cm_slack_webhooks USING btree (monitor) |
| cm_slack_webhooks_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX cm_slack_webhooks_pkey ON cm_slack_webhooks USING btree (id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| cm_slack_webhooks_changed_by_fkey | users | FOREIGN KEY (changed_by) REFERENCES users(id) ON DELETE CASCADE |
| cm_slack_webhooks_created_by_fkey | users | FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE |
| cm_slack_webhooks_monitor_fkey | cm_monitors | FOREIGN KEY (monitor) REFERENCES cm_monitors(id) ON DELETE CASCADE |
### References
| Name | Definition |
| --- | --- |
| cm_action_jobs | cm_action_jobs_slack_webhook_fkey |
# Table "public.cm_trigger_jobs"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| execution_logs | json[] | Yes |  |  |
| failure_message | text | Yes |  |  |
| finished_at | timestamp with time zone | Yes |  |  |
| id | integer | No | nextval('cm_trigger_jobs_id_seq'::regclass) |  |
| last_heartbeat_at | timestamp with time zone | Yes |  |  |
| log_contents | text | Yes |  |  |
| num_failures | integer | No | 0 |  |
| num_resets | integer | No | 0 |  |
| num_results | integer | Yes |  | DEPRECATED: replaced by len(search_results). Can be removed after version 3.37 release cut |
| process_after | timestamp with time zone | Yes |  |  |
| query | bigint | No |  |  |
| query_string | text | Yes |  |  |
| queued_at | timestamp with time zone | Yes | now() |  |
| results | boolean | Yes |  | DEPRECATED: replaced by len(search_results) > 0. Can be removed after version 3.37 release cut |
| search_results | jsonb | Yes |  |  |
| started_at | timestamp with time zone | Yes |  |  |
| state | text | Yes | 'queued'::text |  |
| worker_hostname | text | No | ''::text |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| cm_trigger_jobs_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX cm_trigger_jobs_pkey ON cm_trigger_jobs USING btree (id) |
### Check constraints
| Name | Definition |
| --- | --- |
| search_results_is_array | CHECK (jsonb_typeof(search_results) = 'array'::text) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| cm_trigger_jobs_query_fk | cm_queries | FOREIGN KEY (query) REFERENCES cm_queries(id) ON DELETE CASCADE |
### References
| Name | Definition |
| --- | --- |
| cm_action_jobs | cm_action_jobs_trigger_event_fk |
# Table "public.cm_webhooks"


Webhook actions configured on code monitors

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| changed_at | timestamp with time zone | No | now() |  |
| changed_by | integer | No |  |  |
| created_at | timestamp with time zone | No | now() |  |
| created_by | integer | No |  |  |
| enabled | boolean | No |  | Whether this Slack webhook action is enabled. When not enabled, the action will not be run when its code monitor generates events |
| id | bigint | No | nextval('cm_webhooks_id_seq'::regclass) |  |
| include_results | boolean | No | false |  |
| monitor | bigint | No |  | The code monitor that the action is defined on |
| url | text | No |  | The webhook URL we send the code monitor event to |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| cm_webhooks_monitor | no | no | no | no | CREATE INDEX cm_webhooks_monitor ON cm_webhooks USING btree (monitor) |
| cm_webhooks_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX cm_webhooks_pkey ON cm_webhooks USING btree (id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| cm_webhooks_changed_by_fkey | users | FOREIGN KEY (changed_by) REFERENCES users(id) ON DELETE CASCADE |
| cm_webhooks_created_by_fkey | users | FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE |
| cm_webhooks_monitor_fkey | cm_monitors | FOREIGN KEY (monitor) REFERENCES cm_monitors(id) ON DELETE CASCADE |
### References
| Name | Definition |
| --- | --- |
| cm_action_jobs | cm_action_jobs_webhook_fkey |
# Table "public.critical_and_site_config"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| contents | text | No |  |  |
| created_at | timestamp with time zone | No | now() |  |
| id | integer | No | nextval('critical_and_site_config_id_seq'::regclass) |  |
| type | critical_or_site | No |  |  |
| updated_at | timestamp with time zone | No | now() |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| critical_and_site_config_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX critical_and_site_config_pkey ON critical_and_site_config USING btree (id) |
| critical_and_site_config_unique | no | Yes | no | no | CREATE UNIQUE INDEX critical_and_site_config_unique ON critical_and_site_config USING btree (id, type) |
# Table "public.discussion_comments"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| author_user_id | integer | No |  |  |
| contents | text | No |  |  |
| created_at | timestamp with time zone | No | now() |  |
| deleted_at | timestamp with time zone | Yes |  |  |
| id | bigint | No | nextval('discussion_comments_id_seq'::regclass) |  |
| reports | text[] | No | '{}'::text[] |  |
| thread_id | bigint | No |  |  |
| updated_at | timestamp with time zone | No | now() |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| discussion_comments_author_user_id_idx | no | no | no | no | CREATE INDEX discussion_comments_author_user_id_idx ON discussion_comments USING btree (author_user_id) |
| discussion_comments_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX discussion_comments_pkey ON discussion_comments USING btree (id) |
| discussion_comments_reports_array_length_idx | no | no | no | no | CREATE INDEX discussion_comments_reports_array_length_idx ON discussion_comments USING btree (array_length(reports, 1)) |
| discussion_comments_thread_id_idx | no | no | no | no | CREATE INDEX discussion_comments_thread_id_idx ON discussion_comments USING btree (thread_id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| discussion_comments_author_user_id_fkey | users | FOREIGN KEY (author_user_id) REFERENCES users(id) ON DELETE RESTRICT |
| discussion_comments_thread_id_fkey | discussion_threads | FOREIGN KEY (thread_id) REFERENCES discussion_threads(id) ON DELETE CASCADE |
# Table "public.discussion_mail_reply_tokens"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| deleted_at | timestamp with time zone | Yes |  |  |
| thread_id | bigint | No |  |  |
| token | text | No |  |  |
| user_id | integer | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| discussion_mail_reply_tokens_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX discussion_mail_reply_tokens_pkey ON discussion_mail_reply_tokens USING btree (token) |
| discussion_mail_reply_tokens_user_id_thread_id_idx | no | no | no | no | CREATE INDEX discussion_mail_reply_tokens_user_id_thread_id_idx ON discussion_mail_reply_tokens USING btree (user_id, thread_id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| discussion_mail_reply_tokens_thread_id_fkey | discussion_threads | FOREIGN KEY (thread_id) REFERENCES discussion_threads(id) ON DELETE CASCADE |
| discussion_mail_reply_tokens_user_id_fkey | users | FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT |
# Table "public.discussion_threads"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| archived_at | timestamp with time zone | Yes |  |  |
| author_user_id | integer | No |  |  |
| created_at | timestamp with time zone | No | now() |  |
| deleted_at | timestamp with time zone | Yes |  |  |
| id | bigint | No | nextval('discussion_threads_id_seq'::regclass) |  |
| target_repo_id | bigint | Yes |  |  |
| title | text | Yes |  |  |
| updated_at | timestamp with time zone | No | now() |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| discussion_threads_author_user_id_idx | no | no | no | no | CREATE INDEX discussion_threads_author_user_id_idx ON discussion_threads USING btree (author_user_id) |
| discussion_threads_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX discussion_threads_pkey ON discussion_threads USING btree (id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| discussion_threads_author_user_id_fkey | users | FOREIGN KEY (author_user_id) REFERENCES users(id) ON DELETE RESTRICT |
| discussion_threads_target_repo_id_fk | discussion_threads_target_repo | FOREIGN KEY (target_repo_id) REFERENCES discussion_threads_target_repo(id) ON DELETE CASCADE |
### References
| Name | Definition |
| --- | --- |
| discussion_comments | discussion_comments_thread_id_fkey |
| discussion_mail_reply_tokens | discussion_mail_reply_tokens_thread_id_fkey |
| discussion_threads_target_repo | discussion_threads_target_repo_thread_id_fkey |
# Table "public.discussion_threads_target_repo"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| branch | text | Yes |  |  |
| end_character | integer | Yes |  |  |
| end_line | integer | Yes |  |  |
| id | bigint | No | nextval('discussion_threads_target_repo_id_seq'::regclass) |  |
| lines | text | Yes |  |  |
| lines_after | text | Yes |  |  |
| lines_before | text | Yes |  |  |
| path | text | Yes |  |  |
| repo_id | integer | No |  |  |
| revision | text | Yes |  |  |
| start_character | integer | Yes |  |  |
| start_line | integer | Yes |  |  |
| thread_id | bigint | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| discussion_threads_target_repo_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX discussion_threads_target_repo_pkey ON discussion_threads_target_repo USING btree (id) |
| discussion_threads_target_repo_repo_id_path_idx | no | no | no | no | CREATE INDEX discussion_threads_target_repo_repo_id_path_idx ON discussion_threads_target_repo USING btree (repo_id, path) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| discussion_threads_target_repo_repo_id_fkey | repo | FOREIGN KEY (repo_id) REFERENCES repo(id) ON DELETE CASCADE |
| discussion_threads_target_repo_thread_id_fkey | discussion_threads | FOREIGN KEY (thread_id) REFERENCES discussion_threads(id) ON DELETE CASCADE |
### References
| Name | Definition |
| --- | --- |
| discussion_threads | discussion_threads_target_repo_id_fk |
# Table "public.event_logs"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| anonymous_user_id | text | No |  |  |
| argument | jsonb | No |  |  |
| cohort_id | date | Yes |  |  |
| feature_flags | jsonb | Yes |  |  |
| id | bigint | No | nextval('event_logs_id_seq'::regclass) |  |
| name | text | No |  |  |
| public_argument | jsonb | No | '{}'::jsonb |  |
| source | text | No |  |  |
| timestamp | timestamp with time zone | No |  |  |
| url | text | No |  |  |
| user_id | integer | No |  |  |
| version | text | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| event_logs_anonymous_user_id | no | no | no | no | CREATE INDEX event_logs_anonymous_user_id ON event_logs USING btree (anonymous_user_id) |
| event_logs_name | no | no | no | no | CREATE INDEX event_logs_name ON event_logs USING btree (name) |
| event_logs_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX event_logs_pkey ON event_logs USING btree (id) |
| event_logs_source | no | no | no | no | CREATE INDEX event_logs_source ON event_logs USING btree (source) |
| event_logs_timestamp | no | no | no | no | CREATE INDEX event_logs_timestamp ON event_logs USING btree ("timestamp") |
| event_logs_timestamp_at_utc | no | no | no | no | CREATE INDEX event_logs_timestamp_at_utc ON event_logs USING btree (date(timezone('UTC'::text, "timestamp"))) |
| event_logs_user_id | no | no | no | no | CREATE INDEX event_logs_user_id ON event_logs USING btree (user_id) |
### Check constraints
| Name | Definition |
| --- | --- |
| event_logs_check_has_user | CHECK (user_id = 0 AND anonymous_user_id <> ''::text OR user_id <> 0 AND anonymous_user_id = ''::text OR user_id <> 0 AND anonymous_user_id <> ''::text) |
| event_logs_check_name_not_empty | CHECK (name <> ''::text) |
| event_logs_check_source_not_empty | CHECK (source <> ''::text) |
| event_logs_check_version_not_empty | CHECK (version <> ''::text) |
# Table "public.executor_heartbeats"


Tracks the most recent activity of executors attached to this Sourcegraph instance.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| architecture | text | No |  | The machine architure running the executor. |
| docker_version | text | No |  | The version of Docker used by the executor. |
| executor_version | text | No |  | The version of the executor. |
| first_seen_at | timestamp with time zone | No | now() | The first time a heartbeat from the executor was received. |
| git_version | text | No |  | The version of Git used by the executor. |
| hostname | text | No |  | The uniquely identifying name of the executor. |
| id | integer | No | nextval('executor_heartbeats_id_seq'::regclass) |  |
| ignite_version | text | No |  | The version of Ignite used by the executor. |
| last_seen_at | timestamp with time zone | No | now() | The last time a heartbeat from the executor was received. |
| os | text | No |  | The operating system running the executor. |
| queue_name | text | No |  | The queue name that the executor polls for work. |
| src_cli_version | text | No |  | The version of src-cli used by the executor. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| executor_heartbeats_hostname_key | no | Yes | no | no | CREATE UNIQUE INDEX executor_heartbeats_hostname_key ON executor_heartbeats USING btree (hostname) |
| executor_heartbeats_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX executor_heartbeats_pkey ON executor_heartbeats USING btree (id) |
# Table "public.external_service_repos"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| clone_url | text | No |  |  |
| created_at | timestamp with time zone | No | transaction_timestamp() |  |
| external_service_id | bigint | No |  |  |
| org_id | integer | Yes |  |  |
| repo_id | integer | No |  |  |
| user_id | integer | Yes |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| external_service_repos_clone_url_idx | no | no | no | no | CREATE INDEX external_service_repos_clone_url_idx ON external_service_repos USING btree (clone_url) |
| external_service_repos_idx | no | no | no | no | CREATE INDEX external_service_repos_idx ON external_service_repos USING btree (external_service_id, repo_id) |
| external_service_repos_org_id_idx | no | no | no | no | CREATE INDEX external_service_repos_org_id_idx ON external_service_repos USING btree (org_id) WHERE org_id IS NOT NULL |
| external_service_repos_repo_id_external_service_id_unique | no | Yes | no | no | CREATE UNIQUE INDEX external_service_repos_repo_id_external_service_id_unique ON external_service_repos USING btree (repo_id, external_service_id) |
| external_service_user_repos_idx | no | no | no | no | CREATE INDEX external_service_user_repos_idx ON external_service_repos USING btree (user_id, repo_id) WHERE user_id IS NOT NULL |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| external_service_repos_external_service_id_fkey | external_services | FOREIGN KEY (external_service_id) REFERENCES external_services(id) ON DELETE CASCADE DEFERRABLE |
| external_service_repos_org_id_fkey | orgs | FOREIGN KEY (org_id) REFERENCES orgs(id) ON DELETE CASCADE |
| external_service_repos_repo_id_fkey | repo | FOREIGN KEY (repo_id) REFERENCES repo(id) ON DELETE CASCADE DEFERRABLE |
| external_service_repos_user_id_fkey | users | FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE DEFERRABLE |
# Table "public.external_service_sync_jobs"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| execution_logs | json[] | Yes |  |  |
| external_service_id | bigint | Yes |  |  |
| failure_message | text | Yes |  |  |
| finished_at | timestamp with time zone | Yes |  |  |
| id | integer | No | nextval('external_service_sync_jobs_id_seq'::regclass) |  |
| last_heartbeat_at | timestamp with time zone | Yes |  |  |
| log_contents | text | Yes |  |  |
| num_failures | integer | No | 0 |  |
| num_resets | integer | No | 0 |  |
| process_after | timestamp with time zone | Yes |  |  |
| queued_at | timestamp with time zone | Yes | now() |  |
| started_at | timestamp with time zone | Yes |  |  |
| state | text | No | 'queued'::text |  |
| worker_hostname | text | No | ''::text |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| external_service_sync_jobs_state_idx | no | no | no | no | CREATE INDEX external_service_sync_jobs_state_idx ON external_service_sync_jobs USING btree (state) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| external_services_id_fk | external_services | FOREIGN KEY (external_service_id) REFERENCES external_services(id) ON DELETE CASCADE |
# Table "public.external_services"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| cloud_default | boolean | No | false |  |
| config | text | No |  |  |
| created_at | timestamp with time zone | No | now() |  |
| deleted_at | timestamp with time zone | Yes |  |  |
| display_name | text | No |  |  |
| encryption_key_id | text | No | ''::text |  |
| has_webhooks | boolean | Yes |  |  |
| id | bigint | No | nextval('external_services_id_seq'::regclass) |  |
| kind | text | No |  |  |
| last_sync_at | timestamp with time zone | Yes |  |  |
| namespace_org_id | integer | Yes |  |  |
| namespace_user_id | integer | Yes |  |  |
| next_sync_at | timestamp with time zone | Yes |  |  |
| token_expires_at | timestamp with time zone | Yes |  |  |
| unrestricted | boolean | No | false |  |
| updated_at | timestamp with time zone | No | now() |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| external_services_has_webhooks_idx | no | no | no | no | CREATE INDEX external_services_has_webhooks_idx ON external_services USING btree (has_webhooks) |
| external_services_namespace_org_id_idx | no | no | no | no | CREATE INDEX external_services_namespace_org_id_idx ON external_services USING btree (namespace_org_id) |
| external_services_namespace_user_id_idx | no | no | no | no | CREATE INDEX external_services_namespace_user_id_idx ON external_services USING btree (namespace_user_id) |
| external_services_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX external_services_pkey ON external_services USING btree (id) |
| kind_cloud_default | no | Yes | no | no | CREATE UNIQUE INDEX kind_cloud_default ON external_services USING btree (kind, cloud_default) WHERE cloud_default = true AND deleted_at IS NULL |
### Check constraints
| Name | Definition |
| --- | --- |
| check_non_empty_config | CHECK (btrim(config) <> ''::text) |
| external_services_max_1_namespace | CHECK (namespace_user_id IS NULL AND namespace_org_id IS NULL OR (namespace_user_id IS NULL) <> (namespace_org_id IS NULL)) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| external_services_namepspace_user_id_fkey | users | FOREIGN KEY (namespace_user_id) REFERENCES users(id) ON DELETE CASCADE DEFERRABLE |
| external_services_namespace_org_id_fkey | orgs | FOREIGN KEY (namespace_org_id) REFERENCES orgs(id) ON DELETE CASCADE DEFERRABLE |
### References
| Name | Definition |
| --- | --- |
| external_service_repos | external_service_repos_external_service_id_fkey |
| external_service_sync_jobs | external_services_id_fk |
| webhook_logs | webhook_logs_external_service_id_fkey |
# Table "public.feature_flag_overrides"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| created_at | timestamp with time zone | No | now() |  |
| deleted_at | timestamp with time zone | Yes |  |  |
| flag_name | text | No |  |  |
| flag_value | boolean | No |  |  |
| namespace_org_id | integer | Yes |  |  |
| namespace_user_id | integer | Yes |  |  |
| updated_at | timestamp with time zone | No | now() |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| feature_flag_overrides_org_id | no | no | no | no | CREATE INDEX feature_flag_overrides_org_id ON feature_flag_overrides USING btree (namespace_org_id) WHERE namespace_org_id IS NOT NULL |
| feature_flag_overrides_unique_org_flag | no | Yes | no | no | CREATE UNIQUE INDEX feature_flag_overrides_unique_org_flag ON feature_flag_overrides USING btree (namespace_org_id, flag_name) |
| feature_flag_overrides_unique_user_flag | no | Yes | no | no | CREATE UNIQUE INDEX feature_flag_overrides_unique_user_flag ON feature_flag_overrides USING btree (namespace_user_id, flag_name) |
| feature_flag_overrides_user_id | no | no | no | no | CREATE INDEX feature_flag_overrides_user_id ON feature_flag_overrides USING btree (namespace_user_id) WHERE namespace_user_id IS NOT NULL |
### Check constraints
| Name | Definition |
| --- | --- |
| feature_flag_overrides_has_org_or_user_id | CHECK (namespace_org_id IS NOT NULL OR namespace_user_id IS NOT NULL) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| feature_flag_overrides_flag_name_fkey | feature_flags | FOREIGN KEY (flag_name) REFERENCES feature_flags(flag_name) ON DELETE CASCADE |
| feature_flag_overrides_namespace_org_id_fkey | orgs | FOREIGN KEY (namespace_org_id) REFERENCES orgs(id) ON DELETE CASCADE |
| feature_flag_overrides_namespace_user_id_fkey | users | FOREIGN KEY (namespace_user_id) REFERENCES users(id) ON DELETE CASCADE |
# Table "public.feature_flags"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| bool_value | boolean | Yes |  | Bool value only defined when flag_type is bool |
| created_at | timestamp with time zone | No | now() |  |
| deleted_at | timestamp with time zone | Yes |  |  |
| flag_name | text | No |  |  |
| flag_type | feature_flag_type | No |  |  |
| rollout | integer | Yes |  | Rollout only defined when flag_type is rollout. Increments of 0.01% |
| updated_at | timestamp with time zone | No | now() |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| feature_flags_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX feature_flags_pkey ON feature_flags USING btree (flag_name) |
### Check constraints
| Name | Definition |
| --- | --- |
| feature_flags_rollout_check | CHECK (rollout >= 0 AND rollout <= 10000) |
| required_bool_fields | CHECK (1 =
CASE
    WHEN flag_type = 'bool'::feature_flag_type AND bool_value IS NULL THEN 0
    WHEN flag_type <> 'bool'::feature_flag_type AND bool_value IS NOT NULL THEN 0
    ELSE 1
END) |
| required_rollout_fields | CHECK (1 =
CASE
    WHEN flag_type = 'rollout'::feature_flag_type AND rollout IS NULL THEN 0
    WHEN flag_type <> 'rollout'::feature_flag_type AND rollout IS NOT NULL THEN 0
    ELSE 1
END) |
### References
| Name | Definition |
| --- | --- |
| feature_flag_overrides | feature_flag_overrides_flag_name_fkey |
# Table "public.gitserver_repos"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| clone_status | text | No | 'not_cloned'::text |  |
| last_changed | timestamp with time zone | No | now() |  |
| last_error | text | Yes |  |  |
| last_fetched | timestamp with time zone | No | now() |  |
| repo_id | integer | No |  |  |
| repo_size_bytes | bigint | Yes |  |  |
| shard_id | text | No |  |  |
| updated_at | timestamp with time zone | No | now() |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| gitserver_repos_cloned_status_idx | no | no | no | no | CREATE INDEX gitserver_repos_cloned_status_idx ON gitserver_repos USING btree (repo_id) WHERE clone_status = 'cloned'::text |
| gitserver_repos_cloning_status_idx | no | no | no | no | CREATE INDEX gitserver_repos_cloning_status_idx ON gitserver_repos USING btree (repo_id) WHERE clone_status = 'cloning'::text |
| gitserver_repos_last_error_idx | no | no | no | no | CREATE INDEX gitserver_repos_last_error_idx ON gitserver_repos USING btree (repo_id) WHERE last_error IS NOT NULL |
| gitserver_repos_not_cloned_status_idx | no | no | no | no | CREATE INDEX gitserver_repos_not_cloned_status_idx ON gitserver_repos USING btree (repo_id) WHERE clone_status = 'not_cloned'::text |
| gitserver_repos_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX gitserver_repos_pkey ON gitserver_repos USING btree (repo_id) |
| gitserver_repos_shard_id | no | no | no | no | CREATE INDEX gitserver_repos_shard_id ON gitserver_repos USING btree (shard_id, repo_id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| gitserver_repos_repo_id_fkey | repo | FOREIGN KEY (repo_id) REFERENCES repo(id) ON DELETE CASCADE |
# Table "public.global_state"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| initialized | boolean | No | false |  |
| site_id | uuid | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| global_state_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX global_state_pkey ON global_state USING btree (site_id) |
# Table "public.insights_query_runner_jobs"


See [enterprise/internal/insights/background/queryrunner/worker.go:Job](https://sourcegraph.com/search?q=repo:%5Egithub%5C.com/sourcegraph/sourcegraph%24+file:enterprise/internal/insights/background/queryrunner/worker.go+type+Job&patternType=literal)

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| cost | integer | No | 500 | Integer representing a cost approximation of executing this search query. |
| execution_logs | json[] | Yes |  |  |
| failure_message | text | Yes |  |  |
| finished_at | timestamp with time zone | Yes |  |  |
| id | integer | No | nextval('insights_query_runner_jobs_id_seq'::regclass) |  |
| last_heartbeat_at | timestamp with time zone | Yes |  |  |
| num_failures | integer | No | 0 |  |
| num_resets | integer | No | 0 |  |
| persist_mode | persistmode | No | 'record'::persistmode | The persistence level for this query. This value will determine the lifecycle of the resulting value. |
| priority | integer | No | 1 | Integer representing a category of priority for this query. Priority in this context is ambiguously defined for consumers to decide an interpretation. |
| process_after | timestamp with time zone | Yes |  |  |
| queued_at | timestamp with time zone | Yes | now() |  |
| record_time | timestamp with time zone | Yes |  |  |
| search_query | text | No |  |  |
| series_id | text | No |  |  |
| started_at | timestamp with time zone | Yes |  |  |
| state | text | Yes | 'queued'::text |  |
| worker_hostname | text | No | ''::text |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| insights_query_runner_jobs_cost_idx | no | no | no | no | CREATE INDEX insights_query_runner_jobs_cost_idx ON insights_query_runner_jobs USING btree (cost) |
| insights_query_runner_jobs_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX insights_query_runner_jobs_pkey ON insights_query_runner_jobs USING btree (id) |
| insights_query_runner_jobs_priority_idx | no | no | no | no | CREATE INDEX insights_query_runner_jobs_priority_idx ON insights_query_runner_jobs USING btree (priority) |
| insights_query_runner_jobs_processable_priority_id | no | no | no | no | CREATE INDEX insights_query_runner_jobs_processable_priority_id ON insights_query_runner_jobs USING btree (priority, id) WHERE state = 'queued'::text OR state = 'errored'::text |
| insights_query_runner_jobs_state_btree | no | no | no | no | CREATE INDEX insights_query_runner_jobs_state_btree ON insights_query_runner_jobs USING btree (state) |
### References
| Name | Definition |
| --- | --- |
| insights_query_runner_jobs_dependencies | insights_query_runner_jobs_dependencies_fk_job_id |
# Table "public.insights_query_runner_jobs_dependencies"


Stores data points for a code insight that do not need to be queried directly, but depend on the result of a query at a different point

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| id | integer | No | nextval('insights_query_runner_jobs_dependencies_id_seq'::regclass) |  |
| job_id | integer | No |  | Foreign key to the job that owns this record. |
| recording_time | timestamp without time zone | No |  | The time for which this dependency should be recorded at using the parents value. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| insights_query_runner_jobs_dependencies_job_id_fk_idx | no | no | no | no | CREATE INDEX insights_query_runner_jobs_dependencies_job_id_fk_idx ON insights_query_runner_jobs_dependencies USING btree (job_id) |
| insights_query_runner_jobs_dependencies_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX insights_query_runner_jobs_dependencies_pkey ON insights_query_runner_jobs_dependencies USING btree (id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| insights_query_runner_jobs_dependencies_fk_job_id | insights_query_runner_jobs | FOREIGN KEY (job_id) REFERENCES insights_query_runner_jobs(id) ON DELETE CASCADE |
# Table "public.insights_settings_migration_jobs"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| completed_at | timestamp without time zone | Yes |  |  |
| global | boolean | Yes |  |  |
| id | integer | No | nextval('insights_settings_migration_jobs_id_seq'::regclass) |  |
| migrated_dashboards | integer | No | 0 |  |
| migrated_insights | integer | No | 0 |  |
| org_id | integer | Yes |  |  |
| runs | integer | No | 0 |  |
| settings_id | integer | No |  |  |
| total_dashboards | integer | No | 0 |  |
| total_insights | integer | No | 0 |  |
| user_id | integer | Yes |  |  |
# Table "public.lsif_configuration_policies"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| id | integer | No | nextval('lsif_configuration_policies_id_seq'::regclass) |  |
| index_commit_max_age_hours | integer | Yes |  | The max age of commits indexed by this configuration policy. If null, the age is unbounded. |
| index_intermediate_commits | boolean | No |  | If the matching Git object is a branch, setting this value to true will also index all commits on the matching branches. Setting this value to false will only consider the tip of the branch. |
| indexing_enabled | boolean | No |  | Whether or not this configuration policy affects auto-indexing schedules. |
| last_resolved_at | timestamp with time zone | Yes |  |  |
| name | text | Yes |  |  |
| pattern | text | No |  | A pattern used to match` names of the associated Git object type. |
| protected | boolean | No | false | Whether or not this configuration policy is protected from modification of its data retention behavior (except for duration). |
| repository_id | integer | Yes |  | The identifier of the repository to which this configuration policy applies. If absent, this policy is applied globally. |
| repository_patterns | text[] | Yes |  | The name pattern matching repositories to which this configuration policy applies. If absent, all repositories are matched. |
| retain_intermediate_commits | boolean | No |  | If the matching Git object is a branch, setting this value to true will also retain all data used to resolve queries for any commit on the matching branches. Setting this value to false will only consider the tip of the branch. |
| retention_duration_hours | integer | Yes |  | The max age of data retained by this configuration policy. If null, the age is unbounded. |
| retention_enabled | boolean | No |  | Whether or not this configuration policy affects data retention rules. |
| type | text | No |  | The type of Git object (e.g., COMMIT, BRANCH, TAG). |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_configuration_policies_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_configuration_policies_pkey ON lsif_configuration_policies USING btree (id) |
| lsif_configuration_policies_repository_id | no | no | no | no | CREATE INDEX lsif_configuration_policies_repository_id ON lsif_configuration_policies USING btree (repository_id) |
# Table "public.lsif_configuration_policies_repository_pattern_lookup"


A lookup table to get all the repository patterns by repository id that apply to a configuration policy.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| policy_id | integer | No |  | The policy identifier associated with the repository. |
| repo_id | integer | No |  | The repository identifier associated with the policy. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_configuration_policies_repository_pattern_lookup_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_configuration_policies_repository_pattern_lookup_pkey ON lsif_configuration_policies_repository_pattern_lookup USING btree (policy_id, repo_id) |
# Table "public.lsif_dependency_indexing_jobs"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| execution_logs | json[] | Yes |  |  |
| external_service_kind | text | No | ''::text | Filter the external services for this kind to wait to have synced. If empty, external_service_sync is ignored and no external services are polled for their last sync time. |
| external_service_sync | timestamp with time zone | Yes |  | The sync time after which external services of the given kind will have synced/created any repositories referenced by the LSIF upload that are resolvable. |
| failure_message | text | Yes |  |  |
| finished_at | timestamp with time zone | Yes |  |  |
| id | integer | No | nextval('lsif_dependency_indexing_jobs_id_seq1'::regclass) |  |
| last_heartbeat_at | timestamp with time zone | Yes |  |  |
| num_failures | integer | No | 0 |  |
| num_resets | integer | No | 0 |  |
| process_after | timestamp with time zone | Yes |  |  |
| queued_at | timestamp with time zone | No | now() |  |
| started_at | timestamp with time zone | Yes |  |  |
| state | text | No | 'queued'::text |  |
| upload_id | integer | Yes |  |  |
| worker_hostname | text | No | ''::text |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_dependency_indexing_jobs_pkey1 | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_dependency_indexing_jobs_pkey1 ON lsif_dependency_indexing_jobs USING btree (id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| lsif_dependency_indexing_jobs_upload_id_fkey1 | lsif_uploads | FOREIGN KEY (upload_id) REFERENCES lsif_uploads(id) ON DELETE CASCADE |
# Table "public.lsif_dependency_repos"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| id | bigint | No | nextval('lsif_dependency_repos_id_seq'::regclass) |  |
| name | text | No |  |  |
| scheme | text | No |  |  |
| version | text | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_dependency_repos_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_dependency_repos_pkey ON lsif_dependency_repos USING btree (id) |
| lsif_dependency_repos_unique_triplet | no | Yes | no | no | CREATE UNIQUE INDEX lsif_dependency_repos_unique_triplet ON lsif_dependency_repos USING btree (scheme, name, version) |
# Table "public.lsif_dependency_syncing_jobs"


Tracks jobs that scan imports of indexes to schedule auto-index jobs.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| execution_logs | json[] | Yes |  |  |
| failure_message | text | Yes |  |  |
| finished_at | timestamp with time zone | Yes |  |  |
| id | integer | No | nextval('lsif_dependency_indexing_jobs_id_seq'::regclass) |  |
| last_heartbeat_at | timestamp with time zone | Yes |  |  |
| num_failures | integer | No | 0 |  |
| num_resets | integer | No | 0 |  |
| process_after | timestamp with time zone | Yes |  |  |
| queued_at | timestamp with time zone | No | now() |  |
| started_at | timestamp with time zone | Yes |  |  |
| state | text | No | 'queued'::text |  |
| upload_id | integer | Yes |  | The identifier of the triggering upload record. |
| worker_hostname | text | No | ''::text |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_dependency_indexing_jobs_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_dependency_indexing_jobs_pkey ON lsif_dependency_syncing_jobs USING btree (id) |
| lsif_dependency_indexing_jobs_upload_id | no | no | no | no | CREATE INDEX lsif_dependency_indexing_jobs_upload_id ON lsif_dependency_syncing_jobs USING btree (upload_id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| lsif_dependency_indexing_jobs_upload_id_fkey | lsif_uploads | FOREIGN KEY (upload_id) REFERENCES lsif_uploads(id) ON DELETE CASCADE |
# Table "public.lsif_dirty_repositories"


Stores whether or not the nearest upload data for a repository is out of date (when update_token > dirty_token).

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| dirty_token | integer | No |  | Set to the value of update_token visible to the transaction that updates the commit graph. Updates of dirty_token during this time will cause a second update. |
| repository_id | integer | No |  |  |
| update_token | integer | No |  | This value is incremented on each request to update the commit graph for the repository. |
| updated_at | timestamp with time zone | Yes |  | The time the update_token value was last updated. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_dirty_repositories_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_dirty_repositories_pkey ON lsif_dirty_repositories USING btree (repository_id) |
# Table "public.lsif_index_configuration"


Stores the configuration used for code intel index jobs for a repository.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| autoindex_enabled | boolean | No | true | Whether or not auto-indexing should be attempted on this repo. Index jobs may be inferred from the repository contents if data is empty. |
| data | bytea | No |  | The raw user-supplied [configuration](https://sourcegraph.com/github.com/sourcegraph/sourcegraph@3.23/-/blob/enterprise/internal/codeintel/autoindex/config/types.go#L3:6) (encoded in JSONC). |
| id | bigint | No | nextval('lsif_index_configuration_id_seq'::regclass) |  |
| repository_id | integer | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_index_configuration_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_index_configuration_pkey ON lsif_index_configuration USING btree (id) |
| lsif_index_configuration_repository_id_key | no | Yes | no | no | CREATE UNIQUE INDEX lsif_index_configuration_repository_id_key ON lsif_index_configuration USING btree (repository_id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| lsif_index_configuration_repository_id_fkey | repo | FOREIGN KEY (repository_id) REFERENCES repo(id) ON DELETE CASCADE |
# Table "public.lsif_indexes"


Stores metadata about a code intel index job.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| commit | text | No |  | A 40-char revhash. Note that this commit may not be resolvable in the future. |
| commit_last_checked_at | timestamp with time zone | Yes |  |  |
| docker_steps | jsonb[] | No |  | An array of pre-index [steps](https://sourcegraph.com/github.com/sourcegraph/sourcegraph@3.23/-/blob/enterprise/internal/codeintel/stores/dbstore/docker_step.go#L9:6) to run. |
| execution_logs | json[] | Yes |  | An array of [log entries](https://sourcegraph.com/github.com/sourcegraph/sourcegraph@3.23/-/blob/internal/workerutil/store.go#L48:6) (encoded as JSON) from the most recent execution. |
| failure_message | text | Yes |  |  |
| finished_at | timestamp with time zone | Yes |  |  |
| id | bigint | No | nextval('lsif_indexes_id_seq'::regclass) |  |
| indexer | text | No |  | The docker image used to run the index command (e.g. sourcegraph/lsif-go). |
| indexer_args | text[] | No |  | The command run inside the indexer image to produce the index file (e.g. ['lsif-node', '-p', '.']) |
| last_heartbeat_at | timestamp with time zone | Yes |  |  |
| local_steps | text[] | No |  | A list of commands to run inside the indexer image prior to running the indexer command. |
| log_contents | text | Yes |  | **Column deprecated in favor of execution_logs.** |
| num_failures | integer | No | 0 |  |
| num_resets | integer | No | 0 |  |
| outfile | text | No |  | The path to the index file produced by the index command relative to the working directory. |
| process_after | timestamp with time zone | Yes |  |  |
| queued_at | timestamp with time zone | No | now() |  |
| repository_id | integer | No |  |  |
| root | text | No |  | The working directory of the indexer image relative to the repository root. |
| started_at | timestamp with time zone | Yes |  |  |
| state | text | No | 'queued'::text |  |
| worker_hostname | text | No | ''::text |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_indexes_commit_last_checked_at | no | no | no | no | CREATE INDEX lsif_indexes_commit_last_checked_at ON lsif_indexes USING btree (commit_last_checked_at) WHERE state <> 'deleted'::text |
| lsif_indexes_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_indexes_pkey ON lsif_indexes USING btree (id) |
| lsif_indexes_repository_id_commit | no | no | no | no | CREATE INDEX lsif_indexes_repository_id_commit ON lsif_indexes USING btree (repository_id, commit) |
| lsif_indexes_state | no | no | no | no | CREATE INDEX lsif_indexes_state ON lsif_indexes USING btree (state) |
### Check constraints
| Name | Definition |
| --- | --- |
| lsif_uploads_commit_valid_chars | CHECK (commit ~ '^[a-z0-9]{40}$'::text) |
# Table "public.lsif_last_index_scan"


Tracks the last time repository was checked for auto-indexing job scheduling.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| last_index_scan_at | timestamp with time zone | No |  | The last time uploads of this repository were considered for auto-indexing job scheduling. |
| repository_id | integer | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_last_index_scan_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_last_index_scan_pkey ON lsif_last_index_scan USING btree (repository_id) |
# Table "public.lsif_last_retention_scan"


Tracks the last time uploads a repository were checked against data retention policies.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| last_retention_scan_at | timestamp with time zone | No |  | The last time uploads of this repository were checked against data retention policies. |
| repository_id | integer | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_last_retention_scan_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_last_retention_scan_pkey ON lsif_last_retention_scan USING btree (repository_id) |
# Table "public.lsif_nearest_uploads"


Associates commits with the complete set of uploads visible from that commit. Every commit with upload data is present in this table.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| commit_bytea | bytea | No |  | A 40-char revhash. Note that this commit may not be resolvable in the future. |
| repository_id | integer | No |  |  |
| uploads | jsonb | No |  | Encodes an {upload_id => distance} map that includes an entry for every upload visible from the commit. There is always at least one entry with a distance of zero. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_nearest_uploads_repository_id_commit_bytea | no | no | no | no | CREATE INDEX lsif_nearest_uploads_repository_id_commit_bytea ON lsif_nearest_uploads USING btree (repository_id, commit_bytea) |
| lsif_nearest_uploads_uploads | no | no | no | no | CREATE INDEX lsif_nearest_uploads_uploads ON lsif_nearest_uploads USING gin (uploads) |
# Table "public.lsif_nearest_uploads_links"


Associates commits with the closest ancestor commit with usable upload data. Together, this table and lsif_nearest_uploads cover all commits with resolvable code intelligence.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| ancestor_commit_bytea | bytea | No |  | The 40-char revhash of the ancestor. Note that this commit may not be resolvable in the future. |
| commit_bytea | bytea | No |  | A 40-char revhash. Note that this commit may not be resolvable in the future. |
| distance | integer | No |  | The distance bewteen the commits. Parent = 1, Grandparent = 2, etc. |
| repository_id | integer | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_nearest_uploads_links_repository_id_ancestor_commit_bytea | no | no | no | no | CREATE INDEX lsif_nearest_uploads_links_repository_id_ancestor_commit_bytea ON lsif_nearest_uploads_links USING btree (repository_id, ancestor_commit_bytea) |
| lsif_nearest_uploads_links_repository_id_commit_bytea | no | no | no | no | CREATE INDEX lsif_nearest_uploads_links_repository_id_commit_bytea ON lsif_nearest_uploads_links USING btree (repository_id, commit_bytea) |
# Table "public.lsif_packages"


Associates an upload with the set of packages they provide within a given packages management scheme.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| dump_id | integer | No |  | The identifier of the upload that provides the package. |
| id | integer | No | nextval('lsif_packages_id_seq'::regclass) |  |
| name | text | No |  | The package name. |
| scheme | text | No |  | The (export) moniker scheme. |
| version | text | Yes |  | The package version. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_packages_dump_id | no | no | no | no | CREATE INDEX lsif_packages_dump_id ON lsif_packages USING btree (dump_id) |
| lsif_packages_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_packages_pkey ON lsif_packages USING btree (id) |
| lsif_packages_scheme_name_version_dump_id | no | no | no | no | CREATE INDEX lsif_packages_scheme_name_version_dump_id ON lsif_packages USING btree (scheme, name, version, dump_id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| lsif_packages_dump_id_fkey | lsif_uploads | FOREIGN KEY (dump_id) REFERENCES lsif_uploads(id) ON DELETE CASCADE |
# Table "public.lsif_references"


Associates an upload with the set of packages they require within a given packages management scheme.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| dump_id | integer | No |  | The identifier of the upload that references the package. |
| filter | bytea | No |  | A [bloom filter](https://sourcegraph.com/github.com/sourcegraph/sourcegraph@3.23/-/blob/enterprise/internal/codeintel/bloomfilter/bloom_filter.go#L27:6) encoded as gzipped JSON. This bloom filter stores the set of identifiers imported from the package. |
| id | integer | No | nextval('lsif_references_id_seq'::regclass) |  |
| name | text | No |  | The package name. |
| scheme | text | No |  | The (import) moniker scheme. |
| version | text | Yes |  | The package version. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_references_dump_id | no | no | no | no | CREATE INDEX lsif_references_dump_id ON lsif_references USING btree (dump_id) |
| lsif_references_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_references_pkey ON lsif_references USING btree (id) |
| lsif_references_scheme_name_version_dump_id | no | no | no | no | CREATE INDEX lsif_references_scheme_name_version_dump_id ON lsif_references USING btree (scheme, name, version, dump_id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| lsif_references_dump_id_fkey | lsif_uploads | FOREIGN KEY (dump_id) REFERENCES lsif_uploads(id) ON DELETE CASCADE |
# Table "public.lsif_retention_configuration"


Stores the retention policy of code intellience data for a repository.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| id | integer | No | nextval('lsif_retention_configuration_id_seq'::regclass) |  |
| max_age_for_non_stale_branches_seconds | integer | No |  | The number of seconds since the last modification of a branch until it is considered stale. |
| max_age_for_non_stale_tags_seconds | integer | No |  | The nujmber of seconds since the commit date of a tagged commit until it is considered stale. |
| repository_id | integer | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_retention_configuration_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_retention_configuration_pkey ON lsif_retention_configuration USING btree (id) |
| lsif_retention_configuration_repository_id_key | no | Yes | no | no | CREATE UNIQUE INDEX lsif_retention_configuration_repository_id_key ON lsif_retention_configuration USING btree (repository_id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| lsif_retention_configuration_repository_id_fkey | repo | FOREIGN KEY (repository_id) REFERENCES repo(id) ON DELETE CASCADE |
# Table "public.lsif_uploads"


Stores metadata about an LSIF index uploaded by a user.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| associated_index_id | bigint | Yes |  |  |
| commit | text | No |  | A 40-char revhash. Note that this commit may not be resolvable in the future. |
| commit_last_checked_at | timestamp with time zone | Yes |  |  |
| committed_at | timestamp with time zone | Yes |  |  |
| execution_logs | json[] | Yes |  |  |
| expired | boolean | No | false | Whether or not this upload data is no longer protected by any data retention policy. |
| failure_message | text | Yes |  |  |
| finished_at | timestamp with time zone | Yes |  |  |
| id | integer | No | nextval('lsif_dumps_id_seq'::regclass) | Used as a logical foreign key with the (disjoint) codeintel database. |
| indexer | text | No |  | The name of the indexer that produced the index file. If not supplied by the user it will be pulled from the index metadata. |
| indexer_version | text | Yes |  | The version of the indexer that produced the index file. If not supplied by the user it will be pulled from the index metadata. |
| last_heartbeat_at | timestamp with time zone | Yes |  |  |
| last_retention_scan_at | timestamp with time zone | Yes |  | The last time this upload was checked against data retention policies. |
| num_failures | integer | No | 0 |  |
| num_parts | integer | No |  | The number of parts src-cli split the upload file into. |
| num_references | integer | Yes |  | Deprecated in favor of reference_count. |
| num_resets | integer | No | 0 |  |
| process_after | timestamp with time zone | Yes |  |  |
| queued_at | timestamp with time zone | Yes |  |  |
| reference_count | integer | Yes |  | The number of references to this upload data from other upload records (via lsif_references). |
| repository_id | integer | No |  |  |
| root | text | No | ''::text | The path for which the index can resolve code intelligence relative to the repository root. |
| started_at | timestamp with time zone | Yes |  |  |
| state | text | No | 'queued'::text |  |
| upload_size | bigint | Yes |  | The size of the index file (in bytes). |
| uploaded_at | timestamp with time zone | No | now() |  |
| uploaded_parts | integer[] | No |  | The index of parts that have been successfully uploaded. |
| worker_hostname | text | No | ''::text |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_uploads_associated_index_id | no | no | no | no | CREATE INDEX lsif_uploads_associated_index_id ON lsif_uploads USING btree (associated_index_id) |
| lsif_uploads_commit_last_checked_at | no | no | no | no | CREATE INDEX lsif_uploads_commit_last_checked_at ON lsif_uploads USING btree (commit_last_checked_at) WHERE state <> 'deleted'::text |
| lsif_uploads_committed_at | no | no | no | no | CREATE INDEX lsif_uploads_committed_at ON lsif_uploads USING btree (committed_at) WHERE state = 'completed'::text |
| lsif_uploads_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_uploads_pkey ON lsif_uploads USING btree (id) |
| lsif_uploads_repository_id_commit | no | no | no | no | CREATE INDEX lsif_uploads_repository_id_commit ON lsif_uploads USING btree (repository_id, commit) |
| lsif_uploads_repository_id_commit_root_indexer | no | Yes | no | no | CREATE UNIQUE INDEX lsif_uploads_repository_id_commit_root_indexer ON lsif_uploads USING btree (repository_id, commit, root, indexer) WHERE state = 'completed'::text |
| lsif_uploads_state | no | no | no | no | CREATE INDEX lsif_uploads_state ON lsif_uploads USING btree (state) |
| lsif_uploads_uploaded_at | no | no | no | no | CREATE INDEX lsif_uploads_uploaded_at ON lsif_uploads USING btree (uploaded_at) |
### Check constraints
| Name | Definition |
| --- | --- |
| lsif_uploads_commit_valid_chars | CHECK (commit ~ '^[a-z0-9]{40}$'::text) |
### References
| Name | Definition |
| --- | --- |
| lsif_dependency_syncing_jobs | lsif_dependency_indexing_jobs_upload_id_fkey |
| lsif_dependency_indexing_jobs | lsif_dependency_indexing_jobs_upload_id_fkey1 |
| lsif_packages | lsif_packages_dump_id_fkey |
| lsif_references | lsif_references_dump_id_fkey |
# Table "public.lsif_uploads_visible_at_tip"


Associates a repository with the set of LSIF upload identifiers that can serve intelligence for the tip of the default branch.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| branch_or_tag_name | text | No | ''::text | The name of the branch or tag. |
| is_default_branch | boolean | No | false | Whether the specified branch is the default of the repository. Always false for tags. |
| repository_id | integer | No |  |  |
| upload_id | integer | No |  | The identifier of the upload visible from the tip of the specified branch or tag. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_uploads_visible_at_tip_repository_id_upload_id | no | no | no | no | CREATE INDEX lsif_uploads_visible_at_tip_repository_id_upload_id ON lsif_uploads_visible_at_tip USING btree (repository_id, upload_id) |
# Table "public.migration_logs"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| error_message | text | Yes |  |  |
| finished_at | timestamp with time zone | Yes |  |  |
| id | integer | No | nextval('migration_logs_id_seq'::regclass) |  |
| migration_logs_schema_version | integer | No |  |  |
| schema | text | No |  |  |
| started_at | timestamp with time zone | No |  |  |
| success | boolean | Yes |  |  |
| up | boolean | No |  |  |
| version | integer | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| migration_logs_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX migration_logs_pkey ON migration_logs USING btree (id) |
# Table "public.names"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| name | citext | No |  |  |
| org_id | integer | Yes |  |  |
| user_id | integer | Yes |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| names_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX names_pkey ON names USING btree (name) |
### Check constraints
| Name | Definition |
| --- | --- |
| names_check | CHECK (user_id IS NOT NULL OR org_id IS NOT NULL) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| names_org_id_fkey | orgs | FOREIGN KEY (org_id) REFERENCES orgs(id) ON UPDATE CASCADE ON DELETE CASCADE |
| names_user_id_fkey | users | FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE |
# Table "public.notebook_stars"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| created_at | timestamp with time zone | No | now() |  |
| notebook_id | integer | No |  |  |
| user_id | integer | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| notebook_stars_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX notebook_stars_pkey ON notebook_stars USING btree (notebook_id, user_id) |
| notebook_stars_user_id_idx | no | no | no | no | CREATE INDEX notebook_stars_user_id_idx ON notebook_stars USING btree (user_id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| notebook_stars_notebook_id_fkey | notebooks | FOREIGN KEY (notebook_id) REFERENCES notebooks(id) ON DELETE CASCADE DEFERRABLE |
| notebook_stars_user_id_fkey | users | FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE DEFERRABLE |
# Table "public.notebooks"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| blocks | jsonb | No | '[]'::jsonb |  |
| blocks_tsvector | tsvector | Yes | generated always as (jsonb_to_tsvector('english'::regconfig, blocks, '["string"]'::jsonb)) stored |  |
| created_at | timestamp with time zone | No | now() |  |
| creator_user_id | integer | Yes |  |  |
| id | bigint | No | nextval('notebooks_id_seq'::regclass) |  |
| namespace_org_id | integer | Yes |  |  |
| namespace_user_id | integer | Yes |  |  |
| public | boolean | No |  |  |
| title | text | No |  |  |
| updated_at | timestamp with time zone | No | now() |  |
| updater_user_id | integer | Yes |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| notebooks_blocks_tsvector_idx | no | no | no | no | CREATE INDEX notebooks_blocks_tsvector_idx ON notebooks USING gin (blocks_tsvector) |
| notebooks_namespace_org_id_idx | no | no | no | no | CREATE INDEX notebooks_namespace_org_id_idx ON notebooks USING btree (namespace_org_id) |
| notebooks_namespace_user_id_idx | no | no | no | no | CREATE INDEX notebooks_namespace_user_id_idx ON notebooks USING btree (namespace_user_id) |
| notebooks_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX notebooks_pkey ON notebooks USING btree (id) |
| notebooks_title_trgm_idx | no | no | no | no | CREATE INDEX notebooks_title_trgm_idx ON notebooks USING gin (title gin_trgm_ops) |
### Check constraints
| Name | Definition |
| --- | --- |
| blocks_is_array | CHECK (jsonb_typeof(blocks) = 'array'::text) |
| notebooks_has_max_1_namespace | CHECK (namespace_user_id IS NULL AND namespace_org_id IS NULL OR (namespace_user_id IS NULL) <> (namespace_org_id IS NULL)) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| notebooks_creator_user_id_fkey | users | FOREIGN KEY (creator_user_id) REFERENCES users(id) ON DELETE SET NULL DEFERRABLE |
| notebooks_namespace_org_id_fkey | orgs | FOREIGN KEY (namespace_org_id) REFERENCES orgs(id) ON DELETE SET NULL DEFERRABLE |
| notebooks_namespace_user_id_fkey | users | FOREIGN KEY (namespace_user_id) REFERENCES users(id) ON DELETE SET NULL DEFERRABLE |
| notebooks_updater_user_id_fkey | users | FOREIGN KEY (updater_user_id) REFERENCES users(id) ON DELETE SET NULL DEFERRABLE |
### References
| Name | Definition |
| --- | --- |
| notebook_stars | notebook_stars_notebook_id_fkey |
# Table "public.org_invitations"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| created_at | timestamp with time zone | No | now() |  |
| deleted_at | timestamp with time zone | Yes |  |  |
| expires_at | timestamp with time zone | Yes |  |  |
| id | bigint | No | nextval('org_invitations_id_seq'::regclass) |  |
| notified_at | timestamp with time zone | Yes |  |  |
| org_id | integer | No |  |  |
| recipient_email | citext | Yes |  |  |
| recipient_user_id | integer | Yes |  |  |
| responded_at | timestamp with time zone | Yes |  |  |
| response_type | boolean | Yes |  |  |
| revoked_at | timestamp with time zone | Yes |  |  |
| sender_user_id | integer | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| org_invitations_org_id | no | no | no | no | CREATE INDEX org_invitations_org_id ON org_invitations USING btree (org_id) WHERE deleted_at IS NULL |
| org_invitations_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX org_invitations_pkey ON org_invitations USING btree (id) |
| org_invitations_recipient_user_id | no | no | no | no | CREATE INDEX org_invitations_recipient_user_id ON org_invitations USING btree (recipient_user_id) WHERE deleted_at IS NULL |
### Check constraints
| Name | Definition |
| --- | --- |
| check_atomic_response | CHECK ((responded_at IS NULL) = (response_type IS NULL)) |
| check_single_use | CHECK (responded_at IS NULL AND response_type IS NULL OR revoked_at IS NULL) |
| either_user_id_or_email_defined | CHECK ((recipient_user_id IS NULL) <> (recipient_email IS NULL)) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| org_invitations_org_id_fkey | orgs | FOREIGN KEY (org_id) REFERENCES orgs(id) |
| org_invitations_recipient_user_id_fkey | users | FOREIGN KEY (recipient_user_id) REFERENCES users(id) |
| org_invitations_sender_user_id_fkey | users | FOREIGN KEY (sender_user_id) REFERENCES users(id) |
# Table "public.org_members"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| created_at | timestamp with time zone | No | now() |  |
| id | integer | No | nextval('org_members_id_seq'::regclass) |  |
| org_id | integer | No |  |  |
| updated_at | timestamp with time zone | No | now() |  |
| user_id | integer | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| org_members_org_id_user_id_key | no | Yes | no | no | CREATE UNIQUE INDEX org_members_org_id_user_id_key ON org_members USING btree (org_id, user_id) |
| org_members_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX org_members_pkey ON org_members USING btree (id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| org_members_references_orgs | orgs | FOREIGN KEY (org_id) REFERENCES orgs(id) ON DELETE RESTRICT |
| org_members_user_id_fkey | users | FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT |
# Table "public.org_members_bkup_1514536731"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| created_at | timestamp with time zone | Yes |  |  |
| id | integer | Yes |  |  |
| org_id | integer | Yes |  |  |
| updated_at | timestamp with time zone | Yes |  |  |
| user_id | integer | Yes |  |  |
| user_id_old | text | Yes |  |  |
# Table "public.org_stats"


Business statistics for organizations

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| code_host_repo_count | integer | Yes | 0 | Count of repositories accessible on all code hosts for this organization. |
| org_id | integer | No |  | Org ID that the stats relate to. |
| updated_at | timestamp with time zone | No | now() |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| org_stats_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX org_stats_pkey ON org_stats USING btree (org_id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| org_stats_org_id_fkey | orgs | FOREIGN KEY (org_id) REFERENCES orgs(id) ON DELETE CASCADE DEFERRABLE |
# Table "public.orgs"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| created_at | timestamp with time zone | No | now() |  |
| deleted_at | timestamp with time zone | Yes |  |  |
| display_name | text | Yes |  |  |
| id | integer | No | nextval('orgs_id_seq'::regclass) |  |
| name | citext | No |  |  |
| slack_webhook_url | text | Yes |  |  |
| updated_at | timestamp with time zone | No | now() |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| orgs_name | no | Yes | no | no | CREATE UNIQUE INDEX orgs_name ON orgs USING btree (name) WHERE deleted_at IS NULL |
| orgs_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX orgs_pkey ON orgs USING btree (id) |
### Check constraints
| Name | Definition |
| --- | --- |
| orgs_display_name_max_length | CHECK (char_length(display_name) <= 255) |
| orgs_name_max_length | CHECK (char_length(name::text) <= 255) |
| orgs_name_valid_chars | CHECK (name ~ '^[a-zA-Z0-9](?:[a-zA-Z0-9]|[-.](?=[a-zA-Z0-9]))*-?$'::citext) |
### References
| Name | Definition |
| --- | --- |
| batch_changes | batch_changes_namespace_org_id_fkey |
| cm_monitors | cm_monitors_org_id_fk |
| cm_recipients | cm_recipients_org_id_fk |
| external_service_repos | external_service_repos_org_id_fkey |
| external_services | external_services_namespace_org_id_fkey |
| feature_flag_overrides | feature_flag_overrides_namespace_org_id_fkey |
| names | names_org_id_fkey |
| notebooks | notebooks_namespace_org_id_fkey |
| org_invitations | org_invitations_org_id_fkey |
| org_members | org_members_references_orgs |
| org_stats | org_stats_org_id_fkey |
| registry_extensions | registry_extensions_publisher_org_id_fkey |
| saved_searches | saved_searches_org_id_fkey |
| search_contexts | search_contexts_namespace_org_id_fk |
| settings | settings_references_orgs |
# Table "public.orgs_open_beta_stats"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| created_at | timestamp with time zone | Yes | now() |  |
| data | jsonb | No | '{}'::jsonb |  |
| id | uuid | No | gen_random_uuid() |  |
| org_id | integer | Yes |  |  |
| user_id | integer | Yes |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| orgs_open_beta_stats_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX orgs_open_beta_stats_pkey ON orgs_open_beta_stats USING btree (id) |
# Table "public.out_of_band_migrations"


Stores metadata and progress about an out-of-band migration routine.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| apply_reverse | boolean | No | false | Whether this migration should run in the opposite direction (to support an upcoming downgrade). |
| component | text | No |  | The name of the component undergoing a migration. |
| created | timestamp with time zone | No | now() | The date and time the migration was inserted into the database (via an upgrade). |
| deprecated_version_major | integer | Yes |  | The lowest Sourcegraph version (major component) that assumes the migration has completed. |
| deprecated_version_minor | integer | Yes |  | The lowest Sourcegraph version (minor component) that assumes the migration has completed. |
| description | text | No |  | A brief description about the migration. |
| id | integer | No | nextval('out_of_band_migrations_id_seq'::regclass) | A globally unique primary key for this migration. The same key is used consistently across all Sourcegraph instances for the same migration. |
| introduced_version_major | integer | No |  | The Sourcegraph version (major component) in which this migration was first introduced. |
| introduced_version_minor | integer | No |  | The Sourcegraph version (minor component) in which this migration was first introduced. |
| is_enterprise | boolean | No | false | When true, these migrations are invisible to OSS mode. |
| last_updated | timestamp with time zone | Yes |  | The date and time the migration was last updated. |
| metadata | jsonb | No | '{}'::jsonb |  |
| non_destructive | boolean | No |  | Whether or not this migration alters data so it can no longer be read by the previous Sourcegraph instance. |
| progress | double precision | No | 0 | The percentage progress in the up direction (0=0%, 1=100%). |
| team | text | No |  | The name of the engineering team responsible for the migration. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| out_of_band_migrations_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX out_of_band_migrations_pkey ON out_of_band_migrations USING btree (id) |
### Check constraints
| Name | Definition |
| --- | --- |
| out_of_band_migrations_component_nonempty | CHECK (component <> ''::text) |
| out_of_band_migrations_description_nonempty | CHECK (description <> ''::text) |
| out_of_band_migrations_progress_range | CHECK (progress >= 0::double precision AND progress <= 1::double precision) |
| out_of_band_migrations_team_nonempty | CHECK (team <> ''::text) |
### References
| Name | Definition |
| --- | --- |
| out_of_band_migrations_errors | out_of_band_migrations_errors_migration_id_fkey |
# Table "public.out_of_band_migrations_errors"


Stores errors that occurred while performing an out-of-band migration.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| created | timestamp with time zone | No | now() | The date and time the error occurred. |
| id | integer | No | nextval('out_of_band_migrations_errors_id_seq'::regclass) | A unique identifer. |
| message | text | No |  | The error message. |
| migration_id | integer | No |  | The identifier of the migration. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| out_of_band_migrations_errors_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX out_of_band_migrations_errors_pkey ON out_of_band_migrations_errors USING btree (id) |
### Check constraints
| Name | Definition |
| --- | --- |
| out_of_band_migrations_errors_message_nonempty | CHECK (message <> ''::text) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| out_of_band_migrations_errors_migration_id_fkey | out_of_band_migrations | FOREIGN KEY (migration_id) REFERENCES out_of_band_migrations(id) ON DELETE CASCADE |
# Table "public.phabricator_repos"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| callsign | citext | No |  |  |
| created_at | timestamp with time zone | No | now() |  |
| deleted_at | timestamp with time zone | Yes |  |  |
| id | integer | No | nextval('phabricator_repos_id_seq'::regclass) |  |
| repo_name | citext | No |  |  |
| updated_at | timestamp with time zone | No | now() |  |
| url | text | No | ''::text |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| phabricator_repos_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX phabricator_repos_pkey ON phabricator_repos USING btree (id) |
| phabricator_repos_repo_name_key | no | Yes | no | no | CREATE UNIQUE INDEX phabricator_repos_repo_name_key ON phabricator_repos USING btree (repo_name) |
# Table "public.product_licenses"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| created_at | timestamp with time zone | No | now() |  |
| id | uuid | No |  |  |
| license_key | text | No |  |  |
| product_subscription_id | uuid | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| product_licenses_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX product_licenses_pkey ON product_licenses USING btree (id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| product_licenses_product_subscription_id_fkey | product_subscriptions | FOREIGN KEY (product_subscription_id) REFERENCES product_subscriptions(id) |
# Table "public.product_subscriptions"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| archived_at | timestamp with time zone | Yes |  |  |
| billing_subscription_id | text | Yes |  |  |
| created_at | timestamp with time zone | No | now() |  |
| id | uuid | No |  |  |
| updated_at | timestamp with time zone | No | now() |  |
| user_id | integer | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| product_subscriptions_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX product_subscriptions_pkey ON product_subscriptions USING btree (id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| product_subscriptions_user_id_fkey | users | FOREIGN KEY (user_id) REFERENCES users(id) |
### References
| Name | Definition |
| --- | --- |
| product_licenses | product_licenses_product_subscription_id_fkey |
# Table "public.query_runner_state"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| exec_duration_ns | bigint | Yes |  |  |
| last_executed | timestamp with time zone | Yes |  |  |
| latest_result | timestamp with time zone | Yes |  |  |
| query | text | Yes |  |  |
# Table "public.registry_extension_releases"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| bundle | text | Yes |  |  |
| created_at | timestamp with time zone | No | now() |  |
| creator_user_id | integer | No |  |  |
| deleted_at | timestamp with time zone | Yes |  |  |
| id | bigint | No | nextval('registry_extension_releases_id_seq'::regclass) |  |
| manifest | jsonb | No |  |  |
| registry_extension_id | integer | No |  |  |
| release_tag | citext | No |  |  |
| release_version | citext | Yes |  |  |
| source_map | text | Yes |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| registry_extension_releases_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX registry_extension_releases_pkey ON registry_extension_releases USING btree (id) |
| registry_extension_releases_registry_extension_id | no | no | no | no | CREATE INDEX registry_extension_releases_registry_extension_id ON registry_extension_releases USING btree (registry_extension_id, release_tag, created_at DESC) WHERE deleted_at IS NULL |
| registry_extension_releases_registry_extension_id_created_at | no | no | no | no | CREATE INDEX registry_extension_releases_registry_extension_id_created_at ON registry_extension_releases USING btree (registry_extension_id, created_at) WHERE deleted_at IS NULL |
| registry_extension_releases_version | no | Yes | no | no | CREATE UNIQUE INDEX registry_extension_releases_version ON registry_extension_releases USING btree (registry_extension_id, release_version) WHERE release_version IS NOT NULL |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| registry_extension_releases_creator_user_id_fkey | users | FOREIGN KEY (creator_user_id) REFERENCES users(id) |
| registry_extension_releases_registry_extension_id_fkey | registry_extensions | FOREIGN KEY (registry_extension_id) REFERENCES registry_extensions(id) ON UPDATE CASCADE ON DELETE CASCADE |
# Table "public.registry_extensions"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| created_at | timestamp with time zone | No | now() |  |
| deleted_at | timestamp with time zone | Yes |  |  |
| id | integer | No | nextval('registry_extensions_id_seq'::regclass) |  |
| manifest | text | Yes |  |  |
| name | citext | No |  |  |
| publisher_org_id | integer | Yes |  |  |
| publisher_user_id | integer | Yes |  |  |
| updated_at | timestamp with time zone | No | now() |  |
| uuid | uuid | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| registry_extensions_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX registry_extensions_pkey ON registry_extensions USING btree (id) |
| registry_extensions_publisher_name | no | Yes | no | no | CREATE UNIQUE INDEX registry_extensions_publisher_name ON registry_extensions USING btree (COALESCE(publisher_user_id, 0), COALESCE(publisher_org_id, 0), name) WHERE deleted_at IS NULL |
| registry_extensions_uuid | no | Yes | no | no | CREATE UNIQUE INDEX registry_extensions_uuid ON registry_extensions USING btree (uuid) |
### Check constraints
| Name | Definition |
| --- | --- |
| registry_extensions_name_length | CHECK (char_length(name::text) > 0 AND char_length(name::text) <= 128) |
| registry_extensions_name_valid_chars | CHECK (name ~ '^[a-zA-Z0-9](?:[a-zA-Z0-9]|[_.-](?=[a-zA-Z0-9]))*$'::citext) |
| registry_extensions_single_publisher | CHECK ((publisher_user_id IS NULL) <> (publisher_org_id IS NULL)) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| registry_extensions_publisher_org_id_fkey | orgs | FOREIGN KEY (publisher_org_id) REFERENCES orgs(id) |
| registry_extensions_publisher_user_id_fkey | users | FOREIGN KEY (publisher_user_id) REFERENCES users(id) |
### References
| Name | Definition |
| --- | --- |
| registry_extension_releases | registry_extension_releases_registry_extension_id_fkey |
# Table "public.repo"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| archived | boolean | No | false |  |
| blocked | jsonb | Yes |  |  |
| created_at | timestamp with time zone | No | now() |  |
| deleted_at | timestamp with time zone | Yes |  |  |
| description | text | Yes |  |  |
| external_id | text | Yes |  |  |
| external_service_id | text | Yes |  |  |
| external_service_type | text | Yes |  |  |
| fork | boolean | Yes |  |  |
| id | integer | No | nextval('repo_id_seq'::regclass) |  |
| metadata | jsonb | No | '{}'::jsonb |  |
| name | citext | No |  |  |
| private | boolean | No | false |  |
| stars | integer | No | 0 |  |
| updated_at | timestamp with time zone | Yes |  |  |
| uri | citext | Yes |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| repo_archived | no | no | no | no | CREATE INDEX repo_archived ON repo USING btree (archived) |
| repo_blocked_idx | no | no | no | no | CREATE INDEX repo_blocked_idx ON repo USING btree ((blocked IS NOT NULL)) |
| repo_created_at | no | no | no | no | CREATE INDEX repo_created_at ON repo USING btree (created_at) |
| repo_external_unique_idx | no | Yes | no | no | CREATE UNIQUE INDEX repo_external_unique_idx ON repo USING btree (external_service_type, external_service_id, external_id) |
| repo_fork | no | no | no | no | CREATE INDEX repo_fork ON repo USING btree (fork) |
| repo_hashed_name_idx | no | no | no | no | CREATE INDEX repo_hashed_name_idx ON repo USING btree (sha256(lower(name::text)::bytea)) WHERE deleted_at IS NULL |
| repo_is_not_blocked_idx | no | no | no | no | CREATE INDEX repo_is_not_blocked_idx ON repo USING btree ((blocked IS NULL)) |
| repo_metadata_gin_idx | no | no | no | no | CREATE INDEX repo_metadata_gin_idx ON repo USING gin (metadata) |
| repo_name_case_sensitive_trgm_idx | no | no | no | no | CREATE INDEX repo_name_case_sensitive_trgm_idx ON repo USING gin ((name::text) gin_trgm_ops) |
| repo_name_idx | no | no | no | no | CREATE INDEX repo_name_idx ON repo USING btree (lower(name::text) COLLATE "C") |
| repo_name_trgm | no | no | no | no | CREATE INDEX repo_name_trgm ON repo USING gin (lower(name::text) gin_trgm_ops) |
| repo_name_unique | no | Yes | no | Yes | CREATE UNIQUE INDEX repo_name_unique ON repo USING btree (name) |
| repo_non_deleted_id_name_idx | no | no | no | no | CREATE INDEX repo_non_deleted_id_name_idx ON repo USING btree (id, name) WHERE deleted_at IS NULL |
| repo_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX repo_pkey ON repo USING btree (id) |
| repo_private | no | no | no | no | CREATE INDEX repo_private ON repo USING btree (private) |
| repo_stars_desc_id_desc_idx | no | no | no | no | CREATE INDEX repo_stars_desc_id_desc_idx ON repo USING btree (stars DESC NULLS LAST, id DESC) WHERE deleted_at IS NULL AND blocked IS NULL |
| repo_stars_idx | no | no | no | no | CREATE INDEX repo_stars_idx ON repo USING btree (stars DESC NULLS LAST) |
| repo_uri_idx | no | no | no | no | CREATE INDEX repo_uri_idx ON repo USING btree (uri) |
### Check constraints
| Name | Definition |
| --- | --- |
| check_name_nonempty | CHECK (name <> ''::citext) |
| repo_metadata_check | CHECK (jsonb_typeof(metadata) = 'object'::text) |
### Triggers
| Name | Definition |
| --- | --- |
| trig_delete_repo_ref_on_external_service_repos | CREATE TRIGGER trig_delete_repo_ref_on_external_service_repos AFTER UPDATE OF deleted_at ON repo FOR EACH ROW EXECUTE FUNCTION delete_repo_ref_on_external_service_repos() |
### References
| Name | Definition |
| --- | --- |
| batch_spec_workspaces | batch_spec_workspaces_repo_id_fkey |
| changeset_specs | changeset_specs_repo_id_fkey |
| changesets | changesets_repo_id_fkey |
| discussion_threads_target_repo | discussion_threads_target_repo_repo_id_fkey |
| external_service_repos | external_service_repos_repo_id_fkey |
| gitserver_repos | gitserver_repos_repo_id_fkey |
| lsif_index_configuration | lsif_index_configuration_repository_id_fkey |
| lsif_retention_configuration | lsif_retention_configuration_repository_id_fkey |
| search_context_repos | search_context_repos_repo_id_fk |
| sub_repo_permissions | sub_repo_permissions_repo_id_fk |
| user_public_repos | user_public_repos_repo_id_fkey |
# Table "public.repo_pending_permissions"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| permission | text | No |  |  |
| repo_id | integer | No |  |  |
| updated_at | timestamp with time zone | No |  |  |
| user_ids_ints | integer[] | No | '{}'::integer[] |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| repo_pending_permissions_perm_unique | no | Yes | no | no | CREATE UNIQUE INDEX repo_pending_permissions_perm_unique ON repo_pending_permissions USING btree (repo_id, permission) |
# Table "public.repo_permissions"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| permission | text | No |  |  |
| repo_id | integer | No |  |  |
| synced_at | timestamp with time zone | Yes |  |  |
| updated_at | timestamp with time zone | No |  |  |
| user_ids_ints | integer[] | No | '{}'::integer[] |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| repo_permissions_perm_unique | no | Yes | no | no | CREATE UNIQUE INDEX repo_permissions_perm_unique ON repo_permissions USING btree (repo_id, permission) |
# Table "public.saved_searches"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| created_at | timestamp with time zone | No | now() |  |
| description | text | No |  |  |
| id | integer | No | nextval('saved_searches_id_seq'::regclass) |  |
| notify_owner | boolean | No |  |  |
| notify_slack | boolean | No |  |  |
| org_id | integer | Yes |  |  |
| query | text | No |  |  |
| slack_webhook_url | text | Yes |  |  |
| updated_at | timestamp with time zone | No | now() |  |
| user_id | integer | Yes |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| saved_searches_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX saved_searches_pkey ON saved_searches USING btree (id) |
### Check constraints
| Name | Definition |
| --- | --- |
| saved_searches_notifications_disabled | CHECK (notify_owner = false AND notify_slack = false) |
| user_or_org_id_not_null | CHECK (user_id IS NOT NULL AND org_id IS NULL OR org_id IS NOT NULL AND user_id IS NULL) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| saved_searches_org_id_fkey | orgs | FOREIGN KEY (org_id) REFERENCES orgs(id) |
| saved_searches_user_id_fkey | users | FOREIGN KEY (user_id) REFERENCES users(id) |
# Table "public.search_context_repos"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| repo_id | integer | No |  |  |
| revision | text | No |  |  |
| search_context_id | bigint | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| search_context_repos_unique | no | Yes | no | no | CREATE UNIQUE INDEX search_context_repos_unique ON search_context_repos USING btree (repo_id, search_context_id, revision) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| search_context_repos_repo_id_fk | repo | FOREIGN KEY (repo_id) REFERENCES repo(id) ON DELETE CASCADE |
| search_context_repos_search_context_id_fk | search_contexts | FOREIGN KEY (search_context_id) REFERENCES search_contexts(id) ON DELETE CASCADE |
# Table "public.search_contexts"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| created_at | timestamp with time zone | No | now() |  |
| deleted_at | timestamp with time zone | Yes |  | This column is unused as of Sourcegraph 3.34. Do not refer to it anymore. It will be dropped in a future version. |
| description | text | No |  |  |
| id | bigint | No | nextval('search_contexts_id_seq'::regclass) |  |
| name | citext | No |  |  |
| namespace_org_id | integer | Yes |  |  |
| namespace_user_id | integer | Yes |  |  |
| public | boolean | No |  |  |
| query | text | Yes |  |  |
| updated_at | timestamp with time zone | No | now() |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| search_contexts_name_namespace_org_id_unique | no | Yes | no | no | CREATE UNIQUE INDEX search_contexts_name_namespace_org_id_unique ON search_contexts USING btree (name, namespace_org_id) WHERE namespace_org_id IS NOT NULL |
| search_contexts_name_namespace_user_id_unique | no | Yes | no | no | CREATE UNIQUE INDEX search_contexts_name_namespace_user_id_unique ON search_contexts USING btree (name, namespace_user_id) WHERE namespace_user_id IS NOT NULL |
| search_contexts_name_without_namespace_unique | no | Yes | no | no | CREATE UNIQUE INDEX search_contexts_name_without_namespace_unique ON search_contexts USING btree (name) WHERE namespace_user_id IS NULL AND namespace_org_id IS NULL |
| search_contexts_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX search_contexts_pkey ON search_contexts USING btree (id) |
| search_contexts_query_idx | no | no | no | no | CREATE INDEX search_contexts_query_idx ON search_contexts USING btree (query) |
### Check constraints
| Name | Definition |
| --- | --- |
| search_contexts_has_one_or_no_namespace | CHECK (namespace_user_id IS NULL OR namespace_org_id IS NULL) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| search_contexts_namespace_org_id_fk | orgs | FOREIGN KEY (namespace_org_id) REFERENCES orgs(id) ON DELETE CASCADE |
| search_contexts_namespace_user_id_fk | users | FOREIGN KEY (namespace_user_id) REFERENCES users(id) ON DELETE CASCADE |
### References
| Name | Definition |
| --- | --- |
| search_context_repos | search_context_repos_search_context_id_fk |
# Table "public.security_event_logs"


Contains security-relevant events with a long time horizon for storage.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| anonymous_user_id | text | No |  | The UUID of the actor associated with the event. |
| argument | jsonb | No |  | An arbitrary JSON blob containing event data. |
| id | bigint | No | nextval('security_event_logs_id_seq'::regclass) |  |
| name | text | No |  | The event name as a CAPITALIZED_SNAKE_CASE string. |
| source | text | No |  | The site section (WEB, BACKEND, etc.) that generated the event. |
| timestamp | timestamp with time zone | No |  |  |
| url | text | No |  | The URL within the Sourcegraph app which generated the event. |
| user_id | integer | No |  | The ID of the actor associated with the event. |
| version | text | No |  | The version of Sourcegraph which generated the event. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| security_event_logs_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX security_event_logs_pkey ON security_event_logs USING btree (id) |
| security_event_logs_timestamp | no | no | no | no | CREATE INDEX security_event_logs_timestamp ON security_event_logs USING btree ("timestamp") |
### Check constraints
| Name | Definition |
| --- | --- |
| security_event_logs_check_has_user | CHECK (user_id = 0 AND anonymous_user_id <> ''::text OR user_id <> 0 AND anonymous_user_id = ''::text OR user_id <> 0 AND anonymous_user_id <> ''::text) |
| security_event_logs_check_name_not_empty | CHECK (name <> ''::text) |
| security_event_logs_check_source_not_empty | CHECK (source <> ''::text) |
| security_event_logs_check_version_not_empty | CHECK (version <> ''::text) |
# Table "public.settings"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| author_user_id | integer | Yes |  |  |
| contents | text | No | '{}'::text |  |
| created_at | timestamp with time zone | No | now() |  |
| id | integer | No | nextval('settings_id_seq'::regclass) |  |
| org_id | integer | Yes |  |  |
| user_id | integer | Yes |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| settings_global_id | no | no | no | no | CREATE INDEX settings_global_id ON settings USING btree (id DESC) WHERE user_id IS NULL AND org_id IS NULL |
| settings_org_id_idx | no | no | no | no | CREATE INDEX settings_org_id_idx ON settings USING btree (org_id) |
| settings_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX settings_pkey ON settings USING btree (id) |
| settings_user_id_idx | no | no | no | no | CREATE INDEX settings_user_id_idx ON settings USING btree (user_id) |
### Check constraints
| Name | Definition |
| --- | --- |
| settings_no_empty_contents | CHECK (contents <> ''::text) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| settings_author_user_id_fkey | users | FOREIGN KEY (author_user_id) REFERENCES users(id) ON DELETE RESTRICT |
| settings_references_orgs | orgs | FOREIGN KEY (org_id) REFERENCES orgs(id) ON DELETE RESTRICT |
| settings_user_id_fkey | users | FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT |
# Table "public.settings_bkup_1514702776"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| author_user_id | integer | Yes |  |  |
| author_user_id_old | text | Yes |  |  |
| contents | text | Yes |  |  |
| created_at | timestamp with time zone | Yes |  |  |
| id | integer | Yes |  |  |
| org_id | integer | Yes |  |  |
| user_id | integer | Yes |  |  |
# Table "public.sub_repo_permissions"


Responsible for storing permissions at a finer granularity than repo

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| path_excludes | text[] | Yes |  |  |
| path_includes | text[] | Yes |  |  |
| repo_id | integer | No |  |  |
| updated_at | timestamp with time zone | No | now() |  |
| user_id | integer | No |  |  |
| version | integer | No | 1 |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| sub_repo_permissions_repo_id_user_id_version_uindex | no | Yes | no | no | CREATE UNIQUE INDEX sub_repo_permissions_repo_id_user_id_version_uindex ON sub_repo_permissions USING btree (repo_id, user_id, version) |
| sub_repo_perms_user_id | no | no | no | no | CREATE INDEX sub_repo_perms_user_id ON sub_repo_permissions USING btree (user_id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| sub_repo_permissions_repo_id_fk | repo | FOREIGN KEY (repo_id) REFERENCES repo(id) ON DELETE CASCADE |
| sub_repo_permissions_users_id_fk | users | FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE |
# Table "public.survey_responses"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| better | text | Yes |  |  |
| created_at | timestamp with time zone | No | now() |  |
| email | text | Yes |  |  |
| id | bigint | No | nextval('survey_responses_id_seq'::regclass) |  |
| reason | text | Yes |  |  |
| score | integer | No |  |  |
| user_id | integer | Yes |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| survey_responses_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX survey_responses_pkey ON survey_responses USING btree (id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| survey_responses_user_id_fkey | users | FOREIGN KEY (user_id) REFERENCES users(id) |
# Table "public.temporary_settings"


Stores per-user temporary settings used in the UI, for example, which modals have been dimissed or what theme is preferred.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| contents | jsonb | Yes |  | JSON-encoded temporary settings. |
| created_at | timestamp with time zone | No | now() |  |
| id | integer | No | nextval('temporary_settings_id_seq'::regclass) |  |
| updated_at | timestamp with time zone | No | now() |  |
| user_id | integer | No |  | The ID of the user the settings will be saved for. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| temporary_settings_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX temporary_settings_pkey ON temporary_settings USING btree (id) |
| temporary_settings_user_id_key | no | Yes | no | no | CREATE UNIQUE INDEX temporary_settings_user_id_key ON temporary_settings USING btree (user_id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| temporary_settings_user_id_fkey | users | FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE |
# Table "public.user_credentials"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| created_at | timestamp with time zone | No | now() |  |
| credential | bytea | No |  |  |
| domain | text | No |  |  |
| encryption_key_id | text | No | ''::text |  |
| external_service_id | text | No |  |  |
| external_service_type | text | No |  |  |
| id | bigint | No | nextval('user_credentials_id_seq'::regclass) |  |
| ssh_migration_applied | boolean | No | false |  |
| updated_at | timestamp with time zone | No | now() |  |
| user_id | integer | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| user_credentials_credential_idx | no | no | no | no | CREATE INDEX user_credentials_credential_idx ON user_credentials USING btree ((encryption_key_id = ANY (ARRAY[''::text, 'previously-migrated'::text]))) |
| user_credentials_domain_user_id_external_service_type_exter_key | no | Yes | no | no | CREATE UNIQUE INDEX user_credentials_domain_user_id_external_service_type_exter_key ON user_credentials USING btree (domain, user_id, external_service_type, external_service_id) |
| user_credentials_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX user_credentials_pkey ON user_credentials USING btree (id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| user_credentials_user_id_fkey | users | FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE DEFERRABLE |
# Table "public.user_emails"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| created_at | timestamp with time zone | No | now() |  |
| email | citext | No |  |  |
| is_primary | boolean | No | false |  |
| last_verification_sent_at | timestamp with time zone | Yes |  |  |
| user_id | integer | No |  |  |
| verification_code | text | Yes |  |  |
| verified_at | timestamp with time zone | Yes |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| user_emails_no_duplicates_per_user | no | Yes | no | no | CREATE UNIQUE INDEX user_emails_no_duplicates_per_user ON user_emails USING btree (user_id, email) |
| user_emails_unique_verified_email | no | no | Yes | no | CREATE INDEX user_emails_unique_verified_email ON user_emails USING btree (email) WHERE verified_at IS NOT NULL |
| user_emails_user_id_is_primary_idx | no | Yes | no | no | CREATE UNIQUE INDEX user_emails_user_id_is_primary_idx ON user_emails USING btree (user_id, is_primary) WHERE is_primary = true |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| user_emails_user_id_fkey | users | FOREIGN KEY (user_id) REFERENCES users(id) |
# Table "public.user_external_accounts"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| account_data | text | Yes |  |  |
| account_id | text | No |  |  |
| auth_data | text | Yes |  |  |
| client_id | text | No |  |  |
| created_at | timestamp with time zone | No | now() |  |
| deleted_at | timestamp with time zone | Yes |  |  |
| encryption_key_id | text | No | ''::text |  |
| expired_at | timestamp with time zone | Yes |  |  |
| id | integer | No | nextval('user_external_accounts_id_seq'::regclass) |  |
| last_valid_at | timestamp with time zone | Yes |  |  |
| service_id | text | No |  |  |
| service_type | text | No |  |  |
| updated_at | timestamp with time zone | No | now() |  |
| user_id | integer | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| user_external_accounts_account | no | Yes | no | no | CREATE UNIQUE INDEX user_external_accounts_account ON user_external_accounts USING btree (service_type, service_id, client_id, account_id) WHERE deleted_at IS NULL |
| user_external_accounts_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX user_external_accounts_pkey ON user_external_accounts USING btree (id) |
| user_external_accounts_user_id | no | no | no | no | CREATE INDEX user_external_accounts_user_id ON user_external_accounts USING btree (user_id) WHERE deleted_at IS NULL |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| user_external_accounts_user_id_fkey | users | FOREIGN KEY (user_id) REFERENCES users(id) |
# Table "public.user_pending_permissions"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| bind_id | text | No |  |  |
| id | integer | No | nextval('user_pending_permissions_id_seq'::regclass) |  |
| object_ids_ints | integer[] | No | '{}'::integer[] |  |
| object_type | text | No |  |  |
| permission | text | No |  |  |
| service_id | text | No |  |  |
| service_type | text | No |  |  |
| updated_at | timestamp with time zone | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| user_pending_permissions_service_perm_object_unique | no | Yes | no | no | CREATE UNIQUE INDEX user_pending_permissions_service_perm_object_unique ON user_pending_permissions USING btree (service_type, service_id, permission, object_type, bind_id) |
# Table "public.user_permissions"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| object_ids_ints | integer[] | No | '{}'::integer[] |  |
| object_type | text | No |  |  |
| permission | text | No |  |  |
| synced_at | timestamp with time zone | Yes |  |  |
| updated_at | timestamp with time zone | No |  |  |
| user_id | integer | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| user_permissions_perm_object_unique | no | Yes | no | no | CREATE UNIQUE INDEX user_permissions_perm_object_unique ON user_permissions USING btree (user_id, permission, object_type) |
# Table "public.user_public_repos"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| repo_id | integer | No |  |  |
| repo_uri | text | No |  |  |
| user_id | integer | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| user_public_repos_user_id_repo_id_key | no | Yes | no | no | CREATE UNIQUE INDEX user_public_repos_user_id_repo_id_key ON user_public_repos USING btree (user_id, repo_id) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| user_public_repos_repo_id_fkey | repo | FOREIGN KEY (repo_id) REFERENCES repo(id) ON DELETE CASCADE |
| user_public_repos_user_id_fkey | users | FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE |
# Table "public.users"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| avatar_url | text | Yes |  |  |
| billing_customer_id | text | Yes |  |  |
| created_at | timestamp with time zone | No | now() |  |
| deleted_at | timestamp with time zone | Yes |  |  |
| display_name | text | Yes |  |  |
| id | integer | No | nextval('users_id_seq'::regclass) |  |
| invalidated_sessions_at | timestamp with time zone | No | now() |  |
| invite_quota | integer | No | 15 |  |
| page_views | integer | No | 0 |  |
| passwd | text | Yes |  |  |
| passwd_reset_code | text | Yes |  |  |
| passwd_reset_time | timestamp with time zone | Yes |  |  |
| search_queries | integer | No | 0 |  |
| searchable | boolean | No | true |  |
| site_admin | boolean | No | false |  |
| tags | text[] | Yes | '{}'::text[] |  |
| tos_accepted | boolean | No | false |  |
| updated_at | timestamp with time zone | No | now() |  |
| username | citext | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| users_billing_customer_id | no | Yes | no | no | CREATE UNIQUE INDEX users_billing_customer_id ON users USING btree (billing_customer_id) WHERE deleted_at IS NULL |
| users_created_at_idx | no | no | no | no | CREATE INDEX users_created_at_idx ON users USING btree (created_at) |
| users_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX users_pkey ON users USING btree (id) |
| users_username | no | Yes | no | no | CREATE UNIQUE INDEX users_username ON users USING btree (username) WHERE deleted_at IS NULL |
### Check constraints
| Name | Definition |
| --- | --- |
| users_display_name_max_length | CHECK (char_length(display_name) <= 255) |
| users_username_max_length | CHECK (char_length(username::text) <= 255) |
| users_username_valid_chars | CHECK (username ~ '^[a-zA-Z0-9](?:[a-zA-Z0-9]|[-.](?=[a-zA-Z0-9]))*-?$'::citext) |
### Triggers
| Name | Definition |
| --- | --- |
| trig_invalidate_session_on_password_change | CREATE TRIGGER trig_invalidate_session_on_password_change BEFORE UPDATE OF passwd ON users FOR EACH ROW EXECUTE FUNCTION invalidate_session_for_userid_on_password_change() |
| trig_soft_delete_user_reference_on_external_service | CREATE TRIGGER trig_soft_delete_user_reference_on_external_service AFTER UPDATE OF deleted_at ON users FOR EACH ROW EXECUTE FUNCTION soft_delete_user_reference_on_external_service() |
### References
| Name | Definition |
| --- | --- |
| access_tokens | access_tokens_creator_user_id_fkey |
| access_tokens | access_tokens_subject_user_id_fkey |
| batch_changes | batch_changes_initial_applier_id_fkey |
| batch_changes | batch_changes_last_applier_id_fkey |
| batch_changes | batch_changes_namespace_user_id_fkey |
| batch_spec_execution_cache_entries | batch_spec_execution_cache_entries_user_id_fkey |
| batch_specs | batch_specs_user_id_fkey |
| changeset_jobs | changeset_jobs_user_id_fkey |
| changeset_specs | changeset_specs_user_id_fkey |
| cm_emails | cm_emails_changed_by_fk |
| cm_emails | cm_emails_created_by_fk |
| cm_monitors | cm_monitors_changed_by_fk |
| cm_monitors | cm_monitors_created_by_fk |
| cm_monitors | cm_monitors_user_id_fk |
| cm_recipients | cm_recipients_user_id_fk |
| cm_slack_webhooks | cm_slack_webhooks_changed_by_fkey |
| cm_slack_webhooks | cm_slack_webhooks_created_by_fkey |
| cm_queries | cm_triggers_changed_by_fk |
| cm_queries | cm_triggers_created_by_fk |
| cm_webhooks | cm_webhooks_changed_by_fkey |
| cm_webhooks | cm_webhooks_created_by_fkey |
| discussion_comments | discussion_comments_author_user_id_fkey |
| discussion_mail_reply_tokens | discussion_mail_reply_tokens_user_id_fkey |
| discussion_threads | discussion_threads_author_user_id_fkey |
| external_service_repos | external_service_repos_user_id_fkey |
| external_services | external_services_namepspace_user_id_fkey |
| feature_flag_overrides | feature_flag_overrides_namespace_user_id_fkey |
| names | names_user_id_fkey |
| notebook_stars | notebook_stars_user_id_fkey |
| notebooks | notebooks_creator_user_id_fkey |
| notebooks | notebooks_namespace_user_id_fkey |
| notebooks | notebooks_updater_user_id_fkey |
| org_invitations | org_invitations_recipient_user_id_fkey |
| org_invitations | org_invitations_sender_user_id_fkey |
| org_members | org_members_user_id_fkey |
| product_subscriptions | product_subscriptions_user_id_fkey |
| registry_extension_releases | registry_extension_releases_creator_user_id_fkey |
| registry_extensions | registry_extensions_publisher_user_id_fkey |
| saved_searches | saved_searches_user_id_fkey |
| search_contexts | search_contexts_namespace_user_id_fk |
| settings | settings_author_user_id_fkey |
| settings | settings_user_id_fkey |
| sub_repo_permissions | sub_repo_permissions_users_id_fk |
| survey_responses | survey_responses_user_id_fkey |
| temporary_settings | temporary_settings_user_id_fkey |
| user_credentials | user_credentials_user_id_fkey |
| user_emails | user_emails_user_id_fkey |
| user_external_accounts | user_external_accounts_user_id_fkey |
| user_public_repos | user_public_repos_user_id_fkey |
# Table "public.versions"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| first_version | text | No |  |  |
| service | text | No |  |  |
| updated_at | timestamp with time zone | No | now() |  |
| version | text | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| versions_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX versions_pkey ON versions USING btree (service) |
### Triggers
| Name | Definition |
| --- | --- |
| versions_insert | CREATE TRIGGER versions_insert BEFORE INSERT ON versions FOR EACH ROW EXECUTE FUNCTION versions_insert_row_trigger() |
# Table "public.webhook_logs"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| encryption_key_id | text | No |  |  |
| external_service_id | integer | Yes |  |  |
| id | bigint | No | nextval('webhook_logs_id_seq'::regclass) |  |
| received_at | timestamp with time zone | No | now() |  |
| request | bytea | No |  |  |
| response | bytea | No |  |  |
| status_code | integer | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| webhook_logs_external_service_id_idx | no | no | no | no | CREATE INDEX webhook_logs_external_service_id_idx ON webhook_logs USING btree (external_service_id) |
| webhook_logs_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX webhook_logs_pkey ON webhook_logs USING btree (id) |
| webhook_logs_received_at_idx | no | no | no | no | CREATE INDEX webhook_logs_received_at_idx ON webhook_logs USING btree (received_at) |
| webhook_logs_status_code_idx | no | no | no | no | CREATE INDEX webhook_logs_status_code_idx ON webhook_logs USING btree (status_code) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| webhook_logs_external_service_id_fkey | external_services | FOREIGN KEY (external_service_id) REFERENCES external_services(id) ON UPDATE CASCADE ON DELETE CASCADE |
# View "public.branch_changeset_specs_and_changesets"

## View query:

```sql
 SELECT changeset_specs.id AS changeset_spec_id,
    COALESCE(changesets.id, (0)::bigint) AS changeset_id,
    changeset_specs.repo_id,
    changeset_specs.batch_spec_id,
    changesets.owned_by_batch_change_id AS owner_batch_change_id,
    repo.name AS repo_name,
    changeset_specs.title AS changeset_name,
    changesets.external_state,
    changesets.publication_state,
    changesets.reconciler_state
   FROM ((changeset_specs
     LEFT JOIN changesets ON (((changesets.repo_id = changeset_specs.repo_id) AND (changesets.current_spec_id IS NOT NULL) AND (EXISTS ( SELECT 1
           FROM changeset_specs changeset_specs_1
          WHERE ((changeset_specs_1.id = changesets.current_spec_id) AND (changeset_specs_1.head_ref = changeset_specs.head_ref)))))))
     JOIN repo ON ((changeset_specs.repo_id = repo.id)))
  WHERE ((changeset_specs.external_id IS NULL) AND (repo.deleted_at IS NULL));
```

# View "public.external_service_sync_jobs_with_next_sync_at"

## View query:

```sql
 SELECT j.id,
    j.state,
    j.failure_message,
    j.queued_at,
    j.started_at,
    j.finished_at,
    j.process_after,
    j.num_resets,
    j.num_failures,
    j.execution_logs,
    j.external_service_id,
    e.next_sync_at
   FROM (external_services e
     JOIN external_service_sync_jobs j ON ((e.id = j.external_service_id)));
```

# View "public.lsif_dumps"

## View query:

```sql
 SELECT u.id,
    u.commit,
    u.root,
    u.queued_at,
    u.uploaded_at,
    u.state,
    u.failure_message,
    u.started_at,
    u.finished_at,
    u.repository_id,
    u.indexer,
    u.indexer_version,
    u.num_parts,
    u.uploaded_parts,
    u.process_after,
    u.num_resets,
    u.upload_size,
    u.num_failures,
    u.associated_index_id,
    u.expired,
    u.last_retention_scan_at,
    u.finished_at AS processed_at
   FROM lsif_uploads u
  WHERE ((u.state = 'completed'::text) OR (u.state = 'deleting'::text));
```

# View "public.lsif_dumps_with_repository_name"

## View query:

```sql
 SELECT u.id,
    u.commit,
    u.root,
    u.queued_at,
    u.uploaded_at,
    u.state,
    u.failure_message,
    u.started_at,
    u.finished_at,
    u.repository_id,
    u.indexer,
    u.indexer_version,
    u.num_parts,
    u.uploaded_parts,
    u.process_after,
    u.num_resets,
    u.upload_size,
    u.num_failures,
    u.associated_index_id,
    u.expired,
    u.last_retention_scan_at,
    u.processed_at,
    r.name AS repository_name
   FROM (lsif_dumps u
     JOIN repo r ON ((r.id = u.repository_id)))
  WHERE (r.deleted_at IS NULL);
```

# View "public.lsif_indexes_with_repository_name"

## View query:

```sql
 SELECT u.id,
    u.commit,
    u.queued_at,
    u.state,
    u.failure_message,
    u.started_at,
    u.finished_at,
    u.repository_id,
    u.process_after,
    u.num_resets,
    u.num_failures,
    u.docker_steps,
    u.root,
    u.indexer,
    u.indexer_args,
    u.outfile,
    u.log_contents,
    u.execution_logs,
    u.local_steps,
    r.name AS repository_name
   FROM (lsif_indexes u
     JOIN repo r ON ((r.id = u.repository_id)))
  WHERE (r.deleted_at IS NULL);
```

# View "public.lsif_uploads_with_repository_name"

## View query:

```sql
 SELECT u.id,
    u.commit,
    u.root,
    u.queued_at,
    u.uploaded_at,
    u.state,
    u.failure_message,
    u.started_at,
    u.finished_at,
    u.repository_id,
    u.indexer,
    u.indexer_version,
    u.num_parts,
    u.uploaded_parts,
    u.process_after,
    u.num_resets,
    u.upload_size,
    u.num_failures,
    u.associated_index_id,
    u.expired,
    u.last_retention_scan_at,
    r.name AS repository_name
   FROM (lsif_uploads u
     JOIN repo r ON ((r.id = u.repository_id)))
  WHERE (r.deleted_at IS NULL);
```

# View "public.reconciler_changesets"

## View query:

```sql
 SELECT c.id,
    c.batch_change_ids,
    c.repo_id,
    c.queued_at,
    c.created_at,
    c.updated_at,
    c.metadata,
    c.external_id,
    c.external_service_type,
    c.external_deleted_at,
    c.external_branch,
    c.external_updated_at,
    c.external_state,
    c.external_review_state,
    c.external_check_state,
    c.diff_stat_added,
    c.diff_stat_changed,
    c.diff_stat_deleted,
    c.sync_state,
    c.current_spec_id,
    c.previous_spec_id,
    c.publication_state,
    c.owned_by_batch_change_id,
    c.reconciler_state,
    c.failure_message,
    c.started_at,
    c.finished_at,
    c.process_after,
    c.num_resets,
    c.closing,
    c.num_failures,
    c.log_contents,
    c.execution_logs,
    c.syncer_error,
    c.external_title,
    c.worker_hostname,
    c.ui_publication_state,
    c.last_heartbeat_at,
    c.external_fork_namespace
   FROM (changesets c
     JOIN repo r ON ((r.id = c.repo_id)))
  WHERE ((r.deleted_at IS NULL) AND (EXISTS ( SELECT 1
           FROM ((batch_changes
             LEFT JOIN users namespace_user ON ((batch_changes.namespace_user_id = namespace_user.id)))
             LEFT JOIN orgs namespace_org ON ((batch_changes.namespace_org_id = namespace_org.id)))
          WHERE ((c.batch_change_ids ? (batch_changes.id)::text) AND (namespace_user.deleted_at IS NULL) AND (namespace_org.deleted_at IS NULL)))));
```

# View "public.site_config"

## View query:

```sql
 SELECT global_state.site_id,
    global_state.initialized
   FROM global_state;
```

# View "public.tracking_changeset_specs_and_changesets"

## View query:

```sql
 SELECT changeset_specs.id AS changeset_spec_id,
    COALESCE(changesets.id, (0)::bigint) AS changeset_id,
    changeset_specs.repo_id,
    changeset_specs.batch_spec_id,
    repo.name AS repo_name,
    COALESCE((changesets.metadata ->> 'Title'::text), (changesets.metadata ->> 'title'::text)) AS changeset_name,
    changesets.external_state,
    changesets.publication_state,
    changesets.reconciler_state
   FROM ((changeset_specs
     LEFT JOIN changesets ON (((changesets.repo_id = changeset_specs.repo_id) AND (changesets.external_id = changeset_specs.external_id))))
     JOIN repo ON ((changeset_specs.repo_id = repo.id)))
  WHERE ((changeset_specs.external_id IS NOT NULL) AND (repo.deleted_at IS NULL));
```

# Type batch_changes_changeset_ui_publication_state

- UNPUBLISHED
- DRAFT
- PUBLISHED

# Type cm_email_priority

- NORMAL
- CRITICAL

# Type critical_or_site

- critical
- site

# Type feature_flag_type

- bool
- rollout

# Type lsif_index_state

- queued
- processing
- completed
- errored
- failed

# Type lsif_upload_state

- uploading
- queued
- processing
- completed
- errored
- deleted
- failed

# Type persistmode

- record
- snapshot

