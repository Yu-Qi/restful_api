package seed

import "github.com/Yu-Qi/restful_api/domain"

func Tasks() []*domain.Task {
	return []*domain.Task{
		{
			ID:     1,
			Name:   "task1",
			Status: 0,
		},
		{
			ID:     2,
			Name:   "task2",
			Status: 1,
		},
		{
			ID:     3,
			Name:   "task3",
			Status: 0,
		},
		{
			ID:     4,
			Name:   "task4",
			Status: 0,
		},
		{
			ID:     5,
			Name:   "task5",
			Status: 1,
		},
	}
}
