package models

// import "encoding/json"

// // ConvertTask _
// type ConvertTask struct {
// 	Identity             string
// 	OrderIdentity        string
// 	Type                 string
// 	StorageClaimIdentity string

// 	Muxer            string
// 	VideoCodec       string
// 	HWAccel          string
// 	VideoBitRate     string
// 	VideoQuality     int
// 	Preset           string
// 	Scale            string
// 	KeyframeInterval int
// }

// // GetID _
// func (ct *ConvertTask) GetID() string {
// 	return ct.Identity
// }

// // GetType _
// func (ct *ConvertTask) GetType() string {
// 	return ConvertV1Type
// }

// // GetOrderID _
// func (ct *ConvertTask) GetOrderID() string {
// 	return ct.OrderIdentity
// }

// // GetStorageClaimIdentity _
// func (ct *ConvertTask) GetStorageClaimIdentity() string {
// 	return ct.OrderIdentity
// }

// // GetPayload _
// func (ct *ConvertTask) GetPayload() (string, error) {
// 	b, err := json.Marshal(ct)

// 	return string(b), err
// }

// // // GetStorageClaim _
// // func (ct *ConvertTask) GetStorageClaim() IStorageClaim {
// // 	return ct.StorageClaim
// // }

// // Failed _
// func (ct *ConvertTask) Failed(err error) {
// 	// TODO:
// 	panic(ErrNotImplemented)
// 	// return
// }
