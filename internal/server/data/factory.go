package data

type DataAccessFactory interface {
	GetUserDataAccess() UserDataAccess
	GetExpenseDataAccess() ExpenseDataAccess
}

func NewDataAccessFactory(params interface{}) DataAccessFactory {
	if params == nil {
		return &memoryDataAccessFactory{}
	}

	switch params.(type) {
	case MongoDataAccessParams:
		return NewMongoDataAccessFactory(params.(MongoDataAccessParams))
	default:
		panic("Error in DAO initialisation")
	}
}
