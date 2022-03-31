package apinto_dashboard

type ZoneName string
const (
	ZhCn  ZoneName = "zh_cn"
	JaJp ZoneName ="ja_jp"
	EnUs ZoneName = "EN_US"
)

func (zone ZoneName)Read(from map[string]string)string  {
	return from[string(zone)]
}