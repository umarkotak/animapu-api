BEGIN;

LOCK TABLE users, animes, mangas, manga_chapters, anime_histories, manga_histories, anime_libraries, manga_libraries IN ACCESS EXCLUSIVE MODE;

SELECT setval(pg_get_serial_sequence('users', 'id'), GREATEST(COALESCE(MAX(id), 1), (SELECT last_value FROM users_id_seq)), COUNT(*) > 0 OR (SELECT is_called FROM users_id_seq)) FROM users;
SELECT setval(pg_get_serial_sequence('animes', 'id'), GREATEST(COALESCE(MAX(id), 1), (SELECT last_value FROM animes_id_seq)), COUNT(*) > 0 OR (SELECT is_called FROM animes_id_seq)) FROM animes;
SELECT setval(pg_get_serial_sequence('mangas', 'id'), GREATEST(COALESCE(MAX(id), 1), (SELECT last_value FROM mangas_id_seq)), COUNT(*) > 0 OR (SELECT is_called FROM mangas_id_seq)) FROM mangas;
SELECT setval(pg_get_serial_sequence('manga_chapters', 'id'), GREATEST(COALESCE(MAX(id), 1), (SELECT last_value FROM manga_chapters_id_seq)), COUNT(*) > 0 OR (SELECT is_called FROM manga_chapters_id_seq)) FROM manga_chapters;
SELECT setval(pg_get_serial_sequence('anime_histories', 'id'), GREATEST(COALESCE(MAX(id), 1), (SELECT last_value FROM anime_histories_id_seq)), COUNT(*) > 0 OR (SELECT is_called FROM anime_histories_id_seq)) FROM anime_histories;
SELECT setval(pg_get_serial_sequence('manga_histories', 'id'), GREATEST(COALESCE(MAX(id), 1), (SELECT last_value FROM manga_histories_id_seq)), COUNT(*) > 0 OR (SELECT is_called FROM manga_histories_id_seq)) FROM manga_histories;
SELECT setval(pg_get_serial_sequence('anime_libraries', 'id'), GREATEST(COALESCE(MAX(id), 1), (SELECT last_value FROM anime_libraries_id_seq)), COUNT(*) > 0 OR (SELECT is_called FROM anime_libraries_id_seq)) FROM anime_libraries;
SELECT setval(pg_get_serial_sequence('manga_libraries', 'id'), GREATEST(COALESCE(MAX(id), 1), (SELECT last_value FROM manga_libraries_id_seq)), COUNT(*) > 0 OR (SELECT is_called FROM manga_libraries_id_seq)) FROM manga_libraries;

COMMIT;
