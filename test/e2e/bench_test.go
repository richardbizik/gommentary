package e2e

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/richardbizik/gommentary/internal/rest/handlers"
	"github.com/stretchr/testify/assert"
)

func BenchmarkCreateComment(b *testing.B) {
	subjectName := fmt.Sprintf("test-subject-%d", rand.Int63())
	port := getRandomPort()
	app := NewApplication(port)
	time.Sleep(1 * time.Second) // wait for server to start
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, respCode := Request(b, "POST", fmt.Sprintf("http://localhost:%d/subject/subjectName/comment", port), handlers.CreateComment{
			SubjectName: &subjectName,
			Text:        fmt.Sprintf("comment-%d", i),
		})
		assert.Equal(b, 200, respCode, resp)
	}
	b.StopTimer()
	app.Stop()
}
