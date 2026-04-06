# Business Requirement Document (BRD)

## Sistem E-Procurement

**Entitas:** PT Victoria Investama, Tbk

**Tanggal:** 5 April 2026

**Versi:** 2.1

**Status:** Draft untuk Persetujuan Manajemen

---

# 1. Pendahuluan

## 1.1 Latar Belakang

Dalam rangka memperkuat tata kelola perusahaan yang baik (Good Corporate Governance/GCG), meningkatkan transparansi, serta memperkuat sistem pengendalian internal, PT. Victoria Investama, Tbk berencana mengimplementasikan Sistem E-Procurement yang terintegrasi dan terdokumentasi secara sistematis. Proses pengadaan saat ini masih dilakukan secara manual atau semi-manual dengan ketergantungan pada _email_ dan dokumen terpisah. Kondisi tersebut menimbulkan beberapa risiko dan keterbatasan, antara lain:

- Tidak adanya sistem terintegrasi _end-to-end_.
- _Monitoring approval_ yang tidak _real-time_.
- Dokumentasi perbandingan vendor yang tidak terstruktur.
- Risiko _conflict of interest_.
- Keterbatasan _audit trail_ sistematis.
- Keterlambatan proses akibat _bottleneck approval_.

Untuk itu diperlukan sistem yang mampu mengintegrasikan seluruh tahapan pengadaan secara digital dan terdokumentasi.

## 1.2 Tujuan Dokumen

Dokumen ini disusun untuk:

- Mendefinisikan kebutuhan bisnis atas Sistem E-Procurement.
- Menjadi dasar persetujuan manajemen dan Direksi.
- Menjadi acuan pengembangan sistem oleh Divisi IT.
- Mendukung kebutuhan audit internal dan eksternal.

## 1.3 Catatan Konsistensi Dokumen

Untuk menjaga konsistensi antar dokumen proyek:

- BRD menjadi acuan utama kebutuhan bisnis, governance, dan target proses procurement.
- FSD menjadi turunan kebutuhan fungsional, use case, lifecycle, validasi, dan aturan operasional sistem.
- TSD menjadi turunan rancangan teknis implementasi, termasuk model autentikasi, arsitektur layanan, integrasi, dan kontrol keamanan.
- Jika pada FSD terdapat lampiran alignment implementasi backend berjalan, lampiran tersebut diperlakukan sebagai konteks implementasi saat ini dan tidak mengubah kebutuhan bisnis target-state yang didefinisikan dalam BRD.

## 1.4 Prinsip Terminologi Dokumen

Untuk mengurangi ambiguitas istilah lintas dokumen proyek:

- Istilah status budget yang digunakan secara baku adalah **Within Budget**, **Over Budget**, dan **Non Budget**.
- Istilah metode pengadaan langsung yang digunakan secara baku adalah **Direct Appointment (DA)**.
- Istilah harga referensi yang digunakan secara baku adalah **Reference Price / eCatalog**.
- Istilah _session_ pada level bisnis/fungsional dipahami sebagai **sesi autentikasi**, sedangkan detail implementasi token dijelaskan pada TSD.
- Jika sebuah dokumen membedakan antara kondisi **implemented**, **partially implemented**, atau **planned**, penanda tersebut hanya menjelaskan status implementasi saat ini dan tidak mengubah target-state bisnis yang didefinisikan pada BRD.

# 2. Tujuan Proyek

Pengembangan Sistem E-Procurement bertujuan untuk:

- Mendigitalisasi proses pengadaan secara menyeluruh.
- Mengimplementasikan _workflow approval_ berbasis struktur organisasi, limit nominal, dan parameter kebijakan dinamis.
- Meningkatkan transparansi proses RFQ dan _bidding_ vendor.
- Memperkuat segregasi tugas (_Segregation of Duties_).
- Menyediakan _audit trail_ lengkap dan tidak dapat dimodifikasi.
- Mengurangi _lead time_ proses pengadaan.
- Memberikan _dashboard monitoring_ kepada manajemen dan Direksi.
- Mengimplementasikan tata kelola pengadaan terpusat (_Multi-Entity Procurement Governance_) di bawah _Holding Company_.

# 3. Multi-Entity Procurement Governance

Sistem E-Procurement didesain untuk menjadi _Group Procurement Platform_ yang akan digunakan oleh seluruh entitas dalam Victoria Financial Group. Konsep ini memerlukan arsitektur sistem yang mendukung _multi-entity_ dan tata kelola terpusat.

## 3.1 Struktur Entitas

Sistem harus mendukung hierarki entitas sebagai berikut:

### Holding Company

- PT Victoria Investama Tbk

### Subsidiaries

- PT Victoria Insurance
- PT Bank Victoria International
- PT Victoria Sekuritas
- PT Victoria Alife Indonesia
- PT Victoria Manajemen Investasi

