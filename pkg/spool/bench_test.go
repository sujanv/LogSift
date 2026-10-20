package spool

import "testing"

func BenchmarkDiskSpoolerPush(b *testing.B) {
	tmpDir := b.TempDir()
	sp, err := NewDiskSpooler(tmpDir, 100000)
	if err != nil {
		b.Fatalf("NewDiskSpooler error: %v", err)
	}
	defer sp.Close()

	line := []byte(`{"level":"info","msg":"benchmark spool line","status":200}`)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = sp.Push(line)
	}
}
