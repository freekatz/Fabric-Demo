package structures

import (
	"encoding/json"
	"fmt"
)

type GenderType string

const (
	Male       GenderType = "male"   // 男性
	Female     GenderType = "female" // 女性
	NoneGender GenderType = ""
)

func getGender(genderStr string) (GenderType, error) {
	switch genderStr {
	case "male":
		return Male, nil
	case "female":
		return Female, nil
	default:
		return NoneGender, fmt.Errorf("No gender type named %s error.", genderStr)
	}
}

type Patient struct {
	Name        string     `json:"name"`        // 名称
	Gender      GenderType `json:"gender"`      // 性别
	Birth       string     `json:"birth"`       // 生日
	IdentifyID  string     `json:"identifyID"`  // 身份证号码
	PhoneNumber string     `json:"phoneNumber"` // 手机号码
	Address     string     `json:"address"`     // 当前住址
	NativePlace string     `json:"nativePlace"` // 籍贯
	CreditCard  string     `json:"creditCard"`  // 银行卡号
}

func (p *Patient) UpdatePatientField(field string, value interface{}) error {
	switch field {
	case "name":
		p.Name = value.(string)
	case "gender":
		genderType, err := getGender(value.(string))
		if err != nil {
			return err
		}
		p.Gender = genderType
	case "birth":
		p.Birth = value.(string)
	case "identifyID":
		p.IdentifyID = value.(string)
	case "phoneNumber":
		p.PhoneNumber = value.(string)
	case "address":
		p.Address = value.(string)
	case "nativePlace":
		p.NativePlace = value.(string)
	case "creditCard":
		p.CreditCard = value.(string)
	default:
		return fmt.Errorf("No field named %s error.", value.(string))
	}
	return nil
}

type PatientInHOS struct {
	Patient
	// 长度为 128，192 或 256
	HealthcareID string `json:"healthcareID"` // 医保卡号, 医院方的链码由病人登记, 医保局方自动生成
}

type PatientInHIB struct {
	Patient
}

func NewPatientInHOS(values []string) PatientInHOS {
	genderType, _ := getGender(values[1])
	return PatientInHOS{
		Patient: Patient{
			Name:        values[0],
			Gender:      genderType,
			Birth:       values[2],
			IdentifyID:  values[3],
			PhoneNumber: values[4],
			Address:     values[5],
			NativePlace: values[6],
			CreditCard:  values[7],
		},
		HealthcareID: values[8],
	}
}

func (p *PatientInHOS) String() string {
	patientAsBytes, _ := json.Marshal(p)
	return string(patientAsBytes)
}

func NewPatientInHIB(values []string) PatientInHIB {
	genderType, _ := getGender(values[1])
	return PatientInHIB{
		Patient: Patient{
			Name:        values[0],
			Gender:      genderType,
			Birth:       values[2],
			IdentifyID:  values[3],
			PhoneNumber: values[4],
			Address:     values[5],
			NativePlace: values[6],
			CreditCard:  values[7],
		},
	}
}

func (p *PatientInHIB) String() string {
	patientAsBytes, _ := json.Marshal(p)
	return string(patientAsBytes)
}
