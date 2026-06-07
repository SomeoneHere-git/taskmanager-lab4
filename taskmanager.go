package taskmanager

import (
	"sync"
	"time"
)

type TaskStatus string

const (
	TODO       TaskStatus = "TODO"
	INPROGRESS TaskStatus = "IN_PROGRESS"
	DONE       TaskStatus = "DONE"
)

type Error string

func (e Error) Error() string { return string(e) }

const (
	ErrUserAlreadyInTeam     Error = "user is already in the team"
	ErrUserNotInTeam         Error = "user not found in the team"
	ErrNegativeElapsedTime   Error = "elapsed time cannot be negative"
	ErrUserAlreadySubscribed Error = "user is already subscribed to the task"
)

type User struct {
	ID    int
	Email string
}

func (u *User) Authenticate() bool {
	return u.Email != ""
}

type Team struct {
	ID      int
	Name    string
	Members []*User
	mu      sync.Mutex
}

type Admin struct {
	User User
}

func (a *Admin) AddUserToTeam(team *Team, user *User) error {
	team.mu.Lock()
	defer team.mu.Unlock()
	for _, m := range team.Members {
		if m.ID == user.ID {
			return ErrUserAlreadyInTeam
		}
	}
	team.Members = append(team.Members, user)
	return nil
}

func (a *Admin) RemoveUserFromTeam(team *Team, user *User) error {
	team.mu.Lock()
	defer team.mu.Unlock()
	for i, m := range team.Members {
		if m.ID == user.ID {
			copy(team.Members[i:], team.Members[i+1:])
			team.Members[len(team.Members)-1] = nil
			team.Members = team.Members[:len(team.Members)-1]
			return nil
		}
	}
	return ErrUserNotInTeam
}

type Timer struct {
	RemainingTime time.Duration
	mu            sync.Mutex
}

func (t *Timer) AutoUpdateTime(elapsed time.Duration) error {
	if elapsed < 0 {
		return ErrNegativeElapsedTime
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.RemainingTime -= elapsed
	if t.RemainingTime < 0 {
		t.RemainingTime = 0
	}
	return nil
}

func (t *Timer) CheckDeadline() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.RemainingTime <= 0
}

type Task struct {
	ID            int
	Description   string
	Status        TaskStatus
	Deadline      time.Time
	AssignedUsers map[int]*User
	Timer         *Timer
}

func (t *Task) UpdateStatus(newStatus TaskStatus) {
	t.Status = newStatus
}

func (t *Task) SubscribeUser(user *User) error {
	if t.AssignedUsers == nil {
		t.AssignedUsers = make(map[int]*User)
	}
	if _, exists := t.AssignedUsers[user.ID]; exists {
		return ErrUserAlreadySubscribed
	}
	t.AssignedUsers[user.ID] = user
	return nil
}

type NotificationService struct{}

func (ns *NotificationService) NotifyDeadlineReached(task *Task) bool {
	if task.Timer != nil && task.Timer.CheckDeadline() {
		return true
	}
	return false
}

type MediaSettings struct {
	SelectedMusic string
}

func (ms *MediaSettings) SelectMusic(track string) {
	ms.SelectedMusic = track
}
