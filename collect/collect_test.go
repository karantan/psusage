package collect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProgramCPU_OldEmpty(t *testing.T) {
	fstats := func(program string) (usages []CPU_Usage) {
		return []CPU_Usage{
			{0.1, 5, "myprogram", 100, "root"},
			{0.4, 61, "myprogram", 200, "worker1"},
			{28.0, 3672, "myprogram", 300, "worker2"},
			{6.1, 3, "myprogram", 400, "worker3"},
		}
	}
	wantRunning := []CPU_Usage{
		{0.1, 5, "myprogram", 100, "root"},
		{0.4, 61, "myprogram", 200, "worker1"},
		{28.0, 3672, "myprogram", 300, "worker2"},
		{6.1, 3, "myprogram", 400, "worker3"},
	}
	wantStopped := []CPU_Usage{}
	gotRunning, gotStopped := ProgramCPU("mydatabase", []CPU_Usage{}, fstats)
	assert.Equal(t, wantRunning, gotRunning)
	assert.Equal(t, wantStopped, gotStopped)
}
func TestProgramCPU_AddUpdate(t *testing.T) {
	fstats := func(program string) (usages []CPU_Usage) {
		return []CPU_Usage{
			{1, 6, "myprogram", 100, "root"},
			{4, 62, "myprogram", 200, "worker1"},
			{8.0, 3673, "myprogram", 300, "worker2"},
			{5, 4, "myprogram", 400, "worker3"},
		}
	}
	wantRunning := []CPU_Usage{
		{1, 6, "myprogram", 100, "root"},
		{4, 62, "myprogram", 200, "worker1"},
		{5, 4, "myprogram", 400, "worker3"},
		{8.0, 3673, "myprogram", 300, "worker2"},
	}
	wantStopped := []CPU_Usage{}
	gotRunning, gotStopped := ProgramCPU("mydatabase", []CPU_Usage{
		{0, 0, "myprogram", 100, "root"},
		{0, 0, "myprogram", 200, "worker1"},
		{0, 0, "myprogram", 400, "worker3"},
	}, fstats)
	assert.Equal(t, wantRunning, gotRunning)
	assert.Equal(t, wantStopped, gotStopped)
}

func TestProgramCPU_AllStopped(t *testing.T) {
	fstats := func(program string) (usages []CPU_Usage) {
		return
	}
	wantStopped := []CPU_Usage{
		{0, 0, "myprogram", 100, "root"},
		{0, 0, "myprogram", 200, "worker1"},
		{0, 0, "myprogram", 400, "worker3"},
	}
	gotRunning, gotStopped := ProgramCPU("mydatabase", []CPU_Usage{
		{0, 0, "myprogram", 100, "root"},
		{0, 0, "myprogram", 200, "worker1"},
		{0, 0, "myprogram", 400, "worker3"},
	}, fstats)
	assert.Equal(t, 0, len(gotRunning))
	assert.Equal(t, wantStopped, gotStopped)
}

func TestProgramCPU_AllEmpty(t *testing.T) {
	fstats := func(program string) (usages []CPU_Usage) {
		return
	}
	gotRunning, gotStopped := ProgramCPU("mydatabase", []CPU_Usage{}, fstats)
	assert.Equal(t, 0, len(gotRunning))
	assert.Equal(t, 0, len(gotStopped))
}

func Test_parseStatPS(t *testing.T) {
	psOut := `
 0.1 00:00:05 3477510 root myprogram
 0.4 00:01:01 3518860 worker1 myprogram
28.0 01:01:12 3520027 worker2 myprogram
 6.1 00:00:03 3519918 worker3 myprogram
`
	got := parseStatPS(psOut)
	want := []CPU_Usage{
		{0.1, 6, "myprogram", 3477510, "root"},
		{0.4, 62, "myprogram", 3518860, "worker1"},
		{28.0, 3673, "myprogram", 3520027, "worker2"},
		{6.1, 4, "myprogram", 3519918, "worker3"},
	}
	assert.Equal(t, want, got)
}

func Test_parseCPUTime(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name         string
		args         args
		wantDuration int
	}{
		{"just seconds", args{"00:00:01"}, 1},
		{"just minutes", args{"00:01:00"}, 60},
		{"just hours", args{"01:00:00"}, 3600},
		{"mix of all", args{"01:01:01"}, 3661},
		{"including days", args{"2-02:02:02"}, 180122},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDuration := parseCPUTime(tt.args.val); gotDuration != tt.wantDuration {
				t.Errorf("parseCPUTime() = %v, want %v", gotDuration, tt.wantDuration)
			}
		})
	}
}
