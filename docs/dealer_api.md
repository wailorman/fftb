# Dealer API

```go
// IContracterDealer _
type IContracterDealer interface {
  AllocatePublisherAuthority(name string) (IAuthor, error)
  // POST /publisher_authorities
  
  AllocatePerformerAuthority(name string) (IAuthor, error)
  // POST /performer_authorities
  
  
  FindFreeSegment(performer IAuthor) (ISegment, error)
  // POST /segments/free | Segment
  
  GetSegmentByID(segmentID string) (ISegment, error)
  // GET /segments/{id} | Segment
  
  
  AllocateSegment(req IDealerRequest) (ISegment, error)
  // POST /segments 
  
  PublishSegment(publisher IAuthor, id string) error
  // POST /segments/{id}/actions/publish
  
  RepublishSegment(publisher IAuthor, id string) error
  // POST /segments/{id}/actions/republish
  
  CancelSegment(publisher IAuthor, id string) error
  // POST /segments/{id}/actions/cancel | { reason: failed }
  
  AcceptSegment(publisher IAuthor, id string) error
  // POST /segments/{id}/actions/accept
  
  FinishSegment(performer IAuthor, id string) error
  // POST /segments/{id}/actions/finish
  
  QuitSegment(performer IAuthor, id string) error
  // POST /segments/{id}/actions/quit
  
  FailSegment(performer IAuthor, id string, err error) error
  // POST /segments/{id}/actions/fail

  GetOutputStorageClaim(publisher IAuthor, segmentID string) (IStorageClaim, error)
  // GET /segments/{id}/output_storage_claim | { storage_claim: http... }
  GetInputStorageClaim(performer IAuthor, segmentID string) (IStorageClaim, error)
  // GET /segments/{id}/input_storage_claim | { storage_claim: http... }
  
  AllocateInputStorageClaim(publisher IAuthor, id string) (IStorageClaim, error)
  // POST /segments/{id}/output_storage_claim | { storage_claim: http... }
  AllocateOutputStorageClaim(performer IAuthor, id string) (IStorageClaim, error)
  // POST /segments/{id}/input_storage_claim | { storage_claim: http... }

  GetQueuedSegmentsCount(fctx context.Context, publisher IAuthor) (int, error)
  // GET /segments?queued=true
  
  GetSegmentsByOrderID(fctx context.Context, orderID string, search ISegmentSearchCriteria) ([]ISegment, error)
  // GET /orders/{order_id}/segments
  // GET /segments?order_id=
  
  GetSegmentsStatesByOrderID(fctx context.Context, orderID string) (map[string]string, error)
  // GET /orders/{order_id}/segments_states

  
  
  
  
  NotifyRawUpload(publisher IAuthor, id string, p Progresser) error
  // POST /segments/{id}/notifications/input_upload | { progress: 0.5 }
  NotifyResultDownload(publisher IAuthor, id string, p Progresser) error
  // POST /segments/{id}/notifications/output_download | { progress: 0.5 }
  NotifyRawDownload(performer IAuthor, id string, p Progresser) error
  // POST /segments/{id}/notifications/input_download | { progress: 0.5 }
  NotifyResultUpload(performer IAuthor, id string, p Progresser) error
  // POST /segments/{id}/notifications/ouput_upload | { progress: 0.5 }
  NotifyProcess(performer IAuthor, id string, p Progresser) error
  // POST /segments/{id}/notifications/process | { progress: 0.5 }



}
```
