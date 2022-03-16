# Table "public.lsif_data_apidocs_num_dumps"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| count | bigint | Yes |  |  |
# Table "public.lsif_data_apidocs_num_dumps_indexed"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| count | bigint | Yes |  |  |
# Table "public.lsif_data_apidocs_num_pages"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| count | bigint | Yes |  |  |
# Table "public.lsif_data_apidocs_num_search_results_private"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| count | bigint | Yes |  |  |
# Table "public.lsif_data_apidocs_num_search_results_public"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| count | bigint | Yes |  |  |
# Table "public.lsif_data_definitions"


Associates (document, range) pairs with the import monikers attached to the range.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| data | bytea | Yes |  | A gob-encoded payload conforming to an array of [LocationData](https://sourcegraph.com/github.com/sourcegraph/sourcegraph@3.26/-/blob/enterprise/lib/codeintel/semantic/types.go#L106:6) types. |
| dump_id | integer | No |  | The identifier of the associated dump in the lsif_uploads table (state=completed). |
| identifier | text | No |  | The moniker identifier. |
| num_locations | integer | No |  | The number of locations stored in the data field. |
| schema_version | integer | No |  | The schema version of this row - used to determine presence and encoding of data. |
| scheme | text | No |  | The moniker scheme. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_definitions_dump_id_schema_version | no | no | no | no | CREATE INDEX lsif_data_definitions_dump_id_schema_version ON lsif_data_definitions USING btree (dump_id, schema_version) |
| lsif_data_definitions_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_definitions_pkey ON lsif_data_definitions USING btree (dump_id, scheme, identifier) |
### Triggers
| Name | Definition |
| --- | --- |
| lsif_data_definitions_schema_versions_insert | CREATE TRIGGER lsif_data_definitions_schema_versions_insert AFTER INSERT ON lsif_data_definitions REFERENCING NEW TABLE AS newtab FOR EACH STATEMENT EXECUTE FUNCTION update_lsif_data_definitions_schema_versions_insert() |
# Table "public.lsif_data_definitions_schema_versions"


Tracks the range of schema_versions for each upload in the lsif_data_definitions table.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| dump_id | integer | No |  | The identifier of the associated dump in the lsif_uploads table. |
| max_schema_version | integer | Yes |  | An upper-bound on the `lsif_data_definitions.schema_version` where `lsif_data_definitions.dump_id = dump_id`. |
| min_schema_version | integer | Yes |  | A lower-bound on the `lsif_data_definitions.schema_version` where `lsif_data_definitions.dump_id = dump_id`. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_definitions_schema_versions_dump_id_schema_version_bo | no | no | no | no | CREATE INDEX lsif_data_definitions_schema_versions_dump_id_schema_version_bo ON lsif_data_definitions_schema_versions USING btree (dump_id, min_schema_version, max_schema_version) |
| lsif_data_definitions_schema_versions_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_definitions_schema_versions_pkey ON lsif_data_definitions_schema_versions USING btree (dump_id) |
# Table "public.lsif_data_docs_search_current_private"


A table indicating the most current search index for a unique repository, root, and language.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| created_at | timestamp with time zone | No | now() | The time this record was inserted. The records with the latest created_at value for the same repository, root, and language is the only visible one and others will be deleted asynchronously. |
| dump_id | integer | No |  | The associated dump identifier. |
| dump_root | text | No |  | The root of the associated dump. |
| id | integer | No | nextval('lsif_data_docs_search_current_private_id_seq'::regclass) |  |
| lang_name_id | integer | No |  | The interned index name of the associated dump. |
| last_cleanup_scan_at | timestamp with time zone | No | now() | The last time this record was checked as part of a data retention scan. |
| repo_id | integer | No |  | The repository identifier of the associated dump. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_docs_search_current_private_last_cleanup_scan_at | no | no | no | no | CREATE INDEX lsif_data_docs_search_current_private_last_cleanup_scan_at ON lsif_data_docs_search_current_private USING btree (last_cleanup_scan_at) |
| lsif_data_docs_search_current_private_lookup | no | no | no | no | CREATE INDEX lsif_data_docs_search_current_private_lookup ON lsif_data_docs_search_current_private USING btree (repo_id, dump_root, lang_name_id, created_at) INCLUDE (dump_id) |
| lsif_data_docs_search_current_private_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_docs_search_current_private_pkey ON lsif_data_docs_search_current_private USING btree (id) |
# Table "public.lsif_data_docs_search_current_public"


A table indicating the most current search index for a unique repository, root, and language.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| created_at | timestamp with time zone | No | now() | The time this record was inserted. The records with the latest created_at value for the same repository, root, and language is the only visible one and others will be deleted asynchronously. |
| dump_id | integer | No |  | The associated dump identifier. |
| dump_root | text | No |  | The root of the associated dump. |
| id | integer | No | nextval('lsif_data_docs_search_current_public_id_seq'::regclass) |  |
| lang_name_id | integer | No |  | The interned index name of the associated dump. |
| last_cleanup_scan_at | timestamp with time zone | No | now() | The last time this record was checked as part of a data retention scan. |
| repo_id | integer | No |  | The repository identifier of the associated dump. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_docs_search_current_public_last_cleanup_scan_at | no | no | no | no | CREATE INDEX lsif_data_docs_search_current_public_last_cleanup_scan_at ON lsif_data_docs_search_current_public USING btree (last_cleanup_scan_at) |
| lsif_data_docs_search_current_public_lookup | no | no | no | no | CREATE INDEX lsif_data_docs_search_current_public_lookup ON lsif_data_docs_search_current_public USING btree (repo_id, dump_root, lang_name_id, created_at) INCLUDE (dump_id) |
| lsif_data_docs_search_current_public_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_docs_search_current_public_pkey ON lsif_data_docs_search_current_public USING btree (id) |
# Table "public.lsif_data_docs_search_lang_names_private"


Each unique language name being stored in the API docs search index.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| id | bigint | No | nextval('lsif_data_docs_search_lang_names_private_id_seq'::regclass) | The ID of the language name. |
| lang_name | text | No |  | The lowercase language name (go, java, etc.) OR, if unknown, the LSIF indexer name. |
| tsv | tsvector | No |  | Indexed tsvector for the lang_name field. Crafted for ordered, case, and punctuation sensitivity, see data_write_documentation.go:textSearchVector. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_docs_search_lang_names_private_lang_name_key | no | Yes | no | no | CREATE UNIQUE INDEX lsif_data_docs_search_lang_names_private_lang_name_key ON lsif_data_docs_search_lang_names_private USING btree (lang_name) |
| lsif_data_docs_search_lang_names_private_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_docs_search_lang_names_private_pkey ON lsif_data_docs_search_lang_names_private USING btree (id) |
| lsif_data_docs_search_lang_names_private_tsv_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_lang_names_private_tsv_idx ON lsif_data_docs_search_lang_names_private USING gin (tsv) |
### References
| Name | Definition |
| --- | --- |
| lsif_data_docs_search_private | lsif_data_docs_search_private_lang_name_id_fk |
# Table "public.lsif_data_docs_search_lang_names_public"


Each unique language name being stored in the API docs search index.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| id | bigint | No | nextval('lsif_data_docs_search_lang_names_public_id_seq'::regclass) | The ID of the language name. |
| lang_name | text | No |  | The lowercase language name (go, java, etc.) OR, if unknown, the LSIF indexer name. |
| tsv | tsvector | No |  | Indexed tsvector for the lang_name field. Crafted for ordered, case, and punctuation sensitivity, see data_write_documentation.go:textSearchVector. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_docs_search_lang_names_public_lang_name_key | no | Yes | no | no | CREATE UNIQUE INDEX lsif_data_docs_search_lang_names_public_lang_name_key ON lsif_data_docs_search_lang_names_public USING btree (lang_name) |
| lsif_data_docs_search_lang_names_public_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_docs_search_lang_names_public_pkey ON lsif_data_docs_search_lang_names_public USING btree (id) |
| lsif_data_docs_search_lang_names_public_tsv_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_lang_names_public_tsv_idx ON lsif_data_docs_search_lang_names_public USING gin (tsv) |
### References
| Name | Definition |
| --- | --- |
| lsif_data_docs_search_public | lsif_data_docs_search_public_lang_name_id_fk |
# Table "public.lsif_data_docs_search_private"


A tsvector search index over API documentation (private repos only)

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| detail | text | No |  | The detail string (e.g. the full function signature and its docs). See protocol/documentation.go:Documentation |
| dump_id | integer | No |  | The identifier of the associated dump in the lsif_uploads table (state=completed). |
| dump_root | text | No |  | Identical to lsif_dumps.root; The working directory of the indexer image relative to the repository root. |
| id | bigint | No | nextval('lsif_data_docs_search_private_id_seq'::regclass) | The row ID of the search result. |
| label | text | No |  | The label string of the result, e.g. a one-line function signature. See protocol/documentation.go:Documentation |
| label_reverse_tsv | tsvector | No |  | Indexed tsvector for the reverse of the label field, for suffix lexeme/word matching. Crafted for ordered, case, and punctuation sensitivity, see data_write_documentation.go:textSearchVector. |
| label_tsv | tsvector | No |  | Indexed tsvector for the label field. Crafted for ordered, case, and punctuation sensitivity, see data_write_documentation.go:textSearchVector. |
| lang_name_id | integer | No |  | The programming language (or indexer name) that produced the result. Foreign key into lsif_data_docs_search_lang_names_private. |
| path_id | text | No |  | The fully qualified documentation page path ID, e.g. including "#section". See GraphQL codeintel.schema:documentationPage for what this is. |
| repo_id | integer | No |  | The repo ID, from the main app DB repo table. Used to search over a select set of repos by ID. |
| repo_name_id | integer | No |  | The repository name that produced the result. Foreign key into lsif_data_docs_search_repo_names_private. |
| search_key | text | No |  | The search key generated by the indexer, e.g. mux.Router.ServeHTTP. It is language-specific, and likely unique within a repository (but not always.) See protocol/documentation.go:Documentation.SearchKey |
| search_key_reverse_tsv | tsvector | No |  | Indexed tsvector for the reverse of the search_key field, for suffix lexeme/word matching. Crafted for ordered, case, and punctuation sensitivity, see data_write_documentation.go:textSearchVector. |
| search_key_tsv | tsvector | No |  | Indexed tsvector for the search_key field. Crafted for ordered, case, and punctuation sensitivity, see data_write_documentation.go:textSearchVector. |
| tags_id | integer | No |  | The tags from the documentation node. Foreign key into lsif_data_docs_search_tags_private. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_docs_search_private_dump_id_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_private_dump_id_idx ON lsif_data_docs_search_private USING btree (dump_id) |
| lsif_data_docs_search_private_dump_root_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_private_dump_root_idx ON lsif_data_docs_search_private USING btree (dump_root) |
| lsif_data_docs_search_private_label_reverse_tsv_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_private_label_reverse_tsv_idx ON lsif_data_docs_search_private USING gin (label_reverse_tsv) |
| lsif_data_docs_search_private_label_tsv_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_private_label_tsv_idx ON lsif_data_docs_search_private USING gin (label_tsv) |
| lsif_data_docs_search_private_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_docs_search_private_pkey ON lsif_data_docs_search_private USING btree (id) |
| lsif_data_docs_search_private_repo_id_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_private_repo_id_idx ON lsif_data_docs_search_private USING btree (repo_id) |
| lsif_data_docs_search_private_search_key_reverse_tsv_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_private_search_key_reverse_tsv_idx ON lsif_data_docs_search_private USING gin (search_key_reverse_tsv) |
| lsif_data_docs_search_private_search_key_tsv_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_private_search_key_tsv_idx ON lsif_data_docs_search_private USING gin (search_key_tsv) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| lsif_data_docs_search_private_lang_name_id_fk | lsif_data_docs_search_lang_names_private | FOREIGN KEY (lang_name_id) REFERENCES lsif_data_docs_search_lang_names_private(id) |
| lsif_data_docs_search_private_repo_name_id_fk | lsif_data_docs_search_repo_names_private | FOREIGN KEY (repo_name_id) REFERENCES lsif_data_docs_search_repo_names_private(id) |
| lsif_data_docs_search_private_tags_id_fk | lsif_data_docs_search_tags_private | FOREIGN KEY (tags_id) REFERENCES lsif_data_docs_search_tags_private(id) |
### Triggers
| Name | Definition |
| --- | --- |
| lsif_data_docs_search_private_delete | CREATE TRIGGER lsif_data_docs_search_private_delete AFTER DELETE ON lsif_data_docs_search_private REFERENCING OLD TABLE AS oldtbl FOR EACH STATEMENT EXECUTE FUNCTION lsif_data_docs_search_private_delete() |
| lsif_data_docs_search_private_insert | CREATE TRIGGER lsif_data_docs_search_private_insert AFTER INSERT ON lsif_data_docs_search_private REFERENCING NEW TABLE AS newtbl FOR EACH STATEMENT EXECUTE FUNCTION lsif_data_docs_search_private_insert() |
# Table "public.lsif_data_docs_search_public"


A tsvector search index over API documentation (public repos only)

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| detail | text | No |  | The detail string (e.g. the full function signature and its docs). See protocol/documentation.go:Documentation |
| dump_id | integer | No |  | The identifier of the associated dump in the lsif_uploads table (state=completed). |
| dump_root | text | No |  | Identical to lsif_dumps.root; The working directory of the indexer image relative to the repository root. |
| id | bigint | No | nextval('lsif_data_docs_search_public_id_seq'::regclass) | The row ID of the search result. |
| label | text | No |  | The label string of the result, e.g. a one-line function signature. See protocol/documentation.go:Documentation |
| label_reverse_tsv | tsvector | No |  | Indexed tsvector for the reverse of the label field, for suffix lexeme/word matching. Crafted for ordered, case, and punctuation sensitivity, see data_write_documentation.go:textSearchVector. |
| label_tsv | tsvector | No |  | Indexed tsvector for the label field. Crafted for ordered, case, and punctuation sensitivity, see data_write_documentation.go:textSearchVector. |
| lang_name_id | integer | No |  | The programming language (or indexer name) that produced the result. Foreign key into lsif_data_docs_search_lang_names_public. |
| path_id | text | No |  | The fully qualified documentation page path ID, e.g. including "#section". See GraphQL codeintel.schema:documentationPage for what this is. |
| repo_id | integer | No |  | The repo ID, from the main app DB repo table. Used to search over a select set of repos by ID. |
| repo_name_id | integer | No |  | The repository name that produced the result. Foreign key into lsif_data_docs_search_repo_names_public. |
| search_key | text | No |  | The search key generated by the indexer, e.g. mux.Router.ServeHTTP. It is language-specific, and likely unique within a repository (but not always.) See protocol/documentation.go:Documentation.SearchKey |
| search_key_reverse_tsv | tsvector | No |  | Indexed tsvector for the reverse of the search_key field, for suffix lexeme/word matching. Crafted for ordered, case, and punctuation sensitivity, see data_write_documentation.go:textSearchVector. |
| search_key_tsv | tsvector | No |  | Indexed tsvector for the search_key field. Crafted for ordered, case, and punctuation sensitivity, see data_write_documentation.go:textSearchVector. |
| tags_id | integer | No |  | The tags from the documentation node. Foreign key into lsif_data_docs_search_tags_public. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_docs_search_public_dump_id_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_public_dump_id_idx ON lsif_data_docs_search_public USING btree (dump_id) |
| lsif_data_docs_search_public_dump_root_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_public_dump_root_idx ON lsif_data_docs_search_public USING btree (dump_root) |
| lsif_data_docs_search_public_label_reverse_tsv_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_public_label_reverse_tsv_idx ON lsif_data_docs_search_public USING gin (label_reverse_tsv) |
| lsif_data_docs_search_public_label_tsv_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_public_label_tsv_idx ON lsif_data_docs_search_public USING gin (label_tsv) |
| lsif_data_docs_search_public_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_docs_search_public_pkey ON lsif_data_docs_search_public USING btree (id) |
| lsif_data_docs_search_public_repo_id_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_public_repo_id_idx ON lsif_data_docs_search_public USING btree (repo_id) |
| lsif_data_docs_search_public_search_key_reverse_tsv_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_public_search_key_reverse_tsv_idx ON lsif_data_docs_search_public USING gin (search_key_reverse_tsv) |
| lsif_data_docs_search_public_search_key_tsv_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_public_search_key_tsv_idx ON lsif_data_docs_search_public USING gin (search_key_tsv) |
### Foreign key constraints
| Name | References | Definition |
| --- | --- | --- |
| lsif_data_docs_search_public_lang_name_id_fk | lsif_data_docs_search_lang_names_public | FOREIGN KEY (lang_name_id) REFERENCES lsif_data_docs_search_lang_names_public(id) |
| lsif_data_docs_search_public_repo_name_id_fk | lsif_data_docs_search_repo_names_public | FOREIGN KEY (repo_name_id) REFERENCES lsif_data_docs_search_repo_names_public(id) |
| lsif_data_docs_search_public_tags_id_fk | lsif_data_docs_search_tags_public | FOREIGN KEY (tags_id) REFERENCES lsif_data_docs_search_tags_public(id) |
### Triggers
| Name | Definition |
| --- | --- |
| lsif_data_docs_search_public_delete | CREATE TRIGGER lsif_data_docs_search_public_delete AFTER DELETE ON lsif_data_docs_search_public REFERENCING OLD TABLE AS oldtbl FOR EACH STATEMENT EXECUTE FUNCTION lsif_data_docs_search_public_delete() |
| lsif_data_docs_search_public_insert | CREATE TRIGGER lsif_data_docs_search_public_insert AFTER INSERT ON lsif_data_docs_search_public REFERENCING NEW TABLE AS newtbl FOR EACH STATEMENT EXECUTE FUNCTION lsif_data_docs_search_public_insert() |
# Table "public.lsif_data_docs_search_repo_names_private"


Each unique repository name being stored in the API docs search index.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| id | bigint | No | nextval('lsif_data_docs_search_repo_names_private_id_seq'::regclass) | The ID of the repository name. |
| repo_name | text | No |  | The fully qualified name of the repository. |
| reverse_tsv | tsvector | No |  | Indexed tsvector for the reverse of the lang_name field, for suffix lexeme/word matching. Crafted for ordered, case, and punctuation sensitivity, see data_write_documentation.go:textSearchVector. |
| tsv | tsvector | No |  | Indexed tsvector for the lang_name field. Crafted for ordered, case, and punctuation sensitivity, see data_write_documentation.go:textSearchVector. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_docs_search_repo_names_private_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_docs_search_repo_names_private_pkey ON lsif_data_docs_search_repo_names_private USING btree (id) |
| lsif_data_docs_search_repo_names_private_repo_name_key | no | Yes | no | no | CREATE UNIQUE INDEX lsif_data_docs_search_repo_names_private_repo_name_key ON lsif_data_docs_search_repo_names_private USING btree (repo_name) |
| lsif_data_docs_search_repo_names_private_reverse_tsv_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_repo_names_private_reverse_tsv_idx ON lsif_data_docs_search_repo_names_private USING gin (reverse_tsv) |
| lsif_data_docs_search_repo_names_private_tsv_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_repo_names_private_tsv_idx ON lsif_data_docs_search_repo_names_private USING gin (tsv) |
### References
| Name | Definition |
| --- | --- |
| lsif_data_docs_search_private | lsif_data_docs_search_private_repo_name_id_fk |
# Table "public.lsif_data_docs_search_repo_names_public"


Each unique repository name being stored in the API docs search index.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| id | bigint | No | nextval('lsif_data_docs_search_repo_names_public_id_seq'::regclass) | The ID of the repository name. |
| repo_name | text | No |  | The fully qualified name of the repository. |
| reverse_tsv | tsvector | No |  | Indexed tsvector for the reverse of the lang_name field, for suffix lexeme/word matching. Crafted for ordered, case, and punctuation sensitivity, see data_write_documentation.go:textSearchVector. |
| tsv | tsvector | No |  | Indexed tsvector for the lang_name field. Crafted for ordered, case, and punctuation sensitivity, see data_write_documentation.go:textSearchVector. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_docs_search_repo_names_public_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_docs_search_repo_names_public_pkey ON lsif_data_docs_search_repo_names_public USING btree (id) |
| lsif_data_docs_search_repo_names_public_repo_name_key | no | Yes | no | no | CREATE UNIQUE INDEX lsif_data_docs_search_repo_names_public_repo_name_key ON lsif_data_docs_search_repo_names_public USING btree (repo_name) |
| lsif_data_docs_search_repo_names_public_reverse_tsv_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_repo_names_public_reverse_tsv_idx ON lsif_data_docs_search_repo_names_public USING gin (reverse_tsv) |
| lsif_data_docs_search_repo_names_public_tsv_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_repo_names_public_tsv_idx ON lsif_data_docs_search_repo_names_public USING gin (tsv) |
### References
| Name | Definition |
| --- | --- |
| lsif_data_docs_search_public | lsif_data_docs_search_public_repo_name_id_fk |
# Table "public.lsif_data_docs_search_tags_private"


Each uniques sequence of space-separated tags being stored in the API docs search index.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| id | bigint | No | nextval('lsif_data_docs_search_tags_private_id_seq'::regclass) | The ID of the tags. |
| tags | text | No |  | The full sequence of space-separated tags. See protocol/documentation.go:Documentation |
| tsv | tsvector | No |  | Indexed tsvector for the tags field. Crafted for ordered, case, and punctuation sensitivity, see data_write_documentation.go:textSearchVector. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_docs_search_tags_private_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_docs_search_tags_private_pkey ON lsif_data_docs_search_tags_private USING btree (id) |
| lsif_data_docs_search_tags_private_tags_key | no | Yes | no | no | CREATE UNIQUE INDEX lsif_data_docs_search_tags_private_tags_key ON lsif_data_docs_search_tags_private USING btree (tags) |
| lsif_data_docs_search_tags_private_tsv_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_tags_private_tsv_idx ON lsif_data_docs_search_tags_private USING gin (tsv) |
### References
| Name | Definition |
| --- | --- |
| lsif_data_docs_search_private | lsif_data_docs_search_private_tags_id_fk |
# Table "public.lsif_data_docs_search_tags_public"


Each uniques sequence of space-separated tags being stored in the API docs search index.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| id | bigint | No | nextval('lsif_data_docs_search_tags_public_id_seq'::regclass) | The ID of the tags. |
| tags | text | No |  | The full sequence of space-separated tags. See protocol/documentation.go:Documentation |
| tsv | tsvector | No |  | Indexed tsvector for the tags field. Crafted for ordered, case, and punctuation sensitivity, see data_write_documentation.go:textSearchVector. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_docs_search_tags_public_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_docs_search_tags_public_pkey ON lsif_data_docs_search_tags_public USING btree (id) |
| lsif_data_docs_search_tags_public_tags_key | no | Yes | no | no | CREATE UNIQUE INDEX lsif_data_docs_search_tags_public_tags_key ON lsif_data_docs_search_tags_public USING btree (tags) |
| lsif_data_docs_search_tags_public_tsv_idx | no | no | no | no | CREATE INDEX lsif_data_docs_search_tags_public_tsv_idx ON lsif_data_docs_search_tags_public USING gin (tsv) |
### References
| Name | Definition |
| --- | --- |
| lsif_data_docs_search_public | lsif_data_docs_search_public_tags_id_fk |
# Table "public.lsif_data_documentation_mappings"


Maps documentation path IDs to their corresponding integral documentationResult vertex IDs, which are unique within a dump.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| dump_id | integer | No |  | The identifier of the associated dump in the lsif_uploads table (state=completed). |
| file_path | text | Yes |  | The document file path for the documentationResult, if any. e.g. the path to the file where the symbol described by this documentationResult is located, if it is a symbol. |
| path_id | text | No |  | The documentation page path ID, see see GraphQL codeintel.schema:documentationPage for what this is. |
| result_id | integer | No |  | The documentationResult vertex ID. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_documentation_mappings_inverse_unique_idx | no | Yes | no | no | CREATE UNIQUE INDEX lsif_data_documentation_mappings_inverse_unique_idx ON lsif_data_documentation_mappings USING btree (dump_id, result_id) |
| lsif_data_documentation_mappings_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_documentation_mappings_pkey ON lsif_data_documentation_mappings USING btree (dump_id, path_id) |
# Table "public.lsif_data_documentation_pages"


Associates documentation pathIDs to their documentation page hierarchy chunk.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| data | bytea | Yes |  | A gob-encoded payload conforming to a `type DocumentationPageData struct` pointer (lib/codeintel/semantic/types.go) |
| dump_id | integer | No |  | The identifier of the associated dump in the lsif_uploads table (state=completed). |
| path_id | text | No |  | The documentation page path ID, see see GraphQL codeintel.schema:documentationPage for what this is. |
| search_indexed | boolean | Yes | false |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_documentation_pages_dump_id_unindexed | no | no | no | no | CREATE INDEX lsif_data_documentation_pages_dump_id_unindexed ON lsif_data_documentation_pages USING btree (dump_id) WHERE NOT search_indexed |
| lsif_data_documentation_pages_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_documentation_pages_pkey ON lsif_data_documentation_pages USING btree (dump_id, path_id) |
### Triggers
| Name | Definition |
| --- | --- |
| lsif_data_documentation_pages_delete | CREATE TRIGGER lsif_data_documentation_pages_delete AFTER DELETE ON lsif_data_documentation_pages REFERENCING OLD TABLE AS oldtbl FOR EACH STATEMENT EXECUTE FUNCTION lsif_data_documentation_pages_delete() |
| lsif_data_documentation_pages_insert | CREATE TRIGGER lsif_data_documentation_pages_insert AFTER INSERT ON lsif_data_documentation_pages REFERENCING NEW TABLE AS newtbl FOR EACH STATEMENT EXECUTE FUNCTION lsif_data_documentation_pages_insert() |
| lsif_data_documentation_pages_update | CREATE TRIGGER lsif_data_documentation_pages_update AFTER UPDATE ON lsif_data_documentation_pages REFERENCING OLD TABLE AS oldtbl NEW TABLE AS newtbl FOR EACH STATEMENT EXECUTE FUNCTION lsif_data_documentation_pages_update() |
# Table "public.lsif_data_documentation_path_info"


Associates documentation page pathIDs to information about what is at that pathID, its immediate children, etc.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| data | bytea | Yes |  | A gob-encoded payload conforming to a `type DocumentationPathInoData struct` pointer (lib/codeintel/semantic/types.go) |
| dump_id | integer | No |  | The identifier of the associated dump in the lsif_uploads table (state=completed). |
| path_id | text | No |  | The documentation page path ID, see see GraphQL codeintel.schema:documentationPage for what this is. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_documentation_path_info_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_documentation_path_info_pkey ON lsif_data_documentation_path_info USING btree (dump_id, path_id) |
# Table "public.lsif_data_documents"


Stores reference, hover text, moniker, and diagnostic data about a particular text document witin a dump.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| data | bytea | Yes |  | A gob-encoded payload conforming to the [DocumentData](https://sourcegraph.com/github.com/sourcegraph/sourcegraph@3.26/-/blob/enterprise/lib/codeintel/semantic/types.go#L13:6) type. This field is being migrated across ranges, hovers, monikers, packages, and diagnostics columns and will be removed in a future release of Sourcegraph. |
| diagnostics | bytea | Yes |  | A gob-encoded payload conforming to the [Diagnostics](https://sourcegraph.com/github.com/sourcegraph/sourcegraph@3.26/-/blob/enterprise/lib/codeintel/semantic/types.go#L18:2) field of the DocumentDatatype. |
| dump_id | integer | No |  | The identifier of the associated dump in the lsif_uploads table (state=completed). |
| hovers | bytea | Yes |  | A gob-encoded payload conforming to the [HoversResults](https://sourcegraph.com/github.com/sourcegraph/sourcegraph@3.26/-/blob/enterprise/lib/codeintel/semantic/types.go#L15:2) field of the DocumentDatatype. |
| monikers | bytea | Yes |  | A gob-encoded payload conforming to the [Monikers](https://sourcegraph.com/github.com/sourcegraph/sourcegraph@3.26/-/blob/enterprise/lib/codeintel/semantic/types.go#L16:2) field of the DocumentDatatype. |
| num_diagnostics | integer | No |  | The number of diagnostics stored in the data field. |
| packages | bytea | Yes |  | A gob-encoded payload conforming to the [PackageInformation](https://sourcegraph.com/github.com/sourcegraph/sourcegraph@3.26/-/blob/enterprise/lib/codeintel/semantic/types.go#L17:2) field of the DocumentDatatype. |
| path | text | No |  | The path of the text document relative to the associated dump root. |
| ranges | bytea | Yes |  | A gob-encoded payload conforming to the [Ranges](https://sourcegraph.com/github.com/sourcegraph/sourcegraph@3.26/-/blob/enterprise/lib/codeintel/semantic/types.go#L14:2) field of the DocumentDatatype. |
| schema_version | integer | No |  | The schema version of this row - used to determine presence and encoding of data. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_documents_dump_id_schema_version | no | no | no | no | CREATE INDEX lsif_data_documents_dump_id_schema_version ON lsif_data_documents USING btree (dump_id, schema_version) |
| lsif_data_documents_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_documents_pkey ON lsif_data_documents USING btree (dump_id, path) |
### Triggers
| Name | Definition |
| --- | --- |
| lsif_data_documents_schema_versions_insert | CREATE TRIGGER lsif_data_documents_schema_versions_insert AFTER INSERT ON lsif_data_documents REFERENCING NEW TABLE AS newtab FOR EACH STATEMENT EXECUTE FUNCTION update_lsif_data_documents_schema_versions_insert() |
# Table "public.lsif_data_documents_schema_versions"


Tracks the range of schema_versions for each upload in the lsif_data_documents table.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| dump_id | integer | No |  | The identifier of the associated dump in the lsif_uploads table. |
| max_schema_version | integer | Yes |  | An upper-bound on the `lsif_data_documents.schema_version` where `lsif_data_documents.dump_id = dump_id`. |
| min_schema_version | integer | Yes |  | A lower-bound on the `lsif_data_documents.schema_version` where `lsif_data_documents.dump_id = dump_id`. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_documents_schema_versions_dump_id_schema_version_boun | no | no | no | no | CREATE INDEX lsif_data_documents_schema_versions_dump_id_schema_version_boun ON lsif_data_documents_schema_versions USING btree (dump_id, min_schema_version, max_schema_version) |
| lsif_data_documents_schema_versions_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_documents_schema_versions_pkey ON lsif_data_documents_schema_versions USING btree (dump_id) |
# Table "public.lsif_data_implementations"


Associates (document, range) pairs with the implementation monikers attached to the range.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| data | bytea | Yes |  | A gob-encoded payload conforming to an array of [LocationData](https://sourcegraph.com/github.com/sourcegraph/sourcegraph@3.26/-/blob/enterprise/lib/codeintel/semantic/types.go#L106:6) types. |
| dump_id | integer | No |  | The identifier of the associated dump in the lsif_uploads table (state=completed). |
| identifier | text | No |  | The moniker identifier. |
| num_locations | integer | No |  | The number of locations stored in the data field. |
| schema_version | integer | No |  | The schema version of this row - used to determine presence and encoding of data. |
| scheme | text | No |  | The moniker scheme. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_implementations_dump_id_schema_version | no | no | no | no | CREATE INDEX lsif_data_implementations_dump_id_schema_version ON lsif_data_implementations USING btree (dump_id, schema_version) |
| lsif_data_implementations_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_implementations_pkey ON lsif_data_implementations USING btree (dump_id, scheme, identifier) |
### Triggers
| Name | Definition |
| --- | --- |
| lsif_data_implementations_schema_versions_insert | CREATE TRIGGER lsif_data_implementations_schema_versions_insert AFTER INSERT ON lsif_data_implementations REFERENCING NEW TABLE AS newtab FOR EACH STATEMENT EXECUTE FUNCTION update_lsif_data_implementations_schema_versions_insert() |
# Table "public.lsif_data_implementations_schema_versions"


Tracks the range of schema_versions for each upload in the lsif_data_implementations table.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| dump_id | integer | No |  | The identifier of the associated dump in the lsif_uploads table. |
| max_schema_version | integer | Yes |  | An upper-bound on the `lsif_data_implementations.schema_version` where `lsif_data_implementations.dump_id = dump_id`. |
| min_schema_version | integer | Yes |  | A lower-bound on the `lsif_data_implementations.schema_version` where `lsif_data_implementations.dump_id = dump_id`. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_implementations_schema_versions_dump_id_schema_versio | no | no | no | no | CREATE INDEX lsif_data_implementations_schema_versions_dump_id_schema_versio ON lsif_data_implementations_schema_versions USING btree (dump_id, min_schema_version, max_schema_version) |
| lsif_data_implementations_schema_versions_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_implementations_schema_versions_pkey ON lsif_data_implementations_schema_versions USING btree (dump_id) |
# Table "public.lsif_data_metadata"


Stores the number of result chunks associated with a dump.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| dump_id | integer | No |  | The identifier of the associated dump in the lsif_uploads table (state=completed). |
| num_result_chunks | integer | Yes |  | A bound of populated indexes in the lsif_data_result_chunks table for the associated dump. This value is used to hash identifiers into the result chunk index to which they belong. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_metadata_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_metadata_pkey ON lsif_data_metadata USING btree (dump_id) |
# Table "public.lsif_data_references"


Associates (document, range) pairs with the export monikers attached to the range.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| data | bytea | Yes |  | A gob-encoded payload conforming to an array of [LocationData](https://sourcegraph.com/github.com/sourcegraph/sourcegraph@3.26/-/blob/enterprise/lib/codeintel/semantic/types.go#L106:6) types. |
| dump_id | integer | No |  | The identifier of the associated dump in the lsif_uploads table (state=completed). |
| identifier | text | No |  | The moniker identifier. |
| num_locations | integer | No |  | The number of locations stored in the data field. |
| schema_version | integer | No |  | The schema version of this row - used to determine presence and encoding of data. |
| scheme | text | No |  | The moniker scheme. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_references_dump_id_schema_version | no | no | no | no | CREATE INDEX lsif_data_references_dump_id_schema_version ON lsif_data_references USING btree (dump_id, schema_version) |
| lsif_data_references_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_references_pkey ON lsif_data_references USING btree (dump_id, scheme, identifier) |
### Triggers
| Name | Definition |
| --- | --- |
| lsif_data_references_schema_versions_insert | CREATE TRIGGER lsif_data_references_schema_versions_insert AFTER INSERT ON lsif_data_references REFERENCING NEW TABLE AS newtab FOR EACH STATEMENT EXECUTE FUNCTION update_lsif_data_references_schema_versions_insert() |
# Table "public.lsif_data_references_schema_versions"


Tracks the range of schema_versions for each upload in the lsif_data_references table.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| dump_id | integer | No |  | The identifier of the associated dump in the lsif_uploads table. |
| max_schema_version | integer | Yes |  | An upper-bound on the `lsif_data_references.schema_version` where `lsif_data_references.dump_id = dump_id`. |
| min_schema_version | integer | Yes |  | A lower-bound on the `lsif_data_references.schema_version` where `lsif_data_references.dump_id = dump_id`. |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_references_schema_versions_dump_id_schema_version_bou | no | no | no | no | CREATE INDEX lsif_data_references_schema_versions_dump_id_schema_version_bou ON lsif_data_references_schema_versions USING btree (dump_id, min_schema_version, max_schema_version) |
| lsif_data_references_schema_versions_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_references_schema_versions_pkey ON lsif_data_references_schema_versions USING btree (dump_id) |
# Table "public.lsif_data_result_chunks"


Associates result set identifiers with the (document path, range identifier) pairs that compose the set.

| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| data | bytea | Yes |  | A gob-encoded payload conforming to the [ResultChunkData](https://sourcegraph.com/github.com/sourcegraph/sourcegraph@3.26/-/blob/enterprise/lib/codeintel/semantic/types.go#L76:6) type. |
| dump_id | integer | No |  | The identifier of the associated dump in the lsif_uploads table (state=completed). |
| idx | integer | No |  | The unique result chunk index within the associated dump. Every result set identifier present should hash to this index (modulo lsif_data_metadata.num_result_chunks). |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| lsif_data_result_chunks_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX lsif_data_result_chunks_pkey ON lsif_data_result_chunks USING btree (dump_id, idx) |
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
# Table "public.rockskip_ancestry"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| ancestor | integer | No |  |  |
| commit_id | character varying(40) | No |  |  |
| height | integer | No |  |  |
| id | integer | No | nextval('rockskip_ancestry_id_seq'::regclass) |  |
| repo_id | integer | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| rockskip_ancestry_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX rockskip_ancestry_pkey ON rockskip_ancestry USING btree (id) |
| rockskip_ancestry_repo_commit_id | no | no | no | no | CREATE INDEX rockskip_ancestry_repo_commit_id ON rockskip_ancestry USING btree (repo_id, commit_id) |
| rockskip_ancestry_repo_id_commit_id_key | no | Yes | no | no | CREATE UNIQUE INDEX rockskip_ancestry_repo_id_commit_id_key ON rockskip_ancestry USING btree (repo_id, commit_id) |
# Table "public.rockskip_repos"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| id | integer | No | nextval('rockskip_repos_id_seq'::regclass) |  |
| last_accessed_at | timestamp with time zone | No |  |  |
| repo | text | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| rockskip_repos_last_accessed_at | no | no | no | no | CREATE INDEX rockskip_repos_last_accessed_at ON rockskip_repos USING btree (last_accessed_at) |
| rockskip_repos_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX rockskip_repos_pkey ON rockskip_repos USING btree (id) |
| rockskip_repos_repo | no | no | no | no | CREATE INDEX rockskip_repos_repo ON rockskip_repos USING btree (repo) |
| rockskip_repos_repo_key | no | Yes | no | no | CREATE UNIQUE INDEX rockskip_repos_repo_key ON rockskip_repos USING btree (repo) |
# Table "public.rockskip_symbols"


| Column | Type | Nullable | Default | Comment |
| --- | --- | --- | --- | --- |
| added | integer[] | No |  |  |
| deleted | integer[] | No |  |  |
| id | integer | No | nextval('rockskip_symbols_id_seq'::regclass) |  |
| name | text | No |  |  |
| path | text | No |  |  |
| repo_id | integer | No |  |  |
### Indexes
| Name | IsPrimaryKey | IsUnique | IsExclusion | IsDeferrable | IndexDefinition |
| --- | --- | --- | --- | --- | --- |
| rockskip_symbols_gin | no | no | no | no | CREATE INDEX rockskip_symbols_gin ON rockskip_symbols USING gin (singleton_integer(repo_id) gin__int_ops, added gin__int_ops, deleted gin__int_ops, name gin_trgm_ops, singleton(name), singleton(lower(name)), path gin_trgm_ops, singleton(path), path_prefixes(path), singleton(lower(path)), path_prefixes(lower(path)), singleton(get_file_extension(path)), singleton(get_file_extension(lower(path)))) |
| rockskip_symbols_pkey | Yes | Yes | no | no | CREATE UNIQUE INDEX rockskip_symbols_pkey ON rockskip_symbols USING btree (id) |
| rockskip_symbols_repo_id_path_name | no | no | no | no | CREATE INDEX rockskip_symbols_repo_id_path_name ON rockskip_symbols USING btree (repo_id, path, name) |