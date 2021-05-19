package structures

import (
	"encoding/json"
	"fmt"
)

type ProcessType string

const (
	Medicine         ProcessType = "medicine"  // 药品
	Equipment        ProcessType = "equipment" // 设备，如射线仪、B 超仪等设备
	Machine          ProcessType = "machine"   // 器械，如手术刀、担架等器械
	Material         ProcessType = "material"  // 耗材，如口罩、针筒等耗材
	NoneProcessType ProcessType = ""
)

func getProcessType(processTypeStr string) (ProcessType, error) {
	switch processTypeStr {
	case "medicine":
		return Medicine, nil
	case "equipment":
		return Equipment, nil
	case "machine":
		return Machine, nil
	case "material":
		return Material, nil
	default:
		return NoneProcessType, fmt.Errorf("No process type named %s error.", processTypeStr)
	}
}

type ProduceRecord struct {
	Name     string `json:"name"` // 名称
	Producer    string `json:"producer"` // 制造商
	Address  string `json:"address"` // 生产地
	Date     string `json:"date"` // 生产日期
	Life     string `json:"life"` // 保质期, -1 时为永不过期
}

type ProcessRecord struct {
	Name     string `json:"name"` // 名称
	Type  ProcessType `json:"type"` // 类型
	Processor    string `json:"processor"` // 加工商
	Address  string `json:"address"` // 加工地
	Date     string `json:"date"` // 加工日期
	Life     string `json:"life"` // 保质期, -1 时为永不过期
}

type TransportRecord struct {
	Transporter string `json:"transporter"` // 运输商
	OriginAddress string `json:"originAddress"` // 起始地
	TargetAddress  string `json:"targetAddress"` // 目的地
	StartDate     string `json:"startDate"` // 运输时间
	EndDate     string `json:"endDate"` // 到达时间
}

type TraceRecord struct {
	ProduceID string `json:"produceID"` // 生产
	ProcessID string `json:"processID"` // 加工
	TransportID string `json:"transportID"` // 运输
}

type TraceHistory struct {
	TxID string `json:"txID"`
	TxValue string `json:"txValue"`
	TxTime string `json:"txTime"`
	TxStatus string `json:"txStatus"`
}

func NewProduceRecord(values []string) ProduceRecord {
	return ProduceRecord{
		Name: values[0],
		Producer: values[1],
		Address: values[2],
		Date: values[3],
		Life: values[4],
	}
}

func NewProcessRecord(values []string) ProcessRecord {
	processType, _ := getProcessType(values[1])
	return ProcessRecord{
		Name: values[0],
		Type: processType,
		Processor: values[2],
		Address: values[3],
		Date: values[4],
		Life: values[5],
	}
}

func NewTransportRecord(values []string) TransportRecord {
	return TransportRecord{
		Transporter: values[0],
		OriginAddress: values[1],
		TargetAddress: values[2],
		StartDate: values[3],
		EndDate: values[4],
	}
}

func NewTraceRecord(values []string) TraceRecord {
	return TraceRecord{
		ProduceID: values[0],
		ProcessID: values[1],
		TransportID: values[2],
	}
}

func NewTraceHistory(values []string) TraceHistory {
	return TraceHistory{
		TxID: values[0],
		TxValue: values[1],
		TxTime: values[2],
		TxStatus: values[3],
	}
}

func (r *TraceRecord) UpdateTraceRecordField(field, value string) error {
	switch field {
	case "produceID":
		r.ProduceID = value
	case "processID":
		r.ProcessID = value
	case "transportID":
		r.TransportID = value
	default:
		return fmt.Errorf("No field named %s error.", field)
	}
	return nil
}

func (r *TraceRecord) GetTraceRecordValue(field string) (string, error) {
	value := ""
	switch field {
	case "produceID":
		value = r.ProduceID
	case "processID":
		value = r.ProcessID
	case "transportID":
		value = r.TransportID
	default:
		return value, fmt.Errorf("No field named %s error.", field)
	}
	return value, nil
}

func (r *ProduceRecord) String() string {
	recordAsBytes, _ := json.Marshal(r)
	return string(recordAsBytes)
}

func (r *ProcessRecord) String() string {
	recordAsBytes, _ := json.Marshal(r)
	return string(recordAsBytes)
}

func (r *TransportRecord) String() string {
	recordAsBytes, _ := json.Marshal(r)
	return string(recordAsBytes)
}

func (r *TraceRecord) String() string {
	recordAsBytes, _ := json.Marshal(r)
	return string(recordAsBytes)
}

func (h *TraceHistory) String() string {
	historyAsBytes, _ := json.Marshal(h)
	return string(historyAsBytes)
}