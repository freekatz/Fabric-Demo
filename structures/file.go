package structures

import (
	"fmt"
)

type FileRecord struct {
	Timestamp   string
	IPFSAddress string
	Fingerprint string
}

func NewFileRecord(args []string) (*FileRecord, error) {
	if args == nil {
		return nil, fmt.Errorf("args is nil.")
	}
	if len(args) != 3 {
		return nil, fmt.Errorf("args len do not match.\n")
	}
	record := new(FileRecord)
	record.Timestamp = args[0]
	record.UpdateIPFSAddress(args[1]) // 检查合法性
	record.UpdateFingerprint(args[2]) // 检查合法性
	return record, nil
}

func (r *FileRecord) UpdateRecordField(field string, value string) error {
	switch field {
	case "address":
		return r.UpdateIPFSAddress(value) // 检查合法性
	case "fingerprint":
		return r.UpdateFingerprint(value) // 检查合法性
	default:
		return fmt.Errorf("No field named %s error.", field)
	}
	return nil
}

func (r *FileRecord) UpdateFingerprint(fp string) error {
	// todo 验证合法性
	r.Fingerprint = fp
	return nil
}

func (r *FileRecord) UpdateIPFSAddress(addr string) error {
	// todo 验证合法性
	r.IPFSAddress = addr
	return nil
}
