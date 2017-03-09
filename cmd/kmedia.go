package cmd

import (
	"database/sql"

	"github.com/Bnei-Baruch/mdb/importer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	command := &cobra.Command{
		Use:   "kmedia",
		Short: "Migrate kmedia to MDB",
		Run: func(cmd *cobra.Command, args []string) {
			importer.ImportKmedia()
		},
	}
	RootCmd.AddCommand(command)

	command = &cobra.Command{
		Use:   "kmedia-fkeys",
		Short: "Add foreign keys to kmedia",
		Long:  "Add foreign keys to kmedia, then run:\n\tsqlboiler -o gmodels_old -p gmodels --no-hooks postgres",
		Run: func(cmd *cobra.Command, args []string) {
			createForeignKeys()
		},
	}
	RootCmd.AddCommand(command)
}

// go run main.go kmedia-fkeys
// sqlboiler -o gmodels_old -p gmodels --no-hooks --no-tests postgres
func createForeignKeys() {
	db, err := sql.Open("postgres", viper.GetString("kmedia.url"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`

ALTER TABLE languages DROP CONSTRAINT IF EXISTS code3_unique CASCADE;
ALTER TABLE languages ADD CONSTRAINT code3_unique UNIQUE (code3);

ALTER TABLE catalogs DROP CONSTRAINT IF EXISTS catalogs_fkey;
ALTER TABLE catalogs ADD CONSTRAINT catalogs_fkey FOREIGN KEY (parent_id) REFERENCES catalogs(id) NOT VALID;
ALTER TABLE catalogs DROP CONSTRAINT IF EXISTS users_fkey;
ALTER TABLE catalogs ADD CONSTRAINT users_fkey FOREIGN KEY (user_id) REFERENCES users(id) NOT VALID;

ALTER TABLE catalogs_containers DROP CONSTRAINT IF EXISTS catalogs_fkey;
ALTER TABLE catalogs_containers ADD CONSTRAINT catalogs_fkey FOREIGN KEY (catalog_id) REFERENCES catalogs(id) NOT VALID;
ALTER TABLE catalogs_containers DROP CONSTRAINT IF EXISTS containers_fkey;
ALTER TABLE catalogs_containers ADD CONSTRAINT containers_fkey FOREIGN KEY (container_id) REFERENCES containers(id) NOT VALID;

DO $$
	BEGIN
		ALTER TABLE catalog_descriptions RENAME COLUMN lang TO lang_id;
		EXCEPTION WHEN OTHERS THEN RAISE NOTICE '%', 'Already Exists';
	END
$$ LANGUAGE plpgsql;
ALTER TABLE catalog_descriptions DROP CONSTRAINT IF EXISTS catalogs_fkey;
ALTER TABLE catalog_descriptions ADD CONSTRAINT catalogs_fkey FOREIGN KEY (catalog_id) REFERENCES catalogs(id) NOT VALID;
ALTER TABLE catalog_descriptions DROP CONSTRAINT IF EXISTS languages_fkey;
ALTER TABLE catalog_descriptions ADD CONSTRAINT languages_fkey FOREIGN KEY (lang_id) REFERENCES languages(code3) NOT VALID;

SELECT DISTINCT * INTO tmp FROM catalogs_container_description_patterns; DROP TABLE catalogs_container_description_patterns; SELECT * INTO catalogs_container_description_patterns FROM tmp; DROP TABLE tmp;
ALTER TABLE catalogs_container_description_patterns DROP CONSTRAINT IF EXISTS catalog_container_description_pattern_pkey;
ALTER TABLE catalogs_container_description_patterns ADD CONSTRAINT catalog_container_description_pattern_pkey PRIMARY KEY (catalog_id, container_description_pattern_id);
ALTER TABLE catalogs_container_description_patterns DROP CONSTRAINT IF EXISTS catalogs_fkey;
ALTER TABLE catalogs_container_description_patterns ADD CONSTRAINT catalogs_fkey FOREIGN KEY (catalog_id) REFERENCES catalogs(id) NOT VALID;
ALTER TABLE catalogs_container_description_patterns DROP CONSTRAINT IF EXISTS container_description_patterns_fkey;
ALTER TABLE catalogs_container_description_patterns ADD CONSTRAINT container_description_patterns_fkey FOREIGN KEY (container_description_pattern_id) REFERENCES container_description_patterns(id) NOT VALID;

DO $$
	BEGIN
		ALTER TABLE containers RENAME COLUMN lang TO lang_id;
		EXCEPTION WHEN OTHERS THEN RAISE NOTICE '%', 'Already Exists';
	END
$$ LANGUAGE plpgsql;
ALTER TABLE containers DROP CONSTRAINT IF EXISTS languages_fkey;
ALTER TABLE containers ADD CONSTRAINT languages_fkey FOREIGN KEY (lang_id) REFERENCES languages(code3) NOT VALID;
ALTER TABLE containers DROP CONSTRAINT IF EXISTS content_types_fkey;
ALTER TABLE containers ADD CONSTRAINT content_types_fkey FOREIGN KEY (content_type_id) REFERENCES content_types(id) NOT VALID;
ALTER TABLE containers DROP CONSTRAINT IF EXISTS virtual_lessons_fkey;
ALTER TABLE containers ADD CONSTRAINT virtual_lessons_fkey FOREIGN KEY (virtual_lesson_id) REFERENCES virtual_lessons(id) NOT VALID;

DO $$
	BEGIN
		ALTER TABLE container_descriptions RENAME COLUMN lang TO lang_id;
		EXCEPTION WHEN OTHERS THEN RAISE NOTICE '%', 'Already Exists';
	END
$$ LANGUAGE plpgsql;
ALTER TABLE container_descriptions DROP CONSTRAINT IF EXISTS languages_fkey;
ALTER TABLE container_descriptions ADD CONSTRAINT languages_fkey FOREIGN KEY (lang_id) REFERENCES languages(code3) NOT VALID;
ALTER TABLE container_descriptions DROP CONSTRAINT IF EXISTS container_descriptions_fkey;
ALTER TABLE container_descriptions ADD CONSTRAINT container_descriptions_fkey FOREIGN KEY (container_id) REFERENCES containers(id) NOT VALID;

ALTER TABLE containers_file_assets DROP CONSTRAINT IF EXISTS lessonfiles_pkey;
ALTER TABLE containers_file_assets DROP CONSTRAINT IF EXISTS containers_file_assets_pkey;
ALTER TABLE containers_file_assets ADD CONSTRAINT containers_file_assets_pkey PRIMARY KEY (container_id, file_asset_id);
ALTER TABLE containers_file_assets DROP CONSTRAINT IF EXISTS containers_fkey;
ALTER TABLE containers_file_assets ADD CONSTRAINT containers_fkey FOREIGN KEY (container_id) REFERENCES containers(id) NOT VALID;
ALTER TABLE containers_file_assets DROP CONSTRAINT IF EXISTS file_assets_fkey;
ALTER TABLE containers_file_assets ADD CONSTRAINT file_assets_fkey FOREIGN KEY (file_asset_id) REFERENCES file_assets(id) NOT VALID;

ALTER TABLE containers_labels DROP CONSTRAINT IF EXISTS container_label_pkey;
ALTER TABLE containers_labels ADD CONSTRAINT container_label_pkey PRIMARY KEY (label_id, container_id);
ALTER TABLE containers_labels DROP CONSTRAINT IF EXISTS containers_fkey;
ALTER TABLE containers_labels ADD CONSTRAINT containers_fkey FOREIGN KEY (container_id) REFERENCES containers(id) NOT VALID;
ALTER TABLE containers_labels DROP CONSTRAINT IF EXISTS labels_fkey;
ALTER TABLE containers_labels ADD CONSTRAINT labels_fkey FOREIGN KEY (label_id) REFERENCES labels(id) NOT VALID;

ALTER TABLE file_asset_descriptions DROP CONSTRAINT IF EXISTS file_asset_descriptions_fkey;
ALTER TABLE file_asset_descriptions ADD CONSTRAINT file_asset_descriptions_fkey FOREIGN KEY (file_id) REFERENCES file_assets(id) NOT VALID;

DO $$
	BEGIN
		ALTER TABLE file_assets RENAME COLUMN lang TO lang_id;
		EXCEPTION WHEN OTHERS THEN RAISE NOTICE '%', 'Already Exists';
	END
$$ LANGUAGE plpgsql;
DO $$
	BEGIN
		ALTER TABLE file_assets RENAME COLUMN servername TO servername_id;
		EXCEPTION WHEN OTHERS THEN RAISE NOTICE '%', 'Already Exists';
	END
$$ LANGUAGE plpgsql;
ALTER TABLE file_assets DROP CONSTRAINT IF EXISTS languages_fkey;
ALTER TABLE file_assets ADD CONSTRAINT languages_fkey FOREIGN KEY (lang_id) REFERENCES languages(code3) NOT VALID;
ALTER TABLE file_assets DROP CONSTRAINT IF EXISTS users_fkey;
ALTER TABLE file_assets ADD CONSTRAINT users_fkey FOREIGN KEY (user_id) REFERENCES users(id) NOT VALID;
ALTER TABLE file_assets DROP CONSTRAINT IF EXISTS servers_fkey;
ALTER TABLE file_assets ADD CONSTRAINT servers_fkey FOREIGN KEY (servername_id) REFERENCES servers(servername) NOT VALID;

SELECT DISTINCT * INTO ru FROM roles_users; DROP TABLE roles_users; SELECT * INTO roles_users FROM ru; DROP TABLE ru;
ALTER TABLE roles_users DROP CONSTRAINT IF EXISTS role_user_pkey;
ALTER TABLE roles_users ADD CONSTRAINT role_user_pkey PRIMARY KEY (role_id, user_id);
ALTER TABLE roles_users DROP CONSTRAINT IF EXISTS roles_fkey;
ALTER TABLE roles_users ADD CONSTRAINT users_fkey FOREIGN KEY (user_id) REFERENCES users(id) NOT VALID;
	`)
	if err != nil {
		panic(err)
	}
}