### Entity Hierarchy

Victoria Investama (Holding)

├── Victoria Insurance

├── Bank Victoria

├── Victoria Sekuritas

├── Victoria Alife Indonesia

└── Victoria Manajemen Investasi

## 3.2 Prinsip Tata Kelola Pengadaan

- **Dukungan Arsitektur Multi-Entity:** Sistem wajib memiliki _multi-entity architecture_ yang memungkinkan pemisahan data, _user_, _approval matrix_, dan _budget_ per entitas secara ketat.
- **Otonomi Pengadaan Entitas:** Setiap entitas memiliki fungsi pengadaan dan _workflow_ yang dapat dikelola secara mandiri oleh _Entity Admin_, namun tetap terikat pada ketentuan Holding.
- **Oversight Holding:** _Holding Company_ (PT Victoria Investama Tbk) memiliki _governance oversight_ melalui peran **Holding Admin** untuk memastikan keseragaman kebijakan dan kepatuhan pengadaan di seluruh Grup.

### Kewenangan Governance Approval

Holding Admin memiliki kewenangan tertinggi untuk menentukan model _governance approval_ bagi masing-masing entitas, yang dapat mencakup:

- _Approval_ cukup sampai level manajemen entitas.
- _Approval_ wajib eskalasi ke holding, terlepas dari nilai pengadaan.
- _Approval_ eskalasi ke holding berdasarkan kondisi atau parameter tertentu (misalnya _Over Budget_ atau nilai pengadaan melebihi batas Entitas).

## 3.3 Prinsip Data Isolation dan Cross-Entity Visibility

Untuk menjaga kerahasiaan dan kepatuhan per entitas, sistem harus menerapkan prinsip isolasi data:

- **User Entitas:** Hanya dapat melihat data pengadaan entitasnya sendiri.
- **Entity Admin:** Hanya dapat mengakses dan mengelola data pengadaan entitas yang menjadi kewenangannya.
- **Holding Admin:** Dapat melihat seluruh aktivitas dan data pengadaan lintas entitas dalam Grup.
- **Management Level Grup:** Dapat memiliki akses _monitoring_ lintas entitas sesuai kewenangan yang diberikan (_role-based access control_).
- **Internal Audit Grup:** Dapat melakukan _review_ dan akses data lintas entitas untuk tujuan audit.

Seluruh akses lintas entitas dikontrol melalui _role-based access control_ yang dikelola oleh Holding Admin.

# 4. Ruang Lingkup Proyek

## 4.1 Ruang Lingkup (In Scope)

### A. Modul Purchase Request (PR)

- Input permintaan barang/jasa oleh user.
- Upload dokumen pendukung.
- Penentuan workflow approval.
- Tracking status permintaan.
- Fitur revisi dan resubmit jika ditolak.

### B. Modul Budget Management

Sistem harus mampu mengelola dan memvalidasi anggaran pengadaan sebelum Purchase Request disetujui.

- Konfigurasi budget per entitas.
- Mode Limited Budget dan Unlimited Budget.
- Validasi budget terhadap permintaan pengadaan.
- Status Over Budget.
- Status Non Budget.
- Escalation approval berdasarkan governance budget.
- Budget dapat dikelola oleh Holding Admin dan Entity Admin sesuai kewenangan.
- Budget dapat diatur berdasarkan entitas, departemen, kategori pengadaan, dan periode.

### C. Modul Kebijakan Pengadaan Dinamis / Dynamic Procurement Policy

Sistem harus mendukung konfigurasi Dynamic Procurement Policy untuk menentukan alur proses berdasarkan parameter yang dapat diubah oleh Holding Admin atau Entity Admin tanpa perubahan sistem.

- Parameter nilai pengadaan.
- Status Budget (Within Budget / Over Budget / Non Budget).
- Kategori Pengadaan (Rutin / Non Rutin).
- Jenis Pengadaan (Barang / Jasa).
- Penentuan workflow.
- Penentuan metode procurement.
- Escalation berdasarkan parameter.

### D. Modul RFQ dan Bidding

- Generate dokumen RFQ otomatis.
- Pengiriman RFQ kepada vendor.
- Input dan pengelolaan quotation vendor.
- Pengaturan batas waktu bidding.
- Pembukaan ulang bidding (jika diperlukan).
- Penutupan bidding.
- Publikasi tender ke Vendor Portal oleh Procurement untuk vendor yang eligible.

### E. Modul Perbandingan dan Evaluasi Vendor

Sistem harus mendukung proses evaluasi yang lebih komprehensif dan transparan.

- Perbandingan harga.
- Technical evaluation.
- Commercial evaluation.
- Weighted scoring.
- Ranking vendor.
- Summary report hasil evaluasi.
- Dokumentasi alasan pemilihan vendor.
- BAFO jika diperlukan.
- Reference Price / eCatalog sebagai referensi kewajaran harga.

