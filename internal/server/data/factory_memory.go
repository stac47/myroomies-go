package data

// Factory to access data from memory

type memoryDataAccessFactory struct{}

func (f memoryDataAccessFactory) GetUserDataAccess() UserDataAccess {
	return GetMemoryUserDataAccess()
}

func (f memoryDataAccessFactory) GetExpenseDataAccess() ExpenseDataAccess {
	return GetMemoryExpenseDataAccess()
}
