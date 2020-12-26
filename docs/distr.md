# Contracter
## PerformTask(Request) : ContracterTasker, err 
## ContracterTasker.DealerTasks() []DealTasker, err

# ContractDealer
## AllocateTask(Request) : DealTasker, err
## FindTaskByID(id) : DealTasker, err
## FindTasksByOrderID(orderId) : []DealTasker, err
## GetStorageClaim(DealTasker) : StorageClaimer
## NotifyRawUpload(IOProgresser) : err 
## NotifyResultDownload(IOProgresser) : err 
## PublishTask(DealTasker) : err
## CancelTask(DealTasker) : err
## Subscription(DealTasker) : Subscriber, err

# WorkDealer
## FindFreeTask() : DealTasker, err
## GetStorageClaim(DealerTask) : StorageClaimer
## NotifyRawDownload(IOProgresser) : err 
## NotifyResultUpload(IOProgresser) : err 
## NotifyProcess(ProcessProgresser) : err 
## FinishTask(ProcessProgresser) : err 

# Subscriber
## Subscribe() : chan Progresser
## Unsubscribe()

# Progresser
## Percent() float64

# Storager
## StorageClaimer.Name() : string, err
## StorageClaimer.Size() : int, err
## StorageClaimer.Writer() : io.WriteCloser, err
## StorageClaimer.Reader() : io.ReadCloser, err

# StorageController
## AllocateStorageClaim(name) : StorageClaimer, err
## PurgeStorageClaim(StorageClaimer) : err


# Worker
## Start()