### F. Modul Purchase Order (PO)

- Generate PO berdasarkan vendor terpilih.
- Workflow approval PO.
- Revisi PO.
- Pengiriman PO ke vendor.
- Konfirmasi vendor.
- Dokumentasi status PO.
- Sinkron dengan governance approval entity/holding bila dipersyaratkan.

### G. Modul Penunjukan Langsung

- Pemilihan metode pengadaan oleh Procurement.
- Penunjukan langsung merupakan pengecualian terhadap minimum vendor bidding dan wajib disertai justifikasi yang terdokumentasi.
- Vendor yang ditunjuk tetap harus melalui pengecekan kelayakan administratif, approved vendor status, dan blacklist check.
- Upload quotation / price list / kontrak sebelumnya / referensi harga lain.
- Dokumentasi untuk audit.

### H. Reporting dan Dashboard

- Monitoring status PR, RFQ, dan PO.
- Analisis lead time end-to-end procurement.
- Rekap pengadaan per periode.
- Monitoring nilai pengadaan per kategori.
- Monitoring proporsi bidding vs Direct Appointment.
- Dashboard per entitas dan dashboard group sesuai kewenangan.
- Monitoring budget usage dan pengadaan Over Budget.

### I. Modul Entity Management

- Pembuatan entitas baru oleh Holding Admin.
- Pengelolaan struktur dan hierarki entitas.
- Pengaturan status entitas aktif/nonaktif.
- Pengaitan governance setting per entitas.
- Pengaitan approval model per entitas.
- Pengaitan budget governance per entitas.

### J. Modul User Management dan User Role Assignment

- Pembuatan user.
- Pengaitan user ke entitas tertentu.
- Penetapan role user.
- Aktivasi/deaktivasi user.
- User role assignment.
- Reset password oleh Holding Admin atau Entity Admin sesuai kewenangan.
- Delegate approver sementara dengan periode berlaku tertentu.
- Holding Admin dapat mengelola user lintas entitas.
- Entity Admin hanya dapat mengelola user di entitasnya sendiri sesuai governance.

### K. Modul Vendor Blacklist dan Vendor Eligibility Control

- Holding Admin dapat menetapkan vendor sebagai blacklist atau menghapus blacklist sesuai kebijakan yang berlaku.
- Status blacklist wajib mempengaruhi eligibility vendor pada proses tender dan pemilihan vendor.
- Blacklist check menjadi kontrol wajib pada vendor prequalification, RFQ, Direct Appointment, dan pembuatan PO.
- Alasan blacklist dan unblacklist harus terdokumentasi untuk kebutuhan audit.

### L. Modul Vendor Portal / Vendor Participation Portal

- Portal eksternal untuk partisipasi vendor dalam tender.
- Vendor dapat melihat tender yang dipublikasikan.
- Vendor dapat melihat requirement tender.
- Vendor dapat submit quotation / penawaran.
- Vendor dapat upload dokumen pendukung.
- Vendor dapat melihat status partisipasi tender.
- Vendor dapat melakukan konfirmasi PO apabila mekanisme tersebut ditetapkan melalui portal.
- Partisipasi vendor tetap tunduk pada approved vendor status, blacklist check, dan mekanisme eligibility yang berlaku.
- Vendor portal bukan full vendor onboarding system.

## 4.2 Di Luar Ruang Lingkup (Out of Scope)

Hal-hal berikut secara eksplisit berada di luar ruang lingkup proyek sistem E-Procurement dan harus ditangani oleh sistem eksternal atau proses manual yang terpisah:

- Proses pembayaran vendor.
- Integrasi langsung dengan sistem perbankan.
- Manajemen kontrak jangka panjang.
- Vendor onboarding system penuh.
- Integrasi ERP penuh.

# 5. Gambaran Proses Bisnis

## 5.1 Proses Purchase Request

- Requestor membuat PR.
- Requestor mengunggah dokumen pendukung.
- Sistem menjalankan workflow approval sesuai governance yang berlaku pada entitas terkait.
- Approval dapat diselesaikan di level entitas atau diekskalasikan ke Holding Approver sesuai setting Holding Admin.
- Jika ditolak, Requestor melakukan revisi dan resubmit.
- Jika disetujui, proses dilanjutkan ke tahap penentuan metode pengadaan oleh Procurement.

## 5.2 Dynamic Approval Workflow

- Workflow approval bersifat dinamis.
- Ditentukan berdasarkan parameter:
  - Nilai pengadaan.
  - Within Budget / Over Budget / Non Budget.
  - Jenis pengadaan.
  - Kategori pengadaan.
  - governance per entitas.
