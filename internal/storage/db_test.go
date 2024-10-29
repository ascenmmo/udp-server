package memoryDB

import (
	"context"
	"fmt"
	"github.com/ascenmmo/udp-server/internal/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
)

func TestDB(t *testing.T) {
	counter := struct {
		c          int
		m          sync.Mutex
		goroutines int
	}{
		0, sync.Mutex{}, 100000,
	}

	t.Log("sending data")
	ctx, cansel := context.WithCancel(context.Background())
	defer cansel()
	db := NewMemoryDb(ctx, 5)
	wg := sync.WaitGroup{}
	for i := 0; i < counter.goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			db.SetData(fmt.Sprintf("%d", i), strconv.Itoa(i))

		}(i)
	}
	t.Log("Wait data")
	wg.Wait()

	t.Log("Getting data")
	for i := 0; i < counter.goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, ok := db.GetData(fmt.Sprintf("%d", i))
			if !ok {
				t.Log("data not found")
				return
			}
			counter.m.Lock()
			counter.c++
			counter.m.Unlock()
		}(i)
	}

	t.Log("Getting  Wait data")
	wg.Wait()

	t.Log("end ")

	if counter.goroutines != counter.c {
		t.Log("len(dataArray) != trice")
		t.Fatal(counter.goroutines, "!=", counter.c)
	}
}

func BenchmarkDB(b *testing.B) {
	db := NewMemoryDb(context.Background(), 5)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.SetData(fmt.Sprintf("%d", i), strconv.Itoa(i))
	}
}

func BenchmarkGetData(b *testing.B) {
	db := NewMemoryDb(context.Background(), 5)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.SetData(fmt.Sprintf("%d", i), strconv.Itoa(i))
	}
}

func TestSetDataTraceCheck(t *testing.T) {
	db := NewMemoryDb(context.Background(), 5)
	setNewRoomData(db)
	data, ok := db.GetData("1")
	assert.True(t, ok)
	room := data.(*entities.Room)

	fmt.Println(room)
	oldID := room.RoomID
	updateData(db)
	fmt.Println(room)

	assert.NotEqual(t, oldID, room.RoomID)
}

func updateData(db *MemoryDb) {
	data, ok := db.GetData("1")
	if !ok {
		panic("data not found")
	}
	room := data.(*entities.Room)

	room.RoomID = uuid.New()
	room.GameID = uuid.New()
	users := room.GetUser()

	for _, v := range users {
		v.ID = uuid.New()
	}
}

func setNewRoomData(db *MemoryDb) {
	room := entities.Room{
		RoomID: uuid.New(),
		GameID: uuid.New(),
	}
	room.SetUser(&entities.User{
		ID: uuid.New(),
	})

	db.SetData("1", &room)
}
