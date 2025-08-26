Roadmap Simple Note App (Next.js + Golang)
Phase 1 — Fondasi (Backend + Frontend Dasar)

🎯 Tujuan: Aplikasi CRUD berjalan dengan login sederhana.

Backend (Golang)

    Setup project Go dengan struktur Clean Architecture:

    /cmd
    /internal/controller
    /internal/service
    /internal/repository
    /internal/model

    Setup database (PostgreSQL) + migration tool (golang-migrate).

    Buat endpoint:

        Auth: register, login (JWT)

        Notes CRUD: create, read, update, delete

    Middleware:

        JWT validation

        Logging request

        Error handler

Frontend (Next.js)

    Setup Next.js project (App Router).

    Halaman:

        Login / Register

        Dashboard (list notes)

        Create / Edit Note

    Integrasi API backend dengan fetch/axios.

    Responsive layout (mobile & desktop).

Phase 2 — Menambah Nilai (Fitur Menengah)

🎯 Tujuan: Aplikasi mulai terasa “nyata” & enak dipakai.

Backend

    Tambahkan kolom tags, pinned, deletedAt di tabel notes.

    Endpoint tambahan:

        Search note

        Filter by tag

        Trash & restore note

    Pagination / limit-offset query.

    Unit test untuk service layer.

Frontend

    Search bar + filter tag.

    UI pin note.

    Trash & restore (folder khusus untuk catatan terhapus).

    Infinite scroll / pagination.

    Loading states & error handling di UI.

Phase 3 — Advanced Features

🎯 Tujuan: Menunjukkan kemampuan di real-time, offline, dan UX modern.

Backend

    Implement WebSocket untuk update note real-time.

    Endpoint upload file (gambar, PDF).

    Endpoint export note (PDF / Markdown).

    Simpan version history note.

    Swagger API documentation.

Frontend

    Integrasi WebSocket untuk update real-time.

    Markdown editor (react-markdown atau editor.js).

    Upload gambar/file di note.

    Export note ke PDF / Markdown.

    Dark mode & theme customizer.

Phase 4 — Production Ready

🎯 Tujuan: Aplikasi siap dilihat recruiter & HR.

    Security

        Validasi input di backend.

        CSRF & XSS protection.

        Rate limiting.

    Deployment

        Dockerfile untuk backend & frontend.

        Docker Compose untuk DB + backend + frontend.

        Deploy di VPS / Render / Railway.

    CI/CD

        GitHub Actions untuk auto build & deploy.

    Performance

        Indexing database.

        Caching (Redis) untuk list notes.

Urutan Pengerjaan Harian

Kalau kamu mau eksekusi cepat:

    Hari 1–3 → Phase 1 (Backend basic + Frontend basic)

    Hari 4–6 → Phase 2 (Search, filter, trash)

    Hari 7–10 → Phase 3 (Real-time, upload, markdown)

    Hari 11–12 → Phase 4 (Security, Docker, Deploy)