- Holding Admin menentukan apakah approval suatu entitas:
  - Cukup di level entitas.
  - Wajib ke holding.
  - Conditional berdasarkan parameter.

## 5.3 Penentuan Metode Pengadaan

- Procurement menentukan metode pengadaan.
- Metode yang tersedia:
  - RFQ / Bidding.
  - Penunjukan Langsung.
- Pemilihan metode harus terdokumentasi dalam sistem.
- Jika metode RFQ/Bidding dipilih, proses dilanjutkan ke RFQ dan Vendor Portal.
- Jika metode Penunjukan Langsung dipilih, proses dilanjutkan ke pemilihan vendor langsung.

## 5.4 Proses RFQ dan Bidding

- Sistem menghasilkan dokumen RFQ.
- Procurement mempublikasikan tender ke Vendor Portal.
- Hanya vendor yang eligible yang dapat mengikuti tender.
- Vendor mengirim quotation dan dokumen pendukung.
- Procurement menerima dan mengelola quotation vendor.
- Procurement menetapkan deadline bidding.
- Jika jumlah quotation belum memenuhi kebijakan, bidding dapat dibuka ulang.
- Jika telah memenuhi, bidding ditutup.

## 5.5 Vendor Portal Participation Flow

- Procurement mempublikasikan tender ke Vendor Portal setelah PR approved dan metode RFQ/Bidding dipilih.
- Vendor yang memenuhi eligibility dapat melihat tender.
- Vendor dapat membaca requirement tender.
- Vendor mengisi penawaran / quotation.
- Vendor mengunggah dokumen pendukung.
- Quotation vendor diterima kembali ke proses evaluasi internal.
- Vendor dapat melihat status partisipasi tender sesuai mekanisme yang ditetapkan.
- Vendor dapat melakukan konfirmasi PO melalui portal apabila kebijakan tersebut digunakan.

## 5.6 Proses Vendor Prequalification dan Evaluasi Vendor

### Tahap 1 - Vendor Prequalification

- Approved vendor status.
- Blacklist check.
- Kelayakan administratif.
- Eligibility.

### Tahap 2 - Technical Evaluation

- Technical capability.
- Experience.
- Compliance.
- Delivery capability.

### Tahap 3 - Commercial Evaluation

- Harga.
- Term of payment.
- Delivery terms.
- Commercial terms lain yang relevan.

### Tahap 4 - Weighted Scoring

- Technical score.
- Commercial score.
- Bobot evaluasi.

### Tahap 5 - BAFO

- BAFO dilakukan bila diperlukan kepada vendor yang memenuhi syarat tahap sebelumnya.

### Tahap 6 - Vendor Selection

- Vendor dipilih dan didokumentasikan alasannya.

## 5.7 Proses Penunjukan Langsung

- Procurement memilih vendor yang akan ditunjuk.
- Penunjukan langsung merupakan pengecualian terhadap minimum vendor bidding.
- Procurement wajib mengisi justifikasi penunjukan langsung.
- Vendor yang ditunjuk tetap harus memenuhi approved vendor status, blacklist check, dan kelayakan minimum governance.
- Procurement mengunggah quotation / price list / kontrak sebelumnya / referensi harga lain.
- Sistem menyimpan dokumentasi pengadaan untuk audit.
- Procurement membuat Purchase Order berdasarkan vendor yang ditunjuk.

## 5.8 Proses Purchase Order

- Procurement membuat PO berdasarkan vendor terpilih.
- PO diajukan untuk approval sesuai governance approval yang berlaku.
- Approval PO dapat berhenti di level entitas atau diekskalasikan ke Holding Approver sesuai governance rule.
- Jika memerlukan revisi, kembali ke Procurement.
- Jika disetujui, PO dikirim ke vendor.

## 5.9 Proses Vendor Confirmation

- Vendor memberikan konfirmasi PO melalui Vendor Portal atau mekanisme resmi lain yang ditetapkan perusahaan.
- Sistem mencatat status konfirmasi vendor.
- Proses procurement dinyatakan selesai setelah tahapan vendor confirmation terpenuhi sesuai kebijakan yang berlaku.

# 6. Struktur Peran dan Tanggung Jawab

Peran dalam sistem E-Procurement dibagi menjadi Peran Administrasi, Peran Bisnis, dan Pihak Eksternal, dengan level kewenangan yang dapat berada pada level entitas maupun group sesuai governance yang berlaku.

