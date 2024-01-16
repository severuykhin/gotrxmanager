# gotrxmanager
A package for implementing transactions distributed through business logic

## Usage

### When the application starts
Compatible with `sql.DB` interface - accepts `*sql.DB` as connection
```
trm := gotrxmanager.NewTransactionManager(sqlconnection)
```

### In business logic
Lets say that out logic contains several operations with data and we want to perfom this operations in one db transaction
```
func SomeUsecase() (any, error) {
    res, err := usecase.trxManager.Do(ctx, func(ctx context.Context) (any, error) {

        res, err := u.someRepo.Create(ctx, entity)
        if err != nil {
            return nil, err
        }

        res2, err := u.someOtherRepo.Store(ctx, someData)
        if err != nil {
            return nil, err
        }

        res3, err := u.someOtherRepo.Update(ctx, someEntity)
        if err != nil {
            return nil, err
        }

        return res3, nil
    })

    if err != nil {
        return nil, err
    }

    return res, nil
}

```
If an error occurs while executing the logic contained in trxManager.Do, the entire transaction will be rolled back. If the logic completed successfully, transaction will be commited


### In repository or data provider
There is not much left to do - all operation whan working with data should be perfomed through a transaction in the context
```
func Create(entity any) error {
    tx, err := gotrxmanager.TxFromContext(ctx)
    if err != nil {
        return err
    }

    _, _ = tx.Exec(query)
}

```
There is no need to create, commit, rollback transaction manually for each operation
