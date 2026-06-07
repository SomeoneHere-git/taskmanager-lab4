package taskmanager

import (
	"testing"
	"time"
)

// Тестування методу AddUserToTeam
func TestAdmin_AddUserToTeam(t *testing.T) {
	// EP: Позитивний. Додавання нового користувача у пусту команду.
	t.Run("Add new user to empty team - Positive EP", func(t *testing.T) {
		// Arrange
		admin := &Admin{User: User{ID: 1}}
		team := &Team{ID: 1, Name: "Dev"}
		user := &User{ID: 2}

		// Act
		err := admin.AddUserToTeam(team, user)

		// Assert
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(team.Members) != 1 || team.Members[0].ID != 2 {
			t.Errorf("User was not added correctly")
		}
	})

	// EP: Негативний. Спроба додати існуючого користувача.
	t.Run("Add existing user - Negative EP", func(t *testing.T) {
		// Arrange
		admin := &Admin{User: User{ID: 1}}
		user := &User{ID: 2}
		team := &Team{ID: 1, Name: "Dev", Members: []*User{user}}

		// Act
		err := admin.AddUserToTeam(team, user)

		// Assert
		if err != ErrUserAlreadyInTeam {
			t.Errorf("Expected ErrUserAlreadyInTeam, got: %v", err)
		}
		if len(team.Members) != 1 {
			t.Errorf("Team members should not duplicate, length: %d", len(team.Members))
		}
	})
}

// Тестування методу RemoveUserFromTeam
func TestAdmin_RemoveUserFromTeam(t *testing.T) {
	// BVA: Позитивний. Видалення користувача, коли в команді лише він (межа масиву: 1 -> 0)
	t.Run("Remove only user from team - Positive BVA", func(t *testing.T) {
		// Arrange
		admin := &Admin{User: User{ID: 1}}
		user := &User{ID: 2}
		team := &Team{ID: 1, Name: "Dev", Members: []*User{user}}

		// Act
		err := admin.RemoveUserFromTeam(team, user)

		// Assert
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(team.Members) != 0 {
			t.Errorf("Expected empty team, got len %d", len(team.Members))
		}
	})

	// EP: Негативний. Спроба видалити користувача, якого немає в команді.
	t.Run("Remove non-existent user - Negative EP", func(t *testing.T) {
		// Arrange
		admin := &Admin{User: User{ID: 1}}
		user1 := &User{ID: 2}
		user2 := &User{ID: 3}
		team := &Team{ID: 1, Name: "Dev", Members: []*User{user1}}

		// Act
		err := admin.RemoveUserFromTeam(team, user2)

		// Assert
		if err != ErrUserNotInTeam {
			t.Errorf("Expected ErrUserNotInTeam, got: %v", err)
		}
		if len(team.Members) != 1 {
			t.Errorf("Team should not be modified")
		}
	})
}

// Тестування методу AutoUpdateTime
func TestTimer_AutoUpdateTime(t *testing.T) {
	// EP: Негативний. Неприпустиме значення часу (elapsed < 0)
	t.Run("Negative elapsed time - Negative EP", func(t *testing.T) {
		// Arrange
		timer := &Timer{RemainingTime: 10 * time.Second}

		// Act
		err := timer.AutoUpdateTime(-5 * time.Second)

		// Assert
		if err != ErrNegativeElapsedTime {
			t.Errorf("Expected ErrNegativeElapsedTime, got %v", err)
		}
		if timer.RemainingTime != 10*time.Second {
			t.Errorf("Timer should not change")
		}
	})

	// EP: Позитивний. Стандартне зменшення часу.
	t.Run("Standard time decrease - Positive EP", func(t *testing.T) {
		// Arrange
		timer := &Timer{RemainingTime: 10 * time.Second}

		// Act
		err := timer.AutoUpdateTime(4 * time.Second)

		// Assert
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if timer.RemainingTime != 6*time.Second {
			t.Errorf("Expected 6s, got %v", timer.RemainingTime)
		}
	})

	// BVA: Позитивний. Зменшення рівно до 0.
	t.Run("Decrease exactly to zero - Positive BVA", func(t *testing.T) {
		// Arrange
		timer := &Timer{RemainingTime: 5 * time.Second}

		// Act
		err := timer.AutoUpdateTime(5 * time.Second)

		// Assert
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if timer.RemainingTime != 0 {
			t.Errorf("Expected 0, got %v", timer.RemainingTime)
		}
	})

	// BVA: Позитивний. Зменшення нижче 0 (має зафіксуватись на 0).
	t.Run("Decrease below zero - Positive BVA", func(t *testing.T) {
		// Arrange
		timer := &Timer{RemainingTime: 5 * time.Second}

		// Act
		err := timer.AutoUpdateTime(10 * time.Second)

		// Assert
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if timer.RemainingTime != 0 {
			t.Errorf("Expected 0, got %v", timer.RemainingTime)
		}
	})
}

// Тестування методу SubscribeUser
func TestTask_SubscribeUser(t *testing.T) {
	// EP: Позитивний. Підписка на завдання.
	t.Run("Subscribe new user - Positive EP", func(t *testing.T) {
		// Arrange
		task := &Task{ID: 1, Description: "Task A"}
		user := &User{ID: 10}

		// Act
		err := task.SubscribeUser(user)

		// Assert
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(task.AssignedUsers) != 1 || task.AssignedUsers[0].ID != 10 {
			t.Errorf("User was not subscribed correctly")
		}
	})

	// EP: Негативний. Спроба підписатись двічі.
	t.Run("Subscribe duplicate user - Negative EP", func(t *testing.T) {
		// Arrange
		user := &User{ID: 10}
		task := &Task{ID: 1, Description: "Task A", AssignedUsers: []*User{user}}

		// Act
		err := task.SubscribeUser(user)

		// Assert
		if err != ErrUserAlreadySubscribed {
			t.Errorf("Expected ErrUserAlreadySubscribed, got %v", err)
		}
		if len(task.AssignedUsers) != 1 {
			t.Errorf("Should not duplicate user")
		}
	})
}

func TestOtherMethods(t *testing.T) {
	// EP: Позитивний
	t.Run("User authentication", func(t *testing.T) {
		// Arrange
		user := &User{Email: "test@test.com"}
		// Act
		ok := user.Authenticate()
		// Assert
		if !ok {
			t.Error("Expected successful auth")
		}
	})

	// EP: Позитивний
	t.Run("Update task status", func(t *testing.T) {
		// Arrange
		task := &Task{Status: TODO}
		// Act
		task.UpdateStatus(INPROGRESS)
		// Assert
		if task.Status != INPROGRESS {
			t.Error("Status was not updated")
		}
	})

	// BVA: Позитивний (границя дедлайну досягнута)
	t.Run("Notification service deadline reached", func(t *testing.T) {
		// Arrange
		ns := &NotificationService{}
		task := &Task{Description: "A", Timer: &Timer{RemainingTime: 0}}
		// Act
		msg := ns.NotifyDeadlineReached(task)
		// Assert
		if msg == "" {
			t.Error("Expected notification message")
		}
	})

	// EP: Позитивний
	t.Run("Media settings selection", func(t *testing.T) {
		// Arrange
		ms := &MediaSettings{}
		// Act
		ms.SelectMusic("Rock")
		// Assert
		if ms.SelectedMusic != "Rock" {
			t.Error("Music was not selected")
		}
	})
}
