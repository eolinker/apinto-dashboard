package v1

type JsonMarshalProxy []byte

func (j *JsonMarshalProxy) MarshalJSON() ([]byte, error) {
	return *j, nil
}
