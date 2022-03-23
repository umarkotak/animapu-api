# ANIMAPU-API
Open source API to fetch latest manga from multiple sources. CORS friendly.

## AVAILABLE SOURCES
| Source Name             | Source ID      | Status      |
| ----------------------- | -------------- | ----------- |
| Mangaupdates X Mangahub | mangaupdates   | |
| Klik Manga              | klikmanga      | |
| Mangabat                | mangabat       | |
| Mangadex                | mangadex       | |
| Maid My                 | maidmy         | |

## USAGES

| Name         | Method      | PATH                                                                         |
| -----------  | ----------- | ---------------------------------------------------------------------------- |
| Latest Manga | GET         | {HOST}/mangas/:manga_source/latest?page=                                     |
| Search Manga | GET         | {HOST}/mangas/:manga_source/search?title=                                    |
| Detail Manga | GET         | {HOST}/mangas/:manga_source/detail/:manga_id?secondary_source_id=            |
| Read Manga   | GET         | {HOST}/mangas/:manga_source/read/:manga_id/:chapter_id? secondary_source_id= |
