package entity

const (
	NewCommers = `Hai, Tampaknya kamu adalah pengguna baru ya.`
	MenuText   = `
Untuk menggunakan aplikasi Finance-Log, Kamu bisa pakai command ini ya:
1. [keluar/masuk] [jumlah] dari/buat [kategori] [deskripsi] -> save keuangan kamu
2. [debit/kredit] enter [jumlah] enter [kategori] enter [deskripsi] -> save keuangan kamu
3. [laporan/report] [hari/bulan] [ini/kemarin]
4. link -> untuk melihat link google sheet kamu
5. bantuan/bantuin -> untuk melihat bantuan

kategori yg tersedia ini ya kak :
1. makan
2. jajan
3. jalan
4. bensin
5. bulanan
6. belanja
7. project
8. kerja
9. lainnya
	`
	ReplyChatSaved           = "Okay kak, aku catet yak. Detailnya gini kak \n\nTime : %s \nCategory : %s \nAmount : %s \nStatus : %s \nDescription : %s"
	ReportTextNotFound       = "Maaf, aku ga nemu data keuangan kamu kak. Kamu udah pernah catat keuangan belum?"
	ReportTextHeader         = "Ini laporan keuangan kamu %s %s kak : \n\n"
	ReportTextPemasukan      = "Total Pemasukan : %s \n"
	ReportTextPengeluaran    = "Total Pengeluaran : %s \n"
	ReportTextCategoryHeader = "\n\n%s by Kategori : \n"
	ReportTextCategory       = "%s : %s \n"
)

type Category string

const (
	Makan   Category = "Makan/Minum"
	Jajan   Category = "Jajan"
	Jalan   Category = "Jalan-Jalan"
	Bensin  Category = "Bensin"
	Bulanan Category = "Bulanan"
	Belanja Category = "Belanja"
	Lainnya Category = "Lain-Lain"
	Project Category = "Project"
	Kerja   Category = "Kerja"
)