| **Peran**        | **Level**     | **Tanggung Jawab Utama**                                                                                                                                                                                                                                                                                                                                                     |
| ---------------- | ------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Holding Admin    | Group         | Membuat entitas baru, mengelola seluruh user di semua entitas, menentukan governance rule (kebijakan grup), menentukan model governance approval per entitas, menentukan approval escalation lintas entitas, dan melihat seluruh aktivitas procurement grup.                                                                                                                 |
| Entity Admin     | Entitas       | Membuat user dalam entitasnya, mengelola workflow procurement (sesuai governance rule), mengatur budget entitas (sesuai governance Holding), dan monitoring procurement di entitasnya.                                                                                                                                                                                       |
| Requestor        | Entitas       | Membuat PR dan melengkapi dokumen.                                                                                                                                                                                                                                                                                                                                           |
| Entity Approver  | Entitas       | Pihak yang menyetujui atau menolak PR/PO sesuai limit entitas.                                                                                                                                                                                                                                                                                                               |
| Holding Approver | Group         | Pihak approver di level Holding, yang umumnya berada pada level direktur atau pejabat berwenang di Holding, dan hanya terlibat jika governance rule untuk entitas tersebut mewajibkan eskalasi ke holding.                                                                                                                                                                   |
| Procurement      | Entitas       | Mengelola RFQ, bidding, comparison, BAFO, dan PO.                                                                                                                                                                                                                                                                                                                            |
| Finance          | Entitas/Group | Memberikan approval tambahan untuk kondisi Over Budget atau Non Budget sesuai governance yang berlaku.                                                                                                                                                                                                                                                                      |
| Management       | Entitas/Group | Monitoring dan evaluasi kinerja pengadaan.                                                                                                                                                                                                                                                                                                                                   |
| Internal Audit   | Entitas/Group | Review dan pengujian kepatuhan sistem dan proses lintas entitas.                                                                                                                                                                                                                                                                                                             |
| Vendor           | Eksternal     | Pihak eksternal / rekanan yang berpartisipasi dalam tender melalui Vendor Portal atau mekanisme resmi lain. Vendor dapat mengirim quotation, mengunggah dokumen pendukung, dan memberikan konfirmasi PO sesuai kebijakan yang berlaku. Vendor bukan role administratif internal dan tidak memiliki akses ke governance, budget, workflow internal, atau data lintas entitas. |

_Catatan: Sebagian peran bisnis melekat pada level entitas, sementara peran tertentu seperti Holding Approver, serta fungsi monitoring dan audit tertentu, dapat berada pada level group sesuai governance yang ditetapkan oleh Holding Admin._

# 7. Business Success Metrics dan KPI

Keberhasilan proyek tidak hanya diukur dari ketersediaan aplikasi, tetapi juga dari dampak bisnis setelah implementasi.

| **KPI / Outcome**                    | **Baseline Awal**                | **Target Tahun 1**                                | **Owner Bisnis**         |
| ------------------------------------ | -------------------------------- | ------------------------------------------------- | ------------------------ |
| Lead time approval PR/PO             | Manual, tidak terukur konsisten  | Penurunan lead time minimum 30%                   | Procurement + Management |
| SLA approval compliance              | Belum terukur sistematis         | >= 90% approval selesai dalam SLA                 | Entity Admin             |
| Ketersediaan audit trail             | Parsial, berbasis email dan file | 100% transaksi utama memiliki jejak audit         | Internal Audit + IT      |
| Transparansi evaluasi vendor         | Tidak seragam                    | 100% vendor selection memiliki evidence evaluasi  | Procurement              |
| Pengendalian budget                  | Manual / reaktif                 | 100% PR memiliki status budget terdokumentasi     | Entity Admin + Finance   |
| Adopsi penggunaan sistem             | Belum ada                        | >= 95% proses in-scope berjalan melalui sistem    | Management + IT          |
| Monitoring lintas entitas            | Manual                           | Dashboard group-level tersedia dan dipakai rutin  | Holding Admin            |

# 8. Prioritas Kebutuhan dan Rencana Implementasi Bertahap

| **Prioritas** | **Makna** | **Contoh Cakupan** |
| ------------- | --------- | ------------------ |
| Must Have | Wajib tersedia agar sistem layak _go-live_ | PR, approval workflow, budget validation, RFQ, quotation, vendor evaluation, PO, audit trail, role access |
| Should Have | Sangat penting, tetapi masih dapat menyusul setelah core stabil | Delegate approver, vendor blacklist, Reference Price / eCatalog, print/export, reminder SLA |
| Could Have | Nilai tambah operasional dan monitoring | Dashboard lanjutan, report agregat, Reference Price yang dihasilkan otomatis, BAFO refinement |
| Future Phase | Direncanakan setelah fase inti berhasil | ERP integration, MFA, SSO, mobile approval, vendor performance scoring historis |

## Rencana Implementasi Bertahap

