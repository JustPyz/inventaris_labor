package jobs

import (
	"log"
	"time"

	"invela-be/internal/repositories"
)

// PeminjamanOverdueJob adalah background job yang berjalan secara periodik
// untuk menandai peminjaman yang melewati batas waktu pengembalian.
type PeminjamanOverdueJob struct {
	repo     *repositories.PeminjamanRepository
	interval time.Duration
	stopCh   chan struct{}
}

// NewPeminjamanOverdueJob membuat instance job baru.
// interval menentukan seberapa sering pengecekan dilakukan (misal: time.Minute).
func NewPeminjamanOverdueJob(repo *repositories.PeminjamanRepository, interval time.Duration) *PeminjamanOverdueJob {
	return &PeminjamanOverdueJob{
		repo:     repo,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

// Start menjalankan job di goroutine terpisah.
// Job langsung dijalankan sekali saat startup, lalu berulang setiap interval.
func (j *PeminjamanOverdueJob) Start() {
	log.Printf("[OverdueJob] started, interval: %s", j.interval)

	go func() {
		// Jalankan sekali saat startup agar tidak perlu tunggu interval pertama
		j.run()

		ticker := time.NewTicker(j.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				j.run()
			case <-j.stopCh:
				log.Println("[OverdueJob] stopped")
				return
			}
		}
	}()
}

// Stop menghentikan job secara graceful.
func (j *PeminjamanOverdueJob) Stop() {
	close(j.stopCh)
}

// run adalah satu siklus pengecekan dan update.
func (j *PeminjamanOverdueJob) run() {
	// Gunakan format tanggal yang sama dengan yang disimpan di database (YYYY-MM-DD)
	today := time.Now().Format("2006-01-02")

	updated, err := j.repo.MarkOverdue(today)
	if err != nil {
		log.Printf("[OverdueJob] error marking overdue: %v", err)
		return
	}

	if updated > 0 {
		log.Printf("[OverdueJob] %d peminjaman ditandai 'melewati batas waktu'", updated)
	}
}
