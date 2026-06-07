package taskmanager

import (
	"errors"
	"sync"
	"time"
)

type TaskStatus string

const (
	TODO       TaskStatus = "TODO"
	INPROGRESS TaskStatus = "IN_PROGRESS"
	DONE       TaskStatus = "DONE"
)

var (
	ErrUserAlreadyInTeam     = errors.New("user is already in the team")
	ErrUserNotInTeam         = errors.New("user not found in the team")
	ErrNegativeElapsedTime   = errors.New("elapsed time cannot be negative")
	ErrUserAlreadySubscribed = errors.New("user is already subscribed to the task")
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
}

type Admin struct {
	User
}

func (a *Admin) AddUserToTeam(team *Team, user *User) error {
	for _, m := range team.Members {
		if m.ID == user.ID {
			return ErrUserAlreadyInTeam
		}
	}
	team.Members = append(team.Members, user)
	return nil
}

func (a *Admin) RemoveUserFromTeam(team *Team, user *User) error {
	for i, m := range team.Members {
		if m.ID == user.ID {
			team.Members = append(team.Members[:i], team.Members[i+1:]...)
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
	t.RemainingTime -= elapsed
	if t.RemainingTime < 0 {
		t.RemainingTime = 0
	}
	return nil
}

func (t *Timer) CheckDeadline() bool {
	return t.RemainingTime <= 0
}

type Task struct {
	ID            int
	Description   string
	Status        TaskStatus
	Deadline      time.Time
	AssignedUsers []*User
	Timer         *Timer
}

func (t *Task) UpdateStatus(newStatus TaskStatus) {
	t.Status = newStatus
}

func (t *Task) SubscribeUser(user *User) error {
	for _, u := range t.AssignedUsers {
		if u.ID == user.ID {
			return ErrUserAlreadySubscribed
		}
	}
	t.AssignedUsers = append(t.AssignedUsers, user)
	return nil
}

type NotificationService struct{}

func (ns *NotificationService) NotifyDeadlineReached(task *Task) string {
	if task.Timer != nil && task.Timer.CheckDeadline() {
		return "Deadline reached for task: " + task.Description
	}
	return ""
}

type MediaSettings struct {
	SelectedMusic string
}

func (ms *MediaSettings) SelectMusic(track string) {
	ms.SelectedMusic = track
}