| **Fase** | **Fokus** | **Output Utama** |
| -------- | --------- | ---------------- |
| Phase 1 | Procurement transaksional inti | Login, PR, approval task, RFQ dasar, quotation, PO, vendor confirmation, user/entity management |
| Phase 2 | Penguatan governance | Budget management penuh, dynamic procurement policy, dynamic approval workflow, vendor blacklist, delegate approver |
| Phase 3 | Monitoring dan optimisasi | Dashboard lanjutan, export/report, pengayaan Reference Price, advanced notification dan escalation |
| Phase 4 | Integrasi enterprise | ERP integration, SSO/MFA, data warehouse / BI, future automation |

# 9. Stakeholder dan RACI Bisnis

| **Area / Keputusan** | **Responsible** | **Accountable** | **Consulted** | **Informed** |
| -------------------- | --------------- | --------------- | ------------- | ------------ |
| Kebijakan procurement group | Holding Admin | Direksi / Manajemen | Procurement, Internal Audit | Seluruh entitas |
| Governance approval per entitas | Holding Admin, Entity Admin | Direksi / Manajemen | Finance, Procurement | Approver terkait |
| Budget governance | Entity Admin, Finance | Direksi / Manajemen | Holding Admin, Procurement | Requestor, Approver |
| Operasional procurement harian | Procurement | Management entitas | Requestor, Approver | Internal Audit |
| Vendor evaluation dan selection | Procurement | Management entitas sesuai kewenangan | User teknis / requestor | Internal Audit |
| User access dan SoD | Holding Admin, Entity Admin | Management / Governance owner | Internal Audit, IT | User terkait |
| Audit review dan compliance | Internal Audit | Direksi / Komite terkait | Holding Admin, IT, Procurement | Management |
| Persetujuan _go-live_ | Divisi IT, Procurement | Direksi / Manajemen | Internal Audit, stakeholder entitas | Seluruh stakeholder proyek |

# 10. Kriteria Penerimaan Bisnis dan UAT Level

| **Area UAT** | **Kriteria Penerimaan Bisnis** |
| ------------ | ------------------------------ |
| Login dan role access | User internal dan vendor hanya dapat mengakses menu dan data sesuai role serta scope entitasnya |
| Purchase Request | Requestor dapat membuat, submit, revisi, dan memonitor PR dengan dokumen pendukung lengkap |
| Budget validation | Setiap PR memiliki hasil status Within Budget / Over Budget / Non Budget yang mempengaruhi alur approval |
| Approval workflow | Approval berjalan berurutan, tidak dapat dilewati, dan sesuai governance entitas / holding |
| Procurement method | Sistem dapat membedakan alur RFQ/Bidding vs Direct Appointment sesuai kebijakan |
| RFQ dan quotation | Vendor eligible dapat menerima tender, melihat detail, dan submit quotation sebelum deadline |
| Vendor evaluation | Procurement dapat melakukan prequalification, evaluasi teknis/komersial, dan mendokumentasikan alasan pemilihan |
| Purchase Order | PO hanya dapat diterbitkan dari proses yang sah dan mengikuti approval yang berlaku |
| Vendor confirmation | Vendor dapat mengkonfirmasi PO dan status transaksi berubah sesuai lifecycle yang didefinisikan |
| Audit trail dan report | Aktivitas utama tercatat, dapat difilter, dan dapat diekspor untuk kebutuhan audit |

# 11. Asumsi, Ketergantungan, dan Constraint Bisnis

- Data master entitas, user, vendor, kategori, dan struktur organisasi tersedia atau dapat disiapkan sebelum UAT.
- Kebijakan approval, budget, dan procurement policy per entitas disepakati sebelum konfigurasi final di sistem.
- Proses pembayaran vendor, invoice verification, dan penerimaan barang/jasa tetap berada di luar ruang lingkup fase ini.
- Notifikasi email bergantung pada ketersediaan layanan SMTP dan kebijakan keamanan infrastruktur perusahaan.
- Keberhasilan adopsi sistem memerlukan sosialisasi, pelatihan user, dan dukungan change management dari manajemen.
- Jika terdapat kebutuhan lintas entitas yang spesifik, keputusan governance tetap ditetapkan oleh Holding Admin dan Direksi sesuai kewenangan.

# 12. Open Issues / Keputusan Bisnis yang Masih Perlu Ditetapkan

| **Topik** | **Dampak** | **Owner Keputusan** |
| --------- | ---------- | ------------------- |
| Batas nilai final per level approver pada tiap entitas | Mempengaruhi konfigurasi approval matrix dan UAT | Holding Admin + Management |
| Apakah escalation SLA bersifat pengingat saja atau dapat auto-escalate pada fase berikutnya | Mempengaruhi governance operasional dan ekspektasi user | Management + Internal Audit |
| Kebijakan final penggunaan Direct Appointment per kategori pengadaan | Mempengaruhi dynamic procurement policy dan kontrol fraud | Procurement + Management |
| Kebutuhan tanda tangan digital formal pada dokumen output | Mempengaruhi legal acceptance dan desain output report | Management + Legal / Compliance |
| Prioritas integrasi ERP / Finance pada fase berikutnya | Mempengaruhi roadmap dan desain transisi proses pasca-_go-live_ | Direksi + IT + Finance |

