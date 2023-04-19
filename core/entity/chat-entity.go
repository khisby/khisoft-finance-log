package entity

const (
	NewCommers = `Hai, Tampaknya kamu adalah pengguna baru ya.`
	MenuText   = `
Untuk menggunakan aplikasi Finance-Log, Kamu bisa pakai command ini ya:
1. [keluar/masuk] [jumlah] dari/buat [kategori] [deskripsi] -> save keuangan kamu
2. [debit/kredit] enter [jumlah] enter [kategori] enter [deskripsi] -> save keuangan kamu
3. [laporan/report] [hari/bulan] [ini/kemarin] -> untuk melihat laporan keuangan kamu
4. hapus terakhir/barusan -> untuk menghapus data terakhir
5. link -> untuk melihat link google sheet kamu
6. bantuan/bantuin -> untuk melihat bantuan

kategori yg tersedia ini ya kak :
1. makan (termasuk: minum)
2. jajan (termasuk: ngopi/nongkrong)
3. jalan (termasuk: liburan)
4. transport (termasuk: bensin, parkir) 
5. bulanan (termasuk: listrik, air, internet)
6. belanja (termasuk: belanja online)
7. project 
8. kerja 
9. Utang
10. Sedekah
11. lainnya
	`
	ReplyChatSaved           = "Okay kak, aku catet yak. Detailnya gini kak \n\nStatus : %s \nNominal : %s \nKategori : %s \nKeterangan : %s"
	ReportTextNotFound       = "Maaf, aku ga nemu data keuangan kamu kak. Kamu udah pernah catat keuangan belum?"
	ReportTextHeader         = "Ini laporan keuangan kamu %s %s kak : \n\n"
	ReportTextPemasukan      = "Total Pemasukan : %s \n"
	ReportTextPengeluaran    = "Total Pengeluaran : %s \n"
	ReportTextCategoryHeader = "\n\n%s by Kategori : \n"
	ReportTextCategory       = "%s : %s \n"
	DeletedText              = "Data %s kamu udah aku hapus kak, Detailnya yg ini kak \n\nStatus : %s \nNominal : %s \nKategori : %s \nKeterangan : %s"
)

type Category string

const (
	Makan     Category = "Makan/Minum"
	Jajan     Category = "Jajan"
	Jalan     Category = "Jalan-Jalan"
	Transport Category = "Transport"
	Bulanan   Category = "Bulanan"
	Belanja   Category = "Belanja"
	Lainnya   Category = "Lain-Lain"
	Project   Category = "Project"
	Kerja     Category = "Kerja"
	Hutang    Category = "Hutang"
	Sedekah   Category = "Sedekah"
)
