package models

// // ConvertContracterRequest _
// type ConvertContracterRequest struct {
// 	ConvertTask convert.ConverterTask
// }

// // ConvertDealerRequest _
// type ConvertDealerRequest struct {
// 	Type             string
// 	Identity         string
// 	VideoCodec       string
// 	HWAccel          string
// 	VideoBitRate     string
// 	VideoQuality     int
// 	Preset           string
// 	Scale            string
// 	KeyframeInterval int
// }

// // GetType _
// func (cdr *ConvertDealerRequest) GetType() string {
// 	return "convert"
// }

// // ID _
// func (cdr *ConvertDealerRequest) ID() string {
// 	return cdr.Identity
// }

// // GetType _
// func (cr *ConvertContracterRequest) GetType() string {
// 	return "convert"
// }

// // ConvertContracterTask _
// type ConvertContracterTask struct {
// 	// ConvertTask convert.ConverterTask
// 	Identity           string
// 	Request            *ConvertContracterRequest
// 	ConvertDealerTasks []*ConvertDealerTask
// 	MessageBus         *MessageBus
// }

// // DealerTasks _
// func (ct *ConvertContracterTask) DealerTasks() []ITask {
// 	dealTaskers := make([]ITask, 0)

// 	for _, task := range ct.ConvertDealerTasks {
// 		dealTaskers = append(dealTaskers, task)
// 	}

// 	// return t.ConvertDealerTasks
// 	return dealTaskers
// }

// // Failed _
// func (ct *ConvertContracterTask) Failed(err error) {
// 	// TODO:
// 	panic(ErrNotImplemented)
// 	// return
// }

// // ConvertDealerTask _
// type ConvertDealerTask struct {
// 	Identity string
// }

// // ID _
// func (dt *ConvertDealerTask) ID() string {
// 	return dt.Identity
// }