# 13. Change Management dan Adoption Readiness

Keberhasilan implementasi sangat dipengaruhi oleh kesiapan organisasi, bukan hanya kesiapan aplikasi.

| **Area** | **Kebutuhan Minimum** | **Owner** |
| -------- | --------------------- | --------- |
| Sosialisasi proses baru | Komunikasi resmi mengenai perubahan alur PR, approval, RFQ, dan PO | Management + Procurement |
| Pelatihan user internal | Training per role untuk Requestor, Approver, Procurement, Admin, dan Audit | IT + Procurement |
| Pelatihan vendor | Panduan penggunaan Vendor Portal dan tata cara submit quotation / konfirmasi PO | Procurement + Vendor Management |
| UAT dan sign-off | Daftar user per entitas untuk UAT, hasil uji, dan sign-off bisnis | Entity Admin + IT |
| Kesiapan transisi produksi | Penetapan tanggal _go-live_, periode freeze, dan dukungan pasca implementasi | IT + Manajemen |
| Masa stabilisasi pasca implementasi | Pendampingan insiden, monitoring issue log, dan evaluasi adopsi awal | IT + Procurement + Internal Audit |

# Lampiran A: Dynamic Approval Matrix

Matriks Approval ini menjadi dasar konfigurasi workflow sistem, yang kini bersifat dinamis dan dapat ditentukan berdasarkan kebijakan pengadaan. Workflow dapat dikonfigurasi melalui sistem tanpa perlu perubahan kode program.

**Penegasan Governance Approval per Entitas:** Approval Matrix ini merupakan kerangka acuan umum. Penerapan aktual pada masing-masing entitas dapat berbeda sesuai governance setting yang ditetapkan oleh Holding Admin, termasuk apakah approval cukup sampai level entitas, wajib eskalasi ke holding, atau bersifat conditional berdasarkan parameter tertentu.

## Contoh Konfigurasi Approval

| **Kondisi Pengadaan (Parameter)**  | **Workflow Approval**                                                      |
| ---------------------------------- | -------------------------------------------------------------------------- |
| ≤ Rp 50.000.000 dan Within Budget  | Entity Approver cukup di entitas (Head of Division)                        |
| \> Rp 50.000.000 - Rp 250.000.000  | Head of Division → Direktur Terkait                                        |
| \> Rp 250.000.000 - Rp 500.000.000 | Direktur → Direktur Utama                                                  |
| \> Rp 500.000.000                  | Escalation ke Holding Approver (Direktur Utama Holding + Komisaris)        |
| Pengadaan Over Budget              | Approval khusus Finance (contoh: CFO Approval) sebelum ke level management |
| Pengadaan Non Budget               | Approval khusus Finance dan escalation ke level Direksi                    |

## Ketentuan Tambahan

- Approval tidak dapat dilewati.
- Approval dilakukan berurutan.
- SLA maksimum 2 hari kerja per level.
- Sistem mengirim reminder otomatis jika melewati SLA.
- Penolakan wajib disertai alasan tertulis.

# Lampiran B: Risk and Fraud Control Matrix (RCM)

Matriks ini mencerminkan pengendalian yang diperkuat melalui Budget Management, Dynamic Procurement Policy, Vendor Prequalification, Weighted Scoring, BAFO, dan Vendor Portal untuk memitigasi risiko fraud, meningkatkan transparansi, dan memperkuat tata kelola pengadaan.

| **No** | **Risiko**                  | **Control Objective**                                 | **Control Activity**                                                                     | **Control Type** | **PIC**                            |
| ------ | --------------------------- | ----------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------- | ---------------------------------- |
| 1      | PR tidak valid              | Memastikan kebutuhan pengadaan sah dan terdokumentasi | Approval berjenjang sesuai workflow                                                      | Preventive       | Entity Approver / Holding Approver |
| 2      | Vendor favorit              | Menjamin transparansi pemilihan vendor                | Minimum vendor bidding atau mandatory justification Direct Appointment; weighted scoring | Preventive       | Procurement                        |
| 3      | Manipulasi harga            | Menjaga integritas data harga                         | Lock data setelah approval; penggunaan Reference Price / eCatalog                        | Preventive       | System                             |
| 4      | Approval bypass             | Menjaga kepatuhan workflow approval                   | Mandatory system approval dan dynamic approval workflow                                  | Preventive       | System                             |
| 5      | Keterlambatan approval      | Menjaga SLA approval                                  | Reminder dan escalation                                                                  | Detective        | System                             |
| 6      | Perubahan tanpa log         | Menjaga auditability                                  | Audit trail permanen                                                                     | Preventive       | System                             |
| 7      | Vendor collusion            | Memastikan harga wajar dan evaluasi objektif          | Comparison report terdokumentasi; technical evaluation; commercial evaluation            | Detective        | Procurement                        |
| 8      | Pengadaan melebihi anggaran | Menjaga kepatuhan anggaran                            | Validasi Budget Management Module dan Over Budget Escalation                             | Preventive       | System / Entity Admin              |
| 9      | Penggunaan vendor blacklist | Menjaga kepatuhan terhadap kebijakan vendor           | Vendor Blacklist Check otomatis saat Vendor Prequalification                             | Preventive       | System                             |

