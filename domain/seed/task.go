package seed

import "github.com/Yu-Qi/restful_api/domain"

func Tasks() []*domain.Task {
	return []*domain.Task{
		{
			Name:   "task1",
			Status: 0,
		},
		{
			Name:   "task2",
			Status: 1,
		},
		{
			Name:   "task3",
			Status: 0,
		},
		{
			Name:   "task4",
			Status: 0,
		},
		{
			Name:   "task5",
			Status: 1,
		},
	}
}
