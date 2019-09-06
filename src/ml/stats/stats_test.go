package stats

import "testing"

const testDbName string = "test"

func storeTestResult(input []string, result bool) {
	storeResultIn(testDbName, input, result)
}

func cleanTestDb() {
	cleanDb(testDbName)
}

func testStats() (int, int) {
	return getStatsIn(testDbName)
}

func assertStats(t *testing.T, mutant int, noMutant int) {
	rMutant, rNoMutant := testStats()
	if rMutant != mutant || rNoMutant != noMutant {
		t.Errorf("Stats dont match (%v,%v) != (%v,%v)",
			mutant, noMutant, rMutant, rNoMutant)
	}
}

func TestEmptyStats(t *testing.T) {
	assertStats(t, 0, 0)
}

func TestAllMutants(t *testing.T) {
	defer cleanTestDb()
	storeTestResult([]string{"A"}, true)
	assertStats(t, 1, 0)
	storeTestResult([]string{"B"}, true)
	assertStats(t, 2, 0)
	storeTestResult([]string{"C"}, true)
	assertStats(t, 3, 0)
	storeTestResult([]string{"C"}, true)
	assertStats(t, 3, 0)
}

func TestAllNormal(t *testing.T) {
	defer cleanTestDb()
	storeTestResult([]string{"A"}, false)
	assertStats(t, 0, 1)
	storeTestResult([]string{"B"}, false)
	assertStats(t, 0, 2)
	storeTestResult([]string{"C"}, false)
	assertStats(t, 0, 3)
	storeTestResult([]string{"C"}, false)
	assertStats(t, 0, 3)
}