# Lampiran C: Role Access Matrix

Matriks ini merangkum hak akses fungsional utama untuk setiap peran dalam sistem E-Procurement.

_Catatan: Peran Finance tidak ditampilkan sebagai kolom terpisah pada matriks ini karena keterlibatannya bersifat conditional, yaitu sebagai approver tambahan untuk pengadaan Over Budget atau Non Budget sesuai governance yang berlaku._

| **Capability Utama**           | **Holding Admin** | **Entity Admin** | **Requestor** | **Entity Approver** | **Holding Approver** | **Procurement** | **Management** | **Internal Audit** | **Vendor** |
| ------------------------------ | ----------------- | ---------------- | ------------- | ------------------- | -------------------- | --------------- | -------------- | ------------------ | ---------- |
| Create PR                      | \-                | \-               | R/W           | \-                  | \-                   | \-              | \-             | \-                 | \-         |
| Revise/Resubmit PR             | \-                | \-               | R/W           | \-                  | \-                   | \-              | \-             | \-                 | \-         |
| Approve PR                     | \-                | \-               | \-            | R/W                 | R/W                  | \-              | \-             | \-                 | \-         |
| Select Procurement Method      | \-                | \-               | \-            | \-                  | \-                   | R/W             | \-             | \-                 | \-         |
| Create/Manage RFQ              | \-                | \-               | \-            | \-                  | \-                   | R/W             | \-             | \-                 | \-         |
| Manage Entity                  | R/W               | \-               | \-            | \-                  | \-                   | \-              | \-             | \-                 | \-         |
| Manage User                    | R/W               | R/W              | \-            | \-                  | \-                   | \-              | \-             | \-                 | \-         |
| Assign User Role               | R/W               | R/W              | \-            | \-                  | \-                   | \-              | \-             | \-                 | \-         |
| Reset Password                 | R/W               | R/W              | \-            | \-                  | \-                   | \-              | \-             | \-                 | \-         |
| Delegate Approver              | R/W               | R/W              | \-            | \-                  | \-                   | \-              | \-             | \-                 | \-         |
| Set Entity Governance          | R/W               | R/W              | \-            | \-                  | \-                   | \-              | \-             | \-                 | \-         |
| Set Dynamic Approval Workflow  | R/W               | R/W              | \-            | \-                  | \-                   | \-              | \-             | \-                 | \-         |
| Set Budget Governance          | R/W               | R/W              | \-            | \-                  | \-                   | \-              | \-             | \-                 | \-         |
| Vendor Blacklist Management    | R/W               | \-               | \-            | \-                  | \-                   | \-              | \-             | \-                 | \-         |
| Vendor Comparison / Evaluation | \-                | \-               | \-            | \-                  | \-                   | R/W             | \-             | \-                 | \-         |
| Reference Price / eCatalog     | \-                | \-               | \-            | \-                  | \-                   | R/W             | \-             | \-                 | \-         |
| Create PO                      | \-                | \-               | \-            | \-                  | \-                   | R/W             | \-             | \-                 | \-         |
| Approve PO                     | \-                | \-               | \-            | R/W                 | R/W                  | \-              | \-             | \-                 | \-         |
| View Cross-Entity Procurement  | R                 | \-               | \-            | \-                  | \-                   | \-              | R              | R                  | \-         |
| View Group Dashboard           | R                 | \-               | \-            | \-                  | \-                   | \-              | R              | R                  | \-         |
| View Audit Trail               | R                 | R                | \-            | \-                  | \-                   | \-              | R              | R                  | \-         |
| Participate in Tender          | \-                | \-               | \-            | \-                  | \-                   | \-              | \-             | \-                 | R/W        |
| Submit Quotation               | \-                | \-               | \-            | \-                  | \-                   | \-              | \-             | \-                 | R/W        |
| Upload Supporting Documents    | \-                | \-               | \-            | \-                  | \-                   | \-              | \-             | \-                 | R/W        |
| Confirm PO                     | \-                | \-               | \-            | \-                  | \-                   | \-              | \-             | \-                 | R/W        |

Keterangan: R/W = Read / Write / Modify, R = Read / View Only, - = No Access
